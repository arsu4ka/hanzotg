package main

import (
	"hanzotg/internal/app/config"
	"hanzotg/internal/app/db"
	"hanzotg/internal/app/hanzo"
	"hanzotg/internal/app/repo"
	"hanzotg/internal/app/tg"
	"log"
	"strings"
)

func main() {
	cfg := config.MustLoad()
	admins := strings.Split(cfg.Admins, ",")

	postgresDb, err := db.OpenPostgres(db.PostgresOptions{
		Host:     cfg.DbHost,
		Port:     cfg.DbPort,
		Username: cfg.DbUser,
		Password: cfg.DbPassword,
		DbName:   cfg.DbName,
	})
	if err != nil {
		log.Fatal("can't open db: ", err)
	}

	userRepo := repo.NewUser(postgresDb)
	verificationMessageRepo := repo.NewVerificationMessage(postgresDb)

	var hanzoApiClient tg.IHanzoApiClient
	if cfg.MockHanzoApi {
		hanzoApiClient = &hanzo.MockApiClient{}
	} else {
		hanzoApiClient = hanzo.NewApiClient(hanzo.ApiConfig{
			SecretKey: cfg.HanzoApiSecret,
			BaseUrl:   cfg.HanzoApiUrl,
		})
	}

	tgBot := tg.NewBot(tg.BotConfig{
		BotToken: cfg.BotToken,
		Admins:   admins,
		Payment: tg.PaymentConfig{
			Price: cfg.Price,
			Wallets: tg.PaymentWallets{
				Solana: cfg.SolanaWallet,
				Tron:   cfg.TronWallet,
				Eth:    cfg.EthWallet,
				Ton:    cfg.TonWallet,
			},
		},
		UserRepo:                userRepo,
		VerificationMessageRepo: verificationMessageRepo,
		HanzoApiClient:          hanzoApiClient,
		GroupInviteLink:         cfg.GroupInviteLink,
	})
	tgBot.Start()
}
