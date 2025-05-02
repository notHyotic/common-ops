package main

import (
	"log"
	"os"
)

func (Ops) Test() {
	args := os.Args
	result, err := ParseSolution(args[2])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Parsed solution: %s\n", args[2])
	log.Printf("Version: %s\n", result.Version)
	log.Printf("Visual Studio Version: %s\n", result.VisualStudioVersion)
	log.Printf("Minimum Visual Studio Version: %s\n\n", result.MinimumVisualStudioVersion)

	log.Println("Projects:")
	for _, proj := range result.Projects {
		prefix := ""
		if proj.IsFolder() {
			prefix = "[Folder] "
		}
		log.Printf("  %s%s (%s)\n", prefix, proj.(*SolutionProject).Name, proj.GetID())

		if folder, ok := proj.(*SolutionFolder); ok && len(folder.Items) > 0 {
			log.Println("    Contains:")
			for _, item := range folder.Items {
				log.Printf("      - %s (%s)\n", item.(*SolutionProject).Name, item.GetID())
			}
		}
	}
}