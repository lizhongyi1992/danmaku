# danmaku

## API

```
type define:

VIDEO_ID      int
USER_ID       int
MESSAGE       string
```

- GET /danmaku?video_id=[VIDEO_ID]&uid=[USER_ID]

- POST /pub_danmaku?video_id=[VIDEO_ID]&uid=[USER_ID]&msg=[MESSAGE]

- POST /danmaku/like?video_id=[VIDEO_ID]&uid=[USER_ID]

