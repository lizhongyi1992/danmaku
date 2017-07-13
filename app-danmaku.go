package main

import (
	"danmaku/syncer"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	_log(config)
	app := &App{config: config}

	redis_connect := func() (*redis.Pool, error) {
		_log(config.Syncer.RedisAddress)
		network, address := "tcp", config.Syncer.RedisAddress
		test, e := redis.Dial(network, address)
		if e != nil {
			_err(e)
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
		db, e := sql.Open("mysql", dsn)
		_dbg(dsn, db, e)
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

	app.syncer_pub = syncer.NewRedisMysqlSyncer(syncer.RedisMysqlSyncerOption{})
	e := app.syncer_pub.Init(redis_connect, mysql_connect)
	if e != nil {
		return nil, e
	}

	app.syncer_pub.
		SetSyncMysqlHandle(
			config.Syncer.FlushIntervalSecond,
			app.insert_danmaku_to_mysql)

	return app, nil
}

func (p *App) Run() {
	go p.syncer_pub.Run()
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
	var e error
	var comment, avatar, nickname string
	var video_id, uid, type_, offset int
	var date int64

	comment, avatar, nickname = c.Query("comment"), c.Query("avatar"), c.Query("nickname")
	video_id, _ = strconv.Atoi(c.Query("video_id"))
	uid, _ = strconv.Atoi(c.Query("uid"))
	type_, _ = strconv.Atoi(c.Query("type_"))
	offset, _ = strconv.Atoi(c.Query("offset"))
	date, e = strconv.ParseInt(c.Query("date"), 10, 64)
	if e != nil {
		c.Status(400)
		return
	}

	record := DanmakuRecord{
		VideoID:   video_id,
		Userno:    uid,
		Avatar:    avatar,
		Nickname:  nickname,
		Type:      type_,
		Heat:      0,
		Offset:    offset,
		Action:    0,
		Timestamp: date,
		Comment:   comment,
	}
	p.syncer_pub.SyncRedis(func(conn redis.Conn) {
		p.write_danmaku_to_redis(record, conn)
	})
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

//
func (p *App) write_danmaku_to_redis(record DanmakuRecord, conn redis.Conn) {
	s, e := json.Marshal(record)
	_dbg(string(s), e)
	_, e = conn.Do("lpush", p.config.Syncer.PubDanmakuListName, string(s))
	if e != nil {
		_err(e)
	}
}

func redis_key_exsits(key string, conn redis.Conn) bool {
	reply, e := conn.Do("exists", key)
	if e != nil {
		_err(e)
		return false
	}
	b := reply.(int64)
	if b == 0 {
		return false
	} else {
		return true
	}
}

func (p *App) insert_danmaku_to_mysql(conn redis.Conn, db *sql.DB) {
	oldname := p.config.Syncer.PubDanmakuListName
	name := join_string_by("_", p.config.Syncer.PubDanmakuListName, p.config.Syncer.RedisShuffleSuffix, fmt.Sprint(os.Getpid()), fmt.Sprint(time.Now().Unix()))

	if !redis_key_exsits(oldname, conn) {
		return
	}
	_, e := conn.Do("rename", oldname, name)
	if e != nil {
		_err(e)
		return
	}
	reply, e := conn.Do("lrange", name, 0, -1)
	if e != nil {
		_err(e)
		return
	}
	toupdate, ok := reply.([]interface{})
	if !ok {
		_err("conv failure")
		return
	}

	tx, e := db.Begin()
	if e != nil {
		_err(e)
	}
	sqlstr := "insert into " + p.config.Syncer.MysqlTable + " (uid,type,heat,action,offset,date,nickname,avatar,comment) values (?,?,?,?,?,from_unixtime(?),?,?,?);"
	_dbg(sqlstr)
	for _, v := range toupdate {
		var r DanmakuRecord
		e := json.Unmarshal(v.([]byte), &r)
		if e != nil {
			_err(e)
			continue
		}
		_, e = tx.Exec(sqlstr, r.Userno, r.Type, r.Heat, r.Action, r.Offset, r.Timestamp, r.Nickname, r.Avatar, r.Comment)
		if e != nil {
			_err(e)
		}
	}

	e = tx.Commit()
	if e != nil {
		_err(e)
	}

	conn.Do("del", name)
}
