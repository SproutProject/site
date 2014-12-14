package main

import (
    "net/smtp"
    "text/template"
    "fmt"
    "bytes"
)

func MailVerify(target string,code string) error {
    templ := template.Must(template.New("Verify").Parse("From: {{.From}}\r\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n{{.Message}}"))
    param := struct {
	From string
	To string
	Subject string
	Message string
    }{
	"sprout@csie.ntu.edu.tw",
	target,
	"資訊之芽驗證信",
	fmt.Sprintf(
	    "<html><body>親愛的資訊之芽報名者：<p>歡迎報名 2015 資訊之芽，您的驗證碼是: %s</p></body></html>",
	    code,
	),
    }
    buf := new(bytes.Buffer)
    templ.Execute(buf,&param)

    err := smtp.SendMail(
	"smtp.csie.ntu.edu.tw:25",
	smtp.PlainAuth(
	    "",
	    Mail_User,
	    Mail_Passwd,
	    "smtp.csie.ntu.edu.tw"),
	"sprout@csie.ntu.edu.tw",
	[]string{target},
	buf.Bytes(),
    )

    return err
}
func MailDone(target string) error {
    templ := template.Must(template.New("Verify").Parse("From: {{.From}}\r\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n{{.Message}}"))
    param := struct {
	From string
	To string
	Subject string
	Message string
    }{
	"sprout@csie.ntu.edu.tw",
	target,
	"資訊之芽驗證信",
	"<html><body>親愛的資訊之芽報名者：<p>您的報名已成功</p></body></html>",
    }
    buf := new(bytes.Buffer)
    templ.Execute(buf,&param)

    err := smtp.SendMail(
	"smtp.csie.ntu.edu.tw:25",
	smtp.PlainAuth(
	    "",
	    Mail_User,
	    Mail_Passwd,
	    "smtp.csie.ntu.edu.tw"),
	"sprout@csie.ntu.edu.tw",
	[]string{target},
	buf.Bytes(),
    )

    return err
}
