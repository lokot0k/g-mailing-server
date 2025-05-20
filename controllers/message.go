package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lokot0k/mservice/models"
	"github.com/lokot0k/mservice/queue"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type MessageController struct {
	DB  *gorm.DB
	RMQ *amqp.Connection
}

func (m *MessageController) Send(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req struct {
		RecipientEmail string `json:"recipient_email" binding:"required,email"`
		Subject        string `json:"subject"         binding:"required"`
		Body           string `json:"body"            binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var recipient models.User
	if err := m.DB.Where("email = ?", req.RecipientEmail).First(&recipient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipient not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}

	msg := models.Message{
		SenderID:    userID,
		RecipientID: recipient.ID,
		Subject:     req.Subject,
		Body:        req.Body,
	}
	if err := m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}
		payload, _ := json.Marshal(msg)
		if err := queue.Publish(m.RMQ, "mail_notifications", payload); err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func (m *MessageController) Inbox(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var msgs []models.Message
	if err := m.DB.
		Preload("Sender").
		Preload("Recipient").
		Where("recipient_id = ?", userID).
		Order("created_at DESC").
		Find(&msgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

func (m *MessageController) Sent(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var msgs []models.Message
	if err := m.DB.
		Preload("Sender").
		Preload("Recipient").
		Where("sender_id = ?", userID).
		Order("created_at DESC").
		Find(&msgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, msgs)
}
