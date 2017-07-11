# danmaku

## API

```
type define:

VIDEO_ID  Int
UID       Int
MSG       String

DANMAKU   Object
{
  userno    Int
  avatar    String
  nickname  String
  type      Int
  heat      Int     // sort by heat, heat=like-dislike
  pub_at    Int
  date      Int     // unix timestamp
  comment   String
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

