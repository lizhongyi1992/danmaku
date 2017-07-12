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
}

type RedisMysqlSyncer interface {
	Syncer
	SyncRedis(func(redis.Conn))
	SetSyncMysqlHandle(func(redis.Conn, *sql.DB))
}

type RedisMysqlSyncerOption struct {
}

func NewRedisMysqlSyncer(opt RedisMysqlSyncerOption) RedisMysqlSyncer {
	return nil
}

func (p *redis_mysql_syncer) Run() {

}

func (p *redis_mysql_syncer) Stop() {

}

func (p *redis_mysql_syncer) StopChan() {

}

func (p *redis_mysql_syncer) SyncRedis(f func(redis.Conn)) {

}

func (p *redis_mysql_syncer) SetSyncMysqlHandle(f func(redis.Conn, *sql.DB)) {

}
