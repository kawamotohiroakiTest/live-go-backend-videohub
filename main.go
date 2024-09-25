package main

import (
	"flag"
	"fmt"
	"live/ai"
	"live/common"
	"live/videohub"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	common.LogTodo(common.INFO, "Starting server on port!!: "+port)
	if err := http.ListenAndServe(":"+port, common.EnableCors(r)); err != nil {
		common.LogError(fmt.Errorf("Error starting server: %v", err))
	}
}
