# SSE Demo

Explore how OpenAI responds with streaming

## standard-sse

测标准sse

```
go run server/standard-sse/main.go
http://localhost:8080/
```

## standard-sse-retrylimit

测限制重试次数的标准sse

第一步:

```
cd client/standard-sse-retrylimit && python3 -m http.server 8081
http://localhost:8081/
```
F12查看Console

第二步:

in another terminal

```
go run server/standard-sse-retrylimit/main.go
```
再次刷新浏览器

## openai-sse

测类openai的非标sse

### 测正常

```
go test -v ./server/openai-sse -run '^TestOpenaiSSE$'
```

### 测异常

#### 取消

```
go test -v ./server/openai-sse -run '^TestOpenaiSSEReqCancel$'
```

#### 超时

```
go test -v ./server/openai-sse -run '^TestOpenaiSSEReqTimeout$'
```

#### 超时

```
go test -v ./server/openai-sse -run '^TestOpenaiSSEReqTimeout$'
```

## 概念对比

### SSE vs WebSocket

| 特性              | SSE (Server-Sent Events)                                 | WebSocket                                  |
| ----------------- | -------------------------------------------------------- | ------------------------------------------ |
| **协议**          | 基于 HTTP 协议（使用持久的 HTTP 连接）                   | 独立的协议，先握手再升级为双向 TCP 连接    |
| **通信方向**      | 单向（服务器 → 客户端）                                  | 双向（服务器 ↔ 客户端）                    |
| **连接建立方式**  | 客户端发起 HTTP GET 请求，服务器响应 `text/event-stream` | 双向握手，使用 `ws://` 或 `wss://` 协议    |
| **数据格式**      | 纯文本，符合 SSE 格式（`data: ...\n\n`）                 | 二进制或文本，自定义数据格式               |
| **支持浏览器**    | 大多数现代浏览器内置支持，使用 `EventSource` API         | 所有现代浏览器均支持，使用 `WebSocket` API |
| **消息推送机制**  | 服务器推送，客户端监听事件                               | 双向消息收发，实时交互                     |
| **连接重试机制**  | 浏览器自动重试，服务端可指定重试间隔 (`retry`)           | 需要应用层自行实现重连逻辑                 |
| **事件类型支持**  | 支持自定义事件类型，客户端可以监听不同事件               | 无内建事件类型，需自行设计消息结构和类型   |
| **HTTP/2 兼容性** | 支持                                                     | 也支持，但依赖具体实现                     |
| **跨域支持**      | 通过 CORS 机制支持                                       | 通过 CORS 和 WebSocket 协议支持            |
| **使用场景**      | 简单的实时推送，如新闻更新、股票价格、通知               | 复杂实时双向交互，如聊天、游戏、协作编辑   |
| **服务器负载**    | 轻量级，适合大量客户端连接                               | 较重，因为双向通信和持久连接               |
| **安全性**        | 支持 HTTPS                                               | 支持 WSS（加密 WebSocket）                 |
| **连接数限制**    | 受浏览器限制，一般每域名6-10个并发连接                   | 受浏览器限制，但比 SSE 更灵活              |
| **协议复杂度**    | 简单，基于文本，易于调试                                 | 较复杂，需处理二进制帧和状态管理           |

### 标准 SSE vs OpenAI SSE

| 方面                 | 标准 SSE                               | OpenAI SSE (OpenAI API 的流式 SSE)                      |
| -------------------- | -------------------------------------- | ------------------------------------------------------- |
| **协议基础**         | 完全符合 W3C SSE 规范                  | 基于 SSE，但格式上有定制                                |
| **连接方式**         | HTTP GET 请求建立持久连接              | HTTP POST 请求，带有请求体（JSON）                      |
| **数据格式**         | `data: <文本内容>\n\n`                 | `data: {...json...}\n\n`，每条是 JSON 片段              |
| **事件格式**         | 文本数据，支持自定义事件名             | 通常只有 `data:` 字段，内容是 JSON 对象                 |
| **请求方法**         | 只用 GET 请求                          | 用 POST 请求发送请求参数（如 model, messages）          |
| **传输内容**         | 纯文本或简单文本消息                   | 包含结构化 JSON 数据，分片传输模型生成的文本、tokens 等 |
| **流结束标识**       | 一般服务器关闭连接或发送自定义结束事件 | `data: [DONE]` 明确告诉客户端结束                       |
| **重试机制**         | 浏览器自动重试，服务器可指定重试时间   | 同标准 SSE，浏览器自动处理重连                          |
| **客户端处理**       | 浏览器用 EventSource API 解析文本      | 客户端需要解析每条 `data:` 的 JSON 内容                 |
| **用途**             | 适合简单单向消息推送（如新闻、通知）   | 用于传输大模型输出的流式响应                            |
| **服务端实现复杂度** | 简单                                   | 需要在后台边生成边编码 JSON 数据，实时推送              |
| **请求体携带参数**   | 不支持（GET 请求无请求体）             | 支持（POST 请求带 JSON 请求体）                         |
