package main

import (
	"github.com/x0tf/server/internal/api"
	"github.com/x0tf/server/internal/config"
	"github.com/x0tf/server/internal/database/postgres"
	"github.com/x0tf/server/internal/static"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize the application configuration
	cfg, usedDotEnv := config.Load()
	if !usedDotEnv {
		log.Println("NOTE: No .env file was found. This is no error and the application will use the systems environment variables.")
	}

	// Initialize the namespace service
	namespaces, err := postgres.NewNamespaceService(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	if err = namespaces.InitializeTable(); err != nil {
		panic(err)
	}
	defer namespaces.Close()

	// Initialize the element service
	elements, err := postgres.NewElementService(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	if err = elements.InitializeTable(); err != nil {
		panic(err)
	}
	defer elements.Close()

	// Initialize the invite service if invites are activated
	var invites *postgres.InviteService
	if cfg.Invites {
		invites, err = postgres.NewInviteService(cfg.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		if err = invites.InitializeTable(); err != nil {
			panic(err)
		}
		defer invites.Close()
	}

	// Start up the REST API
	restApi := &api.API{
		Address:     cfg.APIAddress,
		Production:  static.ApplicationMode == "PROD",
		Version:     static.ApplicationVersion,
		Namespaces:  namespaces,
		Elements:    elements,
		Invites:     invites,
		AdminTokens: cfg.AdminTokens,
	}
	if invites == nil {
		restApi.Invites = nil
	}
	go func() {
		panic(restApi.Serve())
	}()

	// Wait for the program to exit
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sc
}
