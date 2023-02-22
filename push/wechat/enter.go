package wechat

import (
	"github.com/r2dygo/corgi/push"
)

type Text struct {
	Webhook, AtUserIds, Text string
	MentionedList            []string
	MentionedMobileList      []string
}

func (t Text) Content() push.CommonMessage {
	body := push.H{
		"msgtype": "text",
		"text": push.H{
			"content":               t.Text,
			"mentioned_list":        t.MentionedList,
			"mentioned_mobile_list": t.MentionedMobileList,
		},
	}
	url := t.Webhook
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return push.CommonMessage{
		Webhook: url,
		Body:    body,
		Header:  header,
	}
}

// Markdown markdown消息
type Markdown struct {
	Webhook, Text string
}

func (t Markdown) Content() push.CommonMessage {
	body := push.H{
		"msgtype": "markdown",
		"markdown": push.H{
			"content": t.Text,
		},
	}
	header := map[string]string{
		"Content-Type": "application/json",
	}
	url := t.Webhook
	return push.CommonMessage{
		Webhook: url,
		Body:    body,
		Header:  header,
	}
}

// Image 图片消息
type Image struct {
	Webhook, Base64, Md5 string
}

func (t Image) Content() push.CommonMessage {
	body := push.H{
		"msgtype": "image",
		"image": push.H{
			"base64": t.Base64,
			"md5":    t.Md5,
		},
	}
	url := t.Webhook
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return push.CommonMessage{
		Webhook: url,
		Body:    body,
		Header:  header,
	}
}

// News 图文 消息
type News struct {
	Webhook, Title, Description string
	Url, PicUrl                 string
}

func (t News) Content() push.CommonMessage {
	body := push.H{
		"msgtype": "news",
		"news": push.H{
			"articles": []push.H{
				{
					"title":       t.Title,
					"description": t.Description,
					"url":         t.Url,
					"picurl":      t.PicUrl,
				},
			},
		},
	}
	url := t.Webhook
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return push.CommonMessage{
		Webhook: url,
		Body:    body,
		Header:  header,
	}
}
