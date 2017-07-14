package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ListenAt string
	Syncer   syncer_config
}

type syncer_config struct {
	MysqlAddress  string
	MysqlUser     string
	MysqlPassword string
	MysqlTable    string

	RedisAddress  string
	RedisPassword string

	FlushIntervalSecond int
	RedisShuffleSuffix  string

	PubDanmakuListName     string
	LikeDanmakuHashName    string
	DislikeDanmakuHashName string
	UserLikesDanmakuHash   string
}

func default_config() Config {
	c := Config{
		ListenAt: ":8888",
		Syncer: syncer_config{
			RedisAddress:        "localhost:6379",
			RedisPassword:       "",
			MysqlUser:           "root",
			MysqlPassword:       "root",
			MysqlAddress:        "localhost:3306",
			MysqlTable:          "test.tdanmaku",
			FlushIntervalSecond: 2,

			RedisShuffleSuffix:     "inprogress",
			PubDanmakuListName:     "pub_danmaku_delta",
			LikeDanmakuHashName:    "like_danmaku_delta",
			DislikeDanmakuHashName: "dislike_danmaku_delta",
			UserLikesDanmakuHash:   "user_likes_danmaku",
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
