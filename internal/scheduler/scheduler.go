package scheduler

import (
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/internal/cryptocompare"
	"time"
)

type SchedulerConfig struct {
	Interval int
}

func Run(intrvl SchedulerConfig, c cryptocompare.Client) {
	log.Trace().Msg("Scheduler started")
	ticker := time.NewTicker(time.Duration(intrvl.Interval) * time.Second)
	var items cryptocompare.GetPriceReqItems
	for range ticker.C {
		log.Trace().Msg("Scheduler tick")
		_, err := c.GetPrice(items)
		if err != nil {
			log.Err(err).Msg("Error while getting quotes")
		}
		time.Sleep(time.Duration(intrvl.Interval) * time.Second)
	}
}
