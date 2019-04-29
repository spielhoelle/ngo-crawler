package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// Mail struct
type Mail struct {
	senderID string
	toIds    []string
	subject  string
	body     string
}

type smtpServer struct {
	host string
	port string
}

func (s smtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) buildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderID)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func mailer(jobsArray []string) {
	t, err := template.ParseFiles("email.html") //setp 1
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, jobsArray); err != nil {
		return
	}
	mail := Mail{}
	mail.senderID = os.Getenv("MAILUSER")
	mail.toIds = []string{os.Getenv("MAILRECEIVER")}
	mail.subject = "Here are all new found jobs"
	mail.body = buf.String()

	messageBody := mail.buildMessage()

	smtpServer := smtpServer{host: os.Getenv("MAILHOST"), port: os.Getenv("MAILPORT")}

	log.Println(smtpServer.host)
	//build an auth
	auth := smtp.PlainAuth(os.Getenv("MAILUSER"), mail.senderID, os.Getenv("MAILPW"), smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderID); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}
	// actually sends the mail out

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("mail sent successfully")

}
