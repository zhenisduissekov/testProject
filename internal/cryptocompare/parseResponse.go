package cryptocompare

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
)

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
			err = json.Unmarshal(dataDisplay, &temp.Display)
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
