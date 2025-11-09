package monitoring

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/adaken4/clean-town/internal/database"
)

func StartMetricsMonitoring(ctx context.Context, db *sql.DB) {
	metricsChan := make(chan *database.Metrics, 10)

	// Start background monitoring (every 30 seconds)
	go database.MonitorMetrics(ctx, db, 30*time.Second, metricsChan)

	// Process metrics in a separate goroutine
	go func() {
		for metrics := range metricsChan {
			if !metrics.Healthy {
				log.Printf("⚠️  DATABASE UNHEALTHY - Ping took %v", metrics.PingDuration)
			}

			// Log if too many connections are waiting
			if metrics.WaitCount > 100 {
				log.Printf("⚠️  High wait count: %d connections waiting, avg wait: %v",
					metrics.WaitCount, metrics.WaitDuration/time.Duration(metrics.WaitCount))
			}

			// Log pool utilization
			utilization := float64(metrics.InUse) / float64(metrics.MaxOpenConnections) * 100
			if utilization > 80 {
				log.Printf("⚠️  High pool utilization: %.1f%% (%d/%d connections in use)",
					utilization, metrics.InUse, metrics.MaxOpenConnections)
			}

			log.Printf("📊 DB Metrics - Open: %d, InUse: %d, Idle: %d, WaitCount: %d, Healthy: %v",
				metrics.OpenConnections, metrics.InUse, metrics.Idle, metrics.WaitCount, metrics.Healthy)
		}
	}()
}

