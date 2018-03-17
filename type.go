package main

import (
	"github.com/jinzhu/gorm"
	"fmt"
)


type Type struct {
	TypeID   int    `gorm:"column:typeID;PRIMARY_KEY"`
	TypeName string `gorm:"column:typeName;type:varchar(100)"`
	Volume   float32
}

func (Type) TableName() string {
	return "invtypes"
}

type TypeFetcher interface {
	getTypeById(typeId int) (*Type, error)
	getTypeByName(typeName string) (*Type, error)
}

type DbTypeFetcher struct {
	db          *gorm.DB
	cacheById   map[int]*Type
	cacheByName map[string]*Type
}

func NewDbTypeFetcher(db *gorm.DB) *DbTypeFetcher {
	return &DbTypeFetcher{
		db:          db,
		cacheById:   make(map[int]*Type),
		cacheByName: make(map[string]*Type),
	}
}

func (tf *DbTypeFetcher) getTypeById(typeId int) (*Type, error) {
	ty, exists := tf.cacheById[typeId]

	if !exists {
		ty = &Type{}
		tf.db.First(ty, typeId)

		if tf.db.Error != nil {
			return nil, tf.db.Error
		}
	}

	tf.cacheById[typeId] = ty

	return ty, nil
}

func (tf *DbTypeFetcher) getTypeByName(typeName string) (*Type, error) {
	ty, exists := tf.cacheByName[typeName]

	if !exists {
		ty = &Type{}
		tf.db.Where("typeName = ?", typeName).First(ty)

		if tf.db.Error != nil {
			return nil, tf.db.Error
		}

		tf.cacheByName[typeName] = ty
	}

	return ty, nil
}

type TypeFetcherMock struct {
	cacheById   map[int]*Type
	cacheByName map[string]*Type
}

func (tf *TypeFetcherMock) getTypeByName(typeName string) (*Type, error) {
	ty, exists := tf.cacheByName[typeName]

	if !exists {
		return nil, fmt.Errorf("Type %s not found", typeName)
	}

	return ty, nil
}

func (tf *TypeFetcherMock) getTypeById(typeId int) (*Type, error) {
	ty, exists := tf.cacheById[typeId]

	if !exists {
		return nil, fmt.Errorf("Type %d not found", typeId)
	}

	tf.cacheById[typeId] = ty

	return ty, nil
}
