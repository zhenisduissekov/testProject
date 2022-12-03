package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/pkg/repository"
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

func (h *Handler) Publish(c *websocket.Conn) {
	log.Trace().Msg("publish handler started")
	err := h.socket.Subscribe(context.Background(), c)
	if err != nil {
		log.Err(err).Msg("could not subscribe to client")
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
