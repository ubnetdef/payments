package main

import (
	"github.com/jinzhu/configor"
	"github.com/ubnetdef/ScoreTrak/pkg/storage"
)

type Config struct {
	BankUsername       string `default:""`
	BankPassword       string `default:""`
	ScoreTrakUsername  string `default:""`
	ScoreTrakPassword  string `default:""`
	BankBaseUrl        string `default:"http://localhost"`
	ScoreTrakAddress   string `default:"localhost"`
	BankAmountPerPoint uint64 `default:"10"`
	ClientTimeout      uint64 `default:"5"`
	DB                 storage.Config
}

func NewConfig(f string) (*Config, error) {
	conf := &Config{}
	err := configor.Load(conf, f)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
