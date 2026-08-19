package main

import (
	"testing"
	"time"
)

// TestCreateAuthorsInBatches verifies batched writes persist every record.
func TestCreateAuthorsInBatches(t *testing.T) {
	db := newTestDB(t)
	authors := []Author{{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"}}
	if err := CreateAuthorsInBatches(db, authors, 2); err != nil {
		t.Fatalf("create authors in batches: %v", err)
	}
	var count int64
	if err := db.Model(&Author{}).Count(&count).Error; err != nil {
		t.Fatalf("count authors: %v", err)
	}
	if count != int64(len(authors)) {
		t.Fatalf("author count = %d, want %d", count, len(authors))
	}
}

// TestPreloadUsesBatchAssociationQuery verifies the preload path returns associations.
func TestPreloadUsesBatchAssociationQuery(t *testing.T) {
	db := newTestDB(t)
	if _, err := SeedGraph(db); err != nil {
		t.Fatalf("seed graph: %v", err)
	}
	plugin := &QueryMetricsPlugin{}
	if err := InstallQueryMetrics(db, plugin); err != nil {
		t.Fatalf("install query metrics: %v", err)
	}

	authors, err := FindAuthorsWithPosts(db)
	if err != nil {
		t.Fatalf("find authors with posts: %v", err)
	}
	if len(authors) != 1 || len(authors[0].Posts) != 2 {
		t.Fatalf("authors = %#v", authors)
	}
	if plugin.Queries.Load() < 2 {
		t.Fatalf("query count = %d, want at least main and preload queries", plugin.Queries.Load())
	}
}

// TestConfigureConnectionPool verifies the configured database/sql pool values.
func TestConfigureConnectionPool(t *testing.T) {
	db := newTestDB(t)
	sqlDB, err := ConfigureConnectionPool(db, 8, 4, time.Minute, 30*time.Second)
	if err != nil {
		t.Fatalf("configure pool: %v", err)
	}
	stats := sqlDB.Stats()
	if stats.MaxOpenConnections != 8 {
		t.Fatalf("max open connections = %d, want 8", stats.MaxOpenConnections)
	}
}
