package main

import "github.com/jinzhu/gorm"

type Player struct {
	Name string `gorm:"type:varchar(100);PRIMARY_KEY"`
	Main string
}

type PlayerFetcher interface {
	getOrCreatePlayerByName(PlayerName string) (*Player, error)
}

type DbPlayerFetcher struct {
	db          *gorm.DB
	cacheById   map[int]*Player
	cacheByName map[string]*Player
}

func NewDbPlayerFetcher(db *gorm.DB) *DbPlayerFetcher {
	return &DbPlayerFetcher{
		db:          db,
		cacheByName: make(map[string]*Player),
	}
}

func (tf *DbPlayerFetcher) getOrCreatePlayerByName(PlayerName string) (*Player, error) {
	if ty, exists := tf.cacheByName[PlayerName]; exists {
		return ty, nil
	}

	ty := Player{}
	tf.db.Where("name = ?", PlayerName).First(&ty)

	if &ty == nil {
		ty.Name = PlayerName
		tf.db.Create(ty)
	}

	if tf.db.Error != nil {
		return nil, tf.db.Error
	}

	tf.cacheByName[PlayerName] = &ty

	return &ty, nil
}

