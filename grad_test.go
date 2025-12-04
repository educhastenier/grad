package main

import (
	"testing"
)

func TestSubscriptionPathWithTestJavaClass(t *testing.T) {
	// fmt.Println("TestSubscriptionPath")
	input := "subscription/bonita-integration-tests-sp/bonita-integration-tests-client/src/test/java/com/bonitasoft/engine/process/ProcessManagementIT.java"
	expected := "./gradlew -PcreateTestReports :subscription:bonita-integration-tests-sp:bonita-integration-tests-client:integrationTest --tests \"com.bonitasoft.engine.process.ProcessManagementIT\""

	runTest(t, input, expected)
}

func TestCommunityBuildGradlePath(t *testing.T) {
	// fmt.Println("TestCommunityBuildGradlePath")
	input := "community/some/path/to/build.gradle"
	expected := "./gradlew -PcreateTestReports :some:path:to:build"

	runTest(t, input, expected)
}

func TestCommunityBuildKtsPath(t *testing.T) {
	// fmt.Println("TestCommunityScriptKtsPath")
	input := "community/another/path/to/build.gradle.kts"
	expected := "./gradlew -PcreateTestReports :another:path:to:build"

	runTest(t, input, expected)
}

func TestDefaultTaskOverride(t *testing.T) {
	// Save the original value and restore it after the test
	originalTask := gradleTask
	defer func() { gradleTask = originalTask }()

	// Set a custom task to test override functionality
	gradleTask = "assemble"

	input := "community/another/path/to"
	expected := "./gradlew -PcreateTestReports :another:path:to:assemble"

	runTest(t, input, expected)
}

func runTest(t *testing.T, input string, expected string) {
	// fmt.Println("testing:", t.Name())
	got := transformPath(input)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
