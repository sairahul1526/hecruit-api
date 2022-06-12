package util

import (
	"fmt"
	CONFIG "hecruit-backend/config"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var emailWaitGroup sync.WaitGroup

func SendEmails(emails []map[string]string) {
	// send emails at once
	for _, email := range emails {
		emailWaitGroup.Add(1)
		go sendSESMail(email["from"], email["to"], email["title"], email["body"])
	}
	emailWaitGroup.Wait()
}

func sendSESMail(from, to, title, body string) {
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

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{ // Required
			ToAddresses: []*string{
				aws.String(to), // Required
			},
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
				Html: &ses.Content{
					Data:    aws.String(body), // Required
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &ses.Content{ // Required
				Data:    aws.String(title), // Required
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(from),
	}

	// send email
	svc.SendEmail(params)
}
