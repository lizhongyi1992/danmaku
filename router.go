package main

import (
	"danmaku/syncer"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type App struct {
	syncer_pub     syncer.RedisMysqlSyncer
	syncer_like    syncer.RedisMysqlSyncer
	syncer_dislike syncer.RedisMysqlSyncer
}

func NewApp(config Config) *App {
	app := &App{}
	return app
}

func (p *App) Stop() {
	p.syncer_pub.Stop()
	p.syncer_like.Stop()
	p.syncer_dislike.Stop()
}

func (p *App) WaitForExit() {
	<-p.syncer_pub.StopChan()
	<-p.syncer_like.StopChan()
	<-p.syncer_dislike.StopChan()
}

func (p *App) danmaku_all(c *gin.Context) {
}

func (p *App) danmaku_pub(c *gin.Context) {
	_log(c.Request.URL.RawQuery)
	//	var comment, avatar, nickname string
	//	var video_id, uid, type_, offset int
	//	var date time.Duration
}

func (p *App) danmaku_like(c *gin.Context) {
	p.syncer_like.SyncRedis(func(conn redis.Conn) {
	})
}

func (p *App) danmaku_dislike(c *gin.Context) {
	p.syncer_dislike.SyncRedis(func(conn redis.Conn) {
	})
}
