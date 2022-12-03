package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/internal/cryptocompare"
	"github.com/zhenisduissekov/testProject/internal/repository"
	"time"
)

// HealthCheck godoc
// @Summary Статус сервера
// @Description проверка статуса.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} response{status=string,message=string} "успешный ответ"
// @Router /health [get]
func (h *Handler) HealthCheck(f *fiber.Ctx) error {
	log.Trace().Msg("healthcheck")
	return f.Status(fiber.StatusOK).JSON(&response{
		Status:  "success",
		Message: "Request successfully processed",
	})
}

func (h *Handler) GetPriceWS(c *websocket.Conn) {
	log.Trace().Msg("get prices WS handler started")
	defer log.Trace().Msg("get prices WS handler finished")
	for {
		log.Info().Msg("publishTest waiting for client message")
		var items cryptocompare.GetPriceReqItems
		err := c.ReadJSON(&items)
		if err != nil {
			log.Err(err).Msg("read:")
			break
		}

		fmt.Println("TSYMS:", items.TSYMS, " FSYMS:", items.FSYMS)

		result, err := h.crypto.GetPrice(items)
		if err != nil {
			result, err = h.crypto.ReadFromDB(items)
			if err != nil {
				log.Err(err).Msgf("get quotes could not get price for %v", items)
				result = "could not get price for " + items.FSYMS + " " + items.TSYMS
			}
		}

		if err := c.WriteJSON(result); err != nil {
			log.Err(err).Msg("write:")
			break
		}
	}
}

func (h *Handler) PublishTest(c *websocket.Conn) {
	for {
		log.Info().Msg("publishTest waiting for client message")
		time.Sleep(3 * time.Second)
		msg := repository.IndexedMessages{
			Number:    0,
			Text:      "тестовое сообщение",
			Status:    "sent",
			CreatedOn: time.Now().String(),
			CreatedBy: "test client",
		}

		if err := c.WriteJSON(msg); err != nil {
			log.Err(err).Msg("write:")
			break
		}
	}
}
