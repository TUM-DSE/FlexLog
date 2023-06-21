package benchmark

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type BenchConfig struct {
	Appends   int           `yaml:"appends"`
	Reads     int           `yaml:"reads"`
	Threads   int           `yaml:"threads"`
	Runtime   time.Duration `yaml:"runtime"`
	Endpoints []string      `yaml:"endpoints"`
	Clients   int           `yaml:"clients"`
	Times     int           `yaml:"times"`
}

func GetBenchConfig(configFile string) (*BenchConfig, error) {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	b := &BenchConfig{}
	err = yaml.Unmarshal(buf, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
