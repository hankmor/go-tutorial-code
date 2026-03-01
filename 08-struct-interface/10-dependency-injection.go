package main

import "fmt"

// 实战建议 2：使用接口实现依赖注入

type User struct {
	ID   int
	Name string
}

// 定义接口
type UserRepository interface {
	GetByID(id int) (*User, error)
	Save(user *User) error
}

// 真实实现
type MySQLUserRepository struct{}

func (r *MySQLUserRepository) GetByID(id int) (*User, error) {
	fmt.Printf("MySQL: Getting user %d from database\n", id)
	return &User{ID: id, Name: "Alice"}, nil
}

func (r *MySQLUserRepository) Save(user *User) error {
	fmt.Printf("MySQL: Saving user %d to database\n", user.ID)
	return nil
}

// Mock 实现（用于测试）
type MockUserRepository struct{}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	fmt.Printf("Mock: Getting user %d\n", id)
	return &User{ID: id, Name: "Mock User"}, nil
}

func (m *MockUserRepository) Save(user *User) error {
	fmt.Printf("Mock: Saving user %d\n", user.ID)
	return nil
}

// 服务依赖接口，不依赖具体实现
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) UpdateUser(user *User) error {
	return s.repo.Save(user)
}

func main() {
	fmt.Println("=== 依赖注入示例 ===")

	// 生产环境：使用真实实现
	fmt.Println("\n--- Production ---")
	mysqlRepo := &MySQLUserRepository{}
	prodService := NewUserService(mysqlRepo)
	user, _ := prodService.GetUser(1)
	fmt.Printf("Got user: %+v\n", user)

	// 测试环境：使用 Mock 实现
	fmt.Println("\n--- Testing ---")
	mockRepo := &MockUserRepository{}
	testService := NewUserService(mockRepo)
	user, _ = testService.GetUser(1)
	fmt.Printf("Got user: %+v\n", user)

	// 优点：
	// 1. 易于测试：可以轻松替换实现
	// 2. 解耦：服务不依赖具体实现
	// 3. 灵活：可以在运行时切换实现
}
