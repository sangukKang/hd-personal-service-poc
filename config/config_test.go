package config

import (
	"reflect"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			want: &Config{
				Host:    "127.0.0.1",
				Port:    "12001",
				RPCPort: "12000",
				Develop: true,
				S3Config: S3Config{
					Port:         "9999",
					S3Region:     "eu-central-1",
					S3BucketName: "test",
				},
				Log: LoggerConfig{
					Type:  "stdout",
					Level: "info",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadConfigFile()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
