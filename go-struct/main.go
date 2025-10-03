package main

import "fmt"

type Employee struct {
	Name string
	Age  uint
}

// Nested struct
type Company struct {
	ID        int
	Name      string
	Address   string
	Employees []*Employee // slice of pointer
}

func main() {
	emp := &Employee{Name: "Khoa", Age: 21}
	c := &Company{
		ID:        2004,
		Name:      "VNG",
		Address:   "TP.HCM",
		Employees: []*Employee{emp},
	}

	fmt.Printf("%p\n", emp)            // địa chỉ Employee
	fmt.Printf("%p\n", c.Employees[0]) // cùng địa chỉ Employee

}
