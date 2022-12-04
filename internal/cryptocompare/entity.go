package cryptocompare

import (
	"github.com/zhenisduissekov/testProject/internal/repository"
	"net/http"
	"time"
)

type PriceReqItems struct {
	FSYMS string `json:"fsyms" example:"BTC,ETH"`
	TSYMS string `json:"tsyms" example:"USD,EUR"`
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
	CHANGE24HOUR    string `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR string `json:"CHANGEPCT24HOUR"`
	OPEN24HOUR      string `json:"OPEN24HOUR"`
	VOLUME24HOUR    string `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  string `json:"VOLUME24HOURTO"`
	HIGH24HOUR      string `json:"HIGH24HOUR"`
	PRICE           string `json:"PRICE"`
	FROMSYMBOL      string `json:"FROMSYMBOL"`
	TOSYMBOL        string `json:"TOSYMBOL"`
	LASTUPDATE      string `json:"LASTUPDATE"`
	SUPPLY          string `json:"SUPPLY"`
	MKTCAP          string `json:"MKTCAP"`
}

type Config struct {
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

func New(cnf Config, repo repository.Repository) Client {
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
