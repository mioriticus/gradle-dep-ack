package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"sort"
	"strings"
)

type Dependency struct {
	Name     string
	License  string
	Src      string
	Group    string
	Artifact string
	Version  string
}

func main() {
	rawDeps := getGradleDeps("gradle-dep-list.txt")
	cleanedDeps := cleanDependencies(rawDeps)
	deps := sortDependencies(cleanedDeps)

	for _, dep := range deps {
		fmt.Println(dep.Src)
		if dep.Group != "" && dep.Artifact != "" {
			url := fmt.Sprintf("https://mvnrepository.com/artifact/%s/%s", dep.Group, dep.Artifact)
			doc, err := goquery.NewDocument(url)
			if err != nil {
				log.Fatal(err)
			}
			artifactName := doc.Find("h2.im-title a").Text()
			artifactLicense := doc.Find("span.b.lic").Text()
			if artifactName != "" {
				fmt.Println("Library:", artifactName)
				fmt.Println("License:", artifactLicense)
			}
		}
		fmt.Println("--")
	}
}

func sortDependencies(deps []Dependency) []Dependency {
	// Deduplicate dependencies
	found := make(map[string]bool)
	uniqueDeps := make(map[string]Dependency)
	for _, dep := range deps {
		if !found[dep.Src] {
			found[dep.Src] = true
			uniqueDeps[dep.Src] = dep
		}
	}

	// Sort keys
	var keys []string
	for k := range found {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sortedDeps []Dependency
	for _, k := range keys {
		sortedDeps = append(sortedDeps, uniqueDeps[k])
	}

	return sortedDeps
}

func cleanDependencies(deps []string) []Dependency {
	for i, dep := range deps {
		// Remove tree drawing
		dep = strings.Replace(dep, "+---", "", -1)
		dep = strings.Replace(dep, "\\---", "", -1)
		dep = strings.Replace(dep, "|   ", "", -1)
		// Remove extra spaces
		dep = strings.TrimSpace(dep)
		// Now put it back
		deps[i] = dep
	}

	dependencies := make([]Dependency, 0)
	for _, dep := range deps {
		//fmt.Println(dep) // debug
		// Skip projects
		if strings.HasPrefix(dep, "project") {
			continue
		}

		dependency := Dependency{Src: dep}

		// Split by ":". We get group:artifact:version. If not, ignore it.
		parts := strings.Split(dep, ":")
		if len(parts) == 3 {
			dependency.Group = parts[0]
			dependency.Artifact = parts[1]

			version := parts[2]
			versionParts := strings.Split(version, " -> ")
			if len(versionParts) == 2 {
				version = versionParts[1]
				// Redo Src
				dependency.Src = fmt.Sprintf("%s:%s:%s", dependency.Group, dependency.Artifact, version)
			}
			dependency.Version = version
		}

		// Add it to list
		dependencies = append(dependencies, dependency)
	}

	return dependencies
}

func getGradleDeps(fileName string) []string {
	startLines := false
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "compile - Classpath for compiling the main sources." {
			startLines = true
			continue // skip this line
		}
		// If it gets to the empty line, we're done
		if startLines && text == "" {
			break
		}
		if startLines {
			// If it has "(*)" inside, it's a duplicate
			if !strings.Contains(text, "(*)") {
				lines = append(lines, text)
			}
		}
	}
	return lines
}
