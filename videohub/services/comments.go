package services

import (
	"context"
	"live/videohub/models" // モデルのパッケージ
	pb "live/videohub/pb"  // プロトバッファのパッケージ
	"log"

	"gorm.io/gorm"
)

type CommentsService struct {
	pb.UnimplementedCommentsServiceServer
	DB *gorm.DB // DB接続
}

// NewCommentsService はDB接続を受け取り、サービスを初期化する
func NewCommentsService(db *gorm.DB) *CommentsService {
	return &CommentsService{
		DB: db,
	}
}

// GetComments メソッドの実装
func (s *CommentsService) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	var comments []models.Comment

	// DBからvideo_idに基づいてコメントを取得
	if err := s.DB.Where("video_id = ?", req.VideoId).Find(&comments).Error; err != nil {
		log.Printf("Error retrieving comments: %v", err)
		return nil, err
	}

	// コメントをプロトコルバッファ形式に変換
	var pbComments []*pb.Comment
	for _, comment := range comments {
		pbComments = append(pbComments, &pb.Comment{
			Id:      int64(comment.ID),
			UserId:  int64(comment.UserID),
			VideoId: int64(comment.VideoID),
			Content: comment.Content,
		})
	}

	// レスポンスを返す
	return &pb.GetCommentsResponse{
		Comments: pbComments,
	}, nil
}
