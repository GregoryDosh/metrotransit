package metrotransit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Stop holds the information about a MetroTransit Stop, including Departures,
// StopID, Description, and UpdateTime
type Stop struct {
	Departures  []StopDeparture `json:"departures"`
	StopID      int             `json:"stop_id"`
	Description string          `json:"description"`
	UpdateTime  time.Time       `json:"update_time"`
}

type StopDetails struct {
	StopID             int64
	StopCode           string
	StopName           string
	StopDesc           string
	StopLat            float64
	StopLon            float64
	ZoneID             string
	StopURL            string
	LocationType       int64
	WheelchairBoarding int64
}

type StopDeparture struct {
	Actual           bool              `json:"Actual"`
	BlockNumber      int               `json:"BlockNumber"`
	DepartureText    string            `json:"DepartureText"`
	DepartureTime    jsonDepartureTime `json:"DepartureTime"`
	Description      string            `json:"Description"`
	Gate             string            `json:"Gate"`
	Route            string            `json:"Route"`
	RouteDirection   string            `json:"RouteDirection"`
	Terminal         string            `json:"Terminal"`
	VehicleHeading   int               `json:"VehicleHeading"`
	VehicleLatitude  float32           `json:"VehicleLatitude"`
	VehicleLongitude float32           `json:"VehicleLongitude"`
}

func (d *StopDeparture) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["actual"] = d.Actual
	m["block_number"] = d.BlockNumber
	m["departure_text"] = d.DepartureText
	m["departure_time"] = d.DepartureTime
	m["description"] = d.Description
	m["gate"] = d.Gate
	m["route"] = d.Route
	m["route_direction"] = d.RouteDirection
	m["terminal"] = d.Terminal
	m["vehicle_heading"] = d.VehicleHeading
	m["vehicle_latitude"] = d.VehicleLatitude
	m["vehicle_longitude"] = d.VehicleLongitude
	return json.Marshal(m)
}

type jsonDepartureTime struct {
	time.Time
}

func (d *jsonDepartureTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) != 30 {
		return fmt.Errorf("date is wrong length: %d expected 30", len(s))
	}
	i, err := strconv.ParseInt(s[8:18], 10, 64)
	if err != nil {
		return err
	}
	t := time.Unix(i, 0)
	d.Time = t
	return nil
}
