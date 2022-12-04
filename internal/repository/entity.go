package repository

import (
	"github.com/zhenisduissekov/testProject/internal/connection"
)

const (
	tableName = "crypto.quotes"
)

type IndexedMessages struct {
	Index     int    `json:"index"`
	Number    int32  `json:"number"`
	Text      string `json:"text"`
	Status    string `json:"status"`
	CreatedOn string `json:"createdOn"`
	CreatedBy string `json:"createdBy"`
}

type Repository struct {
	connection connection.Client
}

func New(conn connection.Client) Repository {
	return Repository{
		connection: conn,
	}
}
