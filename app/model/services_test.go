package model

import (
	"context"
	"strings"
	"testing"
)

// Test query indexes

// TestQueryGetServiceByNameAndUserUUID is used to test whether the index is used to query the service by it's name and user uuid
func TestQueryGetServiceByNameAndUserUUID(t *testing.T) {
	setupTest()

	rows, err := db.NamedExplainQuery(context.Background(), appendExplain(queryCheckServiceByNameAndUserUUID), map[string]interface{}{
		"name":      "backend",
		"user_uuid": userUUID,
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

// TestQueryGetServiceByIDAndUserUUID is used to test whether the index is used to query the service by it's service id and user uuid
func TestQueryGetServiceByIDAndUserUUID(t *testing.T) {
	setupTest()

	rows, err := db.NamedExplainQuery(context.Background(), appendExplain(queryCheckServiceByIDAndUserUUID), map[string]interface{}{
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

// TestQueryGetService is used to test whether the index is used to query the service and it's version by it's service id and user uuid
func TestQueryGetService(t *testing.T) {
	setupTest()

	rows, err := db.NamedExplainQuery(context.Background(), appendExplain(queryGetService), map[string]interface{}{
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
