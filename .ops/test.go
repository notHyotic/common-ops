package main

import (
	"fmt"
	"log"
	"path/filepath"
)

func (Ops) Test() {
	slnPath := "MySolution.sln"           // Replace with actual .sln file path
	targetPackage := "Newtonsoft.Json"    // Replace with the package to check for

	projects, err := parseSln(slnPath)
	if err != nil {
		log.Println("Failed to read .sln file:", err)
		return
	}

	slnDir := filepath.Dir(slnPath)

	for _, project := range projects {
		if project.IsSolutionFolder {
			fmt.Printf("Solution Folder: %-50s\n", project.Path)
			continue
		}

		fullPath := filepath.Join(slnDir, project.Path)
		isTest, isWeb, hasPkg, err := analyzeProject(fullPath, targetPackage)
		if err != nil {
			fmt.Printf("Error analyzing %s: %v\n", fullPath, err)
			continue
		}

		fmt.Printf("Project: %-50s | Test: %-5t | Web: %-5t | Has '%s': %-5t\n",
			project.Path, isTest, isWeb, targetPackage, hasPkg)
	}
}
