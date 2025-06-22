package model

// 请求格式
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StreamRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// 响应格式
type StreamChunk struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Choices []ChunkItem `json:"choices"`
}

type ChunkItem struct {
	Delta        Delta   `json:"delta"`
	Index        int     `json:"index"`
	FinishReason *string `json:"finish_reason"` // nil 表示未结束
}

type Delta struct {
	Content string `json:"content,omitempty"`
}
