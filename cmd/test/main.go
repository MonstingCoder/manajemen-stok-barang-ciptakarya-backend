package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"time"
)

func main() {
	fn6()
}

func fn6() {
	withBufio := func(loop uint) (timeTaken time.Duration) {
		now := time.Now()

		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		for i := range loop {
			fmt.Fprintln(w, "loop ke :", i)
		}

		timeTaken = time.Since(now)
		return
	}
	withoutBufio := func(loop uint) (timeTaken time.Duration) {
		now := time.Now()

		for i := range loop {
			fmt.Println("loop ke :", i)
		}

		timeTaken = time.Since(now)
		return
	}

	var loop uint = 1_000_000
	tt1 := withBufio(loop)
	tt2 := withoutBufio(loop)
	fmt.Printf("\nWaktu yang dibutuhkan dengan bufio: %v\nWaktu yang dibutuhkan tanpa bufio: %v\n", tt1, tt2)
}

func fn5() {
	type Orang struct {
		Nama string
		Umur uint
	}

	var (
		a Orang
		b *Orang
		c = Orang{}
		d = new(Orang)
	)

	fmt.Printf("%+v\n%+v\n%+v\n%+v\n", a, b, c, d)
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
