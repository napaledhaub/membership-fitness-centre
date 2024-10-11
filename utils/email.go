package utils

import (
	"net/smtp"
)

var from = ""
var passwordEmail = "gaqlxlrphumweaxw "
var host = "smtp.gmail.com"
var port = "587"

func SendPackageEmails(to string, msg []byte) error {
	auth := smtp.PlainAuth("", from, passwordEmail, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
	return err
}
