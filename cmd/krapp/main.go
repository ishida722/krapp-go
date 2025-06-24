package main

import (
	"os"

	"github.com/ishida722/krapp-go/cmd/krapp/krapp"
)

func main() {
	if err := krapp.Execute(); err != nil {
		os.Exit(1)
	}
}
