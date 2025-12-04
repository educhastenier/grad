package main

import (
	"testing"
)

func TestSubscriptionPathWithTestJavaClass(t *testing.T) {
	// fmt.Println("TestSubscriptionPath")
	input := "subscription/bonita-integration-tests-sp/bonita-integration-tests-client/src/test/java/com/bonitasoft/engine/process/ProcessManagementIT.java"
	expected := "./gradlew -PcreateTestReports :subscription:bonita-integration-tests-sp:bonita-integration-tests-client:integrationTest --tests \"com.bonitasoft.engine.process.ProcessManagementIT\""

	cfg := &Config{}
	runTest(t, input, expected, cfg)
}

func TestCommunityBuildGradlePath(t *testing.T) {
	// fmt.Println("TestCommunityBuildGradlePath")
	input := "community/some/path/to/build.gradle"
	expected := "./gradlew -PcreateTestReports :some:path:to:build"

	cfg := &Config{}
	runTest(t, input, expected, cfg)
}

func TestCommunityBuildKtsPath(t *testing.T) {
	// fmt.Println("TestCommunityScriptKtsPath")
	input := "community/another/path/to/build.gradle.kts"
	expected := "./gradlew -PcreateTestReports :another:path:to:build"

	cfg := &Config{}
	runTest(t, input, expected, cfg)
}

func TestDefaultTaskOverride(t *testing.T) {
	// Test that custom task overrides the default "build" task
	cfg := &Config{
		GradleTask: "assemble",
	}

	input := "community/another/path/to"
	expected := "./gradlew -PcreateTestReports :another:path:to:assemble"

	runTest(t, input, expected, cfg)
}

func runTest(t *testing.T, input string, expected string, cfg *Config) {
	// fmt.Println("testing:", t.Name())
	got := transformPath(input, cfg)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
