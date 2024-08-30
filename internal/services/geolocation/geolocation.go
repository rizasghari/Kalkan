package geolocation

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoLocation struct {
	db *geoip2.Reader
}

func NewGeoLocation(dbType string) *GeoLocation {
	geoLocation := GeoLocation{}

	if err := geoLocation.connectToDatabase(dbType); err != nil {
		log.Fatal(err)
	}
	return &geoLocation
}

func (gl *GeoLocation) connectToDatabase(dbType string) error {
	db, err := geoip2.Open(fmt.Sprintf("./db/GeoLite2-%s.mmdb", dbType))
    if err != nil {
        return err
    }

	gl.db = db
	return nil
}

func (gl *GeoLocation) GetLocation(ipAddr string) (*geoip2.City, error) {
	ip := net.ParseIP(ipAddr)
    record, err := gl.db.City(ip)
    if err != nil {
        log.Fatal(err)
    }
	return record, nil
}