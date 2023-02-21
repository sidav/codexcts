package main

import "fmt"

type unitSpecial struct {
	name  string
	value int
}

func (us *unitSpecial) getFormattedName() string {
	if us.value > 0 {
		return fmt.Sprintf("%s %d", us.name, us.value)
	}
	return us.name
}
