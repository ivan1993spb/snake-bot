package config

import (
	"flag"
	"fmt"
)

// Default values for the server's settings
const (
	defaultAddress    = ":8080"
	defaultJWTSecret  = "/etc/snake-bot/jwt-secret.base64"
	defaultForbidCORS = false
	defaultDebug      = false

	defaultSnakeServer = "localhost:8080"
	defaultWSS         = false

	defaultBotsLimit = 100

	defaultLogEnableJSON = false
	defaultLogLevel      = "info"

	defaultStoragePath = ""
)

// Flag labels
const (
	flagLabelAddress    = "address"
	flagLabelJWTSecret  = "jwt-secret"
	flagLabelForbidCORS = "forbid-cors"
	flagLabelDebug      = "debug"

	flagLabelSnakeServer = "snake-server"
	flagLabelWSS         = "wss"

	flagLabelBotsLimit = "bots-limit"

	flagLabelLogEnableJSON = "log-json"
	flagLabelLogLevel      = "log-level"

	flagLabelStoragePath = "storage"
)

// Flag usage descriptions
const (
	flagUsageAddress    = "address to listen to"
	flagUsageJWTSecret  = "path to a base64 encoded secret for JWT signing"
	flagUsageForbidCORS = "forbid cross-origin resource sharing"
	flagUsageDebug      = "add profiling routes"

	flagUsageSnakeServer = "snake server's address: host:port"
	flagUsageWSS         = "use secure web-socket connection"

	flagUsageBotsLimit = "overall bots limit"

	flagUsageLogEnableJSON = "use json logging format"
	flagUsageLogLevel      = "log level: panic, fatal, error, warning, info or debug"

	flagUsageStoragePath = "path to a state file"
)

// Server structure contains configurations for the server
type Server struct {
	Address    string
	JWTSecret  string
	ForbidCORS bool
	Debug      bool
}

type Target struct {
	Address string
	WSS     bool
}

type Bots struct {
	Limit int
}

// Log structure defines preferences for logging
type Log struct {
	EnableJSON bool
	Level      string
}

// Storage structure defines preferences for storage
type Storage struct {
	Path string
}

// Config is a base server configuration structure
type Config struct {
	Server  Server
	Target  Target
	Log     Log
	Bots    Bots
	Storage Storage
}

// Fields returns a map of all configurations
func (c Config) Fields() map[string]interface{} {
	return map[string]interface{}{
		flagLabelAddress:    c.Server.Address,
		flagLabelJWTSecret:  c.Server.JWTSecret,
		flagLabelForbidCORS: c.Server.ForbidCORS,
		flagLabelDebug:      c.Server.Debug,

		flagLabelSnakeServer: c.Target.Address,
		flagLabelWSS:         c.Target.WSS,

		flagLabelBotsLimit: c.Bots.Limit,

		flagLabelLogEnableJSON: c.Log.EnableJSON,
		flagLabelLogLevel:      c.Log.Level,

		flagLabelStoragePath: c.Storage.Path,
	}
}

// Default settings
var defaultConfig = Config{
	Server: Server{
		Address:    defaultAddress,
		JWTSecret:  defaultJWTSecret,
		ForbidCORS: defaultForbidCORS,
		Debug:      defaultDebug,
	},

	Target: Target{
		Address: defaultSnakeServer,
		WSS:     defaultWSS,
	},

	Bots: Bots{
		Limit: defaultBotsLimit,
	},

	Log: Log{
		EnableJSON: defaultLogEnableJSON,
		Level:      defaultLogLevel,
	},

	Storage: Storage{
		Path: defaultStoragePath,
	},
}

// DefaultConfig returns configuration by default
func DefaultConfig() Config {
	return defaultConfig
}

// ParseFlags parses flags and returns a config based on the default configuration
func ParseFlags(flagSet *flag.FlagSet, args []string, defaults Config) (Config, error) {
	if flagSet.Parsed() {
		panic("program composition error: the provided FlagSet has been parsed")
	}

	config := defaults

	// Address
	flagSet.StringVar(&config.Server.Address, flagLabelAddress,
		defaults.Server.Address, flagUsageAddress)
	flagSet.StringVar(&config.Server.JWTSecret, flagLabelJWTSecret,
		defaults.Server.JWTSecret, flagUsageJWTSecret)
	flagSet.BoolVar(&config.Server.ForbidCORS, flagLabelForbidCORS,
		defaults.Server.ForbidCORS, flagUsageForbidCORS)
	flagSet.BoolVar(&config.Server.Debug, flagLabelDebug,
		defaults.Server.Debug, flagUsageDebug)

	flagSet.StringVar(&config.Target.Address, flagLabelSnakeServer,
		defaults.Target.Address, flagUsageSnakeServer)
	flagSet.BoolVar(&config.Target.WSS, flagLabelWSS,
		defaults.Target.WSS, flagUsageWSS)

	flagSet.IntVar(&config.Bots.Limit, flagLabelBotsLimit,
		defaults.Bots.Limit, flagUsageBotsLimit)

	// Logging
	flagSet.BoolVar(&config.Log.EnableJSON, flagLabelLogEnableJSON,
		defaults.Log.EnableJSON, flagUsageLogEnableJSON)
	flagSet.StringVar(&config.Log.Level, flagLabelLogLevel,
		defaults.Log.Level, flagUsageLogLevel)

	// Storage
	flagSet.StringVar(&config.Storage.Path, flagLabelStoragePath,
		defaults.Storage.Path, flagUsageStoragePath)

	if err := flagSet.Parse(args); err != nil {
		return defaults, fmt.Errorf("cannot parse flags: %s", err)
	}

	return config, nil
}
