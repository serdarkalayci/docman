package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	rest "github.com/serdarkalayci/docman/api/document/adapters/comm/rest"
	"github.com/serdarkalayci/docman/api/document/adapters/tracing"

	"github.com/nicholasjackson/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/docman/api/document/adapters/data/postgres"

	util "github.com/serdarkalayci/docman/api/document/util"
)

var bindAddress = env.String("BASE_URL", false, ":5500", "Bind address for rest server")
var otlpEndpoint = env.String("OTLP_ENDPOINT", false, "localhost:30080", "OpenTelemetry Collector endpoint")

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	util.SetLogLevels()
	env.Parse()
	//dbContext := memory.NewDataContext()
	dbContext, err := postgres.NewDataContext()
	if err != nil {
		log.Fatal().Msgf("error received from data source. Quitting. Error message is %s", err.Error())
		os.Exit(1)
	}
	//s := rest.NewAPIContext(dbContext, bindAddress)
	s := rest.NewAPIContext(bindAddress, dbContext.HealthRepository, dbContext.DocumentRepository)
	tpShutdown, err := tracing.InitProvider(*otlpEndpoint)	
	if err != nil {
		log.Fatal().Msgf("error received from tracing provider. Quitting. Error message is %s", err.Error())
		os.Exit(1)
	}
	// start the http server
	go func() {
		log.Debug().Msgf("Starting server on %s", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			log.Error().Err(err).Msg("Error starting rest server")
			if err := tpShutdown(context.Background()); err != nil {
				log.Error().Err(err)
			}
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Info().Msgf("Got signal: %s", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
