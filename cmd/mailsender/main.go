package main

import (
	"github.com/richard-on/mail-service/pkg/server"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/richard-on/mail-service/config"
	"github.com/richard-on/mail-service/pkg/logger"
)

var (
	version string
	build   string
)

// @title         Mail Service API
// @version       1.0
// @contact.name  Richard Ragusski
// @contact.url   https://richardhere.dev/
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host          localhost/
// @BasePath      mail/v1/
func main() {
	var err error

	log := logger.NewLogger(config.DefaultWriter,
		zerolog.TraceLevel,
		"auth-setup")

	config.GoDotEnv, err = strconv.ParseBool(os.Getenv("GODOTENV"))
	if err != nil {
		log.Infof("can't load env variable. Trying godotenv next. err: %v\n", err)
		config.GoDotEnv = true
	}

	if config.GoDotEnv {
		err = godotenv.Load()
		if err != nil {
			log.Fatal(err, "abort. Cannot load env variables using godotenv.")
		}
	}

	config.Init(log)

	if !fiber.IsChild() {
		log.Info("env and logger setup complete")
		log.Infof("mail-service - version: %v; build: %v; FiberPrefork: %v; MaxCPU: %v", version, build, config.FiberPrefork, config.MaxCPU)
	}

	runtime.GOMAXPROCS(config.MaxCPU)

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryInfo.DSN,
		TracesSampleRate: config.SentryInfo.TSR,
	})
	if err != nil {
		log.Fatal(err, "sentry init failed")
	}
	defer sentry.Flush(2 * time.Second)

	if !fiber.IsChild() {
		log.Info("sentry setup complete")
	}

	// Start Fiber server
	server := server.NewApp()
	server.Run()

	sentry.Recover()
}
