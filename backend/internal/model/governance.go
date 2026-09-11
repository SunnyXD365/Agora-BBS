package model

import "time"

type GovernancePolicy struct {
	CoolingSeconds    int `json:"cooling_seconds"`
	ReplyDwellSeconds int `json:"reply_dwell_seconds"`
	HeartbeatSeconds  int `json:"heartbeat_seconds"`
	LongTopicChars    int `json:"long_topic_chars"`
	Level1ReadSeconds int `json:"level_1_read_seconds"`
	Level2ReadSeconds int `json:"level_2_read_seconds"`
	Level3ReadSeconds int `json:"level_3_read_seconds"`
}

type ReadingSession struct {
	PublicID          string    `json:"id"`
	TopicID           int64     `json:"topic_id"`
	Progress          int       `json:"progress"`
	ReadingSeconds    int       `json:"reading_seconds"`
	ReplyDwellSeconds int       `json:"reply_dwell_seconds"`
	BottomReached     bool      `json:"bottom_reached"`
	Eligible          bool      `json:"eligible"`
	Completed         bool      `json:"completed"`
	LastHeartbeatAt   time.Time `json:"last_heartbeat_at"`
}

type StartReadingReq struct {
	TopicID int64 `json:"topic_id" binding:"required"`
}

type ReadingHeartbeatReq struct {
	Progress     int  `json:"progress" binding:"min=0,max=100"`
	ReplyFocused bool `json:"reply_focused"`
}
