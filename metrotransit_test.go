package metrotransit_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/GregoryDosh/metrotransit"
)

type mockDS struct{}

func (mockDS *mockDS) GetStopDetails(stopID int) (*metrotransit.StopDetails, error) {
	var stop metrotransit.StopDetails
	switch stopID {
	case 5611:
		stop.StopName = "Victoria St & Idaho Ave"
	case 17946:
		stop.StopName = "Hennepin Ave & 7th St S"
	}
	return &stop, nil
}
func (mockDS *mockDS) GetStopDepartures(stopID int) (*[]metrotransit.StopDeparture, error) {
	switch stopID {
	case 5611:
		return &[]metrotransit.StopDeparture{
			{},
		}, nil
	case 17946:
		return &[]metrotransit.StopDeparture{
			{},
		}, nil
	}
	return nil, errors.New("no information found")
}

var _ = Describe("Metrotransit", func() {
	Describe("GetDepartures", func() {
		env := metrotransit.Env{
			DS: &mockDS{},
		}
		Context("with invalid parameters", func() {
			It("-12 should return an error", func() {
				_, err := env.GetDepartures(-12)
				Expect(err).To(MatchError(errors.New("invalid stopID")))
			})
		})

		Context("with valid parameters stopID=5611", func() {
			d, err := env.GetDepartures(5611)
			It("should not error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should have correct description", func() {
				Expect(d.Description).To(Equal("Victoria St & Idaho Ave"))
			})
			It("should have non empty departures", func() {
				Expect(d.Departures).ShouldNot(BeEmpty())
			})
		})

		Context("with valid parameters stopID=17946", func() {
			d, err := env.GetDepartures(17946)
			It("should not error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should have correct description", func() {
				Expect(d.Description).To(Equal("Hennepin Ave & 7th St S"))
			})
			It("should have non empty departures", func() {
				Expect(d.Departures).ShouldNot(BeEmpty())
			})
		})
	})
})