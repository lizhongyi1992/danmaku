package main

type DanmakuRecord struct {
	DanmakuID int    `json:"danmaku_id"`
	VideoID   int    `json:"video_id"`
	Userno    int    `json:"userno"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	Type      int    `json:"type"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	Heat      int    `json:"heat"`
	Offset    int    `json:"offset"`
	Action    int    `json:"action"`
	Timestamp int64  `json:"date"`
	Comment   string `json:"comment"`
}
