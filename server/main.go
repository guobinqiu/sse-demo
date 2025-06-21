package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StreamRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	// 限制只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// 解析 JSON 请求体
	var reqData StreamRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("收到请求:", reqData.Messages[0].Content)

	// 处理客户端断连
	ctx := r.Context()

	// 模拟数据生成
	words := strings.Split("Hi there! How can I help you today? 😊", " ")

	for _, word := range words {
		select {
		case <-ctx.Done():
			fmt.Println("客户端取消连接")
			return
		default:
			fmt.Fprintf(w, "data: %s\n\n", word)
			flusher.Flush() // 立即发送到客户端
			time.Sleep(200 * time.Millisecond)
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
