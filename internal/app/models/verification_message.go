package models

import (
	"log"
	"strconv"
)

type VerificationMessageStatus string

const (
	VerificationMessageStatusPending  VerificationMessageStatus = "pending"
	VerificationMessageStatusAccepted VerificationMessageStatus = "accepted"
	VerificationMessageStatusDeclined VerificationMessageStatus = "declined"
)

type VerificationMessage struct {
	ChatId      string
	MessageId   string
	AboutUserId string
	Status      VerificationMessageStatus
}

func (m *VerificationMessage) MessageSig() (string, int64) {
	id, err := strconv.Atoi(m.ChatId)
	if err != nil {
		log.Fatal("can't convert chat id from str to int: ", err)
	}

	return m.MessageId, int64(id)
}
