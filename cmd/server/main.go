package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/x0tf/server/internal/api"
	"github.com/x0tf/server/internal/config"
	"github.com/x0tf/server/internal/database/postgres"
	"github.com/x0tf/server/internal/gateway"
	"github.com/x0tf/server/internal/static"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize the application configuration
	cfg, usedDotEnv := config.Load()
	if !usedDotEnv {
		log.Info("NOTE: No .env file was found. This is no error and the application will use the systems environment variables.")
	}

	// Initialize the namespace service
	namespaces, err := postgres.NewNamespaceService(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if err = namespaces.InitializeTable(); err != nil {
		log.Fatal(err)
	}
	defer namespaces.Close()

	// Initialize the element service
	elements, err := postgres.NewElementService(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if err = elements.InitializeTable(); err != nil {
		log.Fatal(err)
	}
	defer elements.Close()

	// Initialize the invite service if invites are activated
	var invites *postgres.InviteService
	if cfg.Invites {
		invites, err = postgres.NewInviteService(cfg.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		if err = invites.InitializeTable(); err != nil {
			log.Fatal(err)
		}
		defer invites.Close()
	}

	// Start up the REST API
	restApi := &api.API{
		Settings: &api.Settings{
			Address:           cfg.APIAddress,
			RequestsPerMinute: cfg.APIRequestsPerMinute,
			Production:        static.ApplicationMode == "PROD",
			Version:           static.ApplicationVersion,
			AdminTokens:       cfg.AdminTokens,
			InvitesEnabled:    cfg.Invites,
		},
		Services: &api.Services{
			Namespaces: namespaces,
			Elements:   elements,
			Invites:    invites,
		},
	}
	go func() {
		if err := restApi.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	// Start up the gateway
	gw := &gateway.Gateway{
		Address:      cfg.GatewayAddress,
		Production:   static.ApplicationMode == "PROD",
		Namespaces:   namespaces,
		Elements:     elements,
		RootRedirect: cfg.GatewayRootRedirect,
	}
	go func() {
		if err := gw.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the program to exit
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sc

	// Gracefully shut down the REST API
	if err := restApi.Shutdown(); err != nil {
		log.Error(err)
	}

	// Gracefully shut down the gateway
	if err := gw.Shutdown(); err != nil {
		log.Error(err)
	}
}
