package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"gorm.io/gorm"
)

// Product 是库存商品模型，用于演示事务中的条件扣库存。
type Product struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Stock int `gorm:"not null"`
}

// Purchase 是商品购买记录模型。
type Purchase struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint
	Quantity  int
}

// HookUser 是演示 GORM 创建钩子的用户模型。
type HookUser struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
}

// BeforeCreate 在创建用户前将演示密码转换为哈希值。
// 真实项目应使用专门的密码哈希算法和参数配置，不应直接使用 SHA-256 存储密码。
func (u *HookUser) BeforeCreate(_ *gorm.DB) error {
	if u.PasswordHash == "" {
		return errors.New("password is required")
	}
	if len(u.PasswordHash) == 64 {
		return nil
	}
	sum := sha256.Sum256([]byte(u.PasswordHash))
	u.PasswordHash = hex.EncodeToString(sum[:])
	return nil
}

// CreatePurchase 在一个事务中扣减库存并创建购买记录。
func CreatePurchase(db *gorm.DB, productID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&Product{}).
			Where("id = ? AND stock >= ?", productID, quantity).
			UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Create(&Purchase{ProductID: productID, Quantity: quantity}).Error
	})
}

// ActiveAuthors 返回只包含启用作者条件的查询作用域。
func ActiveAuthors(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", "active")
}

// AdultAuthors 返回年龄不小于 18 岁的查询作用域。
func AdultAuthors(db *gorm.DB) *gorm.DB {
	return db.Where("age >= ?", 18)
}
