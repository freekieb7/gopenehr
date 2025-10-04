package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestSuite represents a collection of benchmark tests
type TestSuite struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tests       []Test `json:"tests"`
}

// Test represents a single benchmark test
type Test struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	AQL         string         `json:"aql"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Timeout     *time.Duration `json:"timeout,omitempty"`
}

// BenchmarkResult represents the result of a single benchmark test
type BenchmarkResult struct {
	Test         Test          `json:"test"`
	Duration     time.Duration `json:"duration"`
	Success      bool          `json:"success"`
	Error        string        `json:"error,omitempty"`
	ResponseSize int           `json:"response_size"`
	StatusCode   int           `json:"status_code"`
	RowCount     *int          `json:"row_count,omitempty"`
}

// BenchmarkReport represents the complete benchmark results
type BenchmarkReport struct {
	Suite     TestSuite         `json:"suite"`
	Results   []BenchmarkResult `json:"results"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Duration  time.Duration     `json:"total_duration"`
	Summary   Summary           `json:"summary"`
}

// Summary provides aggregate statistics
type Summary struct {
	TotalTests      int           `json:"total_tests"`
	PassedTests     int           `json:"passed_tests"`
	FailedTests     int           `json:"failed_tests"`
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`
}

// JUnit XML structures for Azure DevOps compatibility
type JUnitTestSuites struct {
	XMLName  xml.Name         `xml:"testsuites"`
	Name     string           `xml:"name,attr"`
	Tests    int              `xml:"tests,attr"`
	Failures int              `xml:"failures,attr"`
	Time     string           `xml:"time,attr"`
	Suites   []JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Time      string          `xml:"time,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Name      string        `xml:"name,attr"`
	ClassName string        `xml:"classname,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
}

type JUnitFailure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Text    string   `xml:",chardata"`
}

// QueryRequest represents the request to the /query endpoint
type QueryRequest struct {
	AQL        string         `json:"aql"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// QueryResponse represents the response from the /query endpoint
type QueryResponse struct {
	Rows []map[string]any `json:"rows,omitempty"`
	Meta map[string]any   `json:"meta,omitempty"`
}

// BenchmarkRunner executes benchmark tests
type BenchmarkRunner struct {
	APIURL  string
	Timeout time.Duration
	Verbose bool
}

// RunBenchmarks executes all tests in a suite and returns a report
func (r *BenchmarkRunner) RunBenchmarks(ctx context.Context, suite TestSuite) (*BenchmarkReport, error) {
	startTime := time.Now()

	report := &BenchmarkReport{
		Suite:     suite,
		Results:   make([]BenchmarkResult, 0, len(suite.Tests)),
		StartTime: startTime,
	}

	for i, test := range suite.Tests {
		if r.Verbose {
			fmt.Printf("Running test %d/%d: %s\n", i+1, len(suite.Tests), test.Name)
		}

		result := r.runSingleTest(ctx, test)
		report.Results = append(report.Results, result)

		if r.Verbose {
			status := "PASS"
			if !result.Success {
				status = "FAIL"
			}
			fmt.Printf("  %s - %v\n", status, result.Duration)
			if result.Error != "" {
				fmt.Printf("    Error: %s\n", result.Error)
			}
		}
	}

	endTime := time.Now()
	report.EndTime = endTime
	report.Duration = endTime.Sub(startTime)
	report.Summary = r.calculateSummary(report.Results)

	return report, nil
}

// runSingleTest executes a single benchmark test
func (r *BenchmarkRunner) runSingleTest(ctx context.Context, test Test) BenchmarkResult {
	timeout := r.Timeout
	if test.Timeout != nil {
		timeout = *test.Timeout
	}

	testCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()

	// Prepare request
	queryReq := QueryRequest{
		AQL:        test.AQL,
		Parameters: test.Parameters,
	}

	reqBody, err := json.Marshal(queryReq)
	if err != nil {
		return BenchmarkResult{
			Test:     test,
			Duration: time.Since(start),
			Success:  false,
			Error:    fmt.Sprintf("Failed to marshal request: %v", err),
		}
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/openehr/v1/query", strings.TrimSuffix(r.APIURL, "/"))
	req, err := http.NewRequestWithContext(testCtx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return BenchmarkResult{
			Test:     test,
			Duration: time.Since(start),
			Success:  false,
			Error:    fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return BenchmarkResult{
			Test:     test,
			Duration: time.Since(start),
			Success:  false,
			Error:    fmt.Sprintf("HTTP request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return BenchmarkResult{
			Test:         test,
			Duration:     duration,
			Success:      false,
			Error:        fmt.Sprintf("Failed to read response: %v", err),
			StatusCode:   resp.StatusCode,
			ResponseSize: len(respBody),
		}
	}

	result := BenchmarkResult{
		Test:         test,
		Duration:     duration,
		StatusCode:   resp.StatusCode,
		ResponseSize: len(respBody),
	}

	// Check if request was successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Success = true

		// Try to parse response and count rows
		var queryResp QueryResponse
		if err := json.Unmarshal(respBody, &queryResp); err == nil {
			rowCount := len(queryResp.Rows)
			result.RowCount = &rowCount
		}
	} else {
		result.Success = false
		result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return result
}

// calculateSummary computes aggregate statistics from results
func (r *BenchmarkRunner) calculateSummary(results []BenchmarkResult) Summary {
	summary := Summary{
		TotalTests: len(results),
	}

	if len(results) == 0 {
		return summary
	}

	var totalDuration time.Duration
	summary.MinDuration = results[0].Duration
	summary.MaxDuration = results[0].Duration

	for _, result := range results {
		if result.Success {
			summary.PassedTests++
		} else {
			summary.FailedTests++
		}

		totalDuration += result.Duration
		if result.Duration < summary.MinDuration {
			summary.MinDuration = result.Duration
		}
		if result.Duration > summary.MaxDuration {
			summary.MaxDuration = result.Duration
		}
	}

	summary.AverageDuration = totalDuration / time.Duration(len(results))
	return summary
}

// loadTestSuite loads a test suite from a JSON file
func loadTestSuite(filename string) (TestSuite, error) {
	var suite TestSuite

	data, err := os.ReadFile(filename)
	if err != nil {
		return suite, fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, &suite); err != nil {
		return suite, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return suite, nil
}

// getDefaultTestSuite returns a built-in test suite
func getDefaultTestSuite() TestSuite {
	return TestSuite{
		Name:        "Default AQL Benchmarks",
		Description: "Built-in performance tests for AQL queries",
		Tests: []Test{
			{
				Name:        "EHR - Tiny",
				Description: "Very small EHR result set for baseline performance",
				AQL:         "SELECT e FROM EHR e LIMIT 5",
			},
			{
				Name:        "EHR - Small",
				Description: "Small EHR result set",
				AQL:         "SELECT e FROM EHR e LIMIT 25",
			},
			{
				Name:        "EHR - Medium",
				Description: "Medium EHR result set",
				AQL:         "SELECT e FROM EHR e LIMIT 100",
			},
			{
				Name:        "Composition - Small",
				Description: "Small composition result set",
				AQL:         "SELECT c FROM COMPOSITION c LIMIT 10",
			},
			{
				Name:        "EHR Filter",
				Description: "EHR query with simple condition",
				AQL:         "SELECT e FROM EHR e WHERE e/ehr_id/value EXISTS LIMIT 10",
			},
		},
	}
}

// outputResults saves benchmark results in the specified format
func outputResults(report *BenchmarkReport, outputFile, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return outputJSON(report, outputFile)
	case "junit":
		return outputJUnit(report, outputFile)
	case "both":
		if err := outputJSON(report, getOutputFilename(outputFile, "json")); err != nil {
			return err
		}
		return outputJUnit(report, getOutputFilename(outputFile, "xml"))
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// outputJSON saves results as JSON
func outputJSON(report *BenchmarkReport, outputFile string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if outputFile == "" {
		fmt.Println(string(data))
		return nil
	}

	return os.WriteFile(outputFile, data, 0644)
}

// outputJUnit saves results as JUnit XML
func outputJUnit(report *BenchmarkReport, outputFile string) error {
	testSuite := JUnitTestSuite{
		Name:     report.Suite.Name,
		Tests:    len(report.Results),
		Failures: report.Summary.FailedTests,
		Time:     fmt.Sprintf("%.3f", report.Duration.Seconds()),
	}

	for _, result := range report.Results {
		testCase := JUnitTestCase{
			Name:      result.Test.Name,
			ClassName: "AQLBenchmark",
			Time:      fmt.Sprintf("%.3f", result.Duration.Seconds()),
		}

		if !result.Success {
			testCase.Failure = &JUnitFailure{
				Message: "Test failed",
				Text:    result.Error,
			}
		}

		testSuite.TestCases = append(testSuite.TestCases, testCase)
	}

	testSuites := JUnitTestSuites{
		Name:     "AQL Benchmarks",
		Tests:    len(report.Results),
		Failures: report.Summary.FailedTests,
		Time:     fmt.Sprintf("%.3f", report.Duration.Seconds()),
		Suites:   []JUnitTestSuite{testSuite},
	}

	data, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %w", err)
	}

	xmlData := []byte(xml.Header + string(data))

	if outputFile == "" {
		fmt.Println(string(xmlData))
		return nil
	}

	return os.WriteFile(outputFile, xmlData, 0644)
}

// getOutputFilename generates output filename with extension
func getOutputFilename(base, ext string) string {
	if base == "" {
		return ""
	}

	if filepath.Ext(base) == "" {
		return base + "." + ext
	}

	// Replace extension
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return name + "." + ext
}

// printSummary prints a human-readable summary to stdout
func printSummary(report *BenchmarkReport) {
	fmt.Printf("\n=== Benchmark Results ===\n")
	fmt.Printf("Suite: %s\n", report.Suite.Name)
	fmt.Printf("Total Duration: %v\n", report.Duration)
	fmt.Printf("Tests: %d total, %d passed, %d failed\n",
		report.Summary.TotalTests,
		report.Summary.PassedTests,
		report.Summary.FailedTests)

	if report.Summary.TotalTests > 0 {
		fmt.Printf("Performance: avg=%v, min=%v, max=%v\n",
			report.Summary.AverageDuration,
			report.Summary.MinDuration,
			report.Summary.MaxDuration)
	}

	// Print failed tests
	if report.Summary.FailedTests > 0 {
		fmt.Printf("\nFailed Tests:\n")
		for _, result := range report.Results {
			if !result.Success {
				fmt.Printf("  - %s: %s\n", result.Test.Name, result.Error)
			}
		}
	}

	fmt.Printf("\n")
}

func main() {
	var (
		apiURL     = flag.String("url", "http://localhost:8080", "Base URL of the OpenEHR API")
		suiteFile  = flag.String("suite", "", "Path to test suite JSON file (optional)")
		outputFile = flag.String("output", "", "Output file for results (optional)")
		format     = flag.String("format", "json", "Output format: json, junit, or both")
		timeout    = flag.Duration("timeout", 30*time.Second, "Default timeout for each test")
		verbose    = flag.Bool("verbose", false, "Enable verbose output")
		singleAQL  = flag.String("aql", "", "Single AQL query to benchmark")
		testName   = flag.String("name", "Single AQL Test", "Name for single AQL test")
	)
	flag.Parse()

	ctx := context.Background()

	// Validate inputs
	if *apiURL == "" {
		fmt.Fprintf(os.Stderr, "Error: API URL is required\n")
		os.Exit(1)
	}

	var suite TestSuite
	var err error

	if *singleAQL != "" {
		// Create a single test from command line argument
		suite = TestSuite{
			Name:        "Command Line Test",
			Description: "Single AQL query provided via command line",
			Tests: []Test{
				{
					Name:        *testName,
					Description: "AQL query provided via --aql flag",
					AQL:         *singleAQL,
					Timeout:     timeout,
				},
			},
		}
	} else if *suiteFile != "" {
		// Load test suite from file
		suite, err = loadTestSuite(*suiteFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading test suite: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Use default built-in suite
		suite = getDefaultTestSuite()
	}

	if len(suite.Tests) == 0 {
		fmt.Fprintf(os.Stderr, "No tests to run\n")
		os.Exit(1)
	}

	// Run benchmarks
	runner := &BenchmarkRunner{
		APIURL:  *apiURL,
		Timeout: *timeout,
		Verbose: *verbose,
	}

	report, err := runner.RunBenchmarks(ctx, suite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running benchmarks: %v\n", err)
		os.Exit(1)
	}

	// Output results
	if err := outputResults(report, *outputFile, *format); err != nil {
		fmt.Fprintf(os.Stderr, "Error outputting results: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	printSummary(report)

	// Exit with error code if any tests failed
	if report.Summary.FailedTests > 0 {
		os.Exit(1)
	}
}
