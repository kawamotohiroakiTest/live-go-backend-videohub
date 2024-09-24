package main

import (
	"encoding/json"
	"fmt"
	"live/common"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// 動画データのサンプル
var videos = []map[string]string{
	{"id": "1", "title": "動画1", "description": "これは動画1です"},
	{"id": "2", "title": "動画2", "description": "これは動画2です"},
	{"id": "3", "title": "これは動画3です", "description": "これは動画3です"},
}

// 動画一覧を返すハンドラ
func listVideosHandler(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(videos); err != nil {
		http.Error(w, "Failed to encode videos to JSON", http.StatusInternalServerError)
		return
	}
}

// メイン関数
func main() {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// ルーターを作成し、CORSミドルウェアを登録
	r := mux.NewRouter()
	r.Use(common.EnableCors)
	fmt.Println("CORS ミドルウェアが登録されました")

	// 動画一覧のエンドポイントを設定
	r.HandleFunc("/api/v1/videos/list", listVideosHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origin := os.Getenv("API_ALLOWED_ORIGIN")
	fmt.Println("許可するオリジン: " + origin)

	fmt.Println("サーバーをポート" + port + "で開始します")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Println(err)
	}
}
