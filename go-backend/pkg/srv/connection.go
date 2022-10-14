package srv

import (
	aapb "agones.dev/agones/pkg/allocation/go"
	"context"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	utilService "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type connectionServiceServer struct {
	pb.UnimplementedConnectionServiceServer
	jwtService utilService.JWTService
	allocator  aapb.AllocationServiceClient
	characters pb.CharactersServiceClient
	tracer     trace.Tracer

	// localhostMode used to tell whether to search for a server on kubernetes or return a constant localhost connection
	localhostMode bool

	// namespace kubernetes namespace to search for gameservers in
	namespace string
}

func NewConnectionServiceServer(
	jwtService utilService.JWTService,
	allocator aapb.AllocationServiceClient,
	characters pb.CharactersServiceClient,
	namespace string,
	localHostMode bool,
) pb.ConnectionServiceServer {

	return &connectionServiceServer{
		jwtService:    jwtService,
		allocator:     allocator,
		characters:    characters,
		localhostMode: localHostMode,
		namespace:     namespace,
		tracer:        otel.Tracer("connection"),
	}
}

func (s *connectionServiceServer) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	if s.localhostMode {
		return &pb.ConnectResponse{
			Address: "127.0.0.1",
			Port:    7777,
		}, nil
	}

	// If the current user can't get the character, then deny the request
	//character, err := s.characters.GetCharacter(
	//    ctx,
	//    &pb.CharacterTarget{CharacterId: request.CharacterId},
	//)
	//if err != nil {
	//
	//    fmt.Println("err 1")
	//    return nil, err
	//}

	//world := "Scene_Demo"
	//if character.Location != nil && character.Location.World != "" {
	//    world = character.Location.World
	//}

	allocatorReq := &aapb.AllocationRequest{
		Namespace: s.namespace,
		GameServerSelectors: []*aapb.GameServerSelector{
			{
				//MatchLabels: map[string]string{
				//    "world": world,
				//},
				GameServerState: aapb.GameServerSelector_ALLOCATED,
				Players: &aapb.PlayerSelector{
					MinAvailable: 1,
					MaxAvailable: 1000,
				},
			},
			{
				//MatchLabels: map[string]string{
				//    "world": world,
				//},
				GameServerState: aapb.GameServerSelector_READY,
				Players: &aapb.PlayerSelector{
					MinAvailable: 1,
					MaxAvailable: 1000,
				},
			},
		},
	}

	allocatorResp, err := s.allocator.Allocate(serverAuthContext(s.jwtService, "sro.com/gamebackend/v1/"), allocatorReq)
	if err != nil {
		fmt.Println("err 2")
		//fmt.Printf("world: %s", world)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ConnectResponse{
		Address: allocatorResp.Address,
		Port:    uint32(allocatorResp.Ports[0].Port),
	}, nil
}
