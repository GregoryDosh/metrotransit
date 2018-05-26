package metrotransit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Stop holds the information about a MetroTransit Stop
type Stop struct {
	Departures []Departure `json:"departures"`
	Details    Details     `json:"stop_details"`
	StopID     int         `json:"stop_id"`
	UpdateTime time.Time   `json:"update_time"`
}

// Details holds the stop specific information.
type Details struct {
	ID                 int64   `json:"id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	ZoneID             string  `json:"zone_id"`
	URL                string  `json:"url"`
	LocationType       int64   `json:"location_type"`
	WheelchairBoarding int64   `json:"wheelchair_boarding"`
}

// Departure holds information about one particular bus
type Departure struct {
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
	VehicleLatitude  float64           `json:"VehicleLatitude"`
	VehicleLongitude float64           `json:"VehicleLongitude"`
}

// MarshalJSON implemented to correctly snake_case the json output.
func (d *Departure) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON used to convert odd date format into a date format Go can recognize
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
