package main

import (
	"flag"
	"fmt"
	"live/ai"
	"live/common"
	"live/videohub"
	"live/videohub/services"
	"net"
	"net/http"
	"os"

	pb "live/videohub/services/comments_proto"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// フラグのパース
	flag.Parse()

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

	r := mux.NewRouter()

	videohub.RegisterRoutes(r, dbConn)
	ai.RegisterRoutes(r)

	// gRPC サーバーのセットアップ

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			common.LogError(fmt.Errorf("Failed to listen on gRPC port: %v", err))
			return
		}

		grpcServer := grpc.NewServer()

		// コメントサービスを gRPC サーバーに登録
		commentService := services.NewCommentsServiceServer(dbConn)
		pb.RegisterCommentsServiceServer(grpcServer, commentService)

		fmt.Println("Starting gRPC server on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			common.LogError(fmt.Errorf("Failed to serve gRPC server: %v", err))
		}
	}()

	common.LogTodo(common.INFO, "Starting server on port!!: "+port)
	if err := http.ListenAndServe(":"+port, common.EnableCors(r)); err != nil {
		common.LogError(fmt.Errorf("Error starting server: %v", err))
	}
}
