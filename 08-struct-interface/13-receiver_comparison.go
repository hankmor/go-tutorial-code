package main

import (
	"fmt"
)

type User struct {
	Name string
	Age  int
}

// 值接收者 (Value Receiver)
// 会复制一份 User 结构体，对它的修改不会影响原始数据
func (u User) UpdateNameValue(newName string) {
	fmt.Printf("[值接收者] 修改前: %v (地址: %p)\n", u, &u)
	u.Name = newName
	fmt.Printf("[值接收者] 修改后: %v (地址: %p)\n", u, &u)
}

// 指针接收者 (Pointer Receiver)
// 修改会影响原始数据，且避免了大结构体的内存拷贝
func (u *User) UpdateNamePointer(newName string) {
	fmt.Printf("[指针接收者] 修改前: %v (地址: %p)\n", *u, u)
	u.Name = newName
	fmt.Printf("[指针接收者] 修改后: %v (地址: %p)\n", *u, u)
}

func main() {
	u := User{Name: "老墨", Age: 18}
	fmt.Printf("--- 原始数据: %v (地址: %p) ---\n\n", u, &u)

	fmt.Println("1. 调用值接收者方法:")
	u.UpdateNameValue("极客老墨-值")
	fmt.Printf("调用后原始数据: %v\n\n", u)

	fmt.Println("2. 调用指针接收者方法:")
	u.UpdateNamePointer("极客老墨-指针")
	fmt.Printf("调用后原始数据: %v\n\n", u)

	fmt.Println("避坑总结: 如果需要修改结构体成员，或者结构体很大，务必使用指针接收者。")
}
