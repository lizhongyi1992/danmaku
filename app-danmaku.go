package main

import (
	"danmaku/syncer"
	"database/sql"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type App struct {
	config Config

	syncer_pub     syncer.RedisMysqlSyncer
	syncer_like    syncer.RedisMysqlSyncer
	syncer_dislike syncer.RedisMysqlSyncer
}

func NewApp(config Config) (*App, error) {
	app := &App{}

	redis_connect := func() (*redis.Pool, error) {
		network, address := "tcp", config.Syncer.RedisAddress
		test, e := redis.Dial(network, address)
		if e != nil {
			return nil, e
		}
		test.Close()
		return &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.Dial(network, address) },
		}, nil
	}

	mysql_connect := func() (*sql.DB, error) {
		dsn := config.Syncer.MysqlUser + ":" + config.Syncer.MysqlPassword + "@tcp(" + config.Syncer.MysqlAddress + ")/"
		db, e := sql.Open("mysql", connstr)
		_dbg(connstr, db, e)
		if e != nil {
			_err(e)
			return nil, e
		}
		e = db.Ping()
		if e != nil {
			_err(e)
			return nil, e
		}
		return db, nil
	}

	app.syncer_pub.Init(redis_connect, mysql_connect)
	if e != nil {
		return nil, e
	}

	app.syncer_pub.
		SetSyncMysqlHandle(
			config.Syncer.FlushIntervalSecond,
			app.insert_danmaku_to_mysql)

	return app
}

func (p *App) Stop() {
	p.syncer_pub.Stop()
	//p.syncer_like.Stop()
	//p.syncer_dislike.Stop()
}

func (p *App) WaitForExit() {
	<-p.syncer_pub.StopChan()
	//<-p.syncer_like.StopChan()
	//<-p.syncer_dislike.StopChan()
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
	var danmakuid, uid, videoid string
	var like_hset_name string
	key := join_string_by("_", videoid, uid, danmakuid)
	p.syncer_like.SyncRedis(func(conn redis.Conn) {
		_, e := conn.Do("hincrby", like_hset_name, key, 1)
		if e != nil {
			_err(e)
			return
		}
	})
}

func (p *App) danmaku_dislike(c *gin.Context) {
	p.syncer_dislike.SyncRedis(func(conn redis.Conn) {
	})
}

//  async function
func (p *App) insert_danmaku_to_mysql(redis.Conn, *sql.DB) {
}
