package push

type CommonMessage struct {
	Webhook string
	Body    interface{}
	Header  map[string]string
}

// Common 通用消息接口
type Common interface {
	Content() CommonMessage
}

// Result 发送消息请求结果
type Result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type H map[string]interface{}
