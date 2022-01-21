package domain

import "time"

type ChatInfo struct {
	IsWorking  bool
	Done       chan struct{}
	PollInfo   *PollInfo
	HeadPerson *HeadPerson
}

type HeadPerson struct {
	Username string
	ChatID   int64
}

type PollInfo struct {
	ID    int64
	Times []time.Time
}
