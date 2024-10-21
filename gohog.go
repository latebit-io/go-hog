// Package gohog is a client library to interact with and consume mailhog's api
// https://github.com/mailhog/MailHog mailhog is very useful when testing email integration
package gohog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Path used to represent an SMTP forward or return path
type Path struct {
	Relays  []string `json:"relays"`
	Mailbox string   `json:"mailbox"`
	Domain  string   `json:"domain"`
	Params  string   `json:"params"`
}

// Mime email content parts, can be text or rich content
type Mime struct {
	Parts []Content
}

// Headers metadata about the email message
type Headers struct {
	ContentId               []string `json:"content-id"`
	ContentDisposition      []string `json:"content-disposition"`
	ContentTransferEncoding []string `json:"content-transfer-encoding"`
	ContentType             []string `json:"content-type"`
	Date                    []string `json:"date"`
	From                    []string `json:"from"`
	MimeVersion             []string `json:"mime-version"`
	MessageId               []string `json:"message-id"`
	Received                []string `json:"received"`
	ReturnPath              []string `json:"return-path"`
	Subject                 []string `json:"subject"`
	To                      []string `json:"to"`
}

// Content of the email the Mime parts contain the content
type Content struct {
	Headers Headers `json:"headers"`
	Body    string  `json:"body"`
	Size    int     `json:"size"`
	Mime    Mime    `json:"mime"`
}

// Raw information about the email
type Raw struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Date string   `json:"date"`
	Helo string   `json:"helo"`
}

// Message is used to represent a single email message
type Message struct {
	Id      string    `json:"id"`
	From    Path      `json:"from"`
	To      []Path    `json:"to"`
	Content Content   `json:"content"`
	Created time.Time `json:"created"`
	Mime    Mime      `json:"mime"`
	Raw     Raw       `json:"raw"`
}

// Messages a struct the will contain a list of messages used for search and messages endpoint
type Messages struct {
	Total int       `json:"total"`
	Count int       `json:"count"`
	Start int       `json:"start"`
	Items []Message `json:"items"`
}

// Base end points for mailhog
const (
	messagesUrl = "/api/v2/messages"
	messageUrl  = "/api/v1/messages"
	searchUrl   = "/api/v2/search"
	deleteUrl   = "/api/v1/messages"
)

// Search types used by the search function
const (
	// searchFrom used to search the "from" field on emails
	searchFrom = "from"
	// searchTo used to search the "to" field on emails
	searchTo = "to"
	// searchContains used to search the content on emails
	searchContains = "containing"
)

type GoHog struct {
	url    string
	client *http.Client
}

// NewGoHogClient creates a mailhog client
func NewGoHogClient(url string, client *http.Client) *GoHog {
	return &GoHog{
		url:    url,
		client: client,
	}
}

// Messages will return email messages between the start and limit
func (h *GoHog) Messages(start, limit int) (Messages, error) {
	messages := Messages{}
	mailHogUrl := fmt.Sprintf("%s%s?start=%d&limit=%d", h.url, messagesUrl, start, limit)
	err := doGet(mailHogUrl, &messages, *h.client)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

// Message will return a single email message by id
// Example id: eZVH3mSvQl9oWzIc4U1j1zWO8TWFWNv123iPrS0sOkE=@mailhog.example
func (h *GoHog) Message(id string) (Message, error) {
	message := Message{}
	mailHogUrl := fmt.Sprintf("%s%s/%s", h.url, messageUrl, id)
	err := doGet(mailHogUrl, &message, *h.client)
	if err != nil {
		return message, err
	}
	return message, nil
}

// Search will search for messages by searchType using constants searchFrom, searchTo, and searchContains
// you can limit how many email messages will be returned
// query field is string that will be searched for in the email message
func (h *GoHog) Search(searchType, query string, start, limit int) (Messages, error) {
	messages := Messages{}
	mailHogUrl := fmt.Sprintf("%s%s?kind=%s&query=%s&start=%d&limit=%d", h.url, searchUrl, searchType, url.QueryEscape(query),
		start, limit)
	err := doGet(mailHogUrl, &messages, *h.client)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

// Delete this will delete a email message from mailhog by id
// Example id: eZVH3mSvQl9oWzIc4U1j1zWO8TWFWNv123iPrS0sOkE=@mailhog.example
func (h *GoHog) Delete(id string) error {
	mailHogUrl := fmt.Sprintf("%s%s/%s", h.url, deleteUrl, id)
	req, err := http.NewRequest(http.MethodDelete, mailHogUrl, nil)
	if err != nil {
		return err
	}
	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Email could not be deleted - %s", res.Status))
	}
	return nil
}

// DeleteAll deletes all email messages on the server
func (h *GoHog) DeleteAll() error {
	mailHogUrl := fmt.Sprintf("%s%s", h.url, deleteUrl)
	req, err := http.NewRequest(http.MethodDelete, mailHogUrl, nil)
	if err != nil {
		return err
	}
	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("All emails could not be deleted - %s", res.Status))
	}
	return nil
}

// Utility method for get request, do not use directly
func doGet(url string, model interface{}, client http.Client) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&model)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
