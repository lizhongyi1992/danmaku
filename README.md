# danmaku

## API

```
type define:

video_id  Int
uid       Int
curr_uid  Int
comment   String

Danmaku   Object
{
  userno    Int     // publish user
  avatar    String  // avatar url
  nickname  String  
  type      Int     // this field use by client side define
  heat      Int     // sort by heat, heat=max(like-dislike,0)
  offset    Int     // offset from video begin(0s)
  action    Int     // 0=None 1=like 2=dislike
  date      Int     // unix timestamp
  comment   String  // danmaku info
}
```

- GET /danmaku/all?video_id=&curr_uid=

success: 200 List{Danmaku[]}

failure: 400 error_string

- POST /danmaku/pub?video_id=&uid=&comment=&avatar=&nickname=&type=&offset=&date=

success: 200

failure: 400

- POST /danmaku/like?video_id=&uid=&danmakuid=

success: 200

failure: 400

- POST /danmaku/dislike?video_id=&uid=&danmakuid=

success: 200

failure: 400

