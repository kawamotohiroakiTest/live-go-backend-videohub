package services

import (
	"context"
	pb "live/videohub/pb" // プロトバッファのパッケージ
)

type CommentsService struct {
	pb.UnimplementedCommentsServiceServer
}

func NewCommentsService() *CommentsService {
	return &CommentsService{}
}

// GetComments メソッドの実装
func (s *CommentsService) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	// コメントの取得処理をここに実装
	return &pb.GetCommentsResponse{
		Comments: []*pb.Comment{
			{Id: 1, UserId: 123, VideoId: req.VideoId, Content: "This is a comment"}, // 仮のコメント
		},
	}, nil
}
