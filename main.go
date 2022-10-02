package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/beykansen/internet-speed-monitor/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "internet_speed_monitor"
const fatalErrorCountThreshold = 5

var (
	errorCount    uint64
	hostname      = getHostName()
	downloadSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "download_speed",
	}, []string{"host"})
	uploadSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "upload_speed",
	}, []string{"host"})
	jitter = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "jitter",
	}, []string{"host"})
	latency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "latency",
	}, []string{"host"})
)

var args struct {
	Port        int    `arg:"required, env:PORT, -p,--port" help:"port to listen"`
	CallbackUrl string `arg:"env:CALLBACK, -c,--callback" help:"callback url"`
	Interval    int    `arg:"required, env:INTERVAL, -i,--interval" help:"interval as minutes"`
}

func main() {
	arg.MustParse(&args)
	log.Printf("Server started on :%d. Callback url: %s Interval: %d\n", args.Port, args.CallbackUrl, args.Interval)
	go collectMetrics()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/metrics", 301)
	})
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", args.Port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func collectMetrics() {
	for {
		rand.Seed(time.Now().UnixNano())
		if result, err := runSpeedTest(); err != nil {
			atomic.AddUint64(&errorCount, 1)
			if errorCount > fatalErrorCountThreshold {
				log.Fatalf("Got error while runnig speedtest and threshold exceeded. Terminating... %s \n", err.Error())
			}
			log.Printf("Got error while runnig speedtest %s \n", err.Error())
		} else {
			downloadSpeed.WithLabelValues(hostname).Set(result.DownloadSpeed)
			uploadSpeed.WithLabelValues(hostname).Set(result.UploadSpeed)
			jitter.WithLabelValues(hostname).Set(result.Jitter)
			latency.WithLabelValues(hostname).Set(result.Latency)

			if len(strings.TrimSpace(args.CallbackUrl)) > 0 {
				if err := pkg.Callback(args.CallbackUrl, result); err != nil {
					log.Printf("Got error while callback %s\n", err.Error())
				}
			}
		}
		time.Sleep(time.Duration(args.Interval)*time.Minute + time.Duration(rand.Intn(120))*time.Second)
	}
}

func runSpeedTest() (*pkg.Result, error) {
	speedTestPath := "/usr/bin/speedtest"

	if output, err := exec.Command("which", "speedtest").Output(); err != nil {
		return nil, err
	} else {
		speedTestPath = strings.ReplaceAll(string(output), "\n", "")
		speedTestPath = strings.ReplaceAll(speedTestPath, "\r\n", "")
	}
	output, err := exec.Command(speedTestPath, "--accept-license", "--accept-gdpr", "--f", "json-pretty").Output()
	if err != nil {
		return nil, err
	}

	result, err := pkg.FromBytesToResult(output)
	if err != nil {
		return nil, err
	}

	log.Printf("Speed Test Finished. Download Speed: %fmbps Upload Speed: %fmbps Jitter: %fms Latency: %fms\n", result.DownloadSpeed, result.UploadSpeed, result.Jitter, result.Latency)
	return result, nil
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hostname
}
