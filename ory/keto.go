package ory

import (
	keto "github.com/ory/keto-client-go"
)

var ketoClient *keto.APIClient

type KetoConfig struct {
    AdminURL  string `env:"Keto_ADMIN_URL"`
    PublicURL string `env:"Keto_PUBLIC_URL"`
}

func InitKeto(cfg *KratosConfig) {
    config := keto.NewConfiguration()
    config.Servers = keto.ServerConfigurations{
        {
            URL: cfg.AdminURL, // استفاده از Admin API برای عملیات مدیریتی
        },
        {
            URL: cfg.PublicURL, // استفاده از Public API برای عملیات عمومی
        },
    }
    ketoClient = keto.NewAPIClient(config)
}

func GetKeto() *keto.APIClient {
    return ketoClient
}