package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// 实战建议 3：结构体标签（Tags）

type User struct {
	ID        int       `json:"id" db:"user_id"`
	Name      string    `json:"name" db:"user_name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Password  string    `json:"-" db:"password"` // json:"-" 表示不序列化
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"` // omitempty: 空值时不序列化
	InStock     bool    `json:"in_stock"`
}

func main() {
	fmt.Println("=== 结构体标签示例 ===")

	// JSON 序列化
	user := User{
		ID:        1,
		Name:      "Alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
		Password:  "secret123", // 不会被序列化
	}

	data, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println("\n--- User JSON ---")
	fmt.Println(string(data))
	// 注意：Password 字段不会出现在 JSON 中

	// omitempty 示例
	product1 := Product{
		ID:      1,
		Name:    "Laptop",
		Price:   999.99,
		InStock: true,
	}

	product2 := Product{
		ID:          2,
		Name:        "Mouse",
		Price:       29.99,
		Description: "Wireless mouse", // 有值
		InStock:     true,
	}

	data1, _ := json.MarshalIndent(product1, "", "  ")
	data2, _ := json.MarshalIndent(product2, "", "  ")

	fmt.Println("\n--- Product 1 (no description) ---")
	fmt.Println(string(data1))
	// description 字段不会出现

	fmt.Println("\n--- Product 2 (with description) ---")
	fmt.Println(string(data2))
	// description 字段会出现

	// 常用标签：
	// - json:"field_name" - JSON 序列化
	// - xml:"field_name" - XML 序列化
	// - db:"column_name" - 数据库映射
	// - validate:"required" - 数据验证
	// - form:"field_name" - 表单绑定
}
