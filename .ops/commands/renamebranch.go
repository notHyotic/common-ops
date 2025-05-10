package commands

import (
	"log"
	"os"

	"lesiw.io/cmdio/sys"
)

func (Ops) Renamebranch() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": "./",
	})
	args := os.Args
	newbranchname := args[2]
	defer rnr.Close()

	err := rnr.Run("git", "branch", "-M", newbranchname)
	if err != nil {
		log.Fatal(err)
	}
}
