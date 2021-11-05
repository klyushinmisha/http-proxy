package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	Host         *string  `mapstructure:"host"`
	Port         *int     `mapstructure:"port"`
	Hosts        []string `mapstructure:"hosts"`
	BalancerType *string  `mapstructure:"balancer-type"`
	MaxHeaderKb  *int     `mapstructure:"max-header-kb"`
	BufferSizeKb *int     `mapstructure:"buffer-size-kb"`
}

type Config struct {
	Host         string
	Port         int
	Hosts        []string
	BalancerType string
	MaxHeaderKb  int
	BufferSizeKb int
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"host: %s\n"+
			"port: %d\n"+
			"hosts: %v\n"+
			"balancer-type: %s\n"+
			"max-header: %d bytes\n"+
			"buffer-size: %d bytes\n",
		c.Host,
		c.Port,
		c.Hosts,
		c.BalancerType,
		c.MaxHeaderKb,
		c.BufferSizeKb,
	)
}

func LoadConfig(path string) (*Config, error) {
	var (
		Kilobyte            = 1024
		DefaultHost         = "0.0.0.0"
		DefaultPort         = 8080
		DefaultBalancerType = RoundRobinBalancer
		DefaultMaxHeaderKb  = Kilobyte
		DefaultBufferSizeKb = os.Getpagesize()
	)

	conf := new(config)

	loader := viper.New()
	loader.SetConfigType("yaml")
	loader.SetConfigFile(path)

	if err := loader.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := loader.Unmarshal(&conf); err != nil {
		return nil, err
	}

	if conf.Host == nil {
		conf.Host = &DefaultHost
	}

	if conf.Port == nil {
		conf.Port = &DefaultPort
	}

	if conf.BalancerType == nil {
		conf.BalancerType = &DefaultBalancerType
	}

	if conf.MaxHeaderKb == nil {
		conf.MaxHeaderKb = &DefaultMaxHeaderKb
	}

	if conf.BufferSizeKb == nil {
		conf.BufferSizeKb = &DefaultBufferSizeKb
	}

	if *conf.Port < 0 || *conf.Port > 1<<16 {
		return nil, errors.New("port value must be in ranges [0; 65535]")
	}

	if conf.Hosts == nil {
		return nil, errors.New("hosts cannot be empty")
	}

	if *conf.BufferSizeKb < 1 {
		return nil, errors.New("buffer-size-kb must be positive")
	}

	if *conf.MaxHeaderKb < 1 {
		return nil, errors.New("max-header-kb must be positive")
	}

	*conf.BufferSizeKb *= Kilobyte
	*conf.MaxHeaderKb *= Kilobyte

	return &Config{
		Host:         *conf.Host,
		Port:         int(*conf.Port),
		Hosts:        conf.Hosts,
		BalancerType: *conf.BalancerType,
		MaxHeaderKb:  int(*conf.MaxHeaderKb),
		BufferSizeKb: int(*conf.BufferSizeKb),
	}, nil
}
