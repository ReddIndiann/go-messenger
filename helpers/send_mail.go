package helpers

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

func sendMail(email string, subject string, body string) error {
	auth := smtp.PlainAuth("",os.Getenv("SMTP_USER"),os.Getenv("SMTP_PASS"),os.Getenv("SMPT_HOST"))

	msg := []byte("Subject: " + subject + "\r\n" +
	"To: " + email + "\r\n" +
	"MIME-Version: 1.0\r\n" +
	"Content-Type: text/html; charset=\"utf-8\"\r\n" +
	"\r\n" + body + "\r\n")

for {
	err := smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, os.Getenv("SMTP_USER"), []string{email}, msg)
	if err == nil {
		fmt.Println("Email sent successfully to", email)
		return nil
	}

	
	fmt.Println("Error sending email:", err)
	fmt.Println("Retrying in 5 seconds...")
	time.Sleep(5 * time.Second)
}

}