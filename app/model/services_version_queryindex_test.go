package model

import (
	"context"
	"strings"
	"testing"
)

// TestQueryCheckServiceVersionUsingVersion is used to test whether the index is used to query the service version by it's version and service id
func TestQueryCheckServiceVersionUsingVersion(t *testing.T) {
	setupTest()

	rows, err := db.NamedExplainQuery(context.Background(), appendExplain(queryCheckServiceVersionUsingVersion), map[string]interface{}{
		"version":    "v1.0.0",
		"service_id": 1,
		"user_uuid":  userUUID,
	})
	if err != nil {
		t.Fatal("Failed to execute query:", err)
	}
	defer rows.Close()

	// Analyze the query execution plan
	var plan string
	var indexUsed bool
	for rows.Next() {
		if err := rows.Scan(&plan); err != nil {
			t.Fatal("Failed to scan row:", err)
		}

		// Check if the index is being used
		if strings.Contains(plan, "Index Scan") {
			log.Info("Index scan being used")
			indexUsed = true
			break
		}
	}

	if !indexUsed {
		t.Error("Expected index scan but index is not being used")
	}

	if err := rows.Err(); err != nil {
		t.Fatal("Error iterating over rows:", err)
	}
}

// TestQueryCheckServiceVersionUsingSVIR is used to test whether the index is used to query the service version by it's service version id and service id
func TestQueryCheckServiceVersionUsingSVIR(t *testing.T) {
	setupTest()

	rows, err := db.NamedExplainQuery(context.Background(), appendExplain(queryCheckServiceVersionUsingSVID), map[string]interface{}{
		"sv_id":      1,
		"service_id": 1,
		"user_uuid":  userUUID,
	})
	if err != nil {
		t.Fatal("Failed to execute query:", err)
	}
	defer rows.Close()

	// Analyze the query execution plan
	var plan string
	var indexUsed bool
	for rows.Next() {
		if err := rows.Scan(&plan); err != nil {
			t.Fatal("Failed to scan row:", err)
		}

		// Check if the index is being used
		if strings.Contains(plan, "Index Scan") {
			log.Info("Index scan being used")
			indexUsed = true
			break
		}
	}

	if !indexUsed {
		t.Error("Expected index scan but index is not being used")
	}

	if err := rows.Err(); err != nil {
		t.Fatal("Error iterating over rows:", err)
	}
}
