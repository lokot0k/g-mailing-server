package models

import "time"

type Message struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SenderID    uint      `json:"sender_id"`
	RecipientID uint      `json:"recipient_id"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	IsRead      bool      `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
	Sender      User      `gorm:"foreignKey:SenderID"    json:"sender"`
	Recipient   User      `gorm:"foreignKey:RecipientID" json:"recipient"`
}
