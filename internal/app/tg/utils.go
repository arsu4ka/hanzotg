package tg

import (
	"encoding/json"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"hanzotg/internal/app/models"
	"log"
)

func parseUserFromContext(c tele.Context) (string, string) {
	userId := fmt.Sprint(c.Sender().ID)
	username := c.Sender().Username
	return userId, username
}

func addPaymentDataToButton(btn tele.InlineButton, data paymentVerificationData) tele.InlineButton {
	dataString, err := json.Marshal(data)
	if err != nil {
		log.Fatal("can't marshal payment data: ", err)
	}

	btn.Data = string(dataString)
	return btn
}

func getMarkupForPaymentVerificationMsg(userId string) *tele.ReplyMarkup {
	paymentData := paymentVerificationData{
		UserId: userId,
	}
	acceptButton := addPaymentDataToButton(btnAcceptPayment, paymentData)
	//declineButton := addPaymentDataToButton(btnDeclinePayment, paymentData)

	return &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{acceptButton},
		},
	}
}

func parsePaymentData(data string) (paymentVerificationData, error) {
	var paymentData paymentVerificationData
	err := json.Unmarshal([]byte(data), &paymentData)
	return paymentData, err
}

func mustGetUser(c tele.Context) *models.User {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		log.Fatal("user is not in the context")
	}

	return user
}
