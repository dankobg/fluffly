package keto

import (
	"fmt"

	opl "github.com/ory/keto/proto/ory/keto/opl/v1alpha1"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Read   rts.ReadServiceClient
	Write  rts.WriteServiceClient
	Check  rts.CheckServiceClient
	Expand rts.ExpandServiceClient
	Syntax opl.SyntaxServiceClient
}

func NewClient() (*Client, error) {
	rConn, err := grpc.NewClient("localhost:4466", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to keto grpc read server: %w", err)
	}

	wConn, err := grpc.NewClient("localhost:4467", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to keto grpc write server: %w", err)
	}

	c := &Client{
		Read:   rts.NewReadServiceClient(rConn),
		Write:  rts.NewWriteServiceClient(wConn),
		Check:  rts.NewCheckServiceClient(rConn),
		Expand: rts.NewExpandServiceClient(rConn),
		Syntax: opl.NewSyntaxServiceClient(rConn),
	}

	return c, nil
}
