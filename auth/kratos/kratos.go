package kratos

import orykratos "github.com/ory/client-go"

type Client struct {
	Public *orykratos.APIClient
	Admin  *orykratos.APIClient
}

func NewClient(publicURL, adminURL string) *Client {
	publicConf := orykratos.NewConfiguration()
	publicConf.Servers = orykratos.ServerConfigurations{{
		URL:         publicURL,
		Description: "Kratos public server",
	}}

	adminConf := orykratos.NewConfiguration()
	adminConf.Servers = orykratos.ServerConfigurations{{
		URL:         adminURL,
		Description: "Kratos admin server",
	}}

	c := &Client{
		Public: orykratos.NewAPIClient(publicConf),
		Admin:  orykratos.NewAPIClient(adminConf),
	}

	return c
}
