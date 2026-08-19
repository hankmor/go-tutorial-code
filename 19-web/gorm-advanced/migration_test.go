package main

import (
	"testing"
)

// TestJSONMapRoundTrip 验证自定义 JSON 类型可以写入并读回数据库。
func TestJSONMapRoundTrip(t *testing.T) {
	db := newTestDB(t)
	record := MetadataRecord{
		Name:     "Ada",
		Metadata: JSONMap{"role": "admin", "active": true},
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create metadata record: %v", err)
	}

	var loaded MetadataRecord
	if err := db.First(&loaded, record.ID).Error; err != nil {
		t.Fatalf("load metadata record: %v", err)
	}
	if loaded.Metadata["role"] != "admin" || loaded.Metadata["active"] != true {
		t.Fatalf("metadata = %#v", loaded.Metadata)
	}
}

// TestQueryMetricsPlugin verifies that the plugin is installed through GORM callbacks.
func TestQueryMetricsPlugin(t *testing.T) {
	db := newTestDB(t)
	plugin := &QueryMetricsPlugin{}
	if err := InstallQueryMetrics(db, plugin); err != nil {
		t.Fatalf("install plugin: %v", err)
	}

	var authors []Author
	if err := db.Find(&authors).Error; err != nil {
		t.Fatalf("query authors: %v", err)
	}
	if plugin.Queries.Load() == 0 {
		t.Fatal("query callback was not invoked")
	}
}
