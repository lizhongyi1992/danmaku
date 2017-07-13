package main

type DanmakuRecord struct {
	DanmakuID int    `json:"danmakuid"`
	VideoID   int    `json:"videoid"`
	Userno    int    `json:"userno"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	Type      int    `json:"type"`
	Heat      int    `json:"heat"`
	Offset    int    `json:"offset"`
	Action    int    `json:"action"`
	Timestamp int64  `json:"date"`
	Comment   string `json:"comment"`
}
