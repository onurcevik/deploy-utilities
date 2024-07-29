package prometheus

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

// CollectDockerContainerMetrics collects Docker container monitoring
func CollectDockerContainerMetrics(interval time.Duration) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for {
		containers, err := cli.ContainerList(ctx, container.ListOptions{})
		if err != nil {
			log.Printf("error listing containers: %s", err)
			continue
		}

		for _, c := range containers {
			stats, err := cli.ContainerStats(ctx, c.ID, false)
			if err != nil {
				log.Printf("Error inspecting container %s: %s", c.ID, err)
				continue
			}

			defer stats.Body.Close()

			var statsData container.Stats
			err = json.NewDecoder(stats.Body).Decode(&statsData)
			if err != nil {
				log.Printf("failed to decode stats for container %s: %v", c.ID, err)
				continue
			}

			// Collect CPU usage
			cpuDelta := statsData.CPUStats.CPUUsage.TotalUsage - statsData.PreCPUStats.CPUUsage.TotalUsage
			systemDelta := statsData.CPUStats.SystemUsage - statsData.PreCPUStats.SystemUsage
			cpuUsage := float64(cpuDelta) / float64(systemDelta) * float64(len(statsData.CPUStats.CPUUsage.PercpuUsage))
			containerCPUUsage.WithLabelValues(c.ID).Set(cpuUsage)

			// Collect memory usage
			memUsage := statsData.MemoryStats.Usage
			containerMemUsage.WithLabelValues(c.ID).Set(float64(memUsage))

			// You can collect other monitoring similarly

		}

		time.Sleep(interval)
	}
}

// StartMetricsServer starts an HTTP server to expose monitoring to Prometheus
func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
