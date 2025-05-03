package dotnet

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Represents the .csproj XML structure we're interested in
type Project struct {
	Sdk        string `xml:"Sdk,attr"`
	ItemGroups []struct {
		PackageReferences []struct {
			Include string `xml:"Include,attr"`
		} `xml:"PackageReference"`
	} `xml:"ItemGroup"`
}

// Represents a project or solution folder from a .sln file
type SlnProject struct {
	Path              string
	IsSolutionFolder  bool
}

// Parses the .sln file to extract project paths and solution folders
func ParseSln(slnPath string) ([]SlnProject, error) {
	content, err := os.ReadFile(slnPath)
	if err != nil {
		return nil, err
	}

	var projects []SlnProject
	re := regexp.MustCompile(`Project\("\{(.*?)\}"\) = ".*?", "(.*?)", "\{.*?\}"`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) >= 3 {
			projectTypeGUID := strings.ToUpper(match[1])
			projectPath := filepath.FromSlash(match[2])
			isSolutionFolder := projectTypeGUID == "66A26720-8FB5-11D2-AA7E-00C04F688DDE"
			projects = append(projects, SlnProject{
				Path:             projectPath,
				IsSolutionFolder: isSolutionFolder,
			})
		}
	}
	return projects, nil
}

// Analyzes a .csproj file to determine project type and package usage
func AnalyzeProject(csprojPath, targetPackage string) (isTest bool, isWeb bool, hasPackage bool, err error) {
	info, err := os.Stat(csprojPath)
	if err != nil {
		return
	}
	if info.IsDir() {
		log.Println("path "+csprojPath+" is a directory, not a file")
		return
	}

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
			name := strings.ToLower(pr.Include)
			if strings.Contains(name, "test") || strings.Contains(name, "xunit") || strings.Contains(name, "nunit") {
				isTest = true
			}
			if pr.Include == targetPackage {
				hasPackage = true
			}
		}
	}

	if strings.Contains(strings.ToLower(project.Sdk), "web") {
		isWeb = true
	}

	return
}
