package srv

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

    authInterceptor := interceptor.NewAuthInterceptor(
        jwt,
        publicRPCs,
        accountssrv.GetPermissions(accountspb.NewAuthorizationServiceClient(authorizationServer), jwt, "shatteredrealmsonline.com/gamebackend/v1"),
    )

    go processRoleUpdates(accountspb.NewAuthorizationServiceClient(authorizationServer), authInterceptor, jwt, logger)
    go processUserUpdates(accountspb.NewAuthorizationServiceClient(authorizationServer), authInterceptor, jwt, logger)

    grpcServer := grpc.NewServer(
        grpc.ChainUnaryInterceptor(authInterceptor.Unary(), logger.UnaryLogRequest(log.Info)),
        grpc.ChainStreamInterceptor(authInterceptor.Stream(), logger.StreamLogRequest(log.Info)),
    )

    gwmux := runtime.NewServeMux()
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    }

}

func dialAuthorizationServer(
    config option.Config,
) (*grpc.ClientConn, error) {
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    return grpc.Dial(config.AccountsAddress(), opts...)
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
