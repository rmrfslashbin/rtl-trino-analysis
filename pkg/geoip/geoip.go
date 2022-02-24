package geoip

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

// Used to manage varidic options
type Option func(c *Config)

// database configs
type Config struct {
	geodb string
	db    *maxminddb.Reader
}

var ()

// Record defines the fields to fetch from the GeoIP database.
type Record struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`

	Continent struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"continent"`

	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`

	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`

	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`

	Subdivisions []struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"subdivisions"`
}

// GeoIPData represents the data returned.
type GeoIPData struct {
	IP          net.IP  `json:"ip"`
	City        string  `json:"city_name"`
	Continent   string  `json:"continent_code"`
	Country     string  `json:"country_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   uint    `json:"metro_code"`
	TimeZone    string  `json:"time_zone"`
	PostalCode  string  `json:"postal_code"`
	Subdivision string  `json:"subdivision_code"`
}

// New returns a new Config with the given options
func New(opts ...func(*Config)) (*Config, error) {
	config := &Config{}

	// apply options
	for _, opt := range opts {
		opt(config)
	}

	// dns must be set
	if config.geodb == "" {
		return nil, fmt.Errorf("geodb is required")
	}

	db, err := maxminddb.Open(config.geodb)
	if err != nil {
		log.Fatal(err)
	}
	config.db = db

	return config, nil
}

// SetGeoDB sets the database location
func SetGeoDB(geodb string) Option {
	return func(c *Config) {
		c.geodb = geodb
	}
}

// Close closes the database connection
func (c *Config) Close() error {
	return c.db.Close()
}

// Lookip GeoIP data for the given IP address.
func (c *Config) Lookup(ip net.IP) (*GeoIPData, error) {
	var data Record = Record{}

	err := c.db.Lookup(ip, &data)
	if err != nil {
		return nil, err
	}

	subdivs := []string{}
	for _, s := range data.Subdivisions {
		subdivs = append(subdivs, s.IsoCode)
	}

	return &GeoIPData{
		IP:          ip,
		City:        data.City.Names["en"],
		Continent:   data.Continent.Code,
		Country:     data.Country.IsoCode,
		Latitude:    data.Location.Latitude,
		Longitude:   data.Location.Longitude,
		MetroCode:   data.Location.MetroCode,
		TimeZone:    data.Location.TimeZone,
		PostalCode:  data.Postal.Code,
		Subdivision: strings.Join(subdivs, ";"),
	}, nil
}
