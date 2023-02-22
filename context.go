package bot

import (
	"github.com/r2dygo/corgi/push"
	"github.com/r2dygo/corgi/push/ding"
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 >> 1

type Context struct {
	index   int8
	pusher  push.Pusher
	Message ding.Ding
	// 消息内容
	Content string
	Webhook string
	// 1：单聊
	// 2：群聊
	ConversationType string
	IsAdmin          bool
	// 发送者ID
	SenderStaffId  string
	handlers       []func(*Context)
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	ConversationId string
}

func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) addHandler(f func(*Context)) {
	c.handlers = append(c.handlers, f)
}

func (c *Context) Run() {
	if len(c.handlers) == 0 {
		return
	}
	c.handlers[0](c)
	c.Next()
}

func (c *Context) SendText(text string) {
	content := ding.Text{
		Text:      text,
		Webhook:   c.Webhook,
		AtUserIds: c.SenderStaffId,
	}
	c.pusher.Send(content)
}

func (c *Context) SendMarkDown(title, text string) {
	content := ding.Markdown{
		Title:     title,
		Text:      text,
		Webhook:   c.Webhook,
		AtUserIds: c.SenderStaffId,
	}
	c.pusher.Send(content)
}

func (c *Context) SendActionCard(title, text, singleTitle, singleURL string) {
	content := ding.ActionCard{
		Title:       title,
		SingleURL:   singleURL,
		SingleTitle: singleTitle,
		Text:        text,
		Webhook:     c.Webhook,
		AtUserIds:   c.SenderStaffId,
	}
	c.pusher.Send(content)
}

func (c *Context) Send(msg push.Common) {
	c.pusher.Send(msg)
}

func newContext() *Context {
	return &Context{}
}
