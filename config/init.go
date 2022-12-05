package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/richard-on/mail-service/pkg/logger"
)

var LogInfo struct {
	Output        string
	Level         zerolog.Level
	File          string
	ConsoleWriter bool
}

var SentryInfo struct {
	DSN string
	TSR float64
}

var Mailgun struct {
	Host string
	Pass string
}

var SMTP struct {
	Host string
	Port string
}

var Env string
var GoDotEnv bool
var FiberPrefork bool
var MaxCPU int
var Host string
var SecureCookie bool

func Init(log logger.Logger) {
	var err error

	Env = os.Getenv("ENV")

	Host = os.Getenv("HOST")

	Mailgun.Host = os.Getenv("MAILGUN_HOST")
	Mailgun.Pass = os.Getenv("MAILGUN_PASS")

	SMTP.Host = os.Getenv("SMTP_HOST")
	SMTP.Port = os.Getenv("SMTP_PORT")

	SecureCookie, err = strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	if err != nil {
		log.Infof("SECURE_COOKIE init: %v", err)
	}

	GoDotEnv, err = strconv.ParseBool(os.Getenv("GODOTENV"))
	if err != nil {
		log.Infof("GODOTENV init: %v", err)
	}

	FiberPrefork, err = strconv.ParseBool(os.Getenv("FIBER_PREFORK"))
	if err != nil {
		log.Infof("FIBER_PREFORK init: %v", err)
	}

	MaxCPU, err = strconv.Atoi(os.Getenv("MAX_CPU"))
	if err != nil {
		log.Infof("MAX_CPU init: %v", err)
	}

	LogInfo.Output = os.Getenv("LOG_OUTPUT")

	LogInfo.Level, err = zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Infof("LOG_LEVEL init: %v", err)
	}

	LogInfo.File = os.Getenv("LOG_FILE")

	LogInfo.ConsoleWriter, err = strconv.ParseBool(os.Getenv("LOG_CW"))
	if err != nil {
		log.Infof("LOG_CW init: %v", err)
	}

	SentryInfo.DSN = os.Getenv("SENTRY_DSN")

	SentryInfo.TSR, err = strconv.ParseFloat(os.Getenv("SENTRY_TSR"), 64)
	if err != nil {
		log.Infof("SENTRY_TSR init: %v", err)
	}

}
