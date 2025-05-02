package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Project struct {
	Sdk      string `xml:"Sdk,attr"`
	PropertyGroup []struct {
		OutputType string `xml:"OutputType"`
	} `xml:"PropertyGroup"`
	ItemGroups []struct {
		PackageReferences []struct {
			Include string `xml:"Include,attr"`
		} `xml:"PackageReference"`
	} `xml:"ItemGroup"`
}

func parseSln(file string) ([]string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var projects []string
	re := regexp.MustCompile(`Project\(".*?"\) = ".*?", "(.*?)", ".*?"`)
	matches := re.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			projects = append(projects, filepath.FromSlash(match[1]))
		}
	}
	return projects, nil
}

func analyzeProject(csprojPath string, targetPackage string) (isTest bool, isWeb bool, hasPackage bool, err error) {
	file, err := os.Open(csprojPath)
	if err != nil {
		return
	}
	defer file.Close()

	var project Project
	if err = xml.NewDecoder(file).Decode(&project); err != nil {
		return
	}

	for _, ig := range project.ItemGroups {
		for _, pr := range ig.PackageReferences {
			if strings.Contains(strings.ToLower(pr.Include), "test") {
				isTest = true
			}
			if pr.Include == targetPackage {
				hasPackage = true
			}
		}
	}
	if strings.Contains(project.Sdk, "Web") {
		isWeb = true
	}
	return
}
