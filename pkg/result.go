package pkg

import (
	"encoding/json"
	"time"
)

type Result struct {
	DownloadSpeed float64 `json:"downloadSpeed"`
	UploadSpeed   float64 `json:"uploadSpeed"`
	Jitter        float64 `json:"jitter"`
	Latency       float64 `json:"latency"`
}

type RawResult struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
		Low     float64 `json:"low"`
		High    float64 `json:"high"`
	} `json:"ping"`
	Download struct {
		Bandwidth int `json:"bandwidth"`
		Bytes     int `json:"bytes"`
		Elapsed   int `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"download"`
	Upload struct {
		Bandwidth int `json:"bandwidth"`
		Bytes     int `json:"bytes"`
		Elapsed   int `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"upload"`
	Isp       string `json:"isp"`
	Interface struct {
		InternalIp string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIp string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		Id       int    `json:"id"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		Ip       string `json:"ip"`
	} `json:"server"`
	Result struct {
		Id        string `json:"id"`
		Url       string `json:"url"`
		Persisted bool   `json:"persisted"`
	} `json:"result"`
}

func FromBytesToResult(data []byte) (*Result, error) {
	var rawResult RawResult
	if err := json.Unmarshal(data, &rawResult); err != nil {
		return nil, err
	}
	return rawResult.toResult(), nil

}

func (r RawResult) toResult() *Result {
	return &Result{
		DownloadSpeed: float64(r.Download.Bandwidth) / 125000,
		UploadSpeed:   float64(r.Upload.Bandwidth) / 125000,
		Jitter:        r.Ping.Jitter,
		Latency:       r.Ping.Latency,
	}
}
