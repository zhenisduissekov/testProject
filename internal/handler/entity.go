package handler

import (
	"github.com/zhenisduissekov/testProject/internal/cryptocompare"
)

type Handler struct {
	crypto cryptocompare.Client
}

func New(crypto cryptocompare.Client) *Handler {
	return &Handler{
		crypto: crypto,
	}
}

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Results interface{} `json:"results"`
}
