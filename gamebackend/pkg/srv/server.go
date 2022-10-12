package srv

import (
    aapb "agones.dev/agones/pkg/allocation/go"
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    accountspb "github.com/ShatteredRealms/Accounts/pkg/pb"
    accountssrv "github.com/ShatteredRealms/Accounts/pkg/srv"
    characterspb "github.com/ShatteredRealms/Characters/pkg/pb"
    "github.com/ShatteredRealms/GoUtils/pkg/interceptor"
    utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
    "github.com/ShatteredRealms/gamebackend/internal/log"
    "github.com/ShatteredRealms/gamebackend/internal/option"
    "github.com/ShatteredRealms/gamebackend/pkg/pb"
    "github.com/golang-jwt/jwt"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/pkg/errors"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/metadata"
    "google.golang.org/protobuf/types/known/emptypb"
    "io"
    "io/ioutil"
    "time"
)

const (
    retryAfter = time.Second * 10
)

func NewServer(
    jwt utilService.JWTService,
    config option.Config,
    logger log.LoggerService,
) (*grpc.Server, *runtime.ServeMux, error) {
    ctx := context.Background()

    authorizationServer, err := dialAuthorizationServer(config)
    if err != nil {
        return nil, nil, err
    }

    charactersServer, err := dialCharactersServer(config)
    if err != nil {
        return nil, nil, err
    }

    var publicRPCs = map[string]struct{}{
        "/sro.gamebackend.ConnectionService/Connect": {},
    }

    authInterceptor := interceptor.NewAuthInterceptor(
        jwt,
        publicRPCs,
        accountssrv.GetPermissions(accountspb.NewAuthorizationServiceClient(authorizationServer), jwt, "shatteredrealmsonline.com/gamebackend/v1"),
    )

    go processRoleUpdates(accountspb.NewAuthorizationServiceClient(authorizationServer), authInterceptor, jwt, logger)
    go processUserUpdates(accountspb.NewAuthorizationServiceClient(authorizationServer), authInterceptor, jwt, logger)

    localHostMode := config.Mode.Value == option.LocalhostMode

    var allocatorServer *grpc.ClientConn
    if !localHostMode {
        allocatorServer, err = dialAgonesAllocatorServer(config)
        if err != nil {
            return nil, nil, err
        }
    }

    grpcServer := grpc.NewServer(
        grpc.ChainUnaryInterceptor(authInterceptor.Unary(), logger.UnaryLogRequest(log.Info)),
        grpc.ChainStreamInterceptor(authInterceptor.Stream(), logger.StreamLogRequest(log.Info)),
    )

    gwmux := runtime.NewServeMux()
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    }

    connectionServiceServer := NewConnectionServiceServer(
        jwt,
        aapb.NewAllocationServiceClient(allocatorServer),
        characterspb.NewCharactersServiceClient(charactersServer),
        localHostMode,
        config.AgonesNamespace.Value,
    )

    pb.RegisterConnectionServiceServer(grpcServer, connectionServiceServer)
    err = pb.RegisterConnectionServiceHandlerFromEndpoint(
        ctx,
        gwmux,
        config.Address(),
        opts,
    )

    return grpcServer, gwmux, err
}

func dialAuthorizationServer(
    config option.Config,
) (*grpc.ClientConn, error) {
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    return grpc.Dial(config.AccountsAddress(), opts...)
}

func dialCharactersServer(
    config option.Config,
) (*grpc.ClientConn, error) {
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    return grpc.Dial(config.CharactersAddress(), opts...)
}

func dialAgonesAllocatorServer(
    config option.Config,
) (*grpc.ClientConn, error) {
    clientKey, err := ioutil.ReadFile(config.AgonesKeyFile.Value)
    if err != nil {
        return nil, err
    }

    clientCert, err := ioutil.ReadFile(config.AgonesCertFile.Value)
    if err != nil {
        return nil, err
    }

    caCert, err := ioutil.ReadFile(config.AgonesCaCertFile.Value)
    if err != nil {
        return nil, err
    }

    opt, err := createRemoteClusterDialOption(clientCert, clientKey, caCert)
    if err != nil {
        return nil, err
    }

    return grpc.Dial(config.AgonesAllocatorAddress(), opt)
}

func processUserUpdates(
    authorizationClient accountspb.AuthorizationServiceClient,
    interceptor *interceptor.AuthInterceptor,
    jwtService utilService.JWTService,
    logger log.LoggerService,
) {
    userUpdatesClient, err := authorizationClient.SubscribeUserUpdates(serverAuthContext(jwtService), &emptypb.Empty{})
    if err != nil {
        logger.Logf(log.Error, "Unable to subscribe to user updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
        time.Sleep(retryAfter)
        processUserUpdates(authorizationClient, interceptor, jwtService, logger)
        return
    }
    logger.Info("Successfully subscribed to user updates from authorization server.")
    for {
        msg, err := userUpdatesClient.Recv()

        if err == nil {
            logger.Logf(log.Debug, "Update to user %d permissions. Clearing permissions cache for that user.", msg.Id)
            err = interceptor.ClearUserCache(uint(msg.Id))
            if err != nil {
                logger.Logf(log.Warning, "Clearing cache: %v", err)
            }
        } else if err == io.EOF {
            logger.Infof("User updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
            time.Sleep(retryAfter)
            processUserUpdates(authorizationClient, interceptor, jwtService, logger)
            return
        } else {
            logger.Logf(log.Error, "User updates: %v.", err)
            logger.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
            time.Sleep(retryAfter)
            processUserUpdates(authorizationClient, interceptor, jwtService, logger)
            return
        }
    }
}

func processRoleUpdates(
    authorizationClient accountspb.AuthorizationServiceClient,
    interceptor *interceptor.AuthInterceptor,
    jwtService utilService.JWTService,
    logger log.LoggerService,
) {
    roleUpdatesClient, err := authorizationClient.SubscribeRoleUpdates(serverAuthContext(jwtService), &emptypb.Empty{})
    if err != nil {
        logger.Logf(log.Error, "Unable to subscribe to role updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
        time.Sleep(retryAfter)
        processRoleUpdates(authorizationClient, interceptor, jwtService, logger)
        return
    }
    logger.Info("Successfully subscribed to role updates from authorization server.")
    for {
        msg, err := roleUpdatesClient.Recv()
        if err == nil {
            logger.Logf(log.Debug, "Update to role %d permissions. Clearing permissions cache for all users.", msg.Id)
            err = interceptor.ClearCache()
            if err != nil {
                logger.Logf(log.Warning, "Clearing cache: %v", err)
            }
        } else if err == io.EOF {
            logger.Infof("Role updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
            time.Sleep(retryAfter)
            processRoleUpdates(authorizationClient, interceptor, jwtService, logger)
            return
        } else {
            logger.Logf(log.Error, "Role Updates: %v.", err)
            logger.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
            time.Sleep(retryAfter)
            processRoleUpdates(authorizationClient, interceptor, jwtService, logger)
            return
        }
    }
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

func serverAuthContext(jwtService utilService.JWTService) context.Context {
    md := metadata.New(
        map[string]string{
            "authorization": fmt.Sprintf(
                "Bearer %s", generateTemporaryServerToken(jwtService, "shatteredrealmsonline.com/gamebackend/v1"),
            ),
        },
    )
    return metadata.NewOutgoingContext(context.Background(), md)
}

func generateTemporaryServerToken(jwtService utilService.JWTService, requestingHost string) string {
    out, _ := jwtService.Create(time.Second*10, requestingHost, jwt.MapClaims{"sub": 0})
    return out
}
