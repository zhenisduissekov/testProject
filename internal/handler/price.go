package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/internal/cryptocompare"
)

// GetPrices godoc
// @Summary Статус сервера
// @Description проверка статуса.
// @Tags price
// @Accept */*
// @Produce json
// @Param	queryparams		query 		cryptocompare.GetPriceReqItems	false "model get quotes"
// @Success 200 {object} response{status=string,message=string,results=object}
// @Failure 400 {object} response{status=string,message=string,results=object}
// @Failure 500 {object} response{status=string,message=string,results=object}
// @Router /service/price [get]
func (h *Handler) GetPrice(f *fiber.Ctx) (err error) {
	log.Info().Msg("get quotes")
	var items cryptocompare.GetPriceReqItems
	err = f.QueryParser(&items)
	if err != nil {
		log.Err(err).Msg("could not parse query")
		return f.Status(fiber.StatusBadRequest).JSON(&response{
			Status:  "error",
			Message: "Request error",
			Results: err.Error(),
		})
	}

	result, err := h.crypto.GetPrice(items)
	if err != nil {
		result, err = h.crypto.ReadFromDB(items)
		if err != nil {
			log.Err(err).Msg("get quotes could not get price")
			return f.Status(fiber.StatusInternalServerError).JSON(&response{
				Status:  "error",
				Message: "Could not get quotes",
				Results: result,
			})
		}
	}

	return f.Status(fiber.StatusOK).JSON(&response{
		Status:  "success",
		Message: "Request successfully processed",
		Results: result,
	})
}
