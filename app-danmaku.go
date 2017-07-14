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

	syncer_pub          syncer.RedisMysqlSyncer
	syncer_like_dislike syncer.RedisMysqlSyncer

	datasource struct {
		pool *redis.Pool
		db   *sql.DB
	}
}

func NewApp(config Config) (*App, error) {
	_log(config)
	var e error
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
	app.datasource.pool, e = redis_connect()
	if e != nil {
		return nil, e
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
	app.datasource.db, e = mysql_connect()
	if e != nil {
		return nil, e
	}

	app.syncer_pub = syncer.NewRedisMysqlSyncer(syncer.RedisMysqlSyncerOption{})
	e = app.syncer_pub.Init(redis_connect, mysql_connect)
	if e != nil {
		return nil, e
	}

	app.syncer_pub.
		SetSyncMysqlHandle(
			config.Syncer.FlushIntervalSecond,
			app.insert_danmaku_to_mysql)

	// like
	app.syncer_like_dislike = syncer.NewRedisMysqlSyncer(syncer.RedisMysqlSyncerOption{})
	e = app.syncer_like_dislike.Init(redis_connect, mysql_connect)
	if e != nil {
		return nil, e
	}
	app.syncer_like_dislike.
		SetSyncMysqlHandle(
			config.Syncer.FlushIntervalSecond,
			app.update_danmaku_like_dislike_to_mysql)

	return app, nil
}

func (p *App) Run() {
	go p.syncer_pub.Run()
}

func (p *App) Stop() {
	p.syncer_pub.Stop()
	p.syncer_like_dislike.Stop()
}

func (p *App) WaitForExit() {
	<-p.syncer_pub.StopChan()
	<-p.syncer_like_dislike.StopChan()
}

func (p *App) danmaku_all(c *gin.Context) {
	var video_id, uid int
	video_id, _ = strconv.Atoi(c.Query("video_id"))
	uid, _ = strconv.Atoi(c.Query("curr_uid"))
	_dbg(video_id, uid, c.Request.URL.RawQuery)

	table := p.config.Syncer.MysqlTable
	sqlstr := "select * from " + table + " where video_id=" + fmt.Sprint(video_id)
	_dbg(sqlstr)
	rows, e := p.datasource.db.Query(sqlstr)
	if e != nil {
		_err(e)
		c.Status(400)
		return
	}
	defer rows.Close()

	l := []DanmakuRecord{}
	for rows.Next() {
		var DanmakuID int
		var VideoID int
		var Userno int
		var Avatar string
		var Nickname string
		var Type int
		var Likes int
		var Dislikes int
		var Heat int
		var Offset int
		var Action int
		var Date []byte
		var Comment string

		e := rows.Scan(&DanmakuID, &Userno, &VideoID, &Type, &Likes, &Dislikes, &Heat, &Action, &Offset, &Date, &Nickname, &Avatar, &Comment)
		t, e := time.Parse("2006-01-02 15:04:05", string(Date))
		if e != nil {
			_err(e)
			c.Status(400)
			return
		}

		conn := p.datasource.pool.Get()
		user_hash := p.config.Syncer.UserLikesDanmakuHash
		key := join_string_by("_", fmt.Sprint(uid), fmt.Sprint(VideoID), fmt.Sprint(DanmakuID))
		reply, e := conn.Do("hget", user_hash, key)
		_err(key, reply, e)
		if e == nil {
			if action, ok := reply.([]byte); ok {
				switch string(action) {
				case "1":
					Action = 1
				case "2":
					Action = 2
				}
			}
		}

		record := DanmakuRecord{
			DanmakuID: DanmakuID,
			VideoID:   VideoID,
			Userno:    Userno,
			Avatar:    Avatar,
			Nickname:  Nickname,
			Type:      Type,
			Likes:     Likes,
			Dislikes:  Dislikes,
			Heat:      Heat,
			Offset:    Offset,
			Action:    Action,
			Timestamp: t.Unix(),
			Comment:   Comment,
		}
		l = append(l, record)
	}
	_dbg(e, l)

	// TODO:sort by heat

	b, e := json.Marshal(l)
	if e != nil {
		c.Status(400)
		return
	}

	c.String(200, string(b))
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
	type_, _ = strconv.Atoi(c.Query("type"))
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
	_log(c.Request.URL.RawQuery)

	video_id, uid, danmaku_id := c.Query("video_id"), c.Query("uid"), c.Query("danmaku_id")
	if video_id == "" || uid == "" || danmaku_id == "" {
		c.Status(400)
		return
	}

	p.syncer_like_dislike.SyncRedis(func(conn redis.Conn) {
		func() {
			like_hash := p.config.Syncer.LikeDanmakuHashName
			key := danmaku_id
			_, e := conn.Do("hincrby", like_hash, key, 1)
			if e != nil {
				_err(e)
			}
		}()

		func() {
			user_hash := p.config.Syncer.UserLikesDanmakuHash
			key := join_string_by("_", uid, video_id, danmaku_id)
			_, e := conn.Do("hset", user_hash, key, 1)
			if e != nil {
				_err(e)
			}
		}()
	})
}

func (p *App) danmaku_dislike(c *gin.Context) {
	_log(c.Request.URL.RawQuery)

	video_id, uid, danmaku_id := c.Query("video_id"), c.Query("uid"), c.Query("danmaku_id")
	if video_id == "" || uid == "" || danmaku_id == "" {
		c.Status(400)
		return
	}

	p.syncer_like_dislike.SyncRedis(func(conn redis.Conn) {
		func() {
			like_hash := p.config.Syncer.LikeDanmakuHashName
			key := danmaku_id
			_, e := conn.Do("hincrby", like_hash, key, -1) // *decr*
			if e != nil {
				_err(e)
			}
		}()

		func() {
			user_hash := p.config.Syncer.UserLikesDanmakuHash
			key := join_string_by("_", uid, video_id, danmaku_id)
			_, e := conn.Do("hset", user_hash, key, 2)
			if e != nil {
				_err(e)
			}
		}()
	})
}

// pub
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
	sqlstr := "insert into " + p.config.Syncer.MysqlTable + " (uid,video_id,type,likes,dislikes,heat,action,offset,date,nickname,avatar,comment) values (?,?,?,0,0,?,?,?,from_unixtime(?),?,?,?);"
	_dbg(sqlstr)
	for _, v := range toupdate {
		var r DanmakuRecord
		e := json.Unmarshal(v.([]byte), &r)
		if e != nil {
			_err(e)
			continue
		}
		_dbg(r)
		_, e = tx.Exec(sqlstr, r.Userno, r.VideoID, r.Type, r.Heat, r.Action, r.Offset, r.Timestamp, r.Nickname, r.Avatar, r.Comment)
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

// like / dislike
func (p *App) update_danmaku_like_dislike_to_mysql(conn redis.Conn, db *sql.DB) {
}
