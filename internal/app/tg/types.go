package tg

import "hanzotg/internal/app/models"

type PaymentWallets struct {
	Solana string
	Tron   string
	Eth    string
	Ton    string
}

type PaymentConfig struct {
	Price   uint32
	Wallets PaymentWallets
}

type BotConfig struct {
	GroupInviteLink         string
	BotToken                string
	Admins                  []string
	Payment                 PaymentConfig
	UserRepo                IUserRepo
	VerificationMessageRepo IVerificationMessageRepo
	HanzoApiClient          IHanzoApiClient
}

type paymentVerificationData struct {
	UserId string `json:"userId"`
}

type IUserRepo interface {
	Create(id, username string) error
	GetById(id string) (*models.User, error)
	InsertPassword(userId, password string) error
}

type IVerificationMessageRepo interface {
	Create(chatId, messageId, aboutUserId string) error
	UpdateStatuses(aboutUserId string, newStatus models.VerificationMessageStatus) error
	GetMessagesAboutUserId(aboutUserId string) ([]*models.VerificationMessage, error)
	IsUserWaitingForPaymentAccept(userId string) (bool, error)
}

type IHanzoApiClient interface {
	AddUser(userId, username string) error
	GenerateUserPassword(userId string) (string, error)
	DeleteUser(userId string) error
}
