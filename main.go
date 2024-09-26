package main

import (
	"flag"
	"fmt"
	"live/ai"
	"live/common"
	"live/videohub"
	pb "live/videohub/pb"
	"log"
	"net"
	"net/http"
	"os"

	"live/videohub/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// フラグのパース
	flag.Parse()

	// 環境変数の読み込み
	err := godotenv.Load()
	if err != nil {
		common.LogError(fmt.Errorf("Error loading .env file: %v", err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// データベースの初期化
	dbConn, err := common.InitDB()
	if err != nil {
		common.LogError(fmt.Errorf("Error initializing database: %v", err))
		return
	}

	// HTTP ルーターの初期化
	r := mux.NewRouter()

	// ルートの登録
	videohub.RegisterRoutes(r, dbConn)
	ai.RegisterRoutes(r)

	// gRPC サーバーのリスナーを作成
	lis, err := net.Listen("tcp", ":50051") // gRPC サーバーをポート 50051 でリッスン
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// gRPC サーバーの作成
	grpcServer := grpc.NewServer()

	// コメントサービスの初期化
	commentsService := services.NewCommentsService()

	// コメントサービスを gRPC サーバーに登録
	pb.RegisterCommentsServiceServer(grpcServer, commentsService)

	// gRPC サーバーの起動
	go func() {
		fmt.Println("Starting gRPC server on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// HTTP サーバーの起動
	if err := http.ListenAndServe(":"+port, common.EnableCors(r)); err != nil {
		common.LogError(fmt.Errorf("Error starting HTTP server: %v", err))
	}
}
