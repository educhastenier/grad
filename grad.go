package main

import (
	"fmt"
	"os"

	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	// fmt.Printf("%s: starting\n", os.Args[0])
	var path string = ""

	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		pathReadFromCp, err := clipboard.ReadAll()
		if err != nil {
			fmt.Printf("Failed to read from clipboard: %s\n", err)
			return
		}
		if pathReadFromCp != "" {
			path = pathReadFromCp
		}
	}

	foundPath, err := findFile(".", path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if foundPath != "" {
		fmt.Println("Found file in:", foundPath)
		path = foundPath
	}

	cmd := transformPath(path)

	// Copy the command to the clipboard
	err = clipboard.WriteAll(cmd)
	if err != nil {
		fmt.Printf("Failed to copy to clipboard: %s\n", err)
	} else {
		fmt.Println("Command copied to clipboard")
	}

	// fmt.Println(cmd)
	// Execute the command in the terminal
	fmt.Println("Executing command:", cmd)
	command := exec.Command("zsh", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
	}
}

func transformPath(path string) string {
	rep, _ := strings.CutPrefix(path, "community/")
	rep = strings.TrimSuffix(rep, ".kts")
	rep = strings.TrimSuffix(rep, "build.gradle")

	// if rep ends with '.java':
	if strings.HasSuffix(rep, ".java") {
		rep = strings.TrimSuffix(rep, ".java")
		if strings.Contains(rep, "src/test/java") {
			// get the substring after 'src/test/java'
			beforeClassName, className, _ := strings.Cut(rep, "src/test/java/")
			className = strings.ReplaceAll(className, "/", ".")
			beforeClassName = strings.TrimSuffix(beforeClassName, "/")
			rep = fmt.Sprintf("%s:integrationTest --tests \"%s\"", beforeClassName, className)
		}
	} else {
		rep = strings.TrimSuffix(rep, "/")
		rep += ":build"
	}
	// parts := strings.Split(path, ":")
	// className := parts[len(parts)-1]
	// parts = parts[:len(parts)-1]

	rep = strings.ReplaceAll(rep, "/", ":")
	rep = fmt.Sprintf("./gradlew :%s", rep)

	return rep
}

func findFile(root, filename string) (string, error) {
	var foundPath string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (info.Name() == filename || info.Name() == filename+".java") {
			foundPath = path
			return filepath.SkipDir
		}
		return nil
	})
	return foundPath, err
}
