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

success: 200 List{Danmaku[]}

failure: 400 error_string

- POST /danmaku/pub?video_id=&curr_uid=&comment=&avatar=&nickname=&type=&offset=&date=

success: 200

failure: 400

- POST /danmaku/like?video_id=&uid=&danmaku_id=

success: 200

failure: 400

- POST /danmaku/dislike?video_id=&uid=&danmaku_id=

success: 200

failure: 400


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
