package handler

import (
	"github.com/zhenisduissekov/testProject/pkg/cryptocompare"
	"github.com/zhenisduissekov/testProject/pkg/socket"
)

type Handler struct {
	crypto cryptocompare.Client
	socket socket.Client
}

func New(crypto cryptocompare.Client, socket socket.Client) *Handler {
	return &Handler{
		crypto: crypto,
		socket: socket,
	}
}

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Results interface{} `json:"results"`
}
