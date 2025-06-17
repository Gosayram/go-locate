package version

import (
	"runtime"
	"strings"
	"testing"
)

const (
	TestVersion     = "1.2.3"
	TestCommit      = "abc123def456"
	TestDate        = "2024-01-15_12:30:45"
	TestBuiltBy     = "test-user"
	TestBuildNumber = "42"
)

func TestGetVersion(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalBuildNumber := BuildNumber

	// Test without build number
	Version = TestVersion
	BuildNumber = "0"

	result := GetVersion()
	expected := TestVersion
	if result != expected {
		t.Errorf("GetVersion() = %q, want %q", result, expected)
	}

	// Test with build number
	BuildNumber = TestBuildNumber
	result = GetVersion()
	expected = "1.2.3 (build 42)"
	if result != expected {
		t.Errorf("GetVersion() with build = %q, want %q", result, expected)
	}

	// Restore original values
	Version = originalVersion
	BuildNumber = originalBuildNumber
}

func TestGetFullVersionInfo(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalCommit := Commit
	originalDate := Date
	originalBuiltBy := BuiltBy

	// Set test values
	Version = TestVersion
	Commit = TestCommit
	Date = TestDate
	BuiltBy = TestBuiltBy

	result := GetFullVersionInfo()

	// Check if result contains expected components
	if !strings.Contains(result, "glocate") {
		t.Error("GetFullVersionInfo() should contain 'glocate'")
	}

	if !strings.Contains(result, TestVersion) {
		t.Errorf("GetFullVersionInfo() should contain version %q", TestVersion)
	}

	if !strings.Contains(result, TestCommit[:ShortCommitHashLength]) {
		t.Errorf("GetFullVersionInfo() should contain short commit %q", TestCommit[:ShortCommitHashLength])
	}

	if !strings.Contains(result, "built 2024-01-15 12:30:45") {
		t.Error("GetFullVersionInfo() should contain formatted date")
	}

	if !strings.Contains(result, TestBuiltBy) {
		t.Errorf("GetFullVersionInfo() should contain built by %q", TestBuiltBy)
	}

	if !strings.Contains(result, runtime.Version()) {
		t.Error("GetFullVersionInfo() should contain Go version")
	}

	// Restore original values
	Version = originalVersion
	Commit = originalCommit
	Date = originalDate
	BuiltBy = originalBuiltBy
}

func TestGet(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalCommit := Commit
	originalDate := Date
	originalBuiltBy := BuiltBy

	// Set test values
	Version = TestVersion
	Commit = TestCommit
	Date = TestDate
	BuiltBy = TestBuiltBy

	buildInfo := Get()

	if buildInfo.Version != TestVersion {
		t.Errorf("BuildInfo.Version = %q, want %q", buildInfo.Version, TestVersion)
	}

	if buildInfo.Commit != TestCommit {
		t.Errorf("BuildInfo.Commit = %q, want %q", buildInfo.Commit, TestCommit)
	}

	if buildInfo.Date != TestDate {
		t.Errorf("BuildInfo.Date = %q, want %q", buildInfo.Date, TestDate)
	}

	if buildInfo.BuiltBy != TestBuiltBy {
		t.Errorf("BuildInfo.BuiltBy = %q, want %q", buildInfo.BuiltBy, TestBuiltBy)
	}

	if buildInfo.GoVersion != runtime.Version() {
		t.Errorf("BuildInfo.GoVersion = %q, want %q", buildInfo.GoVersion, runtime.Version())
	}

	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if buildInfo.Platform != expectedPlatform {
		t.Errorf("BuildInfo.Platform = %q, want %q", buildInfo.Platform, expectedPlatform)
	}

	// Restore original values
	Version = originalVersion
	Commit = originalCommit
	Date = originalDate
	BuiltBy = originalBuiltBy
}

func TestBuildInfoString(t *testing.T) {
	buildInfo := &BuildInfo{
		Version:   TestVersion,
		Commit:    TestCommit,
		Date:      TestDate,
		BuiltBy:   TestBuiltBy,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}

	result := buildInfo.String()

	expectedComponents := []string{
		"glocate version",
		TestVersion,
		TestCommit,
		TestDate,
		TestBuiltBy,
		runtime.Version(),
		runtime.GOOS + "/" + runtime.GOARCH,
	}

	for _, component := range expectedComponents {
		if !strings.Contains(result, component) {
			t.Errorf("BuildInfo.String() should contain %q, got: %s", component, result)
		}
	}
}

func TestBuildInfoShort(t *testing.T) {
	buildInfo := &BuildInfo{
		Version: TestVersion,
	}

	result := buildInfo.Short()
	expected := "glocate " + TestVersion

	if result != expected {
		t.Errorf("BuildInfo.Short() = %q, want %q", result, expected)
	}
}

func TestShortCommitHash(t *testing.T) {
	// Save original values
	originalCommit := Commit

	// Test with long commit
	Commit = "abcdef1234567890"
	result := GetFullVersionInfo()

	if !strings.Contains(result, "abcdef1") {
		t.Error("Should truncate long commit to 7 characters")
	}

	// Test with short commit
	Commit = "abc123"
	result = GetFullVersionInfo()

	if !strings.Contains(result, "abc123") {
		t.Error("Should preserve short commit as-is")
	}

	// Restore original value
	Commit = originalCommit
}

func TestUnknownValues(t *testing.T) {
	// Save original values
	originalCommit := Commit
	originalDate := Date
	originalBuiltBy := BuiltBy

	// Set unknown values
	Commit = UnknownValue
	Date = UnknownValue
	BuiltBy = UnknownValue

	result := GetFullVersionInfo()

	// Should not contain unknown values in output
	if strings.Contains(result, "(unknown)") {
		t.Error("Should not include unknown commit in parentheses")
	}

	if strings.Contains(result, "built unknown") {
		t.Error("Should not include unknown date")
	}

	if strings.Contains(result, "by unknown") {
		t.Error("Should not include unknown built by")
	}

	// Restore original values
	Commit = originalCommit
	Date = originalDate
	BuiltBy = originalBuiltBy
}
