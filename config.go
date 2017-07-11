package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ListenAt string
	VideoAcc accumulator_config
}

type accumulator_config struct {
	RedisAddress           string
	RedisPassword          string
	RedisHashSetName       string
	RedisHashShuffleSuffix string

	MysqlAddress  string
	MysqlUser     string
	MysqlPassword string
	MysqlDB       string
	MysqlTable    string
	MysqlField    string
	MysqlKey      string

	FlushIntervalSecond int
	MaxKeyCached        int

	WriteDBTimeout int
	FailRetryTimes int
}

func default_config() Config {
	c := Config{
		ListenAt: ":8888",
		VideoAcc: accumulator_config{
			RedisAddress:           "localhost:6379",
			RedisHashSetName:       "acc_views",
			RedisHashShuffleSuffix: "_inprogress",
			MysqlUser:              "root",
			MysqlPassword:          "root",
			MysqlAddress:           "localhost:3306",
			MysqlTable:             "test.tshare",
			MysqlField:             "playtimes",
			MysqlKey:               "no",
			FlushIntervalSecond:    2,
			MaxKeyCached:           1000,
		},
	}
	return c
}

func gen_sample_config() string {
	out, e := yaml.Marshal(default_config())
	if e != nil {
		_err(e)
	}
	return string(out)
}

func ParseConfig(path string) (Config, error) {
	c := Config{}
	b, e := ioutil.ReadFile(path)
	if e != nil {
		return c, e
	}
	e = yaml.Unmarshal(b, &c)
	return c, e
}
