package mailhog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Path struct {
	Relays  []string `json:"relays"`
	Mailbox string   `json:"mailbox"`
	Domain  string   `json:"domain"`
	Params  string   `json:"params"`
}

type Mime struct {
	Parts []Content
}

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

type Content struct {
	Headers Headers `json:"headers"`
	Body    string  `json:"body"`
	Size    int     `json:"size"`
	Mime    Mime    `json:"mime"`
}

type Raw struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Date string   `json:"date"`
	Helo string   `json:"helo"`
}

type Message struct {
	Id      string    `json:"id"`
	From    Path      `json:"from"`
	To      []Path    `json:"to"`
	Content Content   `json:"content"`
	Created time.Time `json:"created"`
	Mime    Mime      `json:"mime"`
	Raw     Raw       `json:"raw"`
}

type Messages struct {
	Total int       `json:"total"`
	Count int       `json:"count"`
	Start int       `json:"start"`
	Items []Message `json:"items"`
}

const (
	messagesUrl = "/api/v2/messages"
	messageUrl  = "/api/v1/messages"
	searchUrl   = "/api/v2/search"
	deleteUrl   = "/api/v1/messages"
)

const (
	searchFrom     = "from"
	searchTo       = "to"
	searchContains = "containing"
)

type MailHog struct {
	url    string
	client *http.Client
}

func NewMailHogClient(url string, client *http.Client) *MailHog {
	return &MailHog{
		url:    url,
		client: client,
	}
}

func (h *MailHog) Messages(start, limit int) (Messages, error) {
	messages := Messages{}
	mailHogUrl := fmt.Sprintf("%s%s?start=%d&limit=%d", h.url, messagesUrl, start, limit)
	err := doGet(mailHogUrl, &messages, *h.client)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

func (h *MailHog) Message(id string) (Message, error) {
	message := Message{}
	mailHogUrl := fmt.Sprintf("%s%s/%s", h.url, messageUrl, id)
	err := doGet(mailHogUrl, &message, *h.client)
	if err != nil {
		return message, err
	}
	return message, nil
}

func (h *MailHog) Search(searchType, query string, start, limit int) (Messages, error) {
	messages := Messages{}
	mailHogUrl := fmt.Sprintf("%s%s?kind=%s&query=%s&start=%d&limit=%d", h.url, searchUrl, searchType, url.QueryEscape(query),
		start, limit)
	err := doGet(mailHogUrl, &messages, *h.client)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

func (h *MailHog) Delete(id string) error {
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
		return errors.New(res.Status)
	}
	return nil
}

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
