package api

import (
	"encoding/json"
	"io"
)

type SpeedStatsResponse struct {
	RX string
	TX string
}

/**
* Get the current rx/tx speed of the given interface
* Tested values:
* * ra0 - wifi module
* * usb0 - 4g module
**/
func (c *Config) SpeedStats(iface string) (resp *SpeedStatsResponse, err error) {
	response, err := c.sessionApiCall("/cgi-bin/luci/admin/status/bandwidth?iface="+iface, nil, nil)
	if err != nil {
		return nil, err
	}

	var responseData [][]int64

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, ErrParseError
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, ErrParseError
	}

	if len(responseData) <= 2 {
		return &SpeedStatsResponse{
			RX: rateLabel(0),
			TX: rateLabel(0),
		}, nil
	}

	// Adapted from response of chart api call
	row := responseData[len(responseData)-1]
	oldRow := responseData[len(responseData)-2]

	timeDelta := row[0] - oldRow[0] // Find time delta between last 2 entries

	rx := rateLabel(float64((row[1] - oldRow[1]) * 1000000 / timeDelta))
	tx := rateLabel(float64((row[3] - oldRow[3]) * 1000000 / timeDelta))

	return &SpeedStatsResponse{
		RX: rx,
		TX: tx,
	}, nil
}
