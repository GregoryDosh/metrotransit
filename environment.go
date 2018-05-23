package metrotransit

import (
	"errors"
	"time"
)

type Env struct {
	DS Datastore
}

// GetDepartures gets the departures from the MetroTransit API
func (env *Env) GetDepartures(stopID int) (*Stop, error) {
	stop := &Stop{
		StopID: stopID,
	}
	if stopID <= 0 {
		return stop, errors.New("invalid stopID")
	}
	stopDetails, err := env.DS.GetStopDetails(stopID)
	if err != nil {
		return nil, err
	}
	stop.Description = stopDetails.StopName

	stopDepartures, err := env.DS.GetStopDepartures(stopID)
	if err != nil {
		return nil, err
	}
	stop.Departures = *stopDepartures

	stop.UpdateTime = time.Now()

	return stop, nil
}
