package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/guobinqiu/sse-demo/model"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data := model.StreamRequest{
		Model: "gpt-4",
		Messages: []model.Message{
			{Role: "user", Content: "Hi"},
		},
		Stream: true,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/stream", bytes.NewReader(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer resp.Body.Close()

	// 逐行读取 HTTP 响应体中的内容
	scanner := bufio.NewScanner(resp.Body)

	// 当 ctx 被取消导致底层连接关闭或读取中断时，scanner.Scan() 会返回 false，从而自动跳出循环
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 6 && line[:6] == "data: " {
			content := line[6:]

			var chunk model.StreamChunk
			if err := json.Unmarshal([]byte(content), &chunk); err != nil {
				fmt.Println("json解析错误:", err)
				continue
			}

			for _, choice := range chunk.Choices {
				if choice.FinishReason != nil && *choice.FinishReason == "stop" {
					fmt.Println("收到结束标志，退出读取")
					break
				}

				fmt.Println(choice.Delta.Content)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Println("请求被取消了")
		} else if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("请求超时了")
		} else {
			fmt.Println("读取错误:", err)
		}
	} else {
		fmt.Println("读取完毕")
	}
}
