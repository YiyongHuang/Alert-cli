package utils

import (
	"fmt"
	"time"
)

const (
	contextType     = "markdown"
	env             = "dev"
	sendAlias       = "研发"
	tosType         = 0
	applicationName = "cluster-monitor"
	expiredTime     = 5
)

type MessageBody struct {
	ApplicationName string   `json:"applicationName"`
	Env             string   `json:"env"`
	ScheduleTime    string   `json:"scheduleTime"`
	ExpireTime      string   `json:"expireTime"`
	SendAlias       string   `json:"sendAlias"`
	ContextType     string   `json:"contextType"`
	Content         string   `json:"content"`
	Tos             []string `json:"tos"`
	TosType         int32    `json:"tosType"`
	ToParties       []string `json:"toParties"`
}

func (msg *MessageBody) BuildMessageBody(content string, tos, toParties []string) error {
	if len(tos) == 0 && len(toParties) == 0 {
		return fmt.Errorf("no member or party to send message")
	}
	msg.ApplicationName = applicationName
	msg.Env = env
	msg.ScheduleTime = time.Now().Format("2006-01-02 15:04:05")
	msg.ExpireTime = time.Now().Add(time.Minute * expiredTime).Format("2006-01-02 15:04:05")
	msg.SendAlias = sendAlias
	msg.ContextType = contextType
	msg.Content = content
	msg.SendAlias = sendAlias
	msg.TosType = tosType
	if len(tos) > 0 {
		msg.Tos = tos
	}
	if len(toParties) > 0 {
		msg.ToParties = toParties
	}
	return nil
}
