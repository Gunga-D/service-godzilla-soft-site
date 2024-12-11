package yandex_mail

import (
	"crypto/tls"
	"net"
	"net/smtp"
)

type client struct {
	addr           string
	yandexLogin    string
	yandexPassword string
}

func NewClient(addr string, yandexLogin string, yandexPassword string) *client {
	return &client{
		addr:           addr,
		yandexLogin:    yandexLogin,
		yandexPassword: yandexPassword,
	}
}

func (c *client) SendMail(to []string, subject string, body string) error {
	conn, host, err := c.createTLSConn(c.addr)
	if err != nil {
		return err
	}

	smtpClient, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer smtpClient.Quit()
	err = smtpClient.Auth(smtp.PlainAuth("", c.yandexLogin, c.yandexPassword, host))
	if err != nil {
		return err
	}

	payload := email{
		sender:  c.yandexLogin,
		to:      to,
		subject: subject,
		body:    body,
	}

	err = smtpClient.Mail(c.yandexLogin)
	if err != nil {
		return err
	}

	for _, addr := range payload.to {
		err = smtpClient.Rcpt(addr)
		if err != nil {
			return err
		}
	}

	writer, err := smtpClient.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write(payload.BuildMessage())
	if err != nil {
		return err
	}
	return nil
}

func (c *client) createTLSConn(addr string) (*tls.Conn, string, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, "", err
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, config)
	return conn, host, err
}
