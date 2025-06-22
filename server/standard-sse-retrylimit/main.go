package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头，启用 SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 告诉客户端断线后重连时间
	// 浏览器内置的 EventSource 会自动解析服务器发送的 retry: 行
	fmt.Fprintf(w, "retry: 3000\n\n")
	flusher.Flush()

	// 模拟数据生成
	words := strings.Split("Hi there! How can I help you today? 😊", " ")

	for _, word := range words {
		fmt.Fprintf(w, "data: %s\n\n", word)
		flusher.Flush() // 立即发送到客户端
		time.Sleep(200 * time.Millisecond)
	}

	// 发送结束信号
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func main() {
	http.HandleFunc("/stream", streamHandler)
	http.ListenAndServe(":8080", nil)
}
