package main

import (
	"fmt"
	"sync"
)

// sync.Once 示例：确保只执行一次

var (
	instance *Singleton
	once     sync.Once
)

type Singleton struct {
	data string
}

func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Println("Creating singleton instance")
		instance = &Singleton{data: "singleton"}
	})
	return instance
}

// 配置加载示例
type Config struct {
	Host string
	Port int
}

var (
	config     *Config
	configOnce sync.Once
)

func LoadConfig() *Config {
	configOnce.Do(func() {
		fmt.Println("Loading config...")
		config = &Config{
			Host: "localhost",
			Port: 8080,
		}
	})
	return config
}

func main() {
	fmt.Println("=== 示例 1: 单例模式 ===")
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			s := GetInstance()
			fmt.Printf("Goroutine %d got instance: %p\n", id, s)
		}(i)
	}

	wg.Wait()

	fmt.Println("\n=== 示例 2: 配置加载 ===")
	var wg2 sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			cfg := LoadConfig()
			fmt.Printf("Goroutine %d: %s:%d\n", id, cfg.Host, cfg.Port)
		}(i)
	}

	wg2.Wait()
}
