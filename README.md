# danmaku

## API

```
type define:

VIDEO_ID  Int
UID       Int
MSG       String

DANMAKU   Object
{
  userno    Int     // publish user
  avatar    String  // avatar url
  nickname  String  
  type      Int     // this field use by client side define
  heat      Int     // sort by heat, heat=max(like-dislike,0)
  pub_at    Int     // offset from video begin(0s)
  action    Int     // 0=None 1=like 2=dislike
  date      Int     // unix timestamp
  comment   String  // danmaku info
}
```

- GET /danmaku/all?video_id=[VIDEO_ID]&uid=[UID]

success: 200 List{DANMAKU[]}

failure: 400 error_string

- POST /danmaku/pub?video_id=[VIDEO_ID]&uid=[UID]&msg=[MSG]

success: 200

failure: 400

- POST /danmaku/like?video_id=[VIDEO_ID]&uid=[UID]

success: 200

failure: 400

- POST /danmaku/dislike?video_id=[VIDEO_ID]&uid=[UID]

success: 200

failure: 400

