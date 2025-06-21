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
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	cancel()
	// }()

	data := map[string]any{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "user", "content": "Hi"},
		},
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/stream", bytes.NewBuffer(jsonBytes))
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
			if content == "[DONE]" {
				fmt.Println("收到结束标志，退出读取")
				break
			}
			fmt.Println(content)
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
