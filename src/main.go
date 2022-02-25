package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dd "github.com/pbkit/lol-champion-server/datadragon"
	proto "github.com/pbkit/lol-champion-server/gen"
)

var champions = dd.LoadChampions()

const (
	httpPort = "8080"
	grpcPort = "8081"
)

type LolChampionServer struct {
	proto.UnimplementedLolChampionServiceServer
}

func (s *LolChampionServer) ListChampions(ctx context.Context, in *proto.ListChampionsRequest) (*proto.ListChampionsResponse, error) {
	return ListChampions(champions), nil
}

func (s *LolChampionServer) GetChampionStory(ctx context.Context, in *proto.GetChampionStoryRequest) (*proto.GetChampionStoryResponse, error) {
	if res := GetChampionStory(champions, in.ChampionKey); res != nil {
		return res, nil
	}
	return &proto.GetChampionStoryResponse{}, status.Errorf(codes.NotFound, "Not Found: %s", in.ChampionKey)
}

func (s *LolChampionServer) GetChampionStats(ctx context.Context, in *proto.GetChampionStatsRequest) (*proto.GetChampionStatsResponse, error) {
	if res := GetChampionStats(champions, in.ChampionKey); res != nil {
		return res, nil
	}
	return &proto.GetChampionStatsResponse{}, status.Errorf(codes.NotFound, "Not Found: %s", in.ChampionKey)
}

func (s *LolChampionServer) GetChampionSkills(ctx context.Context, in *proto.GetChampionSkillsRequest) (*proto.GetChampionSkillsResponse, error) {
	if res := GetChampionSkills(champions, in.ChampionKey); res != nil {
		return res, nil
	}
	return &proto.GetChampionSkillsResponse{}, status.Errorf(codes.NotFound, "Not Found: %s", in.ChampionKey)
}

func main() {
	grpcServer := grpc.NewServer()
	proto.RegisterLolChampionServiceServer(grpcServer, &LolChampionServer{})
	grpcWebServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)

	httpServer := http.Server{
		Addr: ":" + httpPort,
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if grpcWebServer.IsGrpcWebRequest(req) || grpcWebServer.IsAcceptableGrpcCorsRequest(req) {
				grpcWebServer.ServeHTTP(resp, req)
				return
			}
			http.DefaultServeMux.ServeHTTP(resp, req)
		}),
	}

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		log.Printf("gRPC-web server is listening on :%v", httpPort)
		if err := httpServer.ListenAndServe(); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		listener, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			return err
		}
		log.Printf("gRPC server is listening on :%v", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			return err
		}
		return nil
	})

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
		if err := httpServer.Close(); err != nil {
			panic(err)
		}
	}()

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
