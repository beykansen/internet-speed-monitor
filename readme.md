# Internet Speed Monitor
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/beykansen/internet-speed-monitor/build_and_push)
[![CodeQL](https://github.com/beykansen/internet-speed-monitor/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/beykansen/internet-speed-monitor/actions/workflows/codeql-analysis.yml)
![GitHub last commit](https://img.shields.io/github/last-commit/beykansen/internet-speed-monitor)

This app uses ``Speedtest CLI`` to monitor your internet speed within desired interval 
and expose them as prometheus metrics and callback your desired endpoint with results after each run.

### Prerequisites:
If your plan is using the binary, ``Speedtest CLI`` needs to be installed on your machine. Click [here]((https://www.speedtest.net/tr/apps/cli)) to install SPEEDTEST® CLI.


### How to use:
Download binary from releases and run:
```bash
./internet-speed-monitor --interval 15 --port 8080 --callback https://example.com/callback
```
or via docker:
```bash
docker run -d \
  --name internet-speed-monitor \
  -p 8080:8080 \
  -e PORT=8080 \
  -e INTERVAL=15 \
  -e CALLBACK='https://example.com/callback' \
  ghcr.io/beykansen/internet-speed-monitor:latest
```

### Callback Details:
Method: ``POST`` <br />
Example payload:
```json
{
  "downloadSpeed": 1024,
  "uploadSpeed": 512,
  "jitter": 4.5,
  "latency": 7
}
```

### Grafana Dashboard:
[Internet Speed Monitor Dashboard](https://grafana.com/grafana/dashboards/17105)
