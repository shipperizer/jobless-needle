package tasker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/openfga/go-sdk/client"
	"github.com/openfga/go-sdk/credentials"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type OFGAClient struct {
	c OpenFGAInterface
}

func (c *OFGAClient) ReadTuples(ctx context.Context, user, relation, object, continuationToken string) (*client.ClientReadResponse, error) {
	r := c.c.Read(ctx)

	body := client.ClientReadRequest{
		User:     &user,
		Relation: &relation,
		Object:   &object,
	}

	r = r.Body(body).Options(client.ClientReadOptions{ContinuationToken: &continuationToken})
	res, err := c.c.ReadExecute(r)

	return res, err
}

func NewOFGAClient() *OFGAClient {
	c := new(OFGAClient)

	fga, err := client.NewSdkClient(
		&client.ClientConfiguration{
			ApiScheme: "http",
			ApiHost:   "localhost:8080",
			StoreId:   "01HT24CVRRBSF4JQTXP67HVT27",
			Credentials: &credentials.Credentials{
				Method: credentials.CredentialsMethodApiToken,
				Config: &credentials.Config{
					ApiToken: "42",
				},
			},
			AuthorizationModelId: "01HT24DTA107CEYCG3TA70J3KJ",
			Debug:                true,
			HTTPClient:           &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
		},
	)

	if err != nil {
		panic(fmt.Sprintf("issues setting up OpenFGA client %s", err))
	}

	c.c = fga

	return c
}
