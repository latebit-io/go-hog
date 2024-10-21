package gohog

// To run these tests use docker compose up in the project root to start mailhog
import (
	"fmt"
	"gopkg.in/mail.v2"
	"net/http"
	"testing"
)

const mailHogUrl = "http://localhost:8025"

func sendMail() {
	// Create a new message
	message := mail.NewMessage()

	// Set email headers
	message.SetHeader("From", "gohog@mailhog.com")
	message.SetHeader("To", "piggy@mailhog.com")
	message.SetHeader("Subject", "Hello from the MailHog!")

	// Set email body
	message.SetBody("text/plain", "Test mail")

	// Set up the SMTP dialer
	dialer := mail.NewDialer("localhost", 1025, "user", "pass")

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}

func TestGetMessage(t *testing.T) {
	sendMail()
	client := http.Client{}
	gohog := NewGoHogClient(mailHogUrl, &client)
	messages, err := gohog.Search(searchContains, "Hello", 1, 20)
	message, err := gohog.Message(messages.Items[0].Id)
	if err != nil {
		t.Error(err)
	}

	if message.Id != messages.Items[0].Id {
		t.Errorf("Wrong id: %s", message.Id)
	}
}

func TestGetMessages(t *testing.T) {
	sendMail()
	client := http.Client{}
	gohog := NewGoHogClient(mailHogUrl, &client)
	messages, err := gohog.Messages(1, 30)
	if err != nil {
		t.Error(err)
	}

	if messages.Total == 0 {
		t.Errorf("No messages")
	}
}

func TestSearchMessages(t *testing.T) {
	sendMail()
	client := http.Client{}
	gohog := NewGoHogClient(mailHogUrl, &client)
	messages, err := gohog.Search(searchContains, "Hello", 1, 20)
	if err != nil {
		t.Error(err)
	}

	if messages.Total == 0 {
		t.Errorf("No messages")
	}
}

func TestDeleteMessage(t *testing.T) {
	sendMail()
	client := http.Client{}
	gohog := NewGoHogClient(mailHogUrl, &client)
	messages, err := gohog.Search(searchContains, "Hello", 1, 20)
	err = gohog.Delete(messages.Items[0].Id)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteMessages(t *testing.T) {
	sendMail()
	client := http.Client{}
	gohog := NewGoHogClient(mailHogUrl, &client)
	err := gohog.DeleteAll()
	if err != nil {
		t.Error(err)
	}
}
