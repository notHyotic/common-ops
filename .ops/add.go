package main

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Add() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": "./",
	})

	defer rnr.Close()

	err := rnr.Run("git", "add", ".")
	if err != nil {
		log.Fatal(err)
	}
}
