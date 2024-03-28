package tasker

import (
	"context"

	"github.com/openfga/go-sdk/client"
)

// OpenFGAClientInterface is the interface used to decouple the OpenFGA store implementation
type OpenFGAClientInterface interface {
	ReadTuples(context.Context, string, string, string, string) (*client.ClientReadResponse, error)
}

type OpenFGAInterface interface {
	GetAuthorizationModelId() (string, error)
	CreateStore(context.Context) client.SdkClientCreateStoreRequestInterface
	CreateStoreExecute(client.SdkClientCreateStoreRequestInterface) (*client.ClientCreateStoreResponse, error)
	ReadAuthorizationModel(context.Context) client.SdkClientReadAuthorizationModelRequestInterface
	ReadAuthorizationModelExecute(client.SdkClientReadAuthorizationModelRequestInterface) (*client.ClientReadAuthorizationModelResponse, error)
	ReadAuthorizationModels(context.Context) client.SdkClientReadAuthorizationModelsRequestInterface
	ReadAuthorizationModelsExecute(client.SdkClientReadAuthorizationModelsRequestInterface) (*client.ClientReadAuthorizationModelsResponse, error)
	WriteAuthorizationModel(context.Context) client.SdkClientWriteAuthorizationModelRequestInterface
	WriteAuthorizationModelExecute(client.SdkClientWriteAuthorizationModelRequestInterface) (*client.ClientWriteAuthorizationModelResponse, error)
	Read(context.Context) client.SdkClientReadRequestInterface
	ReadExecute(client.SdkClientReadRequestInterface) (*client.ClientReadResponse, error)
	Check(context.Context) client.SdkClientCheckRequestInterface
	CheckExecute(client.SdkClientCheckRequestInterface) (*client.ClientCheckResponse, error)
	BatchCheck(context.Context) client.SdkClientBatchCheckRequestInterface
	BatchCheckExecute(client.SdkClientBatchCheckRequestInterface) (*client.ClientBatchCheckResponse, error)
	Write(context.Context) client.SdkClientWriteRequestInterface
	WriteExecute(client.SdkClientWriteRequestInterface) (*client.ClientWriteResponse, error)
	ListObjects(context.Context) client.SdkClientListObjectsRequestInterface
	ListObjectsExecute(client.SdkClientListObjectsRequestInterface) (*client.ClientListObjectsResponse, error)
}
