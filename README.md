# danmaku

## API

```
type define:

video_id  Int
uid       Int
curr_uid  Int
comment   String

Danmaku   Object
{
  "danmaku_id": Int,                    // global unique danmaku id
  "video_id": Int,                      // video
  "userno": Int,                        // publisher
  "avatar": "String",                   // avatar url 
  "nickname": "String",                 // nickname
  "type": Int,                          // this field use by client side define
  "likes": Int,                         // total likes count
  "dislikes": Int,                      // total dislies count
  "heat": Int,                          // sort by heat, heat=max(like-dislike,0)
  "offset": Int,                        // offset from video begin(0s)
  "action": Int,                        // 0=None 1=like 2=dislike by current user
  "date": Int64,                        // unix timestamp
  "comment": "String",                  // danmaku info
}


```

- GET /danmaku/all?video_id=&curr_uid=

```
success: 200 List{Danmaku[]}
failure: 400 error_string
```

- POST /danmaku/pub?video_id=&curr_uid=&comment=&avatar=&nickname=&type=&offset=&date=

```
success: 200
failure: 400
```

- POST /danmaku/like?video_id=&uid=&danmaku_id=

```
success: 200
failure: 400
```

- POST /danmaku/dislike?video_id=&uid=&danmaku_id=

```
success: 200
failure: 400
```

## data structure explanation 

```
redis:

(hash) like_danmaku_delta
         key:danmaku_id value:count (int)

(hash) dislike_danmaku_delta (hash)
         key:danmaku_id value:count (int)

(list) pub_danmaku_delta     
         []json(danmaku)

***_delta_inprogressing:  renamed by program for update mysql

(hash) user_likes_danmaku
         key:uid_videoid_danmakuid value: 0=None,1=like,2=dislike


mysql:

create table if not exists test.tdanmaku(id int not null auto_increment primary key, uid int,video_id int,type int, likes int, dislikes int, heat int,action int, offset int,date datetime ,nickname varchar(128),avatar text,comment text) default charset=utf8;

```

## build & deploy

```
# build:
liym@desktop:~/gopath/src/danmaku$ make
go build -ldflags "-X main.VERSION=1.0.0 -X 'main.GITHASH=`git rev-parse --short HEAD`' -X 'main.BUILT=`date +"%Y/%m/%d %H:%M:%S"`' -X 'main.GOVERSION=`go version | cut -d" " -f 3`'"
./danmaku -genconf > config.yaml

liym@desktop:~/gopath/src/danmaku$ ./danmaku -help
Usage of ./danmaku:
  -c string
    	config file path (default "config.yaml")
  -genconf
    	generate a sample config
  -version
    	show version string and exit

liym@desktop:~/gopath/src/danmaku$ ./danmaku -version
1.0.0 679a444 2017/07/14 20:06:30 go1.8.3

# configure
liym@desktop:~/gopath/src/danmaku$ ./danmaku -genconf > config.yaml

# startup:
liym@desktop:~/gopath/src/danmaku$ ./danmaku 
  # ...
2017/07/14 20:07:16 [GIN-debug] GET    /ping                     --> main.ping_test (3 handlers)
2017/07/14 20:07:16 [GIN-debug] GET    /danmaku/all              --> main.(*App).(main.danmaku_all)-fm (3 handlers)
2017/07/14 20:07:16 [GIN-debug] POST   /danmaku/pub              --> main.(*App).(main.danmaku_pub)-fm (3 handlers)
2017/07/14 20:07:16 [GIN-debug] POST   /danmaku/like             --> main.(*App).(main.danmaku_like)-fm (3 handlers)
2017/07/14 20:07:16 [GIN-debug] POST   /danmaku/dislike          --> main.(*App).(main.danmaku_dislike)-fm (3 handlers)
2017/07/14 20:07:16 [GIN-debug] Listening and serving HTTP on :8888

# exit:
liym@desktop:~/gopath/src/danmaku$ killall danmaku
```
