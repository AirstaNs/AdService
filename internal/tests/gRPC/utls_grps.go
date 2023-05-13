package gRPC

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"homework10/internal/adapters/repository/adrepo"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/app"
	"homework10/internal/entities"
	grpc2 "homework10/internal/ports/grpc"
	"homework10/internal/util"
	"log"
	"net"
	"time"
)

type gRPCtestClient struct {
	Server grpc2.AdServiceClient
	Stop   func()
}

var (
	errNotFound  = status.Error(codes.NotFound, "not found")
	errForbidden = status.Error(codes.PermissionDenied, "permission denied")
)

const (
	name  = "Oleg"
	email = name + "@mail.ru"
	title = "phone"
	text  = "buy new phone"
)

func getGRPCTestClient() *gRPCtestClient {
	addr := "localhost:50051"
	repo := adrepo.New()
	uRep := userrepo.New()
	formatter := util.NewDateTimeFormatter(time.RFC3339)
	newApp := app.NewApp(repo, uRep, formatter)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(nil))
	grpc2.RegisterAdServiceServer(grpcServer, grpc2.GServer{App: newApp})
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	AdsClient := grpc2.NewAdServiceClient(conn)

	// Функция для остановки сервера gRPC
	stop := func() {
		grpcServer.Stop()
	}
	return &gRPCtestClient{Server: AdsClient, Stop: stop}
}

func setupUsers(client *gRPCtestClient) ([]entities.User, error) {
	user, err := addUser(client, name, email)

	var empty []entities.User
	if err != nil {
		return empty, err
	}
	user1, err1 := addUser(client, email, name)
	if err1 != nil {
		return empty, err
	}

	return []entities.User{user, user1}, nil
}
func setupAds(client *gRPCtestClient, user entities.User, user1 entities.User) ([]entities.Ad, error) {

	ads := make([]entities.Ad, 0)
	var empty []entities.Ad

	ad, err := addAd(client, title, text, user.ID)
	if err != nil {
		return empty, err
	}
	ads = append(ads, ad)

	ad1, err1 := addAd(client, text, title, user1.ID)
	if err1 != nil {
		return empty, err1
	}
	ads = append(ads, ad1)

	return ads, nil
}

func setupUpdateAd(client *gRPCtestClient, sChange *grpc2.ChangeAdStatusRequest) *entities.Ad {
	server := client.Server

	responseAd, _ := server.UpdateAdStatus(context.Background(), sChange)

	newAd := &entities.Ad{
		ID:         responseAd.Id,
		Title:      responseAd.Title,
		Text:       responseAd.Text,
		AuthorID:   responseAd.AuthorId,
		CreateDate: responseAd.CreateDate.AsTime().UTC(),
		UpdateDate: responseAd.UpdateDate.AsTime().UTC(),
		Published:  responseAd.Published,
	}
	return newAd
}

func addUser(client *gRPCtestClient, nickname string, email string) (entities.User, error) {
	server := client.Server

	cUserReq := &grpc2.UserRequest{Nickname: nickname, Email: email}
	res, err := server.AddUser(context.Background(), cUserReq)
	if err != nil {
		empty := entities.User{}
		return empty, err
	}
	user := entities.User{ID: res.Id, Nickname: res.Nickname, Email: res.Email}
	return user, nil
}

func addAd(client *gRPCtestClient, title string, text string, userId int64) (entities.Ad, error) {
	server := client.Server

	adReq := &grpc2.CreateAdRequest{
		Title:  title,
		Text:   text,
		UserId: userId,
	}
	ad, err := server.AddAd(context.Background(), adReq)
	if err != nil {
		empty := entities.Ad{}
		return empty, err
	}
	formatter := util.NewDateTimeFormatter(time.RFC3339)
	cDate, _ := formatter.ToTime(ad.CreateDate.AsTime().UTC())
	uDate, _ := formatter.ToTime(ad.UpdateDate.AsTime().UTC())

	newAd := entities.Ad{
		ID:         ad.Id,
		Title:      ad.Title,
		AuthorID:   ad.AuthorId,
		Text:       ad.Text,
		CreateDate: cDate,
		UpdateDate: uDate,
	}
	return newAd, nil
}
