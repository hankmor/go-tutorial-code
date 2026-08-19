package main

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// TestCreatePurchaseIsAtomic verifies stock deduction and purchase creation share one transaction.
func TestCreatePurchaseIsAtomic(t *testing.T) {
	db := newTestDB(t)
	product := Product{Name: "Keyboard", Stock: 2}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}

	if err := CreatePurchase(db, product.ID, 1); err != nil {
		t.Fatalf("create purchase: %v", err)
	}
	var updated Product
	if err := db.First(&updated, product.ID).Error; err != nil {
		t.Fatalf("reload product: %v", err)
	}
	if updated.Stock != 1 {
		t.Fatalf("stock = %d, want 1", updated.Stock)
	}

	if err := CreatePurchase(db, product.ID, 2); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("insufficient stock error = %v, want gorm.ErrRecordNotFound", err)
	}
	var purchaseCount int64
	if err := db.Model(&Purchase{}).Count(&purchaseCount).Error; err != nil {
		t.Fatalf("count purchases: %v", err)
	}
	if purchaseCount != 1 {
		t.Fatalf("purchase count = %d, want 1", purchaseCount)
	}
}

// TestBeforeCreateHashesPassword verifies hook execution and hook failure.
func TestBeforeCreateHashesPassword(t *testing.T) {
	db := newTestDB(t)
	user := HookUser{Username: "ada", PasswordHash: "plain-password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create hook user: %v", err)
	}
	if user.PasswordHash == "plain-password" || len(user.PasswordHash) != 64 {
		t.Fatalf("password hash = %q, want 64-character hash", user.PasswordHash)
	}

	if err := db.Create(&HookUser{Username: "empty", PasswordHash: ""}).Error; err == nil {
		t.Fatal("empty password was accepted")
	}
}

// TestScopesComposeConditions verifies scopes are composable and parameterized.
func TestScopesComposeConditions(t *testing.T) {
	db := newTestDB(t)
	if err := db.Create(&Author{Name: "Adult active", Age: 30, Status: "active"}).Error; err != nil {
		t.Fatalf("create active author: %v", err)
	}
	if err := db.Create(&Author{Name: "Minor active", Age: 16, Status: "active"}).Error; err != nil {
		t.Fatalf("create minor author: %v", err)
	}
	if err := db.Create(&Author{Name: "Adult inactive", Age: 30, Status: "inactive"}).Error; err != nil {
		t.Fatalf("create inactive author: %v", err)
	}

	var authors []Author
	if err := db.Scopes(ActiveAuthors, AdultAuthors).Find(&authors).Error; err != nil {
		t.Fatalf("find scoped authors: %v", err)
	}
	if len(authors) != 1 || authors[0].Name != "Adult active" {
		t.Fatalf("scoped authors = %#v", authors)
	}
}
