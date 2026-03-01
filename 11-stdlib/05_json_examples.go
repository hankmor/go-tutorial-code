package main

import (
	"encoding/json"
	"fmt"
)

// JSON 处理示例

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"` // 空值时省略
	Age   int    `json:"-"`               // 忽略该字段
}

func main() {
	fmt.Println("=== 示例 1: 结构体转 JSON ===")

	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}

	// 序列化
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(data))

	fmt.Println("\n=== 示例 2: JSON 转结构体 ===")

	jsonStr := `{"id":1,"name":"Alice","email":"alice@example.com"}`

	var user2 User
	err = json.Unmarshal([]byte(jsonStr), &user2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("%+v\n", user2)

	fmt.Println("\n=== 示例 3: 格式化输出 ===")

	user3 := User{ID: 1, Name: "Alice"}

	// 格式化输出（带缩进）
	data, _ = json.MarshalIndent(user3, "", "  ")
	fmt.Println(string(data))

	fmt.Println("\n=== 示例 4: 处理 map 和 slice ===")

	// map 转 JSON
	m := map[string]interface{}{
		"name": "Go",
		"age":  15,
	}
	data, _ = json.Marshal(m)
	fmt.Println("Map:", string(data))

	// slice 转 JSON
	s := []string{"a", "b", "c"}
	data, _ = json.Marshal(s)
	fmt.Println("Slice:", string(data))
}
