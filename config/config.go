package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	confOnce sync.Once
	config   *Config
)

type Config struct {
	Host        string       `json:"host"`
	Port        string       `json:"port"`
	RPCPort     string       `json:"rpcport"`
	Develop     bool         `json:"develop"`
	S3Config    S3Config     `json:"s3Config"`
	KafkaConfig KafkaConfig  `json:"kafkaConfig"`
	Log         LoggerConfig `json:"log"`
}

// logger config
type LoggerConfig struct {
	Type  string `json:"type"`  // options: file, stdout
	Level string `json:"level"` // debug, info, error...
}

type S3Config struct {
	Port         string `json:"port" required:"true" default:"9999"`
	S3Region     string `json:"s3Region" required:"true" default:"eu-central-1"`
	S3BucketName string `json:"s3BucketName" required:"true"`
}

type KafkaConfig struct {
	ConsumerEvent KafkaReadConfig  `json:"consumer_stream" mapstructure:"consumer_stream"`
	ProducerEvent KafkaWriteConfig `json:"producer_stream" mapstructure:"producer_stream"`
}

type KafkaWriteConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type KafkaReadConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
	GroupID string   `json:"groupId"`
	Timeout int64    `json:"timeout"`
}

type FileOutConfig struct {
	Path string `json:"path"`
}

// config 파일(yaml)을 읽고 global struct 에 저장한다.
func LoadConfigFile() *Config {
	filename := "env.yaml"
	confOnce.Do(func() {
		viper.SetConfigType("yaml")
		viper.SetConfigFile(filename)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(err)
		}
		err = viper.Unmarshal(&config)

		if err != nil {
			panic(err)
		}
		fmt.Println(config)
	})
	return config
}
