package scheduler

import (
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/internal/cryptocompare"
	"time"
)

type Config struct {
	Interval int
}

func Run(intrvl Config, c cryptocompare.Client) {
	log.Trace().Msg("Scheduler started")
	ticker := time.NewTicker(time.Duration(intrvl.Interval) * time.Second)
	for range ticker.C {
		log.Info().Msg("updating prices")
		_, err := c.UpdatePrices()
		if err != nil {
			log.Err(err).Msg("could not update prices")
		}
		time.Sleep(time.Duration(intrvl.Interval) * time.Second)
	}
}
