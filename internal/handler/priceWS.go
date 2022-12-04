package handler

import (
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/internal/cryptocompare"
)

func (h *Handler) GetPriceWS(c *websocket.Conn) {
	log.Trace().Msg("get prices WS handler started")
	defer log.Trace().Msg("get prices WS handler finished")
	for {
		log.Info().Msg("publishTest waiting for client message")
		var items cryptocompare.PriceReqItems
		err := c.ReadJSON(&items)
		if err != nil {
			log.Err(err).Msg("could not read json from client, exiting connection")
			break
		}

		result, err := h.crypto.GetPrice(items)
		if err != nil {
			log.Err(err).Msgf("get quotes could not get price for %v", items)
			result = "could not get price for " + items.FSYMS + " " + items.TSYMS
		}

		if err := c.WriteJSON(result); err != nil {
			log.Err(err).Msg("could not write to client, exiting connection")
			break
		}
	}
}
