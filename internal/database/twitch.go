package database

import "sync"

type twitchService struct {
	AuthLock      *sync.Mutex
	RatelimitLock *sync.Mutex
	AccessToken   string
}

func newTwitchService() *twitchService {
	return &twitchService{
		AuthLock:      &sync.Mutex{},
		RatelimitLock: &sync.Mutex{},
		AccessToken:   "",
	}
}

func (t *twitchService) Ping() error {
	return nil
}
