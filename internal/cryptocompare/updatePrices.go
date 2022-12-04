package cryptocompare

import "fmt"

// method for auto updating prices in DB
func (c Client) UpdatePrices() (result interface{}, err error) {
	result, err = c.requestFromCrypto()
	if err != nil {
		return nil, fmt.Errorf("error while getting quotes: %w", err)
	}
	return result, nil
}
