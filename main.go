package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/beykansen/internet-speed-monitor/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "internet_speed_monitor"

var (
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
	go func() {
		for {
			rand.Seed(time.Now().UnixNano())

			result := runSpeedTest()
			downloadSpeed.WithLabelValues(hostname).Set(result.DownloadSpeed)
			uploadSpeed.WithLabelValues(hostname).Set(result.UploadSpeed)
			jitter.WithLabelValues(hostname).Set(result.Jitter)
			latency.WithLabelValues(hostname).Set(result.Latency)

			if len(strings.TrimSpace(args.CallbackUrl)) > 0 {
				pkg.Callback(args.CallbackUrl, result)
			}
			time.Sleep(time.Duration(args.Interval)*time.Minute + time.Duration(rand.Intn(120))*time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", args.Port), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func runSpeedTest() *pkg.Result {
	speedTestPath := "/usr/bin/speedtest"

	if output, err := exec.Command("which", "speedtest").Output(); err != nil {
		panic(err)
	} else {
		speedTestPath = strings.ReplaceAll(string(output), "\n", "")
		speedTestPath = strings.ReplaceAll(speedTestPath, "\r\n", "")
	}
	output, err := exec.Command(speedTestPath, "--accept-license", "--accept-gdpr", "--f", "json-pretty").Output()
	if err != nil {
		panic(err)
	}

	result, err := pkg.FromBytesToResult(output)
	if err != nil {
		panic(err)
	}

	log.Printf("Speed Test Finished. Download Speed: %fmbps Upload Speed: %fmbps Jitter: %fms Latency: %fms", result.DownloadSpeed, result.UploadSpeed, result.Jitter, result.Latency)
	return result
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}
