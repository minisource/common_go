package ory

import(
	hydra "github.com/ory/hydra-client-go"
)

//https://www.ory.sh/docs/hydra/reference/api#tag/oAuth2
var hydraClient *hydra.APIClient

type HydraConfig struct{
    Url string `env:"HYDRA_URL"`
}

func InitHydra(cfg *HydraConfig) {
	config := hydra.NewConfiguration()
    config.Servers = hydra.ServerConfigurations{
        {
            URL: cfg.Url,
        },
    }
    hydraClient = hydra.NewAPIClient(config)
}

func GetHydra() *hydra.APIClient {
	return hydraClient
}