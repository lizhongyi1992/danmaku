package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type App struct {
	video_acc *accumulator
}

func NewApp(config Config) *App {
	app := &App{}
	app.video_acc = NewAccumulator(config.VideoAcc)
	return app
}

func (p *App) Stop() {
	p.video_acc.Stop()
}

func (p *App) incr_video_views(c *gin.Context) {
	video_id, e := strconv.Atoi(c.Query("video_id"))

	if e != nil {
		_log(c.Request.URL.Query())
		c.Status(400)
		return
	}
	p.video_acc.Incr(fmt.Sprint(video_id))
	c.Status(200)
}

func (p *App) danmaku_all(c *gin.Context) {
}
