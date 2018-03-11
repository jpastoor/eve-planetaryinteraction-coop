package main

import (
	"time"
	"crypto/md5"
	"io"
	"fmt"
	"encoding/hex"
)

type Transaction struct {
	Id                string `gorm:"PRIMARY_KEY"`
	CreationDate      time.Time
	Location          string
	SubLocation       string
	PlayerName        string
	Action            string
	Status            string
	TypeName          string
	Quantity          int
	MarkedForCorp     bool
	ProcessedInCommit int // TODO Get a proper int or NULL in here
}

const ACTION_UNLOCK = "Unlock"
const ACTION_LOCK = "Lock"

func (t Transaction) hash() string {
	input := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%d", t.CreationDate, t.SubLocation, t.PlayerName, t.Action, t.Status, t.TypeName, t.Quantity)
	hasher := md5.New()
	io.WriteString(hasher, input)
	return hex.EncodeToString(hasher.Sum(nil))
}
