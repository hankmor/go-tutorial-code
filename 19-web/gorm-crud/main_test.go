package main

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestCRUDLifecycle(t *testing.T) {
	db := newTestDB(t)
	user := &User{Name: "Alice", Age: 20, Email: "alice@example.com"}

	if err := CreateUser(db, user); err != nil {
		t.Fatalf("create: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("create did not populate primary key")
	}

	found, err := FindUser(db, user.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if found.Name != "Alice" {
		t.Fatalf("name = %q, want Alice", found.Name)
	}

	if err := UpdateUserAge(db, user.ID, 0); err != nil {
		t.Fatalf("update zero value: %v", err)
	}
	found, err = FindUser(db, user.ID)
	if err != nil {
		t.Fatalf("find after update: %v", err)
	}
	if found.Age != 0 {
		t.Fatalf("age = %d, want 0", found.Age)
	}

	if err := DeleteUser(db, user.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := FindUser(db, user.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("find deleted user error = %v, want gorm.ErrRecordNotFound", err)
	}
}

func TestCreateRejectsDuplicateEmail(t *testing.T) {
	db := newTestDB(t)
	first := &User{Name: "Alice", Email: "same@example.com"}
	second := &User{Name: "Bob", Email: "same@example.com"}
	if err := CreateUser(db, first); err != nil {
		t.Fatalf("create first: %v", err)
	}
	if err := CreateUser(db, second); err == nil {
		t.Fatal("duplicate email was accepted")
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
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	return db
}
