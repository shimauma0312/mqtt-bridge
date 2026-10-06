package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/go-sql-driver/mysql"

	"readsrv/internal/httpapi"
	"readsrv/internal/store"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if n, err := strconv.Atoi(v); err == nil {
		return time.Duration(n) * time.Second
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("%s: invalid duration %q", key, v)
	}
	return d
}

func main() {
	healthcheck := flag.Bool("healthcheck", false, "GET /healthz して終了コードで返す（コンテナ用）")
	flag.Parse()

	listen := env("LISTEN_ADDR", ":8080")
	if *healthcheck {
		port := listen[strings.LastIndex(listen, ":")+1:]
		resp, err := http.Get("http://127.0.0.1:" + port + "/healthz")
		if err != nil || resp.StatusCode != 200 {
			os.Exit(1)
		}
		return
	}

	cfg := mysql.NewConfig()
	cfg.User = env("DB_USER", "reader")
	cfg.Passwd = env("DB_PASSWORD", "reader")
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%s", env("DB_HOST", "127.0.0.1"), env("DB_PORT", "3306"))
	cfg.DBName = env("DB_NAME", "bridge")
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	srv := httpapi.New(store.New(db, env("GPS_LABEL", "gps")), httpapi.Config{
		OfflineThreshold: envDuration("OFFLINE_THRESHOLD", 300*time.Second),
		TripGap:          envDuration("TRIP_GAP", 10*time.Minute),
	})
	hs := &http.Server{Addr: listen, Handler: srv.Handler(), ReadHeaderTimeout: 5 * time.Second}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		hs.Shutdown(ctx)
	}()
	log.Printf("listening on %s", listen)
	if err := hs.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
