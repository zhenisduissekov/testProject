package cryptocompare

import (
	"github.com/rs/zerolog/log"
	"net/url"
)

// PrepUrl prepares url for request to crypto compare
func (c *Client) prepUrl(fsym, tsym string) (result string, err error) {
	u, err := url.Parse(c.Url)
	if err != nil {
		log.Err(err).Msg("Error while parsing url")
		return "", err
	}
	q := u.Query()
	q.Set("fsyms", fsym)
	q.Set("tsyms", tsym)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
