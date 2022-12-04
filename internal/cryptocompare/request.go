package cryptocompare

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

// requesting from CryptoCompare and parsing response
func (c Client) requestFromCrypto() (result interface{}, err error) {
	url, err := c.prepUrl(c.FSYMS, c.TSYMS)
	if err != nil {
		return nil, err
	}
	log.Trace().Msgf("Getting quote for url: [%s]", url)
	resp, err := c.Http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error while http getting data: %w", err)
	}

	result, err = c.parseResponse(resp, c.FSYMS, c.TSYMS)
	if err != nil {
		return "", fmt.Errorf("error while parsing response: %w", err)
	}
	return result, nil
}
