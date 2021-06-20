package config

import (
	"flag"
	"fmt"
)

// Default values for the server's settings
const (
	defaultAddress    = ":8080"
	defaultForbidCORS = false
	defaultDebug      = false

	defaultSnakeServer = "localhost:8080"
	defaultWSS         = false

	defaultBotsLimit = 100

	defaultLogEnableJSON = false
	defaultLogLevel      = "info"
)

// Flag labels
const (
	flagLabelAddress    = "address"
	flagLabelForbidCORS = "forbid-cors"
	flagLabelDebug      = "debug"

	flagLabelSnakeServer = "snake-server"
	flagLabelWSS         = "wss"

	flagLabelBotsLimit = "bots-limit"

	flagLabelLogEnableJSON = "log-json"
	flagLabelLogLevel      = "log-level"
)

// Flag usage descriptions
const (
	flagUsageAddress    = "address to listen to"
	flagUsageForbidCORS = "forbid cross-origin resource sharing"
	flagUsageDebug      = "add profiling routes"

	flagUsageSnakeServer = "snake server's address: host:port"
	flagUsageWSS         = "use secure web-socket connection"

	flagUsageBotsLimit = "overall bots limit"

	flagUsageLogEnableJSON = "use json logging format"
	flagUsageLogLevel      = "log level: panic, fatal, error, warning, info or debug"
)

// Label names
const (
	fieldLabelAddress    = "address"
	fieldLabelForbidCORS = "forbid-cors"
	fieldLabelDebug      = "debug"

	fieldLabelSnakeServer = "snake server"
	fieldLabelWSS         = "wss"

	fieldLabelBotsLimit = "bots-limit"

	fieldLabelLogEnableJSON = "log-json"
	fieldLabelLogLevel      = "log-level"
)

// Server structure contains configurations for the server
type Server struct {
	Address    string
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

// Config is a base server configuration structure
type Config struct {
	Server Server
	Target Target
	Log    Log
	Bots   Bots
}

// Fields returns a map of all configurations
func (c Config) Fields() map[string]interface{} {
	return map[string]interface{}{
		fieldLabelAddress:    c.Server.Address,
		fieldLabelForbidCORS: c.Server.ForbidCORS,
		fieldLabelDebug:      c.Server.Debug,

		fieldLabelSnakeServer: c.Target.Address,
		fieldLabelWSS:         c.Target.WSS,

		fieldLabelBotsLimit: c.Bots.Limit,

		fieldLabelLogEnableJSON: c.Log.EnableJSON,
		fieldLabelLogLevel:      c.Log.Level,
	}
}

// Default settings
var defaultConfig = Config{
	Server: Server{
		Address:    defaultAddress,
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

	if err := flagSet.Parse(args); err != nil {
		return defaults, fmt.Errorf("cannot parse flags: %s", err)
	}

	return config, nil
}
