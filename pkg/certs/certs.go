package certs

import (
	"fmt"

	"google.golang.org/grpc/credentials"
)

func loadCerts(credsFile string) (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile(credsFile, "")
	if err != nil {
		return nil, fmt.Errorf("gateway cert load error: %s", err)
	}
	return creds, nil
}
