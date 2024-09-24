package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// 動画データのサンプル
var videos = []map[string]string{
	{"id": "1", "title": "動画1", "description": "これは動画1です"},
	{"id": "2", "title": "動画2", "description": "これは動画2です"},
	{"id": "3", "title": "動画3", "description": "これは動画3です"},
}

// 動画一覧を返すハンドラ
func listVideosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	for _, video := range videos {
		fmt.Fprintf(w, "%v\n", video)
	}
}

// メイン関数
func main() {
	r := mux.NewRouter()

	// 動画一覧のエンドポイントを設定
	r.HandleFunc("/videos/list", listVideosHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// サーバーをポート8080で起動
	fmt.Println("サーバーをポート8002開始します!...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
