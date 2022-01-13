package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Scale  float32 `yaml:"scale"`  // 比例缩放（使用此值时宽高失效）
	Width  uint    `yaml:"width"`  // 缩放后的宽
	Height uint    `yaml:"height"` // 缩放后的高
}

func CreateConfFile() error {
	configStr := "scale: 1\nwidth: 0\nheight: 0\n"

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
			Scale:  1,
			Width:  0,
			Height: 0,
		}
		return conf, err
	}

	confStr, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ass_conf.yaml file incorrect, use default config")
		conf = &Config{
			Scale:  1,
			Width:  0,
			Height: 0,
		}
		return conf, err
	}

	config, err := parseConfig(string(confStr))
	if err != nil {
		fmt.Println("ass_conf.yaml file incorrect, use default config")
		conf = &Config{
			Scale:  1,
			Width:  0,
			Height: 0,
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
	if config.Scale > 1 || config.Scale < 0 {
		return errors.New("invalid scale")
	}
	return nil
}
