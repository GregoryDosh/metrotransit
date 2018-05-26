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
	GetStopDetails(stopID int) (*Details, error)
	GetStopDepartures(stopID int) (*[]Departure, error)
}

// DefaultDatastore is an implementation of the Datastore with an underlying data structure of PostgreSQL & an HTTP Client
type DefaultDatastore struct {
	DB         *sql.DB
	HTTPClient *http.Client
}

// InitDefaultDatastore creates and sets up the default PostgreSQL & http client
func InitDefaultDatastore(Host string, Port string, User string, Password string, Database string, SSLMode string) (*DefaultDatastore, error) {
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
	return &DefaultDatastore{
		DB: db,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

// GetStopDetails pulls the stop details from a PostgreSQL table named `mt.stops`.
func (DefaultDatastore *DefaultDatastore) GetStopDetails(stopID int) (*Details, error) {
	var (
		stop                 Details
		dbID                 sql.NullInt64
		dbCode               sql.NullString
		dbName               sql.NullString
		dbDescription        sql.NullString
		dbLatitude           sql.NullFloat64
		dbLongitude          sql.NullFloat64
		dbZoneID             sql.NullString
		dbURL                sql.NullString
		dbLocationType       sql.NullInt64
		dbWheelchairBoarding sql.NullInt64
	)
	// Should parameterize this in the future
	err := DefaultDatastore.DB.QueryRow("select * from mt.stops where stop_id = $1", stopID).Scan(&dbID, &dbCode, &dbName, &dbDescription, &dbLatitude, &dbLongitude, &dbZoneID, &dbURL, &dbLocationType, &dbWheelchairBoarding)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("no stop with that ID")
	case err != nil:
		return nil, err
	}

	stop.ID = dbID.Int64
	stop.Code = dbCode.String
	stop.Name = dbName.String
	stop.Description = dbDescription.String
	stop.Latitude = dbLatitude.Float64
	stop.Longitude = dbLongitude.Float64
	stop.ZoneID = dbZoneID.String
	stop.URL = dbURL.String
	stop.LocationType = dbLocationType.Int64
	stop.WheelchairBoarding = dbWheelchairBoarding.Int64

	return &stop, nil
}

// GetStopDepartures pulls the stop departures from the MetroTransit API.
func (DefaultDatastore *DefaultDatastore) GetStopDepartures(stopID int) (*[]Departure, error) {
	departures := &[]Departure{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://svc.metrotransit.org/NexTrip/%d?format=json", stopID), nil)
	if err != nil {
		return departures, err
	}

	res, err := DefaultDatastore.HTTPClient.Do(req)
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
