package cryptocompare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/internal/repository"
)

type CryptoCompareConfig struct {
	APIKey       string
	TimeOut      int
	FSYMS        string
	TSYMS        string
	WaitInterval int
}

type Client struct {
	Url          string
	Http         http.Client
	FSYMS        string
	TSYMS        string
	WaitInterval int
	repo         repository.Repository
}

func New(cnf CryptoCompareConfig, repo repository.Repository) Client {
	return Client{
		Url: cnf.APIKey,
		Http: http.Client{
			Timeout: time.Duration(cnf.TimeOut) * time.Second,
		},
		FSYMS:        cnf.FSYMS,
		TSYMS:        cnf.TSYMS,
		WaitInterval: cnf.WaitInterval,
		repo:         repo,
	}
}

type GetPriceReqItems struct {
	FSYMS string `json:"fsyms" example:"BTC,ETH"`
	TSYMS string `json:"tsyms" example:"USD,EUR"`
}

// GetQuote returns quotes for specific fsyms and tsyms
func (c Client) GetPrice(reqItems GetPriceReqItems) (result interface{}, err error) {
	if reqItems.FSYMS != "" && reqItems.TSYMS != "" {
		c.FSYMS = reqItems.FSYMS
		c.TSYMS = reqItems.TSYMS
	}

	url, err := c.prepUrl(c.FSYMS, c.TSYMS)
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

	var objmap map[string]map[string]map[string]json.RawMessage
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling response: %w", err)
	}

	for _, v1 := range strings.Split(fsyms, ",") {
		for _, v2 := range strings.Split(tsyms, ",") {
			var temp CryptoResponse
			temp.TYPE = v1 + " >> " + v2
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
	}
	q := u.Query()
	q.Set("fsyms", fsym)
	q.Set("tsyms", tsym)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

type Response struct {
	FSYM           string      `json:"fsym"`
	TSYM           string      `json:"tsym"`
	CryptoResponse interface{} `json:"cryptoresponse"`
}

type CryptoResponse struct {
	TYPE    string                `json:"type"`
	Raw     CryptoResponseRaw     `json:"raw"`
	Display CryptoResponseDisplay `json:"display"`
}

type CryptoResponseRaw struct {
	CHANGE24HOUR    float64 `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR float64 `json:"CHANGEPCT24HOUR"`
	OPEN24HOUR      float64 `json:"OPEN24HOUR"`
	VOLUME24HOUR    float64 `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  float64 `json:"VOLUME24HOURTO"`
	LOW24HOUR       float64 `json:"LOW24HOUR"`
	HIGH24HOUR      float64 `json:"HIGH24HOUR"`
	PRICE           float64 `json:"PRICE"`
	SUPPLY          float64 `json:"SUPPLY"`
	MKTCAP          float64 `json:"MKTCAP"`
}

type CryptoResponseDisplay struct {
	CHANGE24HOUR    float64 `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR float64 `json:"CHANGEPCT24HOUR"`
	OPEN24HOUR      float64 `json:"OPEN24HOUR"`
	VOLUME24HOUR    float64 `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  float64 `json:"VOLUME24HOURTO"`
	HIGH24HOUR      float64 `json:"HIGH24HOUR"`
	PRICE           float64 `json:"PRICE"`
	FROMSYMBOL      string  `json:"FROMSYMBOL"`
	TOSYMBOL        string  `json:"TOSYMBOL"`
	LASTUPDATE      float64 `json:"LASTUPDATE"`
	SUPPLY          float64 `json:"SUPPLY"`
	MKTCAP          float64 `json:"MKTCAP"`
}

// ReadFromDB нужен если api cryptocompare не доступен
func (c Client) ReadFromDB(reqItems GetPriceReqItems) (result []CryptoResponse, err error) {
	if reqItems.FSYMS != "" && reqItems.TSYMS != "" {
		c.FSYMS = reqItems.FSYMS
		c.TSYMS = reqItems.TSYMS
	}
	for _, v1 := range strings.Split(c.FSYMS, ",") {
		for _, v2 := range strings.Split(c.TSYMS, ",") {
			data, err := c.repo.Get(v1, v2)
			if err != nil {
				return result, fmt.Errorf("error while getting data from db: %w", err)
			}
			var temp CryptoResponse
			log.Info().Msgf("data: %s", data)
			err = json.Unmarshal([]byte(data), &temp)
			if err != nil {
				return result, fmt.Errorf("error while unmarshalling data: %w", err)
			}
			temp.TYPE = v1 + " >> " + v2
			result = append(result, temp)
		}
	}
	return result, nil
}
