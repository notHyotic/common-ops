package main

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Revert() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": "./",
	})

	defer rnr.Close()

	err := rnr.Run("git", "reset", "--soft", "HEAD~1")
	if err != nil {
		log.Fatal(err)
	}
}
