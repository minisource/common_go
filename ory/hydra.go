package ory

import(
	hydra "github.com/ory/hydra-client-go"
)

//https://www.ory.sh/docs/hydra/reference/api#tag/oAuth2
var hydraClient *hydra.APIClient

type HydraConfig struct{
    PrivateUrl string `env:"HYDRA_URL_PRIVATE"`
    PublicUrl string `env:"HYDRA_URL_PUBLIC"`
}

func InitHydra(cfg *HydraConfig) {
	config := hydra.NewConfiguration()
    config.Servers = hydra.ServerConfigurations{
        {
            URL: cfg.PrivateUrl,
        },
    }
    hydraClient = hydra.NewAPIClient(config)
}

func GetHydra() *hydra.APIClient {
	return hydraClient
}