package main

import (
    "log"
    "os"

    "github.com/MystiaFin/scalar/backend/internal/db"
    "github.com/MystiaFin/scalar/backend/internal/email"
    "github.com/MystiaFin/scalar/backend/internal/routes"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/robfig/cron/v3"
)

const (
    address  = ":7225"
    syncDays = 14
)

func runSync() {
    summary, err := email.Sync(syncDays, email.Extract)
    if err != nil {
        log.Printf("sync skipped: %v", err)
        return
    }

    log.Printf("sync done: fetched=%d saved=%d skipped=%d failed=%d",
        summary.Fetched, summary.Saved, summary.Skipped, summary.Failed)
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Printf(".env not loaded: %v", err)
    }

    db.Open()
    defer db.DB.Close()

    db.Migrate()

    if os.Getenv("SCALAR_EMAIL") != "" && os.Getenv("SCALAR_APP_PASSWORD") != "" {
        go runSync()

        c := cron.New()
        if _, err := c.AddFunc("@every 30m", runSync); err != nil {
            log.Fatalf("failed to schedule sync: %v", err)
        }
        c.Start()
        defer c.Stop()
    }

    r := gin.Default()
    routes.Register(r)

    log.Printf("scalar running on %s", address)

    if err := r.Run(address); err != nil {
        log.Fatalf("server stopped: %v", err)
    }
}
