package config

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Config Application config parameters
type Config struct {
	HTTP     HTTPConfig
	Log      LogConfig
	Redis    RedisConfig
	Postgres PostgresConfig
	Counter  Counter
}

// HTTPConfig HTTP config parameters
type HTTPConfig struct {
	Address string
}

// LogConfig Logging configuration
type LogConfig struct {
	Level      string
	Format     string
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

// RedisClusterConfig Redis configuration parameters
type RedisClusterConfig struct {
	Master     string
	Replica    string
	Password   string
	DB         int
	MaxRetries int
	Expiration time.Duration
}

// RedisConfig Redis configuration parameters
type RedisConfig struct {
	Address    string
	Password   string
	DB         int
	MaxRetries int
	Expiration time.Duration
}

type PostgresConfig struct {
	Host           string
	Port           int
	Database       string
	User           string
	Password       string
	SSLMode        string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetimConn int
}

type Counter struct {
	MachineID int
	Range     string
}

// InitLogging Initialize logging framework
func InitLogging(lc *LogConfig) {
	switch lc.Format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		fallthrough
	case "text":
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	switch lc.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		fallthrough
	case "info":
		log.SetLevel(log.InfoLevel)
	}

	if lc.Filename == "" {
		log.SetOutput(os.Stdout)
	} else {
		l := &lumberjack.Logger{
			Filename:   lc.Filename,
			MaxSize:    lc.MaxSize,
			MaxAge:     lc.MaxAge,
			MaxBackups: lc.MaxBackups,
			LocalTime:  lc.LocalTime,
			Compress:   lc.Compress,
		}

		log.SetOutput(l)
	}
}
