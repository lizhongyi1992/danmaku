# danmaku

## API

```
type define:

VIDEO_ID  int
UID       int
MSG       string
```

- GET /danmaku?video_id=[VIDEO_ID]&uid=[UID]



- POST /pub_danmaku?video_id=[VIDEO_ID]&uid=[UID]&msg=[MSG]

- POST /danmaku/like?video_id=[VIDEO_ID]&uid=[UID]

- POST /danmaku/dislike?video_id=[VIDEO_ID]&uid=[UID]

