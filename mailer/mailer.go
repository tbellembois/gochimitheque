package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
)

var (
	// MailServerAddress is the SMTP server address
	// such as smtp.univ.fr
	MailServerAddress string
	// MailServerSender is the username used
	// to send mails
	MailServerSender string
	// MailServerPort is the SMTP server port
	MailServerPort string
	// MailServerUseTLS specify if a TLS SMTP connection
	// should be used
	MailServerUseTLS bool
	// MailServerTLSSkipVerify bypass the SMTP TLS verification
	MailServerTLSSkipVerify bool
)

// TestMail send a mail to "to"
func TestMail(to string) error {
	return SendMail(to, "test mail from Chimithèque", "your mail configuration seems ok")
}

// SendMail send a mail
func SendMail(to string, subject string, body string) error {

	var (
		e         error
		tlsconfig *tls.Config
		tlsconn   *tls.Conn
		client    *smtp.Client
		smtpw     io.WriteCloser
		n         int64
		message   string
	)

	// build message
	message += fmt.Sprintf("From: %s\r\n", MailServerSender)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += "Content-Type: text/plain; charset=utf-8\r\n"
	message += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n" + body + "\r\n"

	logger.Log.WithFields(logrus.Fields{
		"globals.MailServerAddress":       MailServerAddress,
		"globals.MailServerPort":          MailServerPort,
		"globals.MailServerSender":        MailServerSender,
		"globals.MailServerUseTLS":        MailServerUseTLS,
		"globals.MailServerTLSSkipVerify": MailServerTLSSkipVerify,
		"subject":                         subject,
		"to":                              to}).Debug("sendMail")

	if MailServerUseTLS {
		// tls
		tlsconfig = &tls.Config{
			InsecureSkipVerify: MailServerTLSSkipVerify,
			ServerName:         MailServerAddress,
		}
		if tlsconn, e = tls.Dial("tcp", MailServerAddress+":"+MailServerPort, tlsconfig); e != nil {
			return e
		}
		defer tlsconn.Close()
		if client, e = smtp.NewClient(tlsconn, MailServerAddress+":"+MailServerPort); e != nil {
			return e
		}
	} else {
		if client, e = smtp.Dial(MailServerAddress + ":" + MailServerPort); e != nil {
			return e
		}
	}
	defer client.Close()

	// to && from
	logger.Log.Debug("setting sender")
	if e = client.Mail(MailServerSender); e != nil {
		return e
	}
	logger.Log.Debug("setting recipient")
	if e = client.Rcpt(to); e != nil {
		return e
	}
	// data
	logger.Log.Debug("setting body")
	if smtpw, e = client.Data(); e != nil {
		return e
	}

	// send message
	logger.Log.Debug("sending message")
	buf := bytes.NewBufferString(message)
	if n, e = buf.WriteTo(smtpw); e != nil {
		return e
	}
	smtpw.Close()
	logger.Log.WithFields(logrus.Fields{"n": n}).Debug("sendMail")

	// send quit command
	logger.Log.Debug("setting quit command")
	_ = client.Quit()

	return nil
}
