package pkg

import (
	"encoding/json"
)

type Result struct {
	DownloadSpeed float64 `json:"downloadSpeed"`
	UploadSpeed   float64 `json:"uploadSpeed"`
	Jitter        float64 `json:"jitter"`
	Latency       float64 `json:"latency"`
}

type RawResult struct {
	Ping struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth int `json:"bandwidth"`
	} `json:"download"`
	Upload struct {
		Bandwidth int `json:"bandwidth"`
	} `json:"upload"`
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
