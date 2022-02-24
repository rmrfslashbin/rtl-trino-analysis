package fetch

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
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
		spew.Dump(config.dsn)
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
