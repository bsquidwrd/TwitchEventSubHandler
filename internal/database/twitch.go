package database

import "sync"

type twitchService struct {
	AuthLock *sync.Mutex
}

func newTwitchService() *twitchService {
	return &twitchService{
		AuthLock: &sync.Mutex{},
	}
}
