package ory

import (
    kratos "github.com/ory/kratos-client-go"
)

var kratosClient *kratos.APIClient

type KratosConfig struct {
    AdminURL  string `env:"KRATOS_ADMIN_URL"`
    PublicURL string `env:"KRATOS_PUBLIC_URL"`
}

// InitKratos برای مقداردهی اولیه کلاینت Kratos
func InitKratos(cfg *KratosConfig) {
    config := kratos.NewConfiguration()
    config.Servers = kratos.ServerConfigurations{
        {
            URL: cfg.AdminURL, // استفاده از Admin API برای عملیات مدیریتی
        },
        {
            URL: cfg.PublicURL, // استفاده از Public API برای عملیات عمومی
        },
    }
    kratosClient = kratos.NewAPIClient(config)
}

func GetKratos() *kratos.APIClient {
    return kratosClient
}