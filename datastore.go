package metrotransit

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	// This is used to register PostgreSQL to the DB packages
	_ "github.com/lib/pq"
)

// Datastore represents the expected interface for the underlying data
type Datastore interface {
	GetStopDetails(stopID int) (*StopDetails, error)
	GetStopDepartures(stopID int) (*[]StopDeparture, error)
}

type defaultDatastore struct {
	db         *sql.DB
	httpclient *http.Client
}

// InitDefaultDatastore creates and sets up the default PostgreSQL & http client
func InitDefaultDatastore(Host string, Port string, User string, Password string, Database string, SSLMode string) (*defaultDatastore, error) {
	if Host == "" || Port == "" || User == "" ||
		Password == "" || Database == "" {
		return nil, errors.New("all fields must be set")
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		User, Password, Database, Host, Port, SSLMode))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &defaultDatastore{
		db: db,
		httpclient: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

func (defaultDatastore *defaultDatastore) GetStopDetails(stopID int) (*StopDetails, error) {
	var (
		stop                 StopDetails
		dbStopID             sql.NullInt64
		dbStopCode           sql.NullString
		dbStopName           sql.NullString
		dbStopDesc           sql.NullString
		dbStopLat            sql.NullFloat64
		dbStopLon            sql.NullFloat64
		dbZoneID             sql.NullString
		dbStopURL            sql.NullString
		dbLocationType       sql.NullInt64
		dbWheelchairBoarding sql.NullInt64
	)
	err := defaultDatastore.db.QueryRow("select * from mt.stops where stop_id = $1", stopID).Scan(&dbStopID, &dbStopCode, &dbStopName, &dbStopDesc, &dbStopLat, &dbStopLon, &dbZoneID, &dbStopURL, &dbLocationType, &dbWheelchairBoarding)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("no stop with that ID")
	case err != nil:
		return nil, err
	}

	stop.StopID = dbStopID.Int64
	stop.StopCode = dbStopCode.String
	stop.StopName = dbStopName.String
	stop.StopDesc = dbStopDesc.String
	stop.StopLat = dbStopLat.Float64
	stop.StopLon = dbStopLon.Float64
	stop.ZoneID = dbZoneID.String
	stop.StopURL = dbStopURL.String
	stop.LocationType = dbLocationType.Int64
	stop.WheelchairBoarding = dbWheelchairBoarding.Int64

	return &stop, nil
}

func (defaultDatastore *defaultDatastore) GetStopDepartures(stopID int) (*[]StopDeparture, error) {
	departures := &[]StopDeparture{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://svc.metrotransit.org/NexTrip/%d?format=json", stopID), nil)
	if err != nil {
		return departures, err
	}

	res, err := defaultDatastore.httpclient.Do(req)
	if err != nil {
		return departures, err
	}

	if res.StatusCode == http.StatusBadRequest {
		return departures, errors.New("bad request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return departures, err
	}

	if err := json.Unmarshal(body, &departures); err != nil {
		return departures, fmt.Errorf("could not parse: %s", err.Error())
	}
	return departures, nil
}
