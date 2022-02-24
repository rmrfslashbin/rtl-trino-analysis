package mysql

import (
	"fmt"

	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/rtl"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Record *rtl.Record `gorm:"embedded"`
}

// Used to manage varidic options
type Option func(c *Config)

// database configs
type Config struct {
	dsn string
	db  *gorm.DB
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
	if db, err := gorm.Open(mysql.Open(config.dsn)); err != nil {
		return nil, err
	} else {
		config.db = db
		db.AutoMigrate(&Record{})
	}

	return config, nil
}

// SetDSN sets the database connection string
func SetDSN(dsn string) Option {
	return func(c *Config) {
		c.dsn = dsn
	}
}

// Insert inserts a record into the database
func (c *Config) Insert(record *Record) error {
	return c.db.Create(record).Error
}
