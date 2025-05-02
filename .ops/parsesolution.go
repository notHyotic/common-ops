package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Project interface {
	GetID() string
	SetParent(*SolutionFolder)
	GetParent() *SolutionFolder
	IsFolder() bool
}

type SolutionProject struct {
	ID     string
	Name   string
	Path   string
	Type   string
	Parent *SolutionFolder
}

func (p *SolutionProject) GetID() string               { return p.ID }
func (p *SolutionProject) SetParent(f *SolutionFolder) { p.Parent = f }
func (p *SolutionProject) GetParent() *SolutionFolder  { return p.Parent }
func (p *SolutionProject) IsFolder() bool              { return false }

type SolutionFolder struct {
	SolutionProject
	Items []Project
}

func (f *SolutionFolder) IsFolder() bool { return true }

type SolutionParserResult struct {
	Version                   string
	VisualStudioVersion       string
	MinimumVisualStudioVersion string
	Projects                  []Project
}

func ParseSolution(path string) (*SolutionParserResult, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("solution file '%s' does not exist", absPath)
	}
	defer file.Close()

	var version, visualStudioVersion, minimumVisualStudioVersion string
	var projects []Project
	inNestedProjects := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Microsoft Visual Studio Solution File, "):
			version = line[39:]

		case strings.HasPrefix(line, "VisualStudioVersion = "):
			visualStudioVersion = line[22:]

		case strings.HasPrefix(line, "MinimumVisualStudioVersion = "):
			minimumVisualStudioVersion = line[29:]

		case strings.HasPrefix(line, "Project(\"{"):
			project, err := parseSolutionProjectLine(line, filepath.Dir(absPath))
			if err != nil {
				continue
			}
			if strings.ToUpper(project.Type) == strings.ToUpper("{66A26720-8FB5-11D2-AA7E-00C04F688DDE}") {
				folder := &SolutionFolder{
					SolutionProject: *project,
				}
				projects = append(projects, folder)
			} else {
				projects = append(projects, project)
			}

		case strings.HasPrefix(line, "GlobalSection(NestedProjects"):
			inNestedProjects = true

		case inNestedProjects && strings.HasPrefix(line, "EndGlobalSection"):
			inNestedProjects = false

		case inNestedProjects:
			parseNestedProjectLine(projects, line)
		}
	}

	return &SolutionParserResult{
		Version:                   version,
		VisualStudioVersion:       visualStudioVersion,
		MinimumVisualStudioVersion: minimumVisualStudioVersion,
		Projects:                  projects,
	}, nil
}

func parseSolutionProjectLine(line, basePath string) (*SolutionProject, error) {
	parts := []string{}
	start := 0
	inQuote := false

	for i := range len(line) {
		if line[i] == '"' {
			if inQuote {
				parts = append(parts, line[start+1:i])
			}
			inQuote = !inQuote
			start = i
		}
	}

	if len(parts) < 4 {
		return nil, errors.New("invalid project line")
	}

	projPath := parts[2]
	if !filepath.IsAbs(projPath) {
		projPath = filepath.Join(basePath, projPath)
	}

	return &SolutionProject{
		Type: parts[0],
		Name: parts[1],
		Path: projPath,
		ID:   parts[3],
	}, nil
}

func parseNestedProjectLine(projects []Project, line string) {
	parts := strings.Split(line, " = ")
	if len(parts) != 2 {
		return
	}
	childID := strings.TrimSpace(parts[0])
	parentID := strings.TrimSpace(parts[1])

	var child Project
	var parent *SolutionFolder

	for _, p := range projects {
		if p.GetID() == childID {
			child = p
		}
		if p.GetID() == parentID {
			if folder, ok := p.(*SolutionFolder); ok {
				parent = folder
			}
		}
	}

	if child != nil && parent != nil {
		parent.Items = append(parent.Items, child)
		child.SetParent(parent)
	}
}
