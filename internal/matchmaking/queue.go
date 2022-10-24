package matchmaking

import (
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

func AddPlayerToQueue(playerID string, role models.Role) error {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	err = db.AddPlayerToQueue(p, role)
	if err != nil {
		return err
	}
	log.Debugf("%s joined the queue as a %s", playerID, role)
	return nil
}

func RemovePlayerFromQueue(playerID string) error {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	err = db.RemovePlayerFromQueue(p)
	if err != nil {
		return err
	}
	log.Debugf("%s left the queue", playerID)
	return nil
}

func IsPlayerInQueue(playerID string) (bool, error) {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return false, err
	}
	return db.IsPlayerInQueue(p)
}