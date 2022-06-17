package util

import (
	"bytes"
	"fmt"
	CONFIG "hecruit-backend/config"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func saveToDiskFile(file []byte, fileName string) (string, bool) {
	newfile := bytes.NewReader(file)

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("saveToDiskFile", err)
		return "", false
	}
	fmt.Println("created")
	defer f.Close()
	_, err = io.Copy(f, newfile)
	if err != nil {
		fmt.Println("saveToDiskFile", err)
		return "", false
	}

	return fileName, true
}

func UploadContentAsFile(fileName string, body []byte) bool {
	savedFileName, saved := saveToDiskFile(body, fileName)
	if !saved {
		return false
	}
	defer os.Remove(fileName)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(CONFIG.S3ID, CONFIG.S3Secret, ""),
		Endpoint:         aws.String(CONFIG.S3Endpoint),
		Region:           aws.String(CONFIG.S3Region),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess := session.New(s3Config)

	svc := s3manager.NewUploader(sess)

	fmt.Println("Uploading content to S3...")

	openedFile, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("UploadContentAsFile "+fileName, err)
		return false
	}
	defer openedFile.Close()

	_, err = svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(CONFIG.S3Bucket),
		Key:    aws.String(fileName),
		Body:   openedFile,
	})
	if err != nil {
		fmt.Println("UploadContentAsFile "+fileName, err)
		return true
	}
	return true
}

func BuildICSFile(meetingLink, organizer, attendees, UID, title, start, end, status, sequence string) string {

	icsContent := `BEGIN:VCALENDAR
VERSION:2.0
CALSCALE:GREGORIAN
METHOD:REQUEST
BEGIN:VEVENT
DTSTART:` + getDateForICS(start) + `
DTEND:` + getDateForICS(end) + `
DTSTAMP:` + getDateForICS(time.Now().UTC().Format(time.RFC3339))
	icsContent += `
LOCATION=` + meetingLink
	icsContent += `
ORGANIZER;CN=` + organizer + `:mailto:` + organizer
	icsContent += `
UID:` + UID
	for _, email := range strings.Split(attendees, ",") {
		icsContent += `
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=TRUE;CN=` + email + `;X-NUM-GUESTS=0:mailto:` + email
	}

	icsContent += `
CREATED:` + getDateForICS(time.Now().UTC().Format(time.RFC3339)) + `
LAST-MODIFIED:` + getDateForICS(time.Now().UTC().Format(time.RFC3339)) + `
LOCATION:
SEQUENCE:` + sequence + `
STATUS:` + status + `
SUMMARY:` + title + `
TRANSP:OPAQUE
END:VEVENT
END:VCALENDAR`

	return icsContent
}

func getDateForICS(input string) string {
	t, _ := time.Parse(time.RFC3339, input)
	return t.Format("20060102T150405Z")
}
