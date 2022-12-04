package cryptocompare

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

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
