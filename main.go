package main

import (
	"fmt"
	cli "github.com/mainak55512/qwe/cli"
	"os"
)

func main() {
	if err := cli.HandleArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
