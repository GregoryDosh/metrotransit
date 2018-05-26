package metrotransit

import (
	"errors"
	"time"
)

// Env is a wrapper around the Datastore so that any struct which implements the Datastore functions can be used.
type Env struct {
	DS Datastore
}

// GetDepartures combines the departure times & stop information into one struct.
func (env *Env) GetDepartures(stopID int) (*Stop, error) {
	stop := &Stop{
		StopID: stopID,
	}
	if stopID <= 0 {
		return stop, errors.New("invalid stopID")
	}
	sd, err := env.DS.GetStopDetails(stopID)
	if err != nil {
		return nil, err
	}
	stop.Details = *sd

	stopDepartures, err := env.DS.GetStopDepartures(stopID)
	if err != nil {
		return nil, err
	}
	stop.Departures = *stopDepartures

	stop.UpdateTime = time.Now()

	return stop, nil
}
