package main

import (
	"fmt"
)

// ===============truyen ham vao ham khac===============
func f(z int, x int, y int, operator func(int, int) int) int {
	return z * operator(x, y)
}

// ===============Gan ham vao mot bien===============
func add(a, b int) int {
	return a + b
}

// ===============Function return another function===============
func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return factor * x
	}
}

// ===============Function return multi value===============
func divMod(a, b int) (int, int) {
	return a / b, a % b
}

type User struct {
	Name string
}

// ===============Value receiver/Pointer receiver===============
func (u User) RenameValue(newName string) {
	u.Name = newName
}

func (u *User) RenamePointer(newName string) {
	u.Name = newName
}

// ===============iota===============
const (
	A = iota + 1
	B
	C
)

func main() {
	// truyen func vao param
	var result1 int = f(C, A, B, add)
	fmt.Println(result1)

	// gan function vao variale
	var paramFunc func(int, int) int = add
	var result2 int = f(C, A, B, paramFunc)
	fmt.Println(result2)

	// function return a function
	double := makeMultiplier(2)
	fmt.Println(double(7)) // 14

	// function return multi value
	div, mod := divMod(7, 3)
	fmt.Println(div, mod)

	// Value receiver/Pointer Receiver
	u := User{Name: "Khoa"}

	u.RenameValue("Alice")
	fmt.Println(u.Name) // vẫn là "Khoa" vì chỉ đổi trên bản copy

	u.RenamePointer("Alice")
	fmt.Println(u.Name) // đổi thành "Alice" vì sửa trực tiếp trên bản gốc

	// ===============Anonymous functions===============
	count := 0
	cnt := func(val int) int {
		count += val
		return count
	}
	cnt(1)
	cnt(7)
	fmt.Println(count)
}
