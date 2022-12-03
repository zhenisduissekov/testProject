package socket

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/pkg/repository"
	"golang.org/x/time/rate"
)

type Client struct {
	Client                  http.Client
	Url                     string
	subscriberMessageBuffer int
	serveMux                http.ServeMux
	subscribersMu           sync.Mutex
	subscribers             map[*subscriber]struct{}
	publishLimiter          *rate.Limiter
	repo                    repository.Repository
}

type subscriber struct {
	userId    string
	msgs      chan repository.IndexedMessages
	closeSlow func()
}

// SendRequest sends request to contact center
func (c *Client) SendRequest(body []byte) (*http.Response, error) {
	log.Trace().Msgf("sending request to %v", c.Url)
	req, err := http.NewRequest(http.MethodPost, c.Url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("could not prepare request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	log.Trace().Msg("request sent successfully")
	return resp, nil
}

func (c *Client) MsgOperator(body []byte) (err error) {

	resp, err := c.SendRequest(body)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d, body: %v", resp.StatusCode, string(respBody))
	}

	fmt.Println("resp status code: ", resp.StatusCode, string(respBody))
	return nil
}
