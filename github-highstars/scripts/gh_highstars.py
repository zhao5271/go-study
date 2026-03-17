#!/usr/bin/env python3
import argparse
import csv
import datetime as dt
import json
import os
import shutil
import ssl
import subprocess
import sys
import urllib.parse
import urllib.request


GITHUB_API = "https://api.github.com"


def _parse_date(value: str) -> str:
    try:
        dt.date.fromisoformat(value)
    except ValueError as e:
        raise argparse.ArgumentTypeError(
            f"Invalid date '{value}'. Expected YYYY-MM-DD."
        ) from e
    return value


def _build_query(args: argparse.Namespace) -> str:
    if args.raw_query:
        return args.raw_query.strip()

    tokens: list[str] = []

    if args.keywords:
        tokens.extend([k.strip() for k in args.keywords if k.strip()])

    if args.topic:
        tokens.append(f"topic:{args.topic}")

    if args.language:
        tokens.append(f"language:{args.language}")

    if args.stars_gte is not None:
        tokens.append(f"stars:>={args.stars_gte}")

    if args.created_gte:
        tokens.append(f"created:>={args.created_gte}")

    if args.pushed_gte:
        tokens.append(f"pushed:>={args.pushed_gte}")

    if args.archived is not None:
        tokens.append(f"archived:{'true' if args.archived else 'false'}")

    if args.fork is not None:
        tokens.append(f"fork:{'true' if args.fork else 'false'}")

    q = " ".join(tokens).strip()
    if not q:
        # A safe default that matches "high star projects".
        q = "stars:>=50000"
    return q


def _request_via_urllib(url: str, token: str | None, user_agent: str) -> tuple[dict, dict]:
    headers = {
        "Accept": "application/vnd.github+json",
        "User-Agent": user_agent,
    }
    if token:
        headers["Authorization"] = f"Bearer {token}"

    req = urllib.request.Request(url, headers=headers)
    try:
        # Use a fresh default context; some Python installs can have a broken
        # cert bundle, which we handle by falling back to curl below.
        ctx = ssl.create_default_context()
        with urllib.request.urlopen(req, timeout=30, context=ctx) as resp:
            data = resp.read().decode("utf-8")
            return json.loads(data), dict(resp.headers)
    except urllib.error.HTTPError as e:
        body = ""
        try:
            body = e.read().decode("utf-8")
        except Exception:
            pass
        msg = f"GitHub API error: HTTP {e.code} {e.reason}"
        if body:
            msg += f"\n{body}"
        raise RuntimeError(msg) from e


def _request_via_curl(url: str, token: str | None, user_agent: str) -> tuple[dict, dict]:
    curl = shutil.which("curl")
    if not curl:
        raise RuntimeError("curl not found, and urllib HTTPS failed.")

    cmd = [
        curl,
        "-fsSL",
        "-D",
        "-",  # dump headers to stdout
        "-o",
        "-",  # write body to stdout
        "-H",
        "Accept: application/vnd.github+json",
        "-H",
        f"User-Agent: {user_agent}",
    ]
    if token:
        cmd.extend(["-H", f"Authorization: Bearer {token}"])
    cmd.append(url)

    proc = subprocess.run(cmd, capture_output=True, text=True)
    if proc.returncode != 0:
        stderr = (proc.stderr or "").strip()
        raise RuntimeError(f"curl request failed (exit {proc.returncode}). {stderr}")

    # curl output is: headers + blank line + body.
    # There can be multiple header blocks (e.g., redirects); take the last.
    parts = proc.stdout.split("\r\n\r\n")
    if len(parts) < 2:
        parts = proc.stdout.split("\n\n")
    if len(parts) < 2:
        raise RuntimeError("Unexpected curl output (no header/body separator).")

    body = parts[-1]
    header_blob = parts[-2]
    headers: dict[str, str] = {}
    for line in header_blob.splitlines():
        if ":" not in line:
            continue
        k, v = line.split(":", 1)
        headers[k.strip()] = v.strip()

    return json.loads(body), headers


def _request(url: str, token: str | None, user_agent: str) -> tuple[dict, dict]:
    try:
        return _request_via_urllib(url, token=token, user_agent=user_agent)
    except urllib.error.URLError as e:
        # Common in locked-down / misconfigured environments:
        # - DNS/network issues
        # - SSL cert verification errors (missing CA bundle)
        if isinstance(getattr(e, "reason", None), ssl.SSLCertVerificationError):
            return _request_via_curl(url, token=token, user_agent=user_agent)
        raise


def _search_repositories(
    *,
    query: str,
    sort: str,
    order: str,
    limit: int,
    per_page: int,
    token: str | None,
    user_agent: str,
) -> tuple[list[dict], dict]:
    results: list[dict] = []
    headers_last: dict = {}

    # GitHub Search API returns at most 1000 results.
    max_limit = min(limit, 1000)
    pages = (max_limit + per_page - 1) // per_page

    for page in range(1, pages + 1):
        params = {
            "q": query,
            "sort": sort,
            "order": order,
            "per_page": str(per_page),
            "page": str(page),
        }
        url = f"{GITHUB_API}/search/repositories?{urllib.parse.urlencode(params)}"
        payload, headers_last = _request(url, token=token, user_agent=user_agent)
        items = payload.get("items", [])
        for item in items:
            results.append(item)
            if len(results) >= max_limit:
                return results, headers_last

        if not items:
            break

    return results, headers_last


def _repo_row(item: dict) -> dict:
    license_info = item.get("license") or {}
    return {
        "full_name": item.get("full_name") or "",
        "html_url": item.get("html_url") or "",
        "stargazers_count": item.get("stargazers_count") or 0,
        "forks_count": item.get("forks_count") or 0,
        "open_issues_count": item.get("open_issues_count") or 0,
        "language": item.get("language") or "",
        "license": license_info.get("spdx_id") or "",
        "updated_at": item.get("updated_at") or "",
        "description": (item.get("description") or "").replace("\n", " ").strip(),
    }


def _print_markdown(rows: list[dict]) -> None:
    headers = [
        "repo",
        "stars",
        "lang",
        "license",
        "updated",
        "description",
    ]
    print("| " + " | ".join(headers) + " |")
    print("| " + " | ".join(["---"] * len(headers)) + " |")
    for r in rows:
        repo_md = f"[{r['full_name']}]({r['html_url']})"
        desc = r["description"]
        if len(desc) > 120:
            desc = desc[:117] + "..."
        print(
            "| "
            + " | ".join(
                [
                    repo_md,
                    str(r["stargazers_count"]),
                    r["language"] or "",
                    r["license"] or "",
                    (r["updated_at"] or "")[:10],
                    desc.replace("|", "\\|"),
                ]
            )
            + " |"
        )


def _print_tsv(rows: list[dict]) -> None:
    fieldnames = [
        "full_name",
        "html_url",
        "stargazers_count",
        "forks_count",
        "open_issues_count",
        "language",
        "license",
        "updated_at",
        "description",
    ]
    writer = csv.DictWriter(sys.stdout, fieldnames=fieldnames, dialect="excel-tab")
    writer.writeheader()
    for r in rows:
        writer.writerow(r)


def main() -> int:
    parser = argparse.ArgumentParser(
        prog="gh-highstars",
        description="Query GitHub for high-star repositories via the Search API.",
    )
    parser.add_argument(
        "--token",
        help="GitHub token. Defaults to env GITHUB_TOKEN.",
        default=None,
    )
    parser.add_argument(
        "--user-agent",
        default="github-highstars-skill",
        help="User-Agent header for GitHub API requests.",
    )

    q = parser.add_argument_group("Query")
    q.add_argument(
        "--q",
        dest="raw_query",
        help="Raw GitHub search query (overrides other query flags).",
        default=None,
    )
    q.add_argument(
        "--keyword",
        dest="keywords",
        action="append",
        help="Keyword (repeatable). Example: --keyword llm --keyword agent",
        default=[],
    )
    q.add_argument("--topic", help="GitHub topic. Example: ai", default=None)
    q.add_argument("--language", help="Language. Example: Go", default=None)
    q.add_argument(
        "--stars",
        dest="stars_gte",
        type=int,
        default=None,
        help="Minimum stars (>=). Example: 50000",
    )
    q.add_argument("--created-after", dest="created_gte", type=_parse_date, default=None)
    q.add_argument("--pushed-after", dest="pushed_gte", type=_parse_date, default=None)
    q.add_argument(
        "--archived",
        type=lambda s: {"true": True, "false": False}[s.lower()],
        choices=[True, False],
        help="Filter archived repos: true/false",
        default=None,
    )
    q.add_argument(
        "--fork",
        type=lambda s: {"true": True, "false": False}[s.lower()],
        choices=[True, False],
        help="Include forks: true/false",
        default=None,
    )

    o = parser.add_argument_group("Output")
    o.add_argument(
        "--format",
        choices=["markdown", "json", "tsv"],
        default="markdown",
        help="Output format.",
    )

    p = parser.add_argument_group("Paging")
    p.add_argument("--limit", type=int, default=20, help="Max results (<= 1000).")
    p.add_argument("--per-page", type=int, default=50, help="Results per page (<= 100).")
    p.add_argument("--sort", default="stars", choices=["stars", "forks", "updated"])
    p.add_argument("--order", default="desc", choices=["desc", "asc"])

    args = parser.parse_args()

    token = args.token or os.environ.get("GITHUB_TOKEN")
    query = _build_query(args)
    limit = max(1, min(int(args.limit), 1000))
    per_page = max(1, min(int(args.per_page), 100))

    items, headers = _search_repositories(
        query=query,
        sort=args.sort,
        order=args.order,
        limit=limit,
        per_page=per_page,
        token=token,
        user_agent=args.user_agent,
    )
    rows = [_repo_row(i) for i in items]

    if args.format == "json":
        json.dump(
            {
                "query": query,
                "count": len(rows),
                "items": rows,
                "rate_limit": {
                    "limit": headers.get("X-RateLimit-Limit"),
                    "remaining": headers.get("X-RateLimit-Remaining"),
                    "reset": headers.get("X-RateLimit-Reset"),
                },
            },
            sys.stdout,
            ensure_ascii=False,
            indent=2,
        )
        print()
    elif args.format == "tsv":
        _print_tsv(rows)
    else:
        _print_markdown(rows)
        remaining = headers.get("X-RateLimit-Remaining")
        reset = headers.get("X-RateLimit-Reset")
        if remaining is not None and reset is not None:
            try:
                reset_dt = dt.datetime.fromtimestamp(int(reset)).astimezone()
                print(f"\nRate limit remaining: {remaining} (resets at {reset_dt:%Y-%m-%d %H:%M:%S %Z})")
            except Exception:
                pass

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
