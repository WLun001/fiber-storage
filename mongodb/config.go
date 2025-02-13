package mongodb

import (
	"crypto/tls"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Config defines the config for storage.
type Config struct {
	//https://docs.mongodb.com/manual/reference/connection-string/
	URI        string
	Database   string
	Collection string

	// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.4.2/mongo/options#ClientOptions
	AppName                  string
	Auth                     options.Credential
	AutoEncryptionOptions    *options.AutoEncryptionOptions
	ConnectTimeout           time.Duration
	Compressors              []string
	Dialer                   options.ContextDialer
	Direct                   bool
	DisableOCSPEndpointCheck bool
	HeartbeatInterval        time.Duration
	Hosts                    []string
	LocalThreshold           time.Duration
	MaxConnIdleTime          time.Duration
	MaxPoolSize              uint64
	MinPoolSize              uint64
	PoolMonitor              *event.PoolMonitor
	Monitor                  *event.CommandMonitor
	ReadConcern              *readconcern.ReadConcern
	ReadPreference           *readpref.ReadPref
	Registry                 *bsoncodec.Registry
	ReplicaSet               string
	RetryReads               bool
	RetryWrites              bool
	ServerSelectionTimeout   time.Duration
	SocketTimeout            time.Duration
	TLSConfig                *tls.Config
	WriteConcern             *writeconcern.WriteConcern
	ZlibLevel                int
	ZstdLevel                int
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	URI:        "mongodb://127.0.0.1:27017",
	Database:   "fiber",
	Collection: "fiber",
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.URI == "" {
		cfg.URI = ConfigDefault.URI
	}
	if cfg.Database == "" {
		cfg.Database = ConfigDefault.Database
	}
	if cfg.Collection == "" {
		cfg.Collection = ConfigDefault.Collection
	}
	return cfg
}
