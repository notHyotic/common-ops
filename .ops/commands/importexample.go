package commands

import (
	"log"

	dotnet "ops/libs/dotnet"

	"lesiw.io/cmdio/sys"
)

// example of how to import libs from libs folder
func (Ops) Importexample() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": "./",
	})

	defer rnr.Close()

	sln, err := dotnet.ParseSln("bn-api-microservices.sln")
	for _, project := range sln {
		log.Println(project.Path)
	}

	if err != nil {
		log.Fatal(err)
	}

}
