package cryptocompare

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

// GetPrice returns quotes for specific fsyms and tsyms or the ones in .env
func (c Client) GetPrice(reqItems PriceReqItems) (result interface{}, err error) {
	if reqItems.FSYMS != "" && reqItems.TSYMS != "" {
		c.FSYMS = reqItems.FSYMS
		c.TSYMS = reqItems.TSYMS
	}

	result, err = c.requestFromCrypto()
	if err != nil {
		log.Err(err).Msg("Error while getting quotes from CryptoCompare")
		result, err = c.readFromDB()
		if err != nil {
			return nil, fmt.Errorf("error while reading from DB: %w", err)
		}
	}
	return result, nil
}
