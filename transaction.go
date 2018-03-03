package main

import (
	"time"
	"crypto/md5"
	"io"
	"fmt"
	"encoding/hex"
)

type Transaction struct {
	Id            string  `gorm:"PRIMARY_KEY"`
	CreationDate  time.Time
	Location      string
	SubLocation   string
	Who           *Player `gorm:"foreignkey:PlayerName"`
	PlayerName    string
	Action        string
	Status        string
	Type          *Type   `gorm:"foreignkey:TypeName"`
	TypeName      string
	Quantity      int
	MarkedForCorp bool
}

func (t Transaction) hash() string {
	input := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%d", t.CreationDate, t.SubLocation, t.PlayerName, t.Action, t.Status, t.TypeName, t.Quantity)
	hasher := md5.New()
	io.WriteString(hasher, input)
	return hex.EncodeToString(hasher.Sum(nil))
}

type Player struct {
	Name string `gorm:"type:varchar(100);PRIMARY_KEY"`
	Main string
}

type Type struct {
	TypeId int    `gorm:"PRIMARY_KEY"`
	Name   string `gorm:"type:varchar(255);UNIQUE_KEY"`
}
