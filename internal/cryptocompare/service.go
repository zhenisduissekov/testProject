package cryptocompare

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetPrice returns quotes for specific fsyms and tsyms
func (c Client) GetPrice(reqItems PriceReqItems) (result interface{}, err error) {
	if reqItems.FSYMS != "" && reqItems.TSYMS != "" {
		c.FSYMS = reqItems.FSYMS
		c.TSYMS = reqItems.TSYMS
	}

	result, err = c.requestFromCrypto()
	if err != nil {
		result, err = c.readFromDB()
		if err != nil {
			return nil, fmt.Errorf("error while reading from DB: %w", err)
		}
	}
	return result, nil
}

func (c Client) UpdatePrices() (result interface{}, err error) {
	result, err = c.requestFromCrypto()
	if err != nil {
		return nil, fmt.Errorf("error while getting quotes: %w", err)
	}
	return result, nil
}

// reqesting from CryptoCompare and parsing response
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

// parses response from CryptoCompare whether it's an error or not
func (c Client) parseResponse(resp *http.Response, fsyms, tsyms string) (result []CryptoResponse, err error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status received: %d", resp.StatusCode)
	}

	if strings.Contains(string(body), "rror") {
		return nil, fmt.Errorf("error response body received: %s", body)
	}

	// parsing response for different pairs, therefore cannot use predefined struct
	var objmap map[string]map[string]map[string]json.RawMessage
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling response: %w", err)
	}

	for _, v1 := range strings.Split(fsyms, ",") {
		for _, v2 := range strings.Split(tsyms, ",") {
			var temp CryptoResponse
			temp.TYPE = v1 + ":" + v2
			dataRaw := objmap["RAW"][v1][v2]
			err := json.Unmarshal(dataRaw, &temp.Raw)
			if err != nil {
				return nil, fmt.Errorf("error while unmarshalling data: %w", err)
			}

			dataDisplay := objmap["DISPLAY"][v1][v2]
			err = json.Unmarshal(dataRaw, &temp.Display)
			if err != nil {
				return nil, fmt.Errorf("error while unmarshalling data: %w", err)
			}

			err = c.repo.Save(dataRaw, dataDisplay, v1, v2, err)
			if err != nil {
				log.Err(err).Msg("Error while saving data")
			}
			result = append(result, temp)
		}
	}
	return result, nil
}

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

// ReadFromDB нужен если api cryptocompare не доступен
func (c Client) readFromDB() (result []CryptoResponse, err error) {
	log.Trace().Msgf("Getting quote from DB %v, %v", c.FSYMS, c.TSYMS)
	for _, v1 := range strings.Split(c.FSYMS, ",") {
		for _, v2 := range strings.Split(c.TSYMS, ",") {
			data, err := c.repo.Read(v1, v2)
			if err != nil {
				return result, fmt.Errorf("error while getting data from db: %w", err)
			}

			var temp CryptoResponse
			err = json.Unmarshal([]byte(data["raw"]), &temp.Raw)
			if err != nil {
				return result, fmt.Errorf("error while unmarshalling data: %w", err)
			}

			err = json.Unmarshal([]byte(data["display"]), &temp.Display)
			if err != nil {
				return result, fmt.Errorf("error while unmarshalling data: %w", err)
			}
			temp.TYPE = v1 + ":" + v2
			result = append(result, temp)
		}
	}
	return result, nil
}
