package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var (
	text = pflag.String("text", "", "text to put on the webpage")
	addr = pflag.String("addr", ":8080", "address to listen on")
)

func main() {
	pflag.Parse()

	if *text == "" {
		log.Fatal("--text option is required!")
	}

	r := gin.Default()
	r.GET("/", TextHandler)
	r.GET("/health", HealthHandler)
	r.NoRoute(TextHandler)

	srv := http.Server{
		Addr:    *addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s\n", err)
	}

	log.Println("Server exiting")
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func TextHandler(c *gin.Context) {
	c.String(http.StatusOK, *text)
}
