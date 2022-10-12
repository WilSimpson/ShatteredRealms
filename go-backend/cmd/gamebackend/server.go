package main

import (
	aapb "agones.dev/agones/pkg/allocation/go"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"time"
)

const (
	retryAfter = time.Second * 10
)

func NewServer(
	jwt service.JWTService,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	authorizationServer, err := grpc.Dial(conf.AccountsAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	charactersServer, err := grpc.Dial(conf.CharactersAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	var publicRPCs = map[string]struct{}{
		"/sro.gamebackend.ConnectionService/Connect": {},
	}

	authInterceptor := interceptor.NewAuthInterceptor(
		jwt,
		publicRPCs,
		srv.GetPermissions(pb.NewAuthorizationServiceClient(authorizationServer), jwt, "sro.com/gamebackend/v1"),
	)

	go srv.ProcessRoleUpdates(pb.NewAuthorizationServiceClient(authorizationServer), authInterceptor, jwt, "sro.com/gamebackend/v1")
	go srv.ProcessUserUpdates(pb.NewAuthorizationServiceClient(authorizationServer), authInterceptor, jwt, "sro.com/gamebackend/v1")

	localHostMode := conf.Mode == ModeLocalHost

	var allocatorServer *grpc.ClientConn
	if !localHostMode {
		allocatorServer, err = dialAgonesAllocatorServer()
		if err != nil {
			return nil, nil, err
		}
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), helpers.UnaryLogRequest()),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), helpers.StreamLogRequest()),
	)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	connectionServiceServer := srv.NewConnectionServiceServer(
		jwt,
		aapb.NewAllocationServiceClient(allocatorServer),
		pb.NewCharactersServiceClient(charactersServer),
		localHostMode,
		conf.AgonesNamespace,
	)

	pb.RegisterConnectionServiceServer(grpcServer, connectionServiceServer)
	err = pb.RegisterConnectionServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Address(),
		opts,
	)

	return grpcServer, gwmux, err
}

func dialAgonesAllocatorServer() (*grpc.ClientConn, error) {
	clientKey, err := ioutil.ReadFile(conf.AgonesKeyFile)
	if err != nil {
		return nil, err
	}

	clientCert, err := ioutil.ReadFile(conf.AgonesCertFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(conf.AgonesCaCertFile)
	if err != nil {
		return nil, err
	}

	opt, err := createRemoteClusterDialOption(clientCert, clientKey, caCert)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(conf.AgonesAllocatorAddress(), opt)
}

// createRemoteClusterDialOption creates a grpc client dial option with TLS configuration.
func createRemoteClusterDialOption(clientCert, clientKey, caCert []byte) (grpc.DialOption, error) {
	// Load client cert
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	if len(caCert) != 0 {
		// Load CA cert, if provided and trust the server certificate.
		// This is required for self-signed certs.
		tlsConfig.RootCAs = x509.NewCertPool()
		if !tlsConfig.RootCAs.AppendCertsFromPEM(caCert) {
			return nil, errors.New("only PEM format is accepted for server CA")
		}
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), nil
}
