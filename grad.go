package main

import (
	"fmt"
	"os"

	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	verbose                   bool
	copyCommandClipboard      bool
	doNotExecuteGradleCommand bool
	gradleTask                string
	rootCmd                   = &cobra.Command{
		Use:   "grad [path] [flags]",
		Short: "Generate Gradle 🐘 command for a given path, passed as argument or from clipboard",
		Long: `Generate Gradle command for a given file/folder path, relative to Gradle root project folder, passed as argument or from clipboard.

The path argument can contain just the name of a class file (with / without .java): it will generate the command to run the integration test Gradle task for that class.
If the path argument is a folder, it will generate the command to build that project ("build" task).`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args)
		},
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "More verbose output")
	rootCmd.Flags().BoolVarP(&copyCommandClipboard, "copy-to-clipboard", "c", false, "Copy the generated command to the clipboard")
	rootCmd.Flags().BoolVarP(&doNotExecuteGradleCommand, "no-execute", "n", false, "Do not automatically run the generated command but simply print it")
	rootCmd.Flags().StringVarP(&gradleTask, "task", "t", "", "Gradle task to run. Default 'integrationTest' for Java files, 'build' for folders")
}

func main() {
	rootCmd.Execute()
}

func runCommand(args []string) {
	if verbose {
		fmt.Println("Starting with parameters:", args)
	}
	var path string = ""

	if len(args) > 0 {
		path = args[0]
	} else {
		if verbose {
			println("No path argument passed, trying to read from clipboard")
		}
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
		if verbose {
			fmt.Println("Found file in:", foundPath)
		}
		path = foundPath
	} else if verbose {
		fmt.Printf("No file found in the current directory (or subdirectories) with name %s. Assuming path is a Gradle path.\n", path)
	}

	cmd := transformPath(path)

	// Copy the command to the clipboard
	if copyCommandClipboard {
		err = clipboard.WriteAll(cmd)
		if err != nil {
			fmt.Printf("Failed to copy to clipboard: %s\n", err)
		} else {
			fmt.Println("Command copied to clipboard")
		}
	}

	fmt.Println("Gradle command:", cmd)

	// Execute the command in the terminal:
	if !doNotExecuteGradleCommand {
		command := exec.Command("zsh", "-c", cmd)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			fmt.Printf("Failed to execute command: %s\n", err)
		}
	} else {
		fmt.Println("Command not executed (flag -n / --no-execute passed)")
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
			task := "integrationTest"
			if gradleTask != "" {
				if verbose {
					fmt.Printf("Overriding default task with '%s'\n", gradleTask)
				}
				task = gradleTask
			}
			rep = fmt.Sprintf("%s:%s --tests \"%s\"", beforeClassName, task, className)
		}
	} else {
		rep = strings.TrimSuffix(rep, "/")
		if gradleTask != "" {
			if verbose {
				fmt.Printf("Overriding default task with '%s'\n", gradleTask)
			}
			rep += ":" + gradleTask
		} else {
			rep += ":build"
		}
	}

	rep = strings.ReplaceAll(rep, "/", ":")
	rep = fmt.Sprintf("./gradlew -PcreateTestReports :%s", rep)

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
