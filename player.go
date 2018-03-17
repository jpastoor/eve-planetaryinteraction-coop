package main

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type Player struct {
	Name string `gorm:"type:varchar(100);PRIMARY_KEY"`
	Main string
}

type PlayerFetcher interface {
	getOrCreatePlayerByName(playerName string) (*Player, error)
}

type DbPlayerFetcher struct {
	db          *gorm.DB
	cacheByName map[string]*Player
}

func NewDbPlayerFetcher(db *gorm.DB) *DbPlayerFetcher {
	return &DbPlayerFetcher{
		db:          db,
		cacheByName: make(map[string]*Player),
	}
}

func (tf *DbPlayerFetcher) getOrCreatePlayerByName(playerName string) (*Player, error) {
	if ty, exists := tf.cacheByName[playerName]; exists {
		return ty, nil
	}

	ty := Player{}
	tf.db.Where("name = ?", playerName).First(&ty)

	if &ty == nil {
		ty.Name = playerName
		tf.db.Create(ty)
	}

	if tf.db.Error != nil {
		return nil, tf.db.Error
	}

	tf.cacheByName[playerName] = &ty

	return &ty, nil
}

/**
Helper when writing tests. Uses only in-memory cache
 */
type DbPlayerFetcherMock struct {
	cacheByName map[string]*Player
}

func (tf *DbPlayerFetcherMock) getOrCreatePlayerByName(playerName string) (*Player, error) {
	if ty, exists := tf.cacheByName[playerName]; exists {
		return ty, nil
	}
	return nil, fmt.Errorf("Player %s not found", playerName)
}
