package ding

import "github.com/r2dygo/corgi/push"

// Ding 钉钉消息
type Ding struct {
	ConversationId string `json:"conversationId"`
	SenderStaffId  string `json:"senderStaffId"`
	AtUsers        []struct {
		DingtalkId string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatbotUserId             string `json:"chatbotUserId"`
	MsgId                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	ConversationType          string `json:"conversationType"`
	SenderId                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      struct {
		Content string `json:"content"`
	} `json:"text"`
	RobotCode string `json:"robotCode"`
	Msgtype   string `json:"msgtype"`
}

// Text 文本消息
type Text struct {
	Webhook, AtUserIds, Text string
}

func (t Text) Content() push.CommonMessage {
	body := push.H{
		"msgtype": "text",
		"text": push.H{
			"content": t.Text,
		},
		"at": push.H{
			"isAtAll":   "False",
			"atUserIds": []string{t.AtUserIds},
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
	Webhook, AtUserIds, Text string
	Title                    string
}

func (t Markdown) Content() push.CommonMessage {
	body := push.H{
		"msgtype": "markdown",
		"markdown": push.H{
			"title": t.Title,
			"text":  t.Text,
		},
		"at": push.H{
			"isAtAll":   "False",
			"atUserIds": []string{t.AtUserIds},
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

// ActionCard actionCard消息
type ActionCard struct {
	Webhook, AtUserIds, Text      string
	Title, SingleTitle, SingleURL string
}

func (t ActionCard) Content() push.CommonMessage {
	body := push.H{
		"msgtype": "actionCard",
		"actionCard": push.H{
			"title":       t.Title,
			"text":        t.Text,
			"singleTitle": t.SingleTitle,
			"singleURL":   t.SingleURL,
		},
		"at": push.H{
			"isAtAll":   "False",
			"atUserIds": []string{t.AtUserIds},
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
