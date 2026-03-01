package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// 完整示例：综合使用标准库

// 用户结构
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// 格式化用户信息
func (u User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s, Created: %s}",
		u.ID, u.Name, u.Email, u.CreatedAt.Format("2006-01-02"))
}

// 保存用户到文件
func saveUser(user User, filename string) error {
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

// 从文件加载用户
func loadUser(filename string) (*User, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &user, nil
}

// 处理用户名
func processName(name string) string {
	// 去除空格并转为标题格式
	name = strings.TrimSpace(name)
	name = strings.Title(strings.ToLower(name))
	return name
}

func main() {
	// 创建用户
	user := User{
		ID:        1,
		Name:      processName("  alice smith  "),
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
	}

	// 打印用户信息
	fmt.Println(user)

	// 保存到文件
	filename := "user.json"
	err := saveUser(user, filename)
	if err != nil {
		fmt.Println("Save error:", err)
		return
	}
	fmt.Println("User saved to", filename)

	// 从文件加载
	loadedUser, err := loadUser(filename)
	if err != nil {
		fmt.Println("Load error:", err)
		return
	}
	fmt.Println("Loaded user:", loadedUser)

	// 计算账户年龄
	age := time.Since(loadedUser.CreatedAt)
	fmt.Printf("Account age: %v\n", age.Round(time.Second))

	// 清理
	os.Remove(filename)
	fmt.Println("Cleaned up")
}
