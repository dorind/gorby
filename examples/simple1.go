package main

import (
	"fmt"

	"github.com/dorind/gorby"
)

// PushString() and validation
func gorby_example1() {
	s := "simple gorby example"
	rb := gorby.NewRuneBuff(0)
	rb.PushString(s)
	srb := rb.String()
	fmt.Println("\ngorby example 1")
	fmt.Println("Pushed String:", `"`+s+`"`)
	fmt.Println("Result String:", `"`+srb+`"`)
	fmt.Println("Valid?", srb == s)
}

// PushRune() by checking the oddness of the rune's value
func gorby_example2() {
	src := "0123456789876543210"
	runes := []rune(src)
	rodd := gorby.NewRuneBuff(0)
	reven := gorby.NewRuneBuff(0)
	for _, r := range runes {
		if r&1 == 1 {
			rodd.PushRune(r)
		} else {
			// do you even push bro?
			reven.PushRune(r)
		}
	}
	fmt.Println("\ngorby example 2")
	fmt.Println("src", src)
	fmt.Println("odd", rodd.String())
	fmt.Println("even", reven.String())
}

// buffer auto-grow
func gorby_example3() {
	src := "0123456789_abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runes := []rune(src)
	len_runes := uint(len(runes))
	rb := gorby.NewRuneBuff(0)
	fmt.Println("\ngorby example 3")
	fmt.Println("src", src)
	fmt.Println("src len", len_runes)
	fmt.Println("rb cap before push", rb.Capacity())
	for _, r := range runes {
		rb.PushRune(r)
	}
	fmt.Println("rb cap after push", rb.Capacity())
	cap_rb := rb.Capacity()
	ncap := cap_rb - gorby.KRUNE_BUFF_MIN
	fmt.Println("rb cap grew by", ncap,
		"runes, about", (float32(cap_rb) / float32(len_runes)),
		"times the required size")
	srb := rb.String()
	fmt.Println("srb", srb)
	fmt.Println("Valid?", srb == src)
}

func main() {
	gorby_example1()
	gorby_example2()
	gorby_example3()
}



