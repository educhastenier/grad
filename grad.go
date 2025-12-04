package main

import (
	"fmt"
	"os"

	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gookit/color"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	LOG_LEVEL_OK      = "ok"
	LOG_LEVEL_WARNING = "warning"

	GRADLE_TASK_INTEGRATION_TEST = "integrationTest"
	GRADLE_TASK_BUILD            = "build"
)

// Config holds all configuration options for the grad command
type Config struct {
	Verbose                   bool
	CopyCommandClipboard      bool
	DoNotExecuteGradleCommand bool
	GradleTask                string
}

// NewConfigFromViper creates a Config instance from viper settings
func NewConfigFromViper() *Config {
	return &Config{
		Verbose:                   viper.GetBool("verbose"),
		CopyCommandClipboard:      viper.GetBool("copy-to-clipboard"),
		DoNotExecuteGradleCommand: viper.GetBool("no-execute"),
		GradleTask:                viper.GetString("task"),
	}
}

var (
	rootCmd = &cobra.Command{
		Use:   "grad [flags] [path]",
		Short: "Generate Gradle 🐘 command for a given path, passed as argument or from clipboard",
		Long: `Generate Gradle command for a given file/folder path, relative to Gradle root project folder, passed as argument or from clipboard.

The path argument can contain just the name of a class file (with / without .java): it will generate the command to run the integration test Gradle task for that class.
If the path argument is a folder, it will generate the command to build that project ("build" task).

You can also use a configuration file to set default values for flags. The configuration file should be named 'config.yaml' and placed in the current directory or in $HOME/.grad

'config.yaml' example:
	verbose: false
	copy-to-clipboard: false
	no-execute: false
	task: integrationTest`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := NewConfigFromViper()
			runCommand(args, cfg)
		},
	}
)

func init() {
	registerBooleanFlag("verbose", "v", "More verbose output")
	registerBooleanFlag("copy-to-clipboard", "c", "Copy the generated command to the clipboard")
	registerBooleanFlag("no-execute", "n", "Do not automatically run the generated command but simply print it")
	registerStringFlag("task", "t", "", "Gradle task to run. Default 'integrationTest' for Java files, 'build' for folders")

	initViper()
}

func registerBooleanFlag(longName, shortName, description string) {
	rootCmd.Flags().BoolP(longName, shortName, false, description)
	viper.BindPFlag(longName, rootCmd.Flags().Lookup(longName))
}

func registerStringFlag(longName, shortName, defaultValue, description string) {
	rootCmd.Flags().StringP(longName, shortName, defaultValue, description)
	viper.BindPFlag(longName, rootCmd.Flags().Lookup(longName))
}

func main() {
	rootCmd.Execute()
}

func runCommand(args []string, cfg *Config) {
	logVerbose(cfg, LOG_LEVEL_WARNING, "Starting with parameters: %s", args)
	var path string = ""

	if len(args) > 0 {
		path = args[0]
	} else {
		logVerbose(cfg, LOG_LEVEL_WARNING, "No path argument passed, trying to read from clipboard")
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
		logVerbose(cfg, LOG_LEVEL_OK, "Found file in: %s", foundPath)
		path = foundPath
	} else {
		logVerbose(cfg, LOG_LEVEL_WARNING, "No file found in the current directory (or subdirectories) with name '%s'. Assuming path is a Gradle path.", path)
	}

	cmd := transformPath(path, cfg)

	// Copy the command to the clipboard
	if cfg.CopyCommandClipboard {
		err = clipboard.WriteAll(cmd)
		if err != nil {
			PrintError("Failed to copy to clipboard: %s", err)
		} else {
			PrintWarning("Command copied to clipboard") // Always logged
		}
	}

	PrintOk("Gradle command: %s\n", color.OpBold.Render(cmd)) // Always logged

	// Execute the command in the terminal:
	if !cfg.DoNotExecuteGradleCommand {
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

func transformPath(path string, cfg *Config) string {
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
			task := GRADLE_TASK_INTEGRATION_TEST
			if cfg.GradleTask != "" {
				logVerbose(cfg, LOG_LEVEL_OK, "Overriding default task with '%s'", cfg.GradleTask)
				task = cfg.GradleTask
			}
			rep = cfmt.Sprintf("%s:%s --tests \"%s\"", beforeClassName, task, className)
		}
	} else {
		rep = strings.TrimSuffix(rep, "/")
		if cfg.GradleTask != "" {
			logVerbose(cfg, LOG_LEVEL_OK, "Overriding default task with '%s'", cfg.GradleTask)
			rep += ":" + cfg.GradleTask
		} else {
			rep += ":" + GRADLE_TASK_BUILD
		}
	}

	rep = strings.ReplaceAll(rep, "/", ":")
	rep = cfmt.Sprintf("./gradlew -PcreateTestReports :%s", rep)

	return rep
}

func findFile(root, filename string) (string, error) {
	// Early return if filename is empty
	if filename == "" {
		return "", nil
	}

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

func initViper() {
	// Set up Viper to read from a configuration file
	viper.SetConfigName("config")      // Name of the config file (without extension)
	viper.SetConfigType("yaml")        // Config file type
	viper.AddConfigPath(".")           // Path to look for the config file
	viper.AddConfigPath("$HOME/.grad") // Fallback path for config file

	// Read the configuration file (optional)
	if err := viper.ReadInConfig(); err != nil {
		// Config file is optional, so this is just a warning
		logVerbose(&Config{Verbose: true}, LOG_LEVEL_WARNING, "No config file found. Using default configuration.")
	}
}

func logVerbose(cfg *Config, level, msg string, a ...interface{}) {
	if !cfg.Verbose {
		return
	}
	switch level {
	case LOG_LEVEL_WARNING:
		PrintWarning(msg, a...)
	case LOG_LEVEL_OK:
		PrintOk(msg, a...)
	default:
		fmt.Printf(msg+"\n", a...)
	}
}
