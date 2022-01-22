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
	Results []Result // Health, Sick, Pass, Negative, Positive
	All     map[int64]bool
}

type Result struct {
	Title string
	Users map[int64]bool
}

type User struct {
	ID        int64
	Username  string
	Firstname string
	Lastname  string
	ChatID    int64
}
