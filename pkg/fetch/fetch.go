package fetch

import (
	"fmt"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/trinodb/trino-go-client/trino"
)

// Entry represents a single Trino entry
type Entry struct {
	Timestamp                string  `db:"timestamp"`
	ClientIP                 string  `db:"client_ip"`
	Status                   int     `db:"status"`
	Bytes                    int64   `db:"bytes"`
	Method                   string  `db:"method"`
	Protocol                 string  `db:"protocol"`
	Host                     string  `db:"host"`
	UriStem                  string  `db:"uri_stem"`
	EdgeLocation             string  `db:"edge_location"`
	EdgeRequestID            string  `db:"edge_request_id"`
	HostHeader               string  `db:"host_header"`
	TimeTaken                float64 `db:"time_taken"`
	ProtoVersion             string  `db:"proto_version"`
	IPVersion                string  `db:"ip_version"`
	UserAgent                string  `db:"user_agent"`
	Referer                  string  `db:"referer"`
	Cookie                   string  `db:"cookie"`
	UriQuery                 string  `db:"uri_query"`
	EdgeResponseResultType   string  `db:"edge_response_result_type"`
	SslProtocol              string  `db:"ssl_protocol"`
	SslCipher                string  `db:"ssl_cipher"`
	EdgeResultType           string  `db:"edge_result_type"`
	ContentType              string  `db:"content_type"`
	ContentLength            int64   `db:"content_length"`
	EdgeDetailedResultType   string  `db:"edge_detailed_result_type"`
	Country                  string  `db:"country"`
	CacheBehaviorPathPattern string  `db:"cache_behavior_path_pattern"`
	Year                     string  `db:"year"`
	Month                    string  `db:"month"`
	Day                      string  `db:"day"`
}

// Record represents a properly formatted and typed Trino entry
type Record struct {
	Timestamp                time.Time `json:"timestamp"`
	ClientIP                 net.IP    `json:"client_ip"`
	Status                   int       `json:"status"`
	Bytes                    int64     `json:"bytes"`
	Method                   string    `json:"method"`
	Protocol                 string    `json:"protocol"`
	Host                     string    `json:"host"`
	UriStem                  string    `json:"uri_stem"`
	EdgeLocation             string    `json:"edge_location"`
	EdgeRequestID            string    `json:"edge_request_id"`
	HostHeader               string    `json:"host_header"`
	TimeTaken                float64   `json:"time_taken"`
	ProtoVersion             string    `json:"proto_version"`
	IPVersion                string    `json:"ip_version"`
	UserAgent                string    `json:"user_agent"`
	Referer                  string    `json:"referer"`
	Cookie                   string    `json:"cookie"`
	UriQuery                 string    `json:"uri_query"`
	EdgeResponseResultType   string    `json:"edge_response_result_type"`
	SslProtocol              string    `json:"ssl_protocol"`
	SslCipher                string    `json:"ssl_cipher"`
	EdgeResultType           string    `json:"edge_result_type"`
	ContentType              string    `json:"content_type"`
	ContentLength            int64     `json:"content_length"`
	EdgeDetailedResultType   string    `json:"edge_detailed_result_type"`
	Country                  string    `json:"country"`
	CacheBehaviorPathPattern string    `json:"cache_behavior_path_pattern"`
	Year                     int       `json:"year"`
	Month                    int       `json:"month"`
	Day                      int       `json:"day"`
}

// Used to manage varidic options
type Option func(c *Config)

// database configs
type Config struct {
	dsn string
	DB  *sqlx.DB
}

// New returns a new Config with the given options
func New(opts ...func(*Config)) (*Config, error) {
	config := &Config{}

	// apply options
	for _, opt := range opts {
		opt(config)
	}

	// dns must be set
	if config.dsn == "" {
		return nil, fmt.Errorf("dsn is required")
	}

	// connect to database
	if db, err := sqlx.Connect("trino", config.dsn); err != nil {
		return nil, err
	} else {
		config.DB = db
	}

	return config, nil
}

// SetDSN sets the database connection string
func SetDSN(dsn string) Option {
	return func(c *Config) {
		c.dsn = dsn
	}
}

// Close closes the database connection
func (c *Config) Close() error {
	return c.DB.Close()
}
