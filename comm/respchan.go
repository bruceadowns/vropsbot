package comm

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

// Request holds incoming messages to vropsbot
type Request struct {
	Channel   string
	User      string
	Arguments []string
}

// Response holds an outgoing slackbot message
type Response struct {
	Channel    string
	Text       string
	IsError    bool
	Parameters slack.PostMessageParameters
}

// Interop is the persistence struct that handles
// saving a slack request and response interaction
type Interop struct {
	Request  *Request
	Response *Response
}

// Save coalates the request and response and persists
func (iop *Interop) Save(db *sql.DB) (err error) {
	_, err = db.Exec(
		`INSERT INTO 'PullActions' (
			'When'
			, 'Channel'
			, 'User'
			, 'Request'
			, 'Response'
			, 'IsError')
		VALUES (?, ?, ?, ?, ?, ?)`,
		time.Now(),
		iop.Request.Channel,
		iop.Request.User,
		strings.Join(iop.Request.Arguments, " "),
		iop.Response.Text,
		iop.Response.IsError)

	return
}

// DefaultMessageParameters returns a default PostMessageParameters for vropsbot
func DefaultMessageParameters() (params slack.PostMessageParameters) {
	params = slack.NewPostMessageParameters()
	params.AsUser = true
	params.EscapeText = false

	return
}

// RespChan creates and returns a buffered channel used to write messages to slack
func RespChan(client *slack.Client, size int) (ch chan *Response) {
	ch = make(chan *Response, size)

	go func() {
		for {
			select {
			case m := <-ch:
				if m.Channel == "" {
					log.Print("Channel is empty")
					continue
				}
				if m.Text == "" {
					log.Print("Text is empty")
					continue
				}
				if _, _, err := client.PostMessage(m.Channel, m.Text, m.Parameters); err != nil {
					log.Printf("Error posting message to %s [%v].", m.Channel, err)
				}
			}
		}
	}()

	return
}
