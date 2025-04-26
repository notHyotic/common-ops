package main

import (
	"log"
	"os"

	"lesiw.io/cmdio/sys"
)

func (Ops) Commit() {
	args := os.Args
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": "./",
	})

	defer rnr.Close()

	err := rnr.Run("git", "commit", "-m", args[2])
	if err != nil {
		log.Fatal(err)
	}
}
