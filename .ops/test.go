package main

import (
	"log"
)

func (Ops) Test() {
	projects, _ := parseSln("bn-api.microservices.sln")
	log.Println(len(projects))
	for _, p := range projects {
		isTest, isWeb, hasPkg, err := analyzeProject(p, "Newtonsoft.Json")
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		log.Printf("Project: %s | Test: %t | Web: %t | HasPkg: %t\n", p, isTest, isWeb, hasPkg)
	}
}
