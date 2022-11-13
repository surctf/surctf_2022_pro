package models

import (
	"crypto/sha1"
	"encoding/base64"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type User struct {
	ID           uint
	PasswordHash []byte
	Pastes       []Paste
}

func (u *User) AddPaste(paste *Paste, db *gorm.DB) error {
	paste.UserID = u.ID

	// Generating paste hash
	hasher := sha1.New()
	salt := strconv.FormatInt(time.Now().UnixNano(), 10)
	hasher.Write([]byte(paste.Title + paste.Content + salt))
	paste.ID = hasher.Sum(nil)

	if tx := db.Create(paste); tx.Error != nil {
		return tx.Error
	}

	if tx := db.Model(u).Preload("Pastes").Find(u, "id = ?", u.ID); tx.Error != nil {
		return tx.Error
	}

	return nil
}

type Paste struct {
	ID      []byte `gorm:"primaryKey"`
	Title   string `gorm:"size:32;"`
	Content string `gorm:"size:140;"`
	UserID  uint
}

func (p *Paste) GetB64URL() string {
	return base64.URLEncoding.EncodeToString(p.ID)
}
