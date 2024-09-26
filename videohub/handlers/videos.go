package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"live/common"
	"live/videohub/models"
	"live/videohub/services"
	"net/http"
	"os"
	"strconv"
	"time"

	pb "live/videohub/pb"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// コメントサービス用gRPCクライアントを初期化する関数
func initCommentsClient() (pb.CommentsServiceClient, *grpc.ClientConn, error) {
	// gRPCサーバーへの接続
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // 実際のアドレスを指定
	if err != nil {
		return nil, nil, err
	}

	// コメントサービスクライアントを生成
	client := pb.NewCommentsServiceClient(conn)
	return client, conn, nil
}

func GetVideoByIDHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// URLからvideo_idを取得
	vars := mux.Vars(r)
	videoIDStr := vars["video_id"]

	// video_idをuintに変換
	videoID, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		http.Error(w, "無効な動画IDです", http.StatusBadRequest)
		return
	}

	fmt.Printf("Debug: Fetching video with ID %d\n", videoID)

	// ストレージサービスの初期化
	var storageService *services.StorageService
	envMode := os.Getenv("ENV_MODE")
	if envMode == "local" {
		storageService, err = services.InitMinioService()
	} else {
		storageService, err = services.NewStorageService()
	}

	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	// 動画情報をデータベースから取得
	video, err := models.GetVideoByID(uint(videoID))
	if err != nil {
		http.Error(w, "動画が見つかりませんでした", http.StatusNotFound)
		return
	}

	// サムネイルや動画ファイルの署名付きURLを生成
	for i, file := range video.Files {
		if file.FilePath != "" {
			video.Files[i].FilePath, err = storageService.GetVideoPresignedURL(file.FilePath)
			if err != nil {
				common.LogVideoHubError(err)
				http.Error(w, "動画URLの取得に失敗しました", http.StatusInternalServerError)
				return
			}
		}
	}

	// GRPCのコメントサービスを呼び出してコメントを取得
	// コメントサービスクライアントの初期化
	commentClient, conn, err := initCommentsClient()
	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "コメントサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// コメントの取得リクエストを作成
	commentReq := &pb.GetCommentsRequest{
		VideoId: int64(videoID),
	}

	// コメントを取得
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	commentsRes, err := commentClient.GetComments(ctx, commentReq)
	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "コメントの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	video.Comments = commentsRes.Comments // ここでコメントを追加

	// 動画情報をJSONで返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(video); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画データのエンコードに失敗しました", http.StatusInternalServerError)
	}
}
