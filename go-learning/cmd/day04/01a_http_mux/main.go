package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	// 健康检查：给部署/容器编排做 readiness/liveness 用
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // Output: HTTP 405 "method not allowed\n"
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "ok") // Output: ok
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :8080 (输出可能变化/不固定：包含时间戳)
	err := http.ListenAndServe(addr, mux)
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :8080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
