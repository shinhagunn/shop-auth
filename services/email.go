package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"github.com/shinhagunn/Shop-Watches/backend/utils"
)

type SendEmail struct{}

func NewSendEmail() *SendEmail {
	return &SendEmail{}
}

func (SendEmail) Process() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "new-user",
		GroupID: "email-new-users",
	})

	log.Println("Service email running ...")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			panic("could not read message " + err.Error())
		}
		userData := msg.Value

		var user models.User

		err = json.Unmarshal(userData, &user)
		if err != nil {
			panic("could not parse userData " + err.Error())
		}

		randomCode := utils.RandomCode()

		UserID := user.ID
		code := &models.Code{
			UserID:         UserID,
			Code:           randomCode,
			CodeExpiration: time.Now().Add(5 * time.Minute),
			State:          "Active",
		}

		if err := code.Collection().Create(code); err != nil {
			log.Println(err)
		}

		subject := "Subject: Verification Account from ShinWatch\n\n"
		body := fmt.Sprintf("Your code: %v", randomCode)
		message := strings.Join([]string{subject, body}, " ")

		SendEmailService(message, user.Email)
	}
}

func SendEmailService(message string, toAddress string) (response bool, err error) {
	fromAddress := "shinhagunzz5@gmail.com"
	fromEmailPassword := "handpum123"
	smtpServer := "smtp.gmail.com"
	smptPort := "587"

	var auth = smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)
	if err := smtp.SendMail(smtpServer+":"+smptPort, auth, fromAddress, []string{toAddress}, []byte(message)); err == nil {
		return true, nil
	} else {
		return false, err
	}
}
