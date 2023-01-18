package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	CompressScale float32 `yaml:"compress_scale"` // 压缩比例缩放（小数）
	ResizeScale   float32 `yaml:"resize_scale"`   // 裁剪比例缩放（小数）
}

func CreateConfFile() error {
	configStr := "scale: 1"

	err := ioutil.WriteFile("./ass_conf.yaml", []byte(configStr), 0666)
	if err != nil {
		return err
	}
	return nil
}

// GetConfig 获取配置
func GetConfig() (*Config, error) {
	f, err := os.Open("./ass_conf.yaml")
	var conf *Config
	if err != nil || f == nil {
		fmt.Println("no ass_conf.yaml file, use default config")
		conf = &Config{
			CompressScale: 1,
			ResizeScale:   1,
		}
		return conf, err
	}

	confStr, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ass_conf.yaml file incorrect, use default config")
		conf = &Config{
			CompressScale: 1,
			ResizeScale:   1,
		}
		return conf, err
	}

	config, err := parseConfig(string(confStr))
	if err != nil {
		fmt.Println("ass_conf.yaml file incorrect, use default config")
		conf = &Config{
			CompressScale: 1,
			ResizeScale:   1,
		}
		return conf, err
	}

	err = checkConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func parseConfig(configStr string) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal([]byte(configStr), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func checkConfig(config *Config) error {
	if config.CompressScale > 1 || config.CompressScale < 0 {
		return errors.New("invalid compress scale")
	}
	if config.ResizeScale > 1 || config.ResizeScale < 0 {
		return errors.New("invalid compress scale")
	}
	return nil
}
