package mailhog

import (
	"net/http"
	"testing"
)

const mailHogUrl = "http://localhost:8025"

func TestGetMessage(t *testing.T) {
	id := "eZVH3mSvQl9oWzIc4U1j1zWO8TWFWNv123iPrS0sOkE=@mailhog.example"
	client := http.Client{}
	mailhog := NewMailHogClient(mailHogUrl, &client)
	message, err := mailhog.Message(id)
	if err != nil {
		t.Error(err)
	}

	if message.Id != id {
		t.Errorf("Wrong id: %s", message.Id)
	}
}

func TestGetMessages(t *testing.T) {
	client := http.Client{}
	mailhog := NewMailHogClient(mailHogUrl, &client)
	messages, err := mailhog.Messages(1, 30)
	if err != nil {
		t.Error(err)
	}

	if messages.Total == 0 {
		t.Errorf("No messages")
	}
}

func TestSearchMessages(t *testing.T) {
	client := http.Client{}
	mailhog := NewMailHogClient(mailHogUrl, &client)
	messages, err := mailhog.Search(searchContains, "084e33e-ebdb-46ee-b085-7a72ec5a65e8", 1, 20)
	if err != nil {
		t.Error(err)
	}

	if messages.Total == 0 {
		t.Errorf("No messages")
	}
}

func TestDeleteMessage(t *testing.T) {
	client := http.Client{}
	mailhog := NewMailHogClient(mailHogUrl, &client)
	err := mailhog.Delete("kLT8D_S0KozRePv7pzcwd-5PnVH3MWw3QY5n3aZLmvc=@mailhog.example")
	if err != nil {
		t.Error(err)
	}
}
