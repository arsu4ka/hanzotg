package tg

import (
	"database/sql"
	"errors"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"hanzotg/internal/app/models"
	"log"
	"slices"
	"sync"
	"time"
)

var (
	wg = sync.WaitGroup{}
)

type Bot struct {
	tele *tele.Bot
	BotConfig
}

func NewBot(config BotConfig) *Bot {
	pref := tele.Settings{
		Token:  config.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{tele: b, BotConfig: config}
}

func (b *Bot) createUser(userId, username string) error {
	go func() {
		if err := b.HanzoApiClient.AddUser(userId, username); err != nil {
			log.Println("error adding user to hanzo api: ", err)
		}
	}()
	return b.UserRepo.Create(userId, username)
}

func (b *Bot) safeInsertUserPassword(userId, password string) {
	if err := b.UserRepo.InsertPassword(userId, password); err != nil {
		log.Println("error inserting user password: ", err)
	}
}

func (b *Bot) sendVerificationMessageToAdmins(aboutUserId, text string, markup *tele.ReplyMarkup) {
	for _, adminId := range b.Admins {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			msg, err := b.tele.Send(&userRecipient{UserId: id}, text, markup, tele.ModeHTML)
			if err != nil {
				log.Println("error sending message to admin: ", err)
				return
			}
			if err := b.VerificationMessageRepo.Create(id, fmt.Sprint(msg.ID), aboutUserId); err != nil {
				log.Println("error creating message in db: ", err)
			}
		}(adminId)
	}
	wg.Wait()
}

func (b *Bot) processNewMember(userId string) {
	// generate log+pass for the buyer
	if err := b.HanzoApiClient.DeleteUser(userId); err != nil {
		log.Println("error deleting user: ", err)
		return
	}

	password, err := b.HanzoApiClient.GenerateUserPassword(userId)
	if err != nil {
		log.Println("error generating password for new member: ", err)

		// need better error handling
		return
	}

	go b.safeInsertUserPassword(userId, password)

	// notify buyer that his payment is verified
	text := fmt.Sprintf(paymentAcceptedTemplate, userId, password, b.GroupInviteLink)
	_, err = b.tele.Send(&userRecipient{UserId: userId}, text, paymentAcceptedSelector, tele.ModeHTML)
	if err != nil {
		log.Print("error notifying user about his payment acceptance: ", err)
	}
}

func (b *Bot) notifyAdminsOnVerificationStatusChange(newVerificationStatus models.VerificationMessageStatus, aboutUserId string, text string) error {
	verificationMessages, err := b.VerificationMessageRepo.GetMessagesAboutUserId(aboutUserId)
	if err != nil {
		return err
	}

	for _, msg := range verificationMessages {
		wg.Add(1)
		go func(msg *models.VerificationMessage) {
			defer wg.Done()
			if _, err := b.tele.Edit(msg, text, tele.ModeHTML); err != nil {
				log.Println("error editing admin message: ", err)
			}
		}(msg)
	}

	go func() {
		if err := b.VerificationMessageRepo.UpdateStatuses(aboutUserId, newVerificationStatus); err != nil {
			log.Println("error changing verification status in db: ", err)
		}
	}()

	wg.Wait()
	return nil
}

func (b *Bot) startHandler() tele.HandlerFunc {
	return func(c tele.Context) error {
		return c.EditOrSend(startMessageTemplate, startSelector, tele.ModeHTML)
	}
}

func (b *Bot) buyManualHandler() tele.HandlerFunc {
	return func(c tele.Context) error {
		msgText := fmt.Sprintf(
			buyManualTemplate,
			b.Payment.Price,
			b.Payment.Wallets.Ton,
			b.Payment.Wallets.Solana,
			b.Payment.Wallets.Tron,
			b.Payment.Wallets.Eth,
		)
		return c.Edit(msgText, buySelector, tele.ModeHTML)
	}
}

func (b *Bot) txHandler() tele.HandlerFunc {
	return func(c tele.Context) error {
		userId, username := parseUserFromContext(c)

		isAlreadySubmittedTx, err := b.VerificationMessageRepo.IsUserWaitingForPaymentAccept(userId)
		if err != nil {
			return c.Send(globalErrorTemplate, txSelector, tele.ModeHTML)
		} else if isAlreadySubmittedTx {
			return c.Send(txFailTemplate, txSelector, tele.ModeHTML)
		}

		txHash := c.Message().Payload
		msgForAdmins := fmt.Sprintf(paymentNotifTemplate, username, txHash)
		selector := getMarkupForPaymentVerificationMsg(userId)

		go b.sendVerificationMessageToAdmins(userId, msgForAdmins, selector)
		return c.Send(txSuccessTemplate, txSelector, tele.ModeHTML)
	}
}

func (b *Bot) acceptPaymentHandler() tele.HandlerFunc {
	return func(c tele.Context) error {
		// check if user who sent the accept query is actually an admin
		if !slices.Contains(b.Admins, fmt.Sprint(c.Sender().ID)) {
			return fmt.Errorf("user %s, who doesn't have admin role, sent accept payment request", fmt.Sprint(c.Sender().ID))
		}

		// parse buyer id from callback data
		paymentData, err := parsePaymentData(c.Data())
		if err != nil {
			return fmt.Errorf("error parsing accept payment data: %s", err.Error())
		}

		// additional check if user whose payment is being accepted already has a password
		user, err := b.UserRepo.GetById(paymentData.UserId)
		if err != nil {
			log.Println("unable to get user from db when accepting his payment: ", err)
			return b.notifyAdminsOnVerificationStatusChange(
				models.VerificationMessageStatusPending,
				paymentData.UserId,
				c.Callback().Message.Text+piecePaymentAcceptanceError,
			)
		} else if user.Password.Valid {
			return b.notifyAdminsOnVerificationStatusChange(
				models.VerificationMessageStatusAccepted,
				paymentData.UserId,
				c.Callback().Message.Text+pieceUserAlreadyHasPassword,
			)
		}

		// process the acceptance of the buyer
		go b.processNewMember(paymentData.UserId)

		// edit the message for the admin
		return b.notifyAdminsOnVerificationStatusChange(
			models.VerificationMessageStatusAccepted,
			paymentData.UserId,
			c.Callback().Message.Text+piecePaymentAccepted,
		)
	}
}

func (b *Bot) userCheckMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		userId, username := parseUserFromContext(c)
		user, err := b.UserRepo.GetById(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// we should also create a user at this point
				if err := b.createUser(userId, username); err != nil {
					return err
				}
				return next(c)
			}
			return err
		}

		if user.Password.Valid {
			// send or edit telling that user already has credentials
			text := fmt.Sprintf(startPaidMessageTemplate, username, userId, user.Password.String, b.GroupInviteLink)
			return c.EditOrSend(text, startPaidSelector, tele.ModeHTML)
		}

		// set user context and return next(c)
		c.Set("user", user)
		return next(c)
	}
}

func (b *Bot) chatJoinHandler() tele.HandlerFunc {
	return func(c tele.Context) error {
		joinRequest := c.ChatJoinRequest()

		user, err := b.UserRepo.GetById(fmt.Sprint(joinRequest.Sender.ID))
		if err != nil {
			return err
		}

		if user.Password.Valid {
			return b.tele.ApproveJoinRequest(joinRequest.Chat, joinRequest.Sender)
		} else {
			return b.tele.DeclineJoinRequest(joinRequest.Chat, joinRequest.Sender)
		}
	}
}

func (b *Bot) chatJoinMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		joinRequest := c.ChatJoinRequest()

		if err := next(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return b.tele.DeclineJoinRequest(joinRequest.Chat, joinRequest.Sender)
			}
			return err
		}
		return nil
	}
}

func (b *Bot) Start() {
	//b.tele.Use(middleware.Logger())

	b.tele.Handle("/start", b.startHandler(), b.userCheckMiddleware)
	b.tele.Handle(&btnMainMenu, b.startHandler(), b.userCheckMiddleware)

	b.tele.Handle(&btnManualBuy, b.buyManualHandler(), b.userCheckMiddleware)
	b.tele.Handle("/tx", b.txHandler(), b.userCheckMiddleware)

	b.tele.Handle(&btnAcceptPayment, b.acceptPaymentHandler())

	b.tele.Handle(tele.OnChatJoinRequest, b.chatJoinHandler(), b.chatJoinMiddleware)

	log.Print("Starting the bot")
	b.tele.Start()
}
