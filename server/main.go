package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// 模拟OpenAI返回的结构体
type ChoiceDelta struct {
	Content string `json:"content,omitempty"`
}

type Choice struct {
	Delta        ChoiceDelta `json:"delta"`
	Index        int         `json:"index"`
	FinishReason *string     `json:"finish_reason"` // 用指针区分null
}

type OpenAIStreamResp struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
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
	data := []string{"Hi", "there", "!", "How", "can", "I", "help", "you", "today", "?", "😊"}

	for _, msg := range data {
		select {
		case <-ctx.Done():
			fmt.Println("客户端取消连接")
			return
		default:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush() // 立即发送到客户端
			time.Sleep(300 * time.Millisecond)
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
