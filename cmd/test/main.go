package main

import (
	"fmt"
	"reflect"
)

func main() {
	fn4()
}

func fn4() {
	a, b := new(int), new(int)
	if _, err := fmt.Scanf("%d>>%d", a, b); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	fmt.Println(*a >> *b)
}

func fn3() {
	fmt.Println("apa itu bufio?")
}

func fn2() {
	type user struct {
		Name string `json:"name" validate:"name"`
	}

	t := reflect.TypeFor[user]()
	for f := range t.Fields() {
		if lookup, ok := f.Tag.Lookup("validate"); ok {
			fmt.Printf(
				"f\t: %+v\n"+
					"get\t: %+v\n"+
					"lookup\t: %+v\n",
				f, f.Tag.Get("validate"), lookup,
			)
		}
	}
}

func fn1() {
	type User struct {
		Name string
		Age  uint8
	}
	var user *User

	fmt.Println("user:", user.Name)
}
