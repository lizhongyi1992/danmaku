package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	VERSION   string
	BUILT     string
	GITHASH   string
	GOVERSION string
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
	log.SetFlags(log.LstdFlags)
}

func callerinfo() string {
	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	return path.Base(file) + "(" + fmt.Sprint(line) + ") " + path.Base(f.Name())
}

func _dbg(v ...interface{}) {
	logger.Println("DBG", v, callerinfo())
}

func _err(v ...interface{}) {
	log.Println("ERR", v, callerinfo())
}

func _exit_if(err error, v ...interface{}) {
	if err != nil {
		log.Println("ERR Exit", err, v, callerinfo())
		os.Exit(-1)
	}
}

func _log(v ...interface{}) {
	log.Println("INF", v, callerinfo())
}

func _exit() {
	os.Exit(0)
}

func waiting_for_interrupt_chan() chan os.Signal {

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	return c
}

func version_string() string {
	return VERSION + " " + GITHASH + " " + BUILT + " " + GOVERSION
}

func ping_test(c *gin.Context) {
	c.String(200, "ok")
}

func join_string_by(sep string, v ...string) string {
	return strings.Join(v, sep)
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
