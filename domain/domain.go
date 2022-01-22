package domain

import "time"

type ChatInfo struct {
	HeadPerson *User
	PollInfo   *PollInfo
	IsWorking  bool
	Done       chan struct{}
	Users      map[int64]User
}

type PollInfo struct {
	ID           string
	Times        []MentionTime
	Results      *PollResults
	CreationDate time.Time
}

type MentionTime struct {
	MenTime time.Time
	Done    bool
}

type PollResults struct {
	Health   map[int64]bool
	Sick     map[int64]bool
	Pass     map[int64]bool
	Negative map[int64]bool
	Positive map[int64]bool
	All      map[int64]bool
}

type User struct {
	ID        int64
	Username  string
	Firstname string
	Lastname  string
	ChatID    int64
}
