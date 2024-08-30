package geolocation

import (
	"testing"

)

func TestGeoLocation(t *testing.T) {
	geolocation := NewGeoLocation("City")
	record, err := geolocation.GetLocation("88.239.138.38")
	if err != nil {
		t.Error(err)
	}
	t.Logf("City: %+v", record.City)
}