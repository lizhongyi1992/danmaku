package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	var genconf, show_version bool
	var config_file string

	flag.BoolVar(&genconf, "genconf", false, "generate a sample config")
	flag.BoolVar(&show_version, "version", false, "show version string and exit")
	flag.StringVar(&config_file, "c", "config.yaml", "config file path")
	flag.Parse()

	if show_version {
		fmt.Println(version_string())
		_exit()
	}

	if genconf {
		fmt.Println(gen_sample_config())
		_exit()
	}

	config, e := ParseConfig(config_file)
	_exit_if(e)
	_log(config)
	_dbg(config)

	r := gin.Default()

	app := NewApp(config)

	r.GET("/ping", ping_test)
	r.POST("/danmaku/all", app.danmaku_all)
	go r.Run(config.ListenAt)

	s := <-waiting_for_interrupt_chan()
	_log("quit when catch signal:", s)

	app.Stop()
	syncer.Stop()
	_log("waiting for syncer exit")
	<-syncer.StopChan()
}
