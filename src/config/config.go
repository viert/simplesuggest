package config

import (
	logging "github.com/op/go-logging"
	"github.com/viert/properties"
	"strings"
)

type Config struct {
	Host       string
	Port       int
	GCPercent  int
	MaxCores   int
	MaxThreads int
	LogLevel   logging.Level
	Log        string
}

var (
	log           *logging.Logger = logging.MustGetLogger("suggest")
	defaultConfig *Config         = &Config{
		Host:       "0.0.0.0",
		Port:       7978,
		GCPercent:  100,
		MaxCores:   8,
		MaxThreads: 10000,
		LogLevel:   logging.DEBUG,
		Log:        "",
	}
)

func Load(filename string) *Config {
	props, err := properties.Load(filename)
	if err != nil {
		log.Error(err.Error())
		log.Notice("Using configuration defaults")
		return defaultConfig
	}
	config := new(Config)

	// Setting Host
	config.Host, err = props.GetString("main.host")
	if err != nil {
		config.Host = defaultConfig.Host
	}

	// Setting Port
	config.Port, err = props.GetInt("main.port")
	if err != nil {
		config.Port = defaultConfig.Port
	}

	// Setting Runtime configuration
	config.GCPercent, err = props.GetInt("runtime.gc_percent")
	if err != nil {
		config.GCPercent = defaultConfig.GCPercent
	}
	config.MaxCores, err = props.GetInt("runtime.max_cores")
	if err != nil {
		config.MaxCores = defaultConfig.MaxCores
	}
	config.MaxThreads, err = props.GetInt("runtime.max_threads")
	if err != nil {
		config.MaxThreads = defaultConfig.MaxThreads
	}

	// Setting logging
	config.Log, err = props.GetString("main.log")
	if err != nil {
		config.Log = defaultConfig.Log
	}

	logLevel, err := props.GetString("main.log_level")
	if err != nil {
		config.LogLevel = defaultConfig.LogLevel
	} else {
		switch strings.ToLower(logLevel) {
		case "debug":
			config.LogLevel = logging.DEBUG
		case "error":
			config.LogLevel = logging.ERROR
		case "info":
			config.LogLevel = logging.INFO
		case "critical":
			config.LogLevel = logging.CRITICAL
		case "notice":
			config.LogLevel = logging.NOTICE
		case "warning":
			config.LogLevel = logging.WARNING
		default:
			config.LogLevel = defaultConfig.LogLevel
		}
	}
	return config
}
