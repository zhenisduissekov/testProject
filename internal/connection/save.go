package connection

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

// saving to DB using connection
func (c *Client) Save(query string, args []interface{}) (err error) {
	resp, err := c.pool.Exec(*c.ctx, query, args...)
	if err != nil {
		return fmt.Errorf("could not exec: %w", err)
	}
	log.Trace().Msgf("rows affected: %s", resp)
	return nil
}
