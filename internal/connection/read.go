package connection

import "fmt"

// reading from DB using connection
func (c *Client) Read(query string, args []interface{}) (map[string]string, error) {
	var raw, display string

	rows, err := c.pool.Query(*c.ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not Query: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&raw, &display)
		if err != nil {
			return nil, fmt.Errorf("could not scan: %w", err)
		}
	}

	data := make(map[string]string)
	data["raw"] = raw
	data["display"] = display
	return data, nil
}
