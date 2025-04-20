package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sugerara/go-age-practice/internal/database"
	"github.com/sugerara/go-age-practice/internal/handlers"
)

func main() {
	// デバッグモードを有効化
	gin.SetMode(gin.DebugMode)

	log.Printf("データベース接続を開始します...")
	// データベース接続の初期化（設定を更新）
	err := database.InitDB("localhost", 5432, "age_user", "age_password", "age_db")
	if err != nil {
		log.Fatalf("データベース初期化エラー: %v", err)
	}
	log.Printf("データベース接続が成功しました")

	// Ginルーターの設定
	r := gin.Default()

	// APIエンドポイントの定義
	r.GET("/api/countries", handlers.GetAllCountries)
	r.GET("/api/countries/:id/capital", handlers.GetCountryCapital)
	r.GET("/api/country-capitals", handlers.GetCountriesWithCapitals)

	log.Printf("サーバーを起動します（ポート:8080）...")
	// サーバーの起動
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}
