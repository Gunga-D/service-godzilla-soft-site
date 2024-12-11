package yandex_mail

type Client interface {
	SendMail(to []string, subject string, body string) error
}
