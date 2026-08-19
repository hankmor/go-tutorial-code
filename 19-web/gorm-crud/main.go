package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 是数据库中的用户模型，DeletedAt 用于实现软删除。
type User struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:100;not null"`
	Age       int            `gorm:"not null;default:18"`
	Email     string         `gorm:"size:255;uniqueIndex;not null"`
}

// NewDB 创建使用内存 SQLite 的 GORM 数据库连接。
func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, fmt.Errorf("migrate user table: %w", err)
	}
	return db, nil
}

// CreateUser 创建用户，并将数据库生成的主键写回 user。
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

// FindUser 按主键查询未被软删除的用户。
func FindUser(db *gorm.DB, id uint) (User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

// UpdateUserAge 更新用户年龄，包括将年龄设置为零的情况。
func UpdateUserAge(db *gorm.DB, id uint, age int) error {
	return db.Model(&User{}).Where("id = ?", id).Update("age", age).Error
}

// DeleteUser 软删除用户，后续普通查询不会返回该用户。
func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}

	user := &User{Name: "Hank", Age: 25, Email: "hank@example.com"}
	if err := CreateUser(db, user); err != nil {
		log.Fatal(err)
	}
	found, err := FindUser(db, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created: id=%d name=%s age=%d\n", found.ID, found.Name, found.Age)

	if err := UpdateUserAge(db, user.ID, 0); err != nil {
		log.Fatal(err)
	}
	if err := DeleteUser(db, user.ID); err != nil {
		log.Fatal(err)
	}
	if _, err := FindUser(db, user.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("find deleted user: %v", err)
	}
	fmt.Println("user deleted softly")
}
