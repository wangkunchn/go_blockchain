package main

import (
	"strconv"
	"fmt"
)

func main() {
	amountInt, _ := strconv.ParseInt("2", 10, 64)
	fmt.Println(amountInt)
}
