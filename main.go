package main

import (
	"fmt"
	cli "github.com/mainak55512/qwe/cli"
)

func main() {
	if err := cli.HandleArgs(); err != nil {
		fmt.Println(err)
	}
}
