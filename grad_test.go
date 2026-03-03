package main

import (
	"testing"
)

func TestIntegrationTestDetection(t *testing.T) {
	input := "subscription/integration-tests-sp/integration-tests-client/src/test/java/com/company/runtime/access/ProcessManagementIT.java"
	expected := "./gradlew -PcreateTestReports :subscription:integration-tests-sp:integration-tests-client:integrationTest --tests \"com.company.runtime.access.ProcessManagementIT\""

	cfg := &Config{}
	runTest(t, input, expected, cfg)
}

func TestUnitTestDetection(t *testing.T) {
	input := "subscription/project-runtime/src/test/java/com/company/runtime/service/MyServiceTest.java"
	expected := "./gradlew -PcreateTestReports :subscription:project-runtime:test --tests \"com.company.runtime.service.MyServiceTest\""

	cfg := &Config{}
	runTest(t, input, expected, cfg)
}

func TestGroovyIntegrationTestDetection(t *testing.T) {
	input := "subscription/integration-tests-sp/src/test/groovy/com/company/runtime/access/ProcessManagementIT.groovy"
	expected := "./gradlew -PcreateTestReports :subscription:integration-tests-sp:integrationTest --tests \"com.company.runtime.access.ProcessManagementIT\""

	cfg := &Config{}
	runTest(t, input, expected, cfg)
}

func TestGroovyUnitTestDetection(t *testing.T) {
	input := "subscription/project-runtime/src/test/groovy/com/company/runtime/service/MyServiceTest.groovy"
	expected := "./gradlew -PcreateTestReports :subscription:project-runtime:test --tests \"com.company.runtime.service.MyServiceTest\""

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
