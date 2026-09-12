package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type Sender interface {
	Configured() bool
	SendAdminCode(context.Context, string, string, time.Duration) error
}

type SMTP struct {
	host, username, password, from string
	port                           int
}

func NewSMTP(host string, port int, username, password, from string) *SMTP {
	return &SMTP{host: strings.TrimSpace(host), port: port, username: username, password: password, from: strings.TrimSpace(from)}
}

func (s *SMTP) Configured() bool { return s.host != "" && s.port > 0 && s.from != "" }

func (s *SMTP) SendAdminCode(_ context.Context, to, code string, validFor time.Duration) error {
	if !s.Configured() {
		return fmt.Errorf("SMTP is not configured")
	}
	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}
	subject := "Agora BBS 管理员登录验证码"
	body := fmt.Sprintf("你的管理员登录验证码是：%s\r\n\r\n验证码 %d 分钟内有效。如非本人操作，请忽略本邮件。", code, int(validFor.Minutes()))
	message := []byte("From: " + s.from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + body)
	return smtp.SendMail(s.host+":"+strconv.Itoa(s.port), auth, s.from, []string{to}, message)
}
