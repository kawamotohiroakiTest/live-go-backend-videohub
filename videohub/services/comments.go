package services

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"live/videohub/models"
	pb "live/videohub/services/comments_proto" // protoファイルの生成物をインポート
)

// CommentsServiceServer は gRPC サービスのサーバーを表します
type CommentsServiceServer struct {
	pb.UnimplementedCommentsServiceServer
	DB *gorm.DB
}

// 新しいコメントを投稿する
func (s *CommentsServiceServer) PostComment(ctx context.Context, req *pb.PostCommentRequest) (*pb.PostCommentResponse, error) {
	// コメントデータの生成
	comment := models.Comment{
		UserID:   uint64(req.UserId),
		VideoID:  uint64(req.VideoId),
		Content:  req.Content,
		Created:  time.Now(),
		Modified: time.Now(),
	}

	// データベースに保存
	if err := s.DB.Create(&comment).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "コメントの作成に失敗しました: %v", err)
	}

	return &pb.PostCommentResponse{
		Success: true,
		Message: "コメントが投稿されました",
	}, nil
}

// 特定の動画に関連するコメントを取得する
func (s *CommentsServiceServer) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	var comments []models.Comment

	// データベースから該当動画IDのコメントを取得
	if err := s.DB.Where("video_id = ?", req.VideoId).Find(&comments).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "コメントの取得に失敗しました: %v", err)
	}

	// gRPC レスポンス用のコメントリストに変換
	var grpcComments []*pb.Comment
	for _, c := range comments {
		grpcComments = append(grpcComments, &pb.Comment{
			Id:         int64(c.ID),
			UserId:     int64(c.UserID),
			VideoId:    int64(c.VideoID),
			Content:    c.Content,
			CreatedAt:  c.Created.Format(time.RFC3339),
			ModifiedAt: c.Modified.Format(time.RFC3339),
		})
	}

	return &pb.GetCommentsResponse{
		Comments: grpcComments,
	}, nil
}

// サーバーの初期化関数
func NewCommentsServiceServer(db *gorm.DB) *CommentsServiceServer {
	return &CommentsServiceServer{
		DB: db,
	}
}
