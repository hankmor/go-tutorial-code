package main

import (
	"fmt"
	"time"
)

// time 包示例

func main() {
	fmt.Println("=== 示例 1: 获取时间 ===")

	// 当前时间
	now := time.Now()
	fmt.Println("当前时间:", now)

	// 指定时间
	t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	fmt.Println("指定时间:", t)

	// Unix 时间戳
	timestamp := time.Now().Unix()
	fmt.Println("时间戳:", timestamp)

	// 从时间戳创建时间
	t = time.Unix(timestamp, 0)
	fmt.Println("从时间戳:", t)

	fmt.Println("\n=== 示例 2: 时间格式化 ===")

	// Go 的格式化模板：2006-01-02 15:04:05
	fmt.Println("日期:", now.Format("2006-01-02"))
	fmt.Println("日期时间:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("时间:", now.Format("15:04:05"))
	fmt.Println("自定义:", now.Format("2006/01/02"))

	// 解析时间字符串
	t, err := time.Parse("2006-01-02", "2024-01-12")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("解析时间:", t)

	fmt.Println("\n=== 示例 3: 时间计算 ===")

	// 加减时间
	later := now.Add(2 * time.Hour)
	earlier := now.Add(-30 * time.Minute)

	fmt.Println("2小时后:", later.Format("15:04:05"))
	fmt.Println("30分钟前:", earlier.Format("15:04:05"))

	// 时间差
	duration := later.Sub(now)
	fmt.Println("时间差:", duration)

	// 比较时间
	if later.After(now) {
		fmt.Println("later is after now")
	}

	if earlier.Before(now) {
		fmt.Println("earlier is before now")
	}

	fmt.Println("\n=== 示例 4: 时间间隔 ===")

	fmt.Println("1秒:", time.Second)
	fmt.Println("1分钟:", time.Minute)
	fmt.Println("1小时:", time.Hour)

	// 自定义间隔
	duration = 2*time.Hour + 30*time.Minute
	fmt.Println("自定义间隔:", duration)

	fmt.Println("\n=== 示例 5: 睡眠 ===")
	fmt.Println("睡眠1秒...")
	time.Sleep(time.Second)
	fmt.Println("醒来了！")
}
