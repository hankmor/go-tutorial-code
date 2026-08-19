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

// TestComplexQueries verifies subqueries, aggregation, joins and raw SQL results.
func TestComplexQueries(t *testing.T) {
	db := newTestDB(t)
	author, err := SeedGraph(db)
	if err != nil {
		t.Fatalf("seed graph: %v", err)
	}
	if err := db.Create(&Author{Name: "No orders"}).Error; err != nil {
		t.Fatalf("create author without orders: %v", err)
	}

	aboveAverage, err := FindOrdersAboveAverage(db)
	if err != nil {
		t.Fatalf("find above average orders: %v", err)
	}
	if len(aboveAverage) != 2 || aboveAverage[0].Amount != 80 || aboveAverage[1].Amount != 120 {
		t.Fatalf("above average orders = %#v", aboveAverage)
	}

	summaries, err := SummarizeAuthorOrders(db, 200)
	if err != nil {
		t.Fatalf("summarize orders: %v", err)
	}
	if len(summaries) != 1 || summaries[0].AuthorID != author.ID || summaries[0].Total != 230 {
		t.Fatalf("summaries = %#v", summaries)
	}

	withOrders, err := FindAuthorsWithOrders(db)
	if err != nil {
		t.Fatalf("find authors with orders: %v", err)
	}
	withoutOrders, err := FindAuthorsWithoutOrders(db)
	if err != nil {
		t.Fatalf("find authors without orders: %v", err)
	}
	if len(withOrders) != 1 || withOrders[0].ID != author.ID {
		t.Fatalf("authors with orders = %#v", withOrders)
	}
	if len(withoutOrders) != 1 || withoutOrders[0].Name != "No orders" {
		t.Fatalf("authors without orders = %#v", withoutOrders)
	}

	rawSummaries, err := RawAuthorOrderTotals(db, 200)
	if err != nil {
		t.Fatalf("raw summaries: %v", err)
	}
	if len(rawSummaries) != 1 || rawSummaries[0].Total != 230 {
		t.Fatalf("raw summaries = %#v", rawSummaries)
	}
}

// TestExecUpdatesOrderStatus verifies parameterized SQL updates and row-count errors.
func TestExecUpdatesOrderStatus(t *testing.T) {
	db := newTestDB(t)
	if _, err := SeedGraph(db); err != nil {
		t.Fatalf("seed graph: %v", err)
	}
	var order Order
	if err := db.Order("id ASC").First(&order).Error; err != nil {
		t.Fatalf("find order: %v", err)
	}
	if err := UpdateOrderStatus(db, order.ID, "cancelled"); err != nil {
		t.Fatalf("update order status: %v", err)
	}
	var updated Order
	if err := db.First(&updated, order.ID).Error; err != nil {
		t.Fatalf("reload order: %v", err)
	}
	if updated.Status != "cancelled" {
		t.Fatalf("status = %q, want cancelled", updated.Status)
	}
	if err := UpdateOrderStatus(db, 999, "cancelled"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("missing order error = %v, want gorm.ErrRecordNotFound", err)
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
