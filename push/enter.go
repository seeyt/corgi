package push

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

var Limit time.Duration = 3

type Pusher interface {
	Send(Common)
	SendByLimit(Common)
	Run()
	AddBeforeSendFunc(f func(Common) Common)
}

type Push struct {
	message       chan Common
	activeMessage chan Common
	beforeFunc    func(Common) Common
}

// Send 无发送限制
func (p *Push) Send(msg Common) {
	p.message <- msg
}

// SendByLimit 3秒钟发送一条
func (p *Push) SendByLimit(msg Common) {
	p.activeMessage <- msg
}

func (p *Push) Run() {
	go func() {
		for mes := range p.message {
			go p.sendFunc(mes)
		}
	}()
	go func() {
		for mes := range p.activeMessage {
			go p.sendFunc(mes)
			time.Sleep(Limit * time.Second)
		}
	}()
}

func (p *Push) AddBeforeSendFunc(f func(Common) Common) {
	p.beforeFunc = f
}

func (p *Push) sendFunc(msg Common) {
	common := msg
	if p.beforeFunc != nil {
		common = p.beforeFunc(common)
	}
	message := msg.Content()
	client := resty.New().R()
	client.SetHeaders(message.Header)
	client.SetBody(message.Body)
	result := &Result{}
	res, _ := client.SetResult(result).Post(message.Webhook)
	if result.Errmsg != "ok" {
		fmt.Println(res)
	}
}

func New() *Push {
	return &Push{
		message:       make(chan Common, 6),
		activeMessage: make(chan Common, 6),
	}
}
