package syncer

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Syncer interface {
	Run()
	Stop()
	StopChan() chan struct{}
}

type redis_mysql_syncer struct {
	pool     *redis.Pool
	db       *sql.DB
	f        func(redis.Conn, *sql.DB)
	task     *TimerTask
	stopchan chan struct{}
}

type RedisMysqlSyncer interface {
	Syncer
	Init(func() (*redis.Pool, error), func() (*sql.DB, error)) error
	SyncRedis(func(redis.Conn))
	SetSyncMysqlHandle(int, func(redis.Conn, *sql.DB))
}

type RedisMysqlSyncerOption struct {
}

func NewRedisMysqlSyncer(opt RedisMysqlSyncerOption) RedisMysqlSyncer {
	return &redis_mysql_syncer{stopchan: make(chan struct{})}
}

func (p *redis_mysql_syncer) Run() {
	go p.task.Start()
	<-p.task.StopChan()
	close(p.stopchan)
}

func (p *redis_mysql_syncer) Stop() {
	p.task.Stop()
}

func (p *redis_mysql_syncer) StopChan() chan struct{} {
	return p.stopchan
}

func (p *redis_mysql_syncer) Init(redis_connect func() (*redis.Pool, error), mysql_connect func() (*sql.DB, error)) error {
	var e error
	p.pool, e = redis_connect()
	if e != nil {
		return e
	}

	p.db, e = mysql_connect()
	if e != nil {
		return e
	}
	return nil
}

func (p *redis_mysql_syncer) SyncRedis(f func(redis.Conn)) {
	f(p.pool.Get())
}

func (p *redis_mysql_syncer) SetSyncMysqlHandle(second int, f func(redis.Conn, *sql.DB)) {
	p.f = f
	p.task = NewTimerTask(second, func() {
		p.f(p.pool.Get(), p.db)
	})
}
