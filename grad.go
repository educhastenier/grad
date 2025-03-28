package main

import (
	"os"

	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gookit/color"
	"github.com/i582/cfmt/cmd/cfmt"
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
		PrintWarning("Starting with parameters: %s", args)
	}
	var path string = ""

	if len(args) > 0 {
		path = args[0]
	} else {
		if verbose {
			PrintWarning("No path argument passed, trying to read from clipboard")
		}
		pathReadFromCp, err := clipboard.ReadAll()
		if err != nil {
			PrintError("Failed to read from clipboard: %s", err)
			return
		}
		if pathReadFromCp != "" {
			path = pathReadFromCp
		}
	}

	foundPath, err := findFile(".", path)
	if err != nil {
		PrintError("Error: %s", err)
		return
	}
	if foundPath != "" {
		if verbose {
			PrintOk("Found file in: %s", foundPath)
		}
		path = foundPath
	} else if verbose {
		PrintWarning("No file found in the current directory (or subdirectories) with name '%s'. Assuming path is a Gradle path.", path)
	}

	cmd := transformPath(path)

	// Copy the command to the clipboard
	if copyCommandClipboard {
		err = clipboard.WriteAll(cmd)
		if err != nil {
			PrintError("Failed to copy to clipboard: %s", err)
		} else {
			PrintWarning("Command copied to clipboard")
		}
	}

	PrintOk("Gradle command: %s\n", color.OpBold.Render(cmd))

	// Execute the command in the terminal:
	if !doNotExecuteGradleCommand {
		command := exec.Command("zsh", "-c", cmd)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			PrintError("Failed to execute command: %s", err)
		}
	} else {
		PrintWarning("Command not executed (flag -n / --no-execute passed)")
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
					PrintOk("Overriding default task with '%s'", gradleTask)
				}
				task = gradleTask
			}
			rep = cfmt.Sprintf("%s:%s --tests \"%s\"", beforeClassName, task, className)
		}
	} else {
		rep = strings.TrimSuffix(rep, "/")
		if gradleTask != "" {
			if verbose {
				PrintOk("Overriding default task with '%s'", gradleTask)
			}
			rep += ":" + gradleTask
		} else {
			rep += ":build"
		}
	}

	rep = strings.ReplaceAll(rep, "/", ":")
	rep = cfmt.Sprintf("./gradlew -PcreateTestReports :%s", rep)

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

func PrintError(msg string, a ...interface{}) (n int, err error) {
	return cfmt.Printf("{{"+msg+"}}::bold|red\n", a...)
}

func PrintOk(msg string, a ...interface{}) (n int, err error) {
	return cfmt.Printf("{{"+msg+"}}::green\n", a...)
}

func PrintWarning(msg string, a ...interface{}) (n int, err error) {
	return cfmt.Printf("{{"+msg+"}}::yellow\n", a...)
}
