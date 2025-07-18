package commands

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Push() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": "./",
	})

	defer rnr.Close()

	err := rnr.Run("git", "push", "origin", "HEAD")
	if err != nil {
		log.Fatal(err)
	}

}
