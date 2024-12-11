package yandex_mail

import (
	"fmt"
	"strings"
)

type email struct {
	sender  string
	to      []string
	subject string
	body    string
}

// Генерирует тело сообщения
func (email *email) BuildMessage() []byte {
	enter := "\r\n"
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	message := ""
	message += fmt.Sprintf("From: %s%s", email.sender, enter)
	if len(email.to) > 0 {
		message += fmt.Sprintf("To: %s%s", strings.Join(email.to, ";"), enter)
	}

	message += fmt.Sprintf("Subject: %s%s", email.subject, enter)
	message += mimeHeaders
	message += enter + email.body
	message += enter + "."

	return []byte(message)
}
