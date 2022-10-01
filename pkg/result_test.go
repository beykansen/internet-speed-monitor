package pkg

import (
	"reflect"
	"testing"
)

var result = `{
    "type": "result",
    "timestamp": "1970-01-01T00:00:00Z",
    "ping": {
        "jitter": 0.5,
        "latency": 4.5,
        "low": 4.435,
        "high": 6.720
    },
    "download": {
        "bandwidth": 375000,
        "bytes": 923134548,
        "elapsed": 8814,
        "latency": {
            "iqm": 5.105,
            "low": 3.498,
            "high": 10.691,
            "jitter": 0.584
        }
    },
    "upload": {
        "bandwidth": 250000,
        "bytes": 83599720,
        "elapsed": 12720,
        "latency": {
            "iqm": 4.123,
            "low": 2.870,
            "high": 343.622,
            "jitter": 6.049
        }
    },
    "isp": "My ISP",
    "interface": {
        "internalIp": "192.168.0.1",
        "name": "en7",
        "macAddr": "11:11:11:11:11:11",
        "isVpn": false,
        "externalIp": "1.1.1.1"
    },
    "server": {
        "id": 12345,
        "host": "speedtest.example.com",
        "port": 8080,
        "name": "name",
        "location": "Mars",
        "country": "Center",
        "ip": "1.2.3.4"
    },
    "result": {
        "id": "1a5b5165-7ebc-4be0-b930-51e96f87b166",
        "url": "https://www.speedtest.net/result/c/1a5b5165-7ebc-4be0-b930-51e96f87b166",
        "persisted": true
    }
}`

func TestFromBytesToResult(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Result
		wantErr bool
	}{
		{name: "should_work", args: struct{ data []byte }{data: []byte(result)}, want: &Result{
			DownloadSpeed: 3,
			UploadSpeed:   2,
			Jitter:        0.5,
			Latency:       4.5,
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromBytesToResult(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBytesToResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromBytesToResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}
