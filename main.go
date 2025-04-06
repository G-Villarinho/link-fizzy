package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/g-villarinho/link-fizz-api/config"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/router"
)

func main() {
	i := di.NewInjector()

	if err := config.LoadEnv(); err != nil {
		log.Fatalf("❌ Failed to load environment variables: %v", err)
	}

	db := initDeps(i)

	mux := router.SetupRoutes(i)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.APIPort),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("\n" +
			"🔗 Link Fizz API\n" +
			"🚀 Server started successfully!\n" +
			"📦 Listening on port " + config.Env.APIPort + "\n" +
			"⏳ Ready to shorten URLs at " + time.Now().Format("02/01/2006 15:04:05"))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Error during server shutdown: %v", err)
	}
	db.Close()
	log.Println("\n✅ Server and database connection shut down properly.")
}
