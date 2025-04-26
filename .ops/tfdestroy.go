package main

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Tfdestroy() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": "./infrastructure",
	})

	defer rnr.Close()

	err := rnr.Run("terraform", "init")
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("terraform", "validate")
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("terraform", "destroy", "-auto-approve")
	if err != nil {
		log.Fatal(err)
	}
}
