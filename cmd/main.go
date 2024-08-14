package main

import (
	"fmt"
	"log"
	"net/http"
	"nvml-exporter/internal/metrics"
	"os"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.ExitErrHandler = func(context *cli.Context, err error) {
		cli.HandleExitCoder(err)
	}
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Value: 3213,
			Usage: "the port on for HTTP requests",
		},
	}

	app.Action = func(cctx *cli.Context) error {
		// Initialize NVML
		if err := nvml.Init(); err != nvml.SUCCESS {
			log.Fatalf("Failed to initialize NVML: %v", err)
		}
		defer nvml.Shutdown()

		// Register all metrics
		metrics.RegisterMetrics()

		// Start the HTTP server
		listen := fmt.Sprintf(":%d", cctx.Int("port"))
		log.Printf("Starting server on %s\n", listen)
		http.Handle("/metrics", promhttp.Handler())
		return http.ListenAndServe(listen, nil)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
