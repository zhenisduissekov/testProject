package socket

import (
	"context"
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/pkg/repository"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"time"
)

type ChatReqItems struct {
	ClientId   string `json:"client_id"`
	ClientName string `json:"client_name"`
	Type       string `json:"type"`
	ChatId     string `json:"chat_id"`
	MessageId  string `json:"message_id"`
	Payload    string `json:"payload"`
	Filename   string `json:"filename"`
}

func New(repo repository.Repository) Client {
	os.Setenv("CONTACT_CENTER_URL", "https://demoback.crm.kz/test9kalf7854ncvfjsw0aznjhf84hjnchjs673jf92whadq325449ibkn") //todo: to remove in production
	url := os.Getenv("CONTACT_CENTER_URL")
	return Client{
		Client: http.Client{
			Timeout: 1 * time.Second,
		},
		Url:                     url,
		subscriberMessageBuffer: 16,
		subscribers:             make(map[*subscriber]struct{}),
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		repo:                    repo,
	}
}

func (cs *Client) Subscribe(ctx context.Context, c *websocket.Conn) error {
	s := &subscriber{
		msgs: make(chan repository.IndexedMessages, cs.subscriberMessageBuffer),
		closeSlow: func() {
			log.Info().Msg("connection is too slow")
			c.Close()
		},
	}
	cs.addSubscriber(s)
	defer cs.deleteSubscriber(s)

	go cs.readingFromSocket(c)

	for {
		select {
		case msg := <-s.msgs:
			err := c.WriteJSON(msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *Client) readingFromSocket(c *websocket.Conn) {
	for {
		var items ChatReqItems
		err := c.ReadJSON(&items)
		if err != nil {
			log.Err(err).Msg("could not read client message")
			break
		}

		itemsJson, err := json.Marshal(items)
		if err != nil {
			log.Err(err).Msg("could not marshal client message")
			continue
		}

		cs.MsgOperator(itemsJson)
	}
}

func (cs *Client) addSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	cs.subscribers[s] = struct{}{}
	cs.subscribersMu.Unlock()
}

func (cs *Client) deleteSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	delete(cs.subscribers, s)
	cs.subscribersMu.Unlock()
}

type OperatorMessage struct {
	ClientId   string `json:"client_id" validate:"omitempty,len=32,alphanum,lowercase" example:"e1979fed66804cb8bcda9cb1e24334b2"`
	ClientName string `json:"client_name"`
	Type       string `json:"type"`
	ChatId     string `json:"chat_id"`
	MessageId  string `json:"message_id"`
	Payload    string `json:"payload"`
	Filename   string `json:"filename"`
}

func (cs *Client) OperatorMsg(items OperatorMessage) {

	publishMsg := repository.IndexedMessages{
		Index:     0,
		Number:    0,
		CreatedBy: "0c49f1b2c00c4f76aea7bbd527ac5c4e", //todo: to replace in production
		CreatedOn: time.Now().Format("2006-01-02 15:04:05"),
		Text:      items.Payload,
	}

	go cs.publish(publishMsg)

}

func (cs *Client) publish(msg repository.IndexedMessages) {
	cs.subscribersMu.Lock()
	defer cs.subscribersMu.Unlock()

	cs.publishLimiter.Wait(context.Background())

	for s := range cs.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}
