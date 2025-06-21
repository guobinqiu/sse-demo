package main

import (
	"fmt"
	"net/http"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	// 限制只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 设置响应头，启用 SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "close")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 获取请求上下文
	ctx := r.Context()

	// 模拟数据生成
	data := []string{"Hello", "World", "This is a stream!"}

	for _, msg := range data {
		select {
		case <-ctx.Done():
			fmt.Println("客户端取消连接")
			return
		default:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush() // 立即发送到客户端
			time.Sleep(1 * time.Second)
		}
	}

	// 发送结束信号
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func main() {
	http.HandleFunc("/stream", streamHandler)
	http.ListenAndServe(":8080", nil)
}
