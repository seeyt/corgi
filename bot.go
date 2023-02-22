package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/r2dygo/corgi/push"
	"github.com/r2dygo/corgi/push/ding"
	"github.com/r2dygo/corgi/schedule"
	"io"
	"log"
	"net/http"
	"strings"
)

type KeywordItem struct {
	pattern string
	fn      func(ctx *Context)
}

type Bot struct {
	Messages    chan *Context
	commandFns  map[string]func(*Context)
	messageFns  []func(*Context)
	KeywordFns  []*KeywordItem
	middleWares []func(*Context)
	pusher      push.Pusher
	scheduler   schedule.Scheduler
}

func (b *Bot) Use(f ...func(*Context)) {
	b.middleWares = append(b.middleWares, f...)
}

func (b *Bot) BeforeSendFunc(f func(push.Common) push.Common) {
	b.pusher.AddBeforeSendFunc(f)
}

func (b *Bot) Send(content push.Common) {
	b.pusher.Send(content)
}

// SendByLimit 3秒钟发送一条，防止频繁请求
func (b *Bot) SendByLimit(content push.Common) {
	b.pusher.SendByLimit(content)
}

func (b *Bot) SetPush(messageManager push.Pusher) {
	b.pusher = messageManager
}

func (b *Bot) SetScheduler(schedule schedule.Scheduler) {
	b.scheduler = schedule
}

func (b *Bot) AddJob(taskName, spec string, f func(bot *Bot) func()) error {
	if b.scheduler == nil {
		return errors.New("请添加scheduler")
	}
	b.scheduler.AddJob(taskName, spec, f(b))
	return nil
}

func (b *Bot) OnMessage(f ...func(ctx *Context)) {
	b.messageFns = append(b.messageFns, f...)
}

func (b *Bot) onMessagesHandle(c *Context) {
	for _, fn := range b.messageFns {
		fn(c)
	}
}

func (b *Bot) push(msg *Context) {
	b.Messages <- msg
}

func (b *Bot) handle() {
	for msg := range b.Messages {
		msg.Run()
	}
}

func (b *Bot) Run(addr, pattern string) {
	go b.pusher.Run()
	go b.scheduler.Run()
	go b.handle()
	http.HandleFunc(pattern, func(ResponseWriter http.ResponseWriter, Request *http.Request) {
		body, _ := io.ReadAll(Request.Body)
		msg := ding.Ding{}
		_ = json.Unmarshal(body, &msg)
		content := strings.TrimSpace(msg.Text.Content)
		ctx := &Context{
			pusher:           b.pusher,
			Message:          msg,
			Content:          content,
			Webhook:          msg.SessionWebhook,
			ConversationType: msg.ConversationType,
			IsAdmin:          msg.IsAdmin,
			SenderStaffId:    msg.SenderStaffId,
			handlers:         append(b.middleWares, b.onMessagesHandle),
			ResponseWriter:   ResponseWriter,
			Request:          Request,
			ConversationId:   msg.ConversationId,
		}
		b.push(ctx)
	})
	fmt.Printf("runing in %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func New() *Bot {
	botManage := &Bot{
		Messages:    make(chan *Context, 10),
		commandFns:  map[string]func(ctx *Context){},
		messageFns:  []func(ctx *Context){},
		middleWares: []func(*Context){},
		pusher:      push.New(),
		KeywordFns:  []*KeywordItem{},
	}
	return botManage
}
