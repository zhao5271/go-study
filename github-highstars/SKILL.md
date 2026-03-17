---
name: github-highstars
description: Find and summarize high-star GitHub repositories (via GitHub Search API)
---

# github-highstars

Finds high-star GitHub repositories using the GitHub Search API, then returns a short, actionable shortlist (with links and key metadata).

## When to use

- User asks for “GitHub 高星项目 / top repos / trending alternatives”.
- User wants candidates to learn from (best-in-class repos in a language/topic).
- User wants a quick repo shortlist for evaluation or benchmarking.

## Instructions

1. Ask for constraints only if missing (topic/keywords, language, minimum stars, time window, count). Otherwise infer sensible defaults.
2. Run the local helper script to query GitHub:
   - `python3 scripts/gh_highstars.py --stars 50000 --limit 20`
   - Optional filters: `--language Go`, `--topic ai`, `--created-after 2024-01-01`, `--pushed-after 2025-01-01`
   - If the user provided a raw GitHub query, pass it via `--q`.
3. Present results as a Markdown table (default output) and then provide 3–5 bullet takeaways (what’s popular, why it’s relevant, which to start with).
4. If GitHub rate limits are hit, ask the user to set `GITHUB_TOKEN` and rerun.

## Notes

- Token: set `GITHUB_TOKEN` for higher rate limits (recommended).
- The GitHub Search API returns a maximum of 1000 results per query; this tool caps at 1000.

## Examples

- Top Go repos with ≥ 30k stars:
  - `python3 scripts/gh_highstars.py --language Go --stars 30000 --limit 30`
- LLM/agent repos created recently:
  - `python3 scripts/gh_highstars.py --keyword llm --keyword agent --created-after 2024-01-01 --stars 2000 --limit 50`
- Raw query mode (exact GitHub search syntax):
  - `python3 scripts/gh_highstars.py --q "topic:ai language:Python stars:>=5000 pushed:>=2025-01-01" --limit 25`
