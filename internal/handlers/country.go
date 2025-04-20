package handlers

import (
	"net/http"

	"github.com/apache/age/drivers/golang/age"
	"github.com/gin-gonic/gin"
	"github.com/sugerara/go-age-practice/internal/database"
	"github.com/sugerara/go-age-practice/internal/models"
)

// GetAllCountries はすべての国を取得します
func GetAllCountries(c *gin.Context) {
	tx, err := database.GetClient().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	cursor, err := tx.ExecCypher(1, `
		MATCH (c:Country)
		RETURN c
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var countries []models.Country
	for cursor.Next() {
		row, err := cursor.GetRow()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		country := row[0].(*age.Vertex)
		props := country.Props()
		countryName := props["name"].(string)

		countries = append(countries, models.Country{
			Name: countryName,
		})
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, countries)
}

// GetCountryCapital は特定の国の首都を取得します
func GetCountryCapital(c *gin.Context) {
	countryName := c.Param("id")

	tx, err := database.GetClient().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベース接続エラー: " + err.Error()})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// クエリを修正: パラメータの渡し方を改善し、より具体的な条件を設定
	cursor, err := tx.ExecCypher(2, `
		MATCH (c:Country {name: %s})-[:HAS_CAPITAL]->(cap:Capital)
		RETURN c.name as country, cap.name as capital
	`, countryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "クエリ実行エラー: " + err.Error()})
		return
	}

	if !cursor.Next() {
		c.JSON(http.StatusNotFound, gin.H{"error": "指定された国または首都が見つかりません: " + countryName})
		return
	}

	row, err := cursor.GetRow()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データ取得エラー: " + err.Error()})
		return
	}

	// 結果の構造を単純化
	relation := models.CountryCapitalRelation{
		Country: models.Country{
			Name: row[0].String(),
		},
		Capital: models.Capital{
			Name: row[1].String(),
		},
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "コミットエラー: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, relation)
}

// GetCountriesWithCapitals は全ての国と首都の関係を返すハンドラーです
func GetCountriesWithCapitals(c *gin.Context) {
	relations, err := database.GetCountryCapitals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, relations)
}
