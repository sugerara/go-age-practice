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
	countryName := c.Param("id") // URLパラメータはそのままにしていますが、実際にはnameで検索します

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

	cursor, err := tx.ExecCypher(3, `
		MATCH (c:Country)-[r:HAS_CAPITAL]->(cap:Capital)
		WHERE c.name = $1
		RETURN c, cap
	`, countryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !cursor.Next() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	row, err := cursor.GetRow()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	country := row[0].(*age.Vertex)
	capital := row[1].(*age.Vertex)

	countryProps := country.Props()
	capitalProps := capital.Props()

	countryName = countryProps["name"].(string)
	capitalName := capitalProps["name"].(string)

	relation := models.CountryCapitalRelation{
		Country: models.Country{
			Name: countryName,
		},
		Capital: models.Capital{
			Name: capitalName,
		},
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
