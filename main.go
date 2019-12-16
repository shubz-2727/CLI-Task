package main

import (
	"fmt"
	"os"
	"path/filepath"

	"task/cmd"
	"task/db"

	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	home, _ := homedir.Dir()
	//fmt.Println(home)
	dbPath := filepath.Join(home, "tasks.db")

	must(db.Init(dbPath))
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
