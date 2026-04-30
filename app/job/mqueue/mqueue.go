package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"amigo-api/app/job/mqueue/internal/config"
	"amigo-api/app/job/mqueue/internal/handler"
	mqhandler "amigo-api/app/job/mqueue/internal/handler/mqueue"
	"amigo-api/app/job/mqueue/internal/svc"
	"amigo-api/common/mqueue"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/mqueue.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// Create Redis client for asynq
	redisOpt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.MQueue.RedisHost, c.MQueue.RedisPort),
		Password: c.MQueue.RedisPass,
		DB:       c.MQueue.RedisDB,
	}

	// Create mqueue config
	mqConfig := &mqueue.QueueConfig{
		RedisOpt:      redisOpt,
		Queues:        c.MQueue.Queues,
		Concurrency:   c.MQueue.Concurrency,
		StrictMode:    false,
		SyncTimeout:   c.MQueue.SyncTimeout,
		RetryDelay:    c.MQueue.RetryDelay,
		MaxRetry:      c.MQueue.MaxRetry,
		Timeout:       c.MQueue.Timeout,
		DeadQueueName: c.MQueue.DeadQueue,
		ServerName:    c.MQueue.ServerName,
	}

	// Initialize global mqueue
	if err := mqueue.InitGlobalMQueue(redisOpt, mqConfig); err != nil {
		logx.Errorf("Failed to initialize global mqueue: %v", err)
	}

	// Create Redis client for handler
	redisClient := redis.NewClient(redisOpt)
	mqhandler.InitRedis(redisClient)

	// Create consumer and register handlers
	consumer := mqueue.NewRedisConsumer(redisOpt, mqConfig)
	registerHandlers(consumer)

	// Start consumer
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		logx.Errorf("Failed to start consumer: %v", err)
	}

	// Start asynq monitoring server in background
	go startMonitoringServer(redisOpt, &c)

	// Setup REST server
	server := rest.MustNewServer(c.RestConf)
	defer func() {
		server.Stop()
		consumer.Stop()
		mqueue.Shutdown()
		redisClient.Close()
	}()

	handlerCtx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, handlerCtx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf("Asynq Monitor available at http://%s:%d/monitor\n", c.Host, c.MQueue.MonitorPort)
	server.Start()
}

func registerHandlers(consumer *mqueue.RedisConsumer) {
	consumer.RegisterHandler("send_sms", mqhandler.NewSendSmsHandler())
}

// startMonitoringServer starts the asynq monitoring web UI
func startMonitoringServer(redisOpt *redis.Options, c *config.Config) {
	// Create Redis client for monitor
	rdb := redis.NewClient(redisOpt)

	// Create HTTP server for monitoring
	monitorMux := http.NewServeMux()

	// Dashboard HTML page
	monitorMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(dashboardHTML)
	})

	// Health check endpoint
	monitorMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Stats endpoint - JSON API using Redis directly
	monitorMux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats, err := getQueueStats(rdb)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(stats)
	})

	// Create server
	monitorAddr := fmt.Sprintf("%s:%d", c.Host, c.MQueue.MonitorPort)
	monitorServer := &http.Server{
		Addr:    monitorAddr,
		Handler: monitorMux,
	}

	logx.Infof("Asynq monitor starting on %s", monitorAddr)

	// Start server in goroutine
	go func() {
		if err := monitorServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.Errorf("Monitor server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown monitor server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := monitorServer.Shutdown(ctx); err != nil {
		logx.Errorf("Monitor server shutdown error: %v", err)
	}

	// Close connections
	rdb.Close()

	logx.Info("Asynq monitor stopped")
}

// getQueueStats returns queue statistics as JSON using Redis directly
func getQueueStats(rdb *redis.Client) ([]byte, error) {
	queues := []string{"critical", "default", "low"}
	stats := make(map[string]interface{})

	ctx := context.Background()
	for _, q := range queues {
		// Get queue length using asynq key patterns
		pendingKey := fmt.Sprintf("asynq:pending:%s", q)

		pendingLen, err := rdb.LLen(ctx, pendingKey).Result()
		if err != nil {
			pendingLen = 0
		}

		stats[q] = map[string]interface{}{
			"pending": pendingLen,
			"active":  0,
		}
	}

	return json.Marshal(map[string]interface{}{"queues": stats})
}

// Dashboard HTML
var dashboardHTML = []byte(`<!DOCTYPE html>
<html>
<head>
    <title>Amigo Asynq Monitor</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; color: #333; }
        .header { background: #2c3e50; color: white; padding: 20px; }
        .header h1 { font-size: 24px; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .card { background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .card-header { padding: 15px 20px; border-bottom: 1px solid #eee; font-weight: bold; }
        .card-body { padding: 20px; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; }
        .stat-card { background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); padding: 20px; }
        .stat-card h3 { color: #666; font-size: 14px; text-transform: uppercase; margin-bottom: 10px; }
        .stat-card .value { font-size: 32px; font-weight: bold; color: #2c3e50; }
        .queue-critical { border-left: 4px solid #e74c3c; }
        .queue-default { border-left: 4px solid #3498db; }
        .queue-low { border-left: 4px solid #95a5a6; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #eee; }
        th { background: #f8f9fa; font-weight: 600; }
        .refresh { display: inline-block; padding: 10px 20px; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer; }
        .refresh:hover { background: #2980b9; }
        .nav { display: flex; gap: 10px; margin-bottom: 20px; }
        .nav a { padding: 10px 20px; background: white; border-radius: 4px; text-decoration: none; color: #333; }
        .nav a:hover { background: #e0e0e0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Amigo Asynq Monitor</h1>
    </div>
    <div class="container">
        <div class="nav">
            <a href="/">Dashboard</a>
            <a href="/stats">Stats JSON</a>
        </div>

        <div class="stats-grid" id="stats"></div>

        <div class="card">
            <div class="card-header">
                <button class="refresh" onclick="loadStats()">Refresh</button>
                <span style="margin-left: 20px; color: #666;">Auto-refresh: 5s</span>
            </div>
            <div class="card-body">
                <h3>Queue Details</h3>
                <div id="queues"></div>
            </div>
        </div>
    </div>

    <script>
        async function loadStats() {
            try {
                const res = await fetch('/stats');
                const data = await res.json();

                const statsDiv = document.getElementById('stats');
                const queuesDiv = document.getElementById('queues');

                let html = '';
                let tableHtml = '<table><thead><tr><th>Queue</th><th>Pending</th><th>Active</th></tr></thead><tbody>';

                for (const [queue, info] of Object.entries(data.queues)) {
                    const queueClass = queue === 'critical' ? 'queue-critical' : queue === 'default' ? 'queue-default' : 'queue-low';
                    const stats = typeof info === 'object' ? info : { pending: 0, active: 0 };

                    html += '<div class="stat-card ' + queueClass + '">' +
                        '<h3>' + queue + ' Queue</h3>' +
                        '<div class="value">' + (stats.pending || 0) + '</div>' +
                        '<p>Pending Tasks</p></div>';

                    tableHtml += '<tr><td>' + queue + '</td><td>' + (stats.pending || 0) + '</td><td>' + (stats.active || 0) + '</td></tr>';
                }

                tableHtml += '</tbody></table>';
                statsDiv.innerHTML = html;
                queuesDiv.innerHTML = tableHtml;
            } catch (e) {
                console.error('Error loading stats:', e);
            }
        }

        // Auto refresh
        loadStats();
        setInterval(loadStats, 5000);
    </script>
</body>
</html>`)
