package main

import (
    "net/smtp"
    "text/template"
    "fmt"
    "bytes"
)

func MailReq() {
    templ := template.Must(template.New("ReqMail").Parse("From: {{.From}}\r\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n{{.Message}}"))
    param := struct {
	From string
	To string
	Subject string
	Message string
    }{
	"sprout@csie.ntu.edu.tw",
	"xxxx@gmail.com",
	"Hello Request",
	"<html><body><h3>This is Request</h3></body></html>",
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
	[]string{"xxxx@gmail.com"},
	buf.Bytes(),
    )
    fmt.Println(err)
}
