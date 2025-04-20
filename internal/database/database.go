package database

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/apache/age/drivers/golang/age"
	_ "github.com/lib/pq" // PostgreSQLドライバー
	"github.com/sugerara/go-age-practice/internal/models"
)

var ag *age.Age

func InitDB(host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Printf("AGEに接続を試みます: %s:%d", host, port)
	var err error
	ag, err = age.ConnectAge("world_graph", dsn)
	if err != nil {
		return fmt.Errorf("AGE接続エラー: %v", err)
	}
	log.Printf("AGE接続が成功しました")

	// データ構造を確認するためのテストクエリを実行
	tx, err := ag.Begin()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	cursor, err := tx.ExecCypher(1, `
		MATCH (c:Country)
		RETURN c
		LIMIT 1
	`)
	if err != nil {
		return fmt.Errorf("テストクエリ実行エラー: %v", err)
	}

	if cursor.Next() {
		row, err := cursor.GetRow()
		if err != nil {
			return fmt.Errorf("テストデータ取得エラー: %v", err)
		}

		country := row[0].(*age.Vertex)
		props := country.Props()
		propsJSON, _ := json.MarshalIndent(props, "", "  ")
		log.Printf("Country vertex properties: %s", string(propsJSON))
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("テストトランザクションのコミットエラー: %v", err)
	}

	return nil
}

func GetCountryCapitals() ([]models.CountryCapitalRelation, error) {
	tx, err := ag.Begin()
	if err != nil {
		return nil, fmt.Errorf("トランザクション開始エラー: %v", err)
	}
	defer func() {
		if err != nil {
			log.Printf("トランザクションをロールバックします")
			tx.Rollback()
		}
	}()

	log.Printf("Cypherクエリを実行します: MATCH (c:Country)-[r:HAS_CAPITAL]->(cap:Capital)")
	cursor, err := tx.ExecCypher(2, `
		MATCH (c:Country)-[r:HAS_CAPITAL]->(cap:Capital)
		RETURN c.name AS country_name, cap.name AS capital_name
	`)
	if err != nil {
		return nil, fmt.Errorf("クエリ実行エラー: %v", err)
	}

	var relations []models.CountryCapitalRelation
	for cursor.Next() {
		row, err := cursor.GetRow()
		if err != nil {
			return nil, fmt.Errorf("行データ取得エラー: %v", err)
		}

		// 行データの型とデータ内容をログに出力
		log.Printf("Debug - row[0] type: %T, value: %v", row[0], row[0])
		log.Printf("Debug - row[1] type: %T, value: %v", row[1], row[1])

		country := row[0].(*age.Vertex)
		capital := row[1].(*age.Vertex)

		countryName := country.Prop("name").(string)
		capitalName := capital.Prop("name").(string)

		log.Printf("取得したデータ: 国=%s, 首都=%s", countryName, capitalName)

		relation := models.CountryCapitalRelation{
			Country: models.Country{
				Name: countryName,
			},
			Capital: models.Capital{
				Name: capitalName,
			},
		}
		relations = append(relations, relation)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("トランザクションのコミットエラー: %v", err)
	}

	log.Printf("合計%d件のデータを取得しました", len(relations))
	return relations, nil
}

// GetClient は互換性のために残していますが、新しい実装では使用しません
func GetClient() *age.Age {
	return ag
}
