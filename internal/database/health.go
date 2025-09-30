package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Metrics holds database connection pool metrics
type Metrics struct {
	// Pool stats
	MaxOpenConnections int           `json:"max_open_connections"`
	OpenConnections    int           `json:"open_connections"`
	InUse              int           `json:"in_use"`
	Idle               int           `json:"idle"`
	WaitCount          int64         `json:"wait_count"`
	WaitDuration       time.Duration `json:"wait_duration"`
	MaxIdleClosed      int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed  int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed  int64         `json:"max_lifetime_closed"`

	// Health check
	Healthy      bool          `json:"healthy"`
	PingDuration time.Duration `json:"ping_duration"`
	LastCheck    time.Time     `json:"last_check"`
}

// CollectMetrics gathers current database pool statistics and health status
func CollectMetrics(db *sql.DB) *Metrics {
	stats := db.Stats()

	// Perform health check with timing
	start := time.Now()
	err := CheckDB(db)
	pingDuration := time.Since(start)

	return &Metrics{
		// Pool configuration and current state
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUse:              stats.InUse,
		Idle:               stats.Idle,

		// Wait statistics
		WaitCount:    stats.WaitCount,
		WaitDuration: stats.WaitDuration,

		// Connection closure counts
		MaxIdleClosed:     stats.MaxIdleClosed,
		MaxIdleTimeClosed: stats.MaxIdleTimeClosed,
		MaxLifetimeClosed: stats.MaxLifetimeClosed,

		// Health status
		Healthy:      err == nil,
		PingDuration: pingDuration,
		LastCheck:    time.Now(),
	}
}

// MonitorMetrics continuously collects metrics at specified intervals
// and sends them to the provided channel. Call cancel() to stop monitoring.
func MonitorMetrics(ctx context.Context, db *sql.DB, interval time.Duration, metricsChan chan<- *Metrics) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics := CollectMetrics(db)
			select {
			case metricsChan <- metrics:
			case <-ctx.Done():
				return
			}
		}
	}
}

func CheckDB(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

// HealthCheck performs comprehensive health check with metrics
type HealthCheck struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Metrics *Metrics      `json:"metrics"`
	Latency time.Duration `json:"latency_ms"`
}

func PerformHealthCheck(db *sql.DB) *HealthCheck {
	start := time.Now()
	metrics := CollectMetrics(db)
	latency := time.Since(start)

	hc := &HealthCheck{
		Metrics: metrics,
		Latency: latency,
	}

	if !metrics.Healthy {
		hc.Status = "unhealthy"
		hc.Message = fmt.Sprintf("database ping failed (took %v)", metrics.PingDuration)
		return hc
	}

	// Check for warning conditions
	utilization := float64(metrics.InUse) / float64(metrics.MaxOpenConnections)
	if utilization > 0.9 {
		hc.Status = "degraded"
		hc.Message = fmt.Sprintf("high connection pool utilization: %.1f%%", utilization*100)
		return hc
	}

	hc.Status = "healthy"
	return hc
}
