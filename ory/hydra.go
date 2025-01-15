package ory

import(
	hydra "github.com/ory/hydra-client-go"
)

//https://www.ory.sh/docs/hydra/reference/api#tag/oAuth2
var hydraClient *hydra.APIClient

type HydraConfig struct{
    AdminURL string `env:"HYDRA_ADMIN_URL"`
    PublicURL string `env:"HYDRA_PUBLIC_URL"`
}

func InitHydra(cfg *HydraConfig) {
	config := hydra.NewConfiguration()
    config.Servers = hydra.ServerConfigurations{
        {
            URL: cfg.AdminURL, // استفاده از Admin API برای عملیات مدیریتی
        },
        {
            URL: cfg.PublicURL, // استفاده از Public API برای عملیات عمومی
        },
    }
    hydraClient = hydra.NewAPIClient(config)
}

func GetHydra() *hydra.APIClient {
	return hydraClient
}