package utils

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	. "github.com/GregoryUnderscore/Mining-Automation-Shared/models"
)

// Generate a URL for a pool.
// @param tx - The active database transaction/connection
// @param algoID - The database ID for the algorithm to use for pool lookup
func generatePoolURL(tx *gorm.DB, algoID uint64) string {
	// Get a pool for the algorithm.
	var pool Pool
	tx.Where("algorithm_id = ?", algoID).Limit(1).Find(&pool)
	if (Pool{}) == pool {
		var algo Algorithm
		tx.Where("id = ?", algoID).Find(&algo)
		log.Fatalf("No pool found for this algorithm: " + algo.Name)
	}
	// Generate the URL. Can use any pool that supports the algorithm.
	url := "stratum+tcp://" + pool.URL + ":" + fmt.Sprint(pool.Port)
	return url
}
