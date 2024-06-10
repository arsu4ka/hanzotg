package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	BotToken     string `env:"BOT_TOKEN" env-required:"true"`
	Price        uint32 `env:"PRICE" env-required:"true"`
	SolanaWallet string `env:"SOLANA_WALLET" env-required:"true"`
	TronWallet   string `env:"TRON_WALLET" env-required:"true"`
	TonWallet    string `env:"TON_WALLET" env-required:"true"`
	EthWallet    string `env:"ETH_WALLET" env-required:"true"`
	Admins       string `env:"ADMINS" env-required:"true"`

	DbHost     string `env:"DB_HOST" env-required:"true"`
	DbPort     string `env:"DB_PORT" env-required:"true"`
	DbUser     string `env:"DB_USER" env-required:"true"`
	DbPassword string `env:"DB_PASSWORD" env-required:"true"`
	DbName     string `env:"DB_NAME" env-required:"true"`

	MockHanzoApi   bool   `env:"MOCK_HANZO_API"`
	HanzoApiUrl    string `env:"HANZO_API_URL" env-required:"true"`
	HanzoApiSecret string `env:"HANZO_API_SECRET" env-required:"true"`

	GroupInviteLink string `env:"GROUP_INVITE_LINK" env-required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatal("error reading config: ", err)
	}

	return cfg
}
