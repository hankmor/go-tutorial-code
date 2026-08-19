package main

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// TestNestedPreload 验证显式嵌套预加载能够读取完整关联图。
func TestNestedPreload(t *testing.T) {
	db := newTestDB(t)
	author, err := SeedGraph(db)
	if err != nil {
		t.Fatalf("seed graph: %v", err)
	}

	loaded, err := LoadAuthorGraph(db, author.ID)
	if err != nil {
		t.Fatalf("load graph: %v", err)
	}
	if len(loaded.Posts) != 2 {
		t.Fatalf("posts = %d, want 2", len(loaded.Posts))
	}
	if len(loaded.Posts[0].Tags) != 1 || len(loaded.Posts[1].Tags) != 2 {
		t.Fatalf("unexpected tag counts: %d, %d", len(loaded.Posts[0].Tags), len(loaded.Posts[1].Tags))
	}
}

// TestAssociationReplace 验证 Replace 只替换指定文章的标签关联。
func TestAssociationReplace(t *testing.T) {
	db := newTestDB(t)
	author, err := SeedGraph(db)
	if err != nil {
		t.Fatalf("seed graph: %v", err)
	}

	var post Post
	if err := db.Where("author_id = ?", author.ID).Order("id ASC").First(&post).Error; err != nil {
		t.Fatalf("find post: %v", err)
	}
	var tag Tag
	if err := db.Where("name = ?", "database").First(&tag).Error; err != nil {
		t.Fatalf("find tag: %v", err)
	}
	if err := ReplacePostTags(db, post.ID, &tag); err != nil {
		t.Fatalf("replace tags: %v", err)
	}

	var loaded Post
	if err := db.Preload("Tags").First(&loaded, post.ID).Error; err != nil {
		t.Fatalf("reload post: %v", err)
	}
	if len(loaded.Tags) != 1 || loaded.Tags[0].Name != "database" {
		t.Fatalf("tags = %#v, want database only", loaded.Tags)
	}
}

// TestTransactionRollsBack 验证事务失败后不会保留部分写入。
func TestTransactionRollsBack(t *testing.T) {
	db := newTestDB(t)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Author{Name: "Rolled back"}).Error; err != nil {
			return err
		}
		return errors.New("stop transaction")
	})
	if err == nil {
		t.Fatal("transaction unexpectedly succeeded")
	}

	var count int64
	if err := db.Model(&Author{}).Where("name = ?", "Rolled back").Count(&count).Error; err != nil {
		t.Fatalf("count authors: %v", err)
	}
	if count != 0 {
		t.Fatalf("rolled back author count = %d, want 0", count)
	}
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := NewDB()
	if err != nil {
		t.Fatalf("new database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql database: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}
