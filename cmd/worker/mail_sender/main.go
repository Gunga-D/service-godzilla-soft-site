package main

import (
	"log"
	"os"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
)

func main() {
	cl := yandex_mail.NewClient(os.Getenv("YANDEX_MAIL_ADDRESS"), os.Getenv("YANDEX_MAIL_LOGIN"), os.Getenv("YANDEX_MAIL_PASSWORD"))
	err := cl.SendMail([]string{
		"dondokov.gunga@mail.ru",
	}, "Test", `
		<html>
			<body>
				<h3>Name:</h3><span>Василий</span><br/><br/>
				<h3>Email:</h3><span>Сообщение</span><br/>
			</body>
		</html>
	`)
	if err != nil {
		log.Fatalln(err)
	}
}
