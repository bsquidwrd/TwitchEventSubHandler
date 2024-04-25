package database

type Service struct {
	Cache    *cacheService
	Database *databaseService
	Twitch   *twitchService
}

func New() *Service {
	return &Service{
		Cache:    newCacheService(),
		Database: newDatabaseService(),
		Twitch:   newTwitchService(),
	}
}
