package util

import (
	"bytes"
	"fmt"
	CONFIG "hecruit-backend/config"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-resty/resty/v2"
	"gopkg.in/gomail.v2"
)

var emailWaitGroup sync.WaitGroup

func SendEmails(emails []map[string]string) {
	// send emails at once
	for _, email := range emails {
		emailWaitGroup.Add(1)
		go sendSESMail(email["from"], email["reply_to"], email["to"], email["title"], email["body"], email["attachment"])
	}
	emailWaitGroup.Wait()
}

func sendSESMail(from, replyTo, to, title, body, attachment string) {
	defer emailWaitGroup.Done()

	// start a new aws session
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	// start a new ses session
	svc := ses.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(CONFIG.AWSAccessKey, CONFIG.AWSSecretKey, ""),
		Region:      aws.String("us-east-2"),
	})

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	if len(replyTo) > 0 {
		msg.SetHeader("Reply-To", replyTo)
	}
	msg.SetHeader("To", strings.Split(to, ",")...)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/html", body)
	if len(attachment) > 0 {
		tempFileName := GenerateRandomID() + "." + getFileExtension(attachment)
		fmt.Println(tempFileName, CONFIG.S3MediaURL+attachment)
		defer os.Remove(tempFileName)
		_, err := resty.New().R().
			SetOutput(tempFileName).
			Get(CONFIG.S3MediaURL + attachment)

		if err != nil {
			fmt.Println("failed to download attachment,", err)
			return
		}

		msg.Attach(tempFileName)
	}

	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	message := ses.RawMessage{Data: emailRaw.Bytes()}

	input := &ses.SendRawEmailInput{RawMessage: &message}

	// send email
	svc.SendRawEmail(input)
}

func getFileExtension(fileName string) string {
	temp := strings.Split(fileName, ".")
	return temp[len(temp)-1]
}
