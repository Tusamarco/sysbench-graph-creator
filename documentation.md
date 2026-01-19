# Sysbench Graph Creator - Technical Documentation

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Core Components](#core-components)
4. [Data Structures](#data-structures)
5. [Processing Pipeline](#processing-pipeline)
6. [Configuration](#configuration)
7. [Usage Guide](#usage-guide)
8. [API Reference](#api-reference)
9. [Development](#development)

---

## Overview

Sysbench Graph Creator is a Go-based application designed to parse, analyze, and visualize performance metrics from sysbench benchmark outputs. The tool transforms raw sysbench test results into interactive HTML charts, static images (JPEG), and CSV data exports.

### Key Features

- **Automated File Processing**: Recursively scans directories for sysbench output files (.csv, .txt)
- **Statistical Analysis**: Calculates standard deviation and distance metrics across multiple test runs
- **Data Filtering**: Supports filtering by producer, version, dimension, test name, and pre/post-write phases
- **Multiple Output Formats**:
  - Interactive HTML charts with embedded graphs
  - Static JPEG images
  - CSV data exports
- **HTTP Server**: Built-in web server to view results locally
- **Progress Tracking**: Real-time progress bars during file processing

### Technology Stack

- **Language**: Go 1.23+
- **Chart Library**: go-echarts (ECharts-based charting)
- **Statistics**: montanaflynn/stats
- **Logging**: sirupsen/logrus
- **Configuration**: TOML format
- **Chart Rendering**: chromedp for headless browser rendering

---

## Architecture

The application follows a modular architecture with clear separation of concerns:

```
sysbench_graph_creator.go (Main Entry Point)
    ├── internal/global/          (Global utilities and configuration)
    │   ├── configuration.go      (TOML config parsing)
    │   ├── help.go               (Help text)
    │   ├── logUtils.go           (Logging setup)
    │   └── utils.go              (Helper functions)
    │
    └── internal/dataObjects/     (Core business logic)
        ├── FileProcessor.go      (File scanning and parsing)
        ├── tests.go              (Data structure definitions)
        ├── Calculator.go         (Statistical calculations)
        ├── GraphMaker.go         (Chart generation)
        ├── bar.go                (Bar chart specifics)
        └── snapshot.go           (Screenshot functionality)
```

### Data Flow

1. **Configuration Loading**: Parse TOML config file and command-line parameters
2. **File Discovery**: Recursively scan source directory for sysbench output files
3. **Data Parsing**: Extract metadata and metrics from each file
4. **Data Aggregation**: Group tests by producer, version, dimension, and action type
5. **Statistical Analysis**: Calculate averages, standard deviations, and error margins
6. **Graph Generation**: Create interactive charts using go-echarts
7. **Output Rendering**: Export to HTML, JPEG, and/or CSV formats
8. **HTTP Server**: Optionally serve results via local web server

---

## Core Components

### 1. Main Application (`sysbench_graph_creator.go`)

The entry point that orchestrates the entire workflow:

- **Flag Parsing**: Handles command-line arguments
- **Configuration Management**: Loads and merges config from file and CLI
- **Initialization**: Sets up logging and validates paths
- **Execution Flow**:
  1. Get file list from source directory
  2. Parse test collections
  3. Calculate statistical results
  4. Generate graphs
  5. Build HTML page
  6. Start HTTP server (optional)

### 2. File Processor (`FileProcessor.go`)

Responsible for discovering and parsing sysbench output files.

**Key Functions**:

- `GetFileList(path)`: Recursively walks directories to find .csv and .txt files
- `GetTestCollectionArray()`: Parses each file into TestCollection objects
- `getTestCollectionData(path)`: Extracts metadata and test data from a single file
- `identifyEndTime()`: Determines the end time of test collections
- `identifyPerPostWrite()`: Distinguishes pre-write and post-write test phases

**File Format Expected**:

Files should contain:
- `META:` lines with semicolon-separated key-value pairs for collection metadata
- `SUBTEST:` sections for individual test runs
- `TEST SUMMARY:` sections with comma-separated metric values

Example:
```
META: testIdentifyer=PS8042_iron_ssd2;dimension=large;actionType=select;runNumber=1;host=10.30.12.4;producer=sysbench;execDate=2024-02-02_12_12_27;engine=innodb;mysqlproducer=Percona Server;mysqlversion=8.0.42
SUBTEST:select_run_inlist
BLOCK: [START] 2024-02-02_12_12_27 (filter: select)
META: testIdentifyer=PS8042_iron_ssd2;dimension=large;actionType=select;runNumber=1;execCommand=run;subtest=select_run_inlist;execDate=2024-02-02_12_12_27;engine=innodb
TEST SUMMARY:
TotalTime,RunningThreads,totalEvents,Events/s,Tot Operations,operations/s,tot reads,reads/s,Tot writes,writes/s,oterOps/s,latencyPct95(μs),Tot errors,errors/s,Tot reconnects,reconnects/s,Latency(ms) min,Latency(ms) max,Latency(ms) avg,Latency(ms) sum
200,1,2642.00,13.21,2642.00,13.21,2642.00,13.21,0.00,0.00,0.00,137.35,0.00,0.00,0.00,0.00,0.04,0.22,0.08,200.00
```

### 3. Calculator (`Calculator.go`)

Performs statistical analysis on collected test data.

**Key Functions**:

- `Init(configuration)`: Initializes the calculator with config
- `BuildResults(testCollections)`: Main orchestrator for building result sets
- `loopCollections()`: Groups collections by matching criteria
- `loopTests()`: Processes tests within collections
- `calculateResults()`: Computes averages and statistics
- `calculateSTD()`: Calculates standard deviation for data sets
- `GroupByProducers()`: Organizes results by database producer

**Statistical Metrics Calculated**:
- Mean values for each metric
- Standard deviation (STD)
- Distance/Error percentage (Gerror)

### 4. Graph Generator (`GraphMaker.go`)

Creates visual representations of benchmark data.

**Key Functions**:

- `Init(config, producers, collections)`: Initializes the graph generator
- `RenderResults()`: Creates all charts based on configured labels
- `BuildPage()`: Generates the HTML page with embedded charts
- `ActivateHTTPServer()`: Starts the HTTP server to view results
- `buildGenericColumn()`: Creates column/bar charts
- `buildSummaryChart()`: Builds aggregate summary charts
- `convertToCSV()`: Exports chart data to CSV format

**Chart Types**:
- Column charts (default)
- Line charts (configurable)
- Summary charts (aggregated metrics)

**Output Formats**:
- HTML with interactive charts
- JPEG static images (using chromedp for rendering)
- CSV data exports

---

## Data Structures

### TestCollection

Represents a complete set of tests from a single file.

```go
type TestCollection struct {
    DateStart        time.Time         // When the test collection started
    DateEnd          time.Time         // When the test collection ended
    Dimension        string            // "small" or "large"
    ExecutionTime    int64             // Total execution time in minutes
    TestName         string            // Identifier for the test
    Producer         string            // Test producer (sysbench, tpcc, dbt3)
    Tests            map[string]Test   // Map of individual tests
    ActionType       int               // READ=0, WRITE=1, READ_AND_WRITE=10
    SelectPostWrites int               // PREWRITE=0, POSTWRITE=1
    HostDB           string            // Database host
    RunNumber        int               // Run iteration number
    Engine           string            // Storage engine (e.g., "innodb")
    Name             string            // File-based name
    MySQLVersion     string            // MySQL version string
    MySQLProducer    string            // MySQL producer (e.g., "Percona Server", "MySQL")
    FileName         string            // Source file name
}
```

### Test

Represents an individual test within a collection.

```go
type Test struct {
    Name          string                 // Test name
    DateStart     time.Time              // Test start time
    DateEnd       time.Time              // Test end time
    Dimension     string                 // "small" or "large"
    ExecutionTime int64                  // Test execution time
    Labels        []string               // Metric labels
    Threads       []int                  // Thread counts used
    ThreadExec    map[int]Execution      // Execution results by thread count
    ActionType    int                    // Test action type
    Filter        string                 // Filter applied
    RunNumber     int                    // Run number
}
```

### Execution

Contains the actual benchmark results for a specific thread count.

```go
type Execution struct {
    Thread    int                   // Number of threads
    Command   string                // Command executed
    Result    map[string]float64    // Metric name -> value mapping
    DateStart time.Time             // Execution start
    DateEnd   time.Time             // Execution end
    Processed bool                  // Processing status flag
}
```

### ResultTest

Aggregated results across multiple test runs.

```go
type ResultTest struct {
    Key        TestKey                         // Unique identifier
    Labels     map[string][]ResultValue        // Results organized by label
    Executions int                             // Number of executions
    STD        float64                         // Standard deviation
    Gerror     float64                         // Error percentage
}
```

### Producer

Aggregates all tests for a specific MySQL producer and version.

```go
type Producer struct {
    MySQLProducer       string         // Producer name
    MySQLVersion        string         // Version string
    TestsResults        []ResultTest   // All test results
    TestsTypes          []TestType     // Test type definitions
    TestCollectionsName string         // Collection name
    STDReadPre          float64        // STD for pre-write reads
    GerrorReadPre       float64        // Error for pre-write reads
    STDReadPost         float64        // STD for post-write reads
    GerrorReadPost      float64        // Error for post-write reads
    STDRWrite           float64        // STD for writes
    GerrorWrite         float64        // Error for writes
    Color               string         // Chart color
}
```

---

## Processing Pipeline

### Phase 1: Initialization

1. **Parse Command-Line Arguments**
   - Required: `--configfile`, `--configpath`
   - Optional: Various filter and output options

2. **Load Configuration**
   - Read TOML configuration file
   - Merge with command-line parameters (CLI takes precedence)

3. **Sanity Checks**
   - Verify source path exists
   - Validate filter combinations
   - Check conflicting options

4. **Path Setup**
   - Create destination directories if needed
   - Set up HTML, CSV, and image output paths

5. **Logging Initialization**
   - Configure log level and target (stdout/file)

### Phase 2: File Processing

1. **File Discovery**
   - Recursively walk source directory
   - Collect paths to .csv and .txt files
   - Skip warmup files (`_warmup_`)

2. **Parse Each File**
   - Extract collection metadata from `META:` lines
   - Parse individual test sections (`SUBTEST:`)
   - Read test execution data from `TEST SUMMARY:` sections
   - Build TestCollection objects

3. **Post-Processing**
   - Identify end times for all collections
   - Determine pre-write vs. post-write phases
   - Apply naming filters

### Phase 3: Statistical Analysis

1. **Group Collections**
   - Match collections by:
     - Test name
     - Dimension (small/large)
     - MySQL producer and version
     - Action type
     - Pre/post-write phase

2. **Aggregate Test Runs**
   - Combine multiple runs of the same test
   - Group by thread count

3. **Calculate Statistics**
   - Compute mean values for each metric
   - Calculate standard deviation
   - Compute error percentages (distance from mean)

4. **Filter Results**
   - Apply configured filters:
     - Producer filter
     - Version filter
     - Dimension filter
     - Title filter
     - Pre/post-write filter

5. **Group by Producers**
   - Organize results by MySQL producer
   - Calculate aggregate statistics per producer

### Phase 4: Visualization

1. **Prepare Chart Data**
   - Extract configured labels (metrics to chart)
   - Create data series for each producer
   - Sort by thread count

2. **Generate Charts**
   - Create column/line charts for each metric
   - Build summary charts (reads/writes)
   - Apply color schemes

3. **Export Outputs**
   - **HTML**: Generate interactive chart page
   - **JPEG** (if enabled): Render charts to images using headless Chrome
   - **CSV** (if enabled): Export raw data to CSV files

4. **Start HTTP Server** (if enabled)
   - Serve generated HTML page
   - Default: http://localhost:8089

---

## Configuration

### TOML Configuration File

Configuration is stored in a TOML file with three main sections:

#### [parser] Section

Controls how files are parsed and processed.

```toml
[parser]
sourceDataPath = "/path/to/sysbench/output"  # Source directory
filterOutliners = true                        # Enable outlier filtering
distanceLabel = "operations/s"                # Metric for distance calculation
```

#### [render] Section

Controls output generation and rendering.

```toml
[render]
graphType = "column"                          # Chart type: "column" or "line"
destinationPath = "/path/to/results"          # Main output directory
csvDestinationPath = "/path/to/csv"           # CSV output directory
printStats = false                            # Print statistical summaries
printData = false                             # Print raw data
printCharts = false                           # Generate JPEG charts
printChartsFormat = "jpeg"                    # Image format
printChartsQuality = 5                        # JPEG quality (1-10)
convertChartsToCsv = true                     # Export to CSV

httpServerPort = 8089                         # HTTP server port
httpServerIp = "localhost"                    # HTTP server IP

# Metrics to chart (comma-separated)
labels = "TotalTime,Events/s,operations/s,writes/s,reads/s,latencyPct95(μs)"
statslabels = "operations/s,latencyPct95(μs)"
readSummaryLabel = "reads/s"                  # Summary metric for reads
writeSummaryLabel = "writes/s"                # Summary metric for writes

chartHeight = 700                             # Chart height in pixels
chartWidth = 1200                             # Chart width in pixels

# Filters
filterTestsByTitle = ""                       # Include only these tests
filterExcludeByTitle = ""                     # Exclude these tests
filterByDimension = "small,large"             # Filter by dimension
filterByVersion = "8.0.37"                    # Filter by MySQL version
filterByProducer = ""                         # Filter by MySQL producer
filterByPrePost = "post"                      # Filter by phase: "pre", "post", or "pre,post"
```

#### [global] Section

Global application settings.

```toml
[global]
testName = "Comparing PS VS MySQL"            # Overall test name/title
logLevel = "info"                             # Log level: debug, info, warn, error
logTarget = "stdout"                          # Log target: "stdout" or "file"
logFile = "/tmp/sysbench_graph_creator.log"   # Log file path
performance = false                           # Performance mode
```

#### [colors] Section

Customize chart colors.

```toml
[colors]
color = ["orange", "blue", "green", "red", "purple"]
```

### Command-Line Parameters

Command-line parameters override TOML configuration values.

**Required**:
- `--configfile`: Config file name
- `--configpath`: Config file directory path

**Optional**:
- `--sourceDataPath`: Source directory (overrides config)
- `--destinationPath`: Destination directory (overrides config)
- `--csvDestinationPath`: CSV destination (overrides config)
- `--testName`: Test name (overrides config)
- `--labels`: Comma-separated list of metrics to chart
- `--printCharts`: Generate JPEG images (true/false)
- `--printData`: Generate HTML output (true/false)
- `--convertCsv`: Export to CSV (true/false)
- `--filterByProducer`: Filter by producer names
- `--filterByVersion`: Filter by version strings
- `--filterByDimension`: Filter by dimension (small/large)
- `--filterByTitle`: Include only matching test names
- `--filterExcludeByTitle`: Exclude matching test names
- `--filterByPrePost`: Filter by pre/post-write phase

---

## Usage Guide

### Basic Usage

1. **Prepare Configuration File**

Create a TOML config file (e.g., `config/mytest.toml`):

```toml
[parser]
sourceDataPath = "/data/sysbench/results"
filterOutliners = true

[render]
destinationPath = "/data/sysbench/graphs"
labels = "operations/s,latencyPct95(μs)"
printData = true

[global]
testName = "MySQL Performance Test"
```

2. **Run the Application**

```bash
./sysbench_graph_creator --configfile=mytest.toml --configpath=./config
```

3. **View Results**

If HTTP server is enabled, navigate to:
```
http://localhost:8089
```

### Advanced Usage Examples

#### Example 1: Generate Only CSV Output

```bash
./sysbench_graph_creator \
  --configfile=mytest.toml \
  --configpath=./config \
  --convertCsv=true \
  --printData=false
```

#### Example 2: Filter by Specific MySQL Version

```bash
./sysbench_graph_creator \
  --configfile=mytest.toml \
  --configpath=./config \
  --filterByVersion="8.0.37"
```

#### Example 3: Compare Pre and Post-Write Performance

```bash
# Only pre-write tests
./sysbench_graph_creator \
  --configfile=mytest.toml \
  --configpath=./config \
  --filterByPrePost="pre"

# Only post-write tests
./sysbench_graph_creator \
  --configfile=mytest.toml \
  --configpath=./config \
  --filterByPrePost="post"
```

#### Example 4: Generate Static Images

```bash
./sysbench_graph_creator \
  --configfile=mytest.toml \
  --configpath=./config \
  --printCharts=true
```

#### Example 5: Focus on Specific Tests

```bash
./sysbench_graph_creator \
  --configfile=mytest.toml \
  --configpath=./config \
  --filterByTitle="select_run_range,write_run_update"
```

### Output Files

After execution, the destination directory will contain:

```
destinationPath/
├── html/
│   └── index.html              # Interactive charts page
├── images/                     # (if printCharts=true)
│   ├── chart_operations_per_s.jpeg
│   ├── chart_latency.jpeg
│   └── ...
└── csv/                        # (if convertChartsToCsv=true)
    ├── operations_per_s.csv
    ├── latency.csv
    └── ...
```

---

## API Reference

### Package: internal/global

#### Configuration

**GetConfig(path string) Configuration**

Loads and parses a TOML configuration file.

- **Parameters**: 
  - `path`: Full path to TOML config file
- **Returns**: `Configuration` struct
- **Exits on error**: Calls `syscall.Exit(2)` if parsing fails

**ParseCommandLine(params Params)**

Merges command-line parameters into configuration. CLI values override config file values.

**CheckPaths()**

Validates and creates necessary output directories.

**SanityChecks()**

Performs validation on configuration parameters. Exits if critical issues found.

#### Utilities

**InitLog(config Configuration) bool**

Initializes the logging system based on configuration.

- **Returns**: `true` if successful, `false` otherwise

**CheckIfPathExists(path string) bool**

Checks if a path exists on the filesystem.

**CreatePath(path string) bool**

Creates a directory path recursively.

**LineCount(path string) (int, error)**

Counts lines in a file (used for progress bars).

**ParsetimeLocal(dateStr, path string) (time.Time, error)**

Parses date strings in various formats.

### Package: internal/dataObjects

#### FileProcessor

**GetFileList(path string) error**

Recursively scans directory for sysbench output files.

- **Parameters**: 
  - `path`: Root directory to scan
- **Returns**: Error if scan fails
- **Side effects**: Populates `arPathFiles` array

**GetTestCollectionArray() ([]TestCollection, error)**

Parses all discovered files into TestCollection objects.

- **Returns**: 
  - Array of `TestCollection` objects
  - Error if parsing fails
- **Side effects**: 
  - Processes files with progress bars
  - Identifies end times and pre/post-write phases

#### Calculator

**Init(configuration Configuration)**

Initializes calculator with configuration.

**BuildResults(testCollections []TestCollection) map[TestKey]ResultTest**

Main processing function. Groups collections, processes tests, calculates statistics.

- **Parameters**: 
  - `testCollections`: Array of parsed test collections
- **Returns**: Map of TestKey to ResultTest with calculated statistics

**GroupByProducers() []Producer**

Organizes results by MySQL producer and version.

- **Returns**: Array of Producer objects with aggregated statistics

#### GraphGenerator

**Init(config Configuration, producers []Producer, collections []TestCollection)**

Initializes graph generator with data and configuration.

**RenderResults() bool**

Generates all configured charts.

- **Returns**: `true` if successful

**BuildPage() bool**

Creates the HTML page with embedded charts.

- **Returns**: `true` if successful

**ActivateHTTPServer()**

Starts the HTTP server to serve generated HTML.

- **Blocks**: This function runs the server and blocks

**convertToCSV() bool**

Exports chart data to CSV files.

- **Returns**: `true` if successful

---

## Development

### Building the Project

```bash
# Clone the repository
git clone https://github.com/Tusamarco/sysbench-graph-creator.git
cd sysbench-graph-creator

# Install dependencies
go mod download

# Build the binary
go build .
```

### Project Structure

```
sysbench-graph-creator/
├── sysbench_graph_creator.go    # Main entry point
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── internal/
│   ├── dataObjects/              # Core data processing
│   │   ├── FileProcessor.go      # File parsing
│   │   ├── Calculator.go         # Statistical analysis
│   │   ├── GraphMaker.go         # Chart generation
│   │   ├── tests.go              # Data structures
│   │   ├── bar.go                # Bar chart utilities
│   │   ├── snapshot.go           # Screenshot functionality
│   │   └── dataFile.go           # File data structures
│   └── global/                   # Global utilities
│       ├── configuration.go      # Config management
│       ├── help.go               # Help text
│       ├── logUtils.go           # Logging utilities
│       └── utils.go              # Helper functions
├── config/                       # Sample config files
│   ├── config-dev-marco.toml
│   └── config-dev-marco_win.toml
└── uml/                          # Architecture diagrams
    ├── concept.puml
    ├── objects.puml
    └── *.png
```

### Key Dependencies

- **github.com/go-echarts/go-echarts/v2**: Chart generation
- **github.com/chromedp/chromedp**: Headless browser for image rendering
- **github.com/montanaflynn/stats**: Statistical calculations
- **github.com/sirupsen/logrus**: Logging framework
- **github.com/Tusamarco/toml**: TOML configuration parsing
- **github.com/schollz/progressbar**: Progress indication

### Extending the Application

#### Adding New Metrics

1. Ensure the metric is present in the `TEST SUMMARY:` section of input files
2. Add the metric name to the `labels` configuration
3. The system will automatically parse and chart it

#### Adding New Chart Types

1. Extend the `charTest` struct in `GraphMaker.go` if needed
2. Implement a new chart building function (e.g., `buildNewChartType()`)
3. Call the function in `RenderResults()`

#### Custom Filters

1. Add filter parameters to `Params` and `Render` structs in `configuration.go`
2. Implement filter logic in `Calculator.go` or `GraphMaker.go`
3. Add command-line flag parsing in `sysbench_graph_creator.go`

### Logging

The application uses structured logging with logrus. Log levels:

- **Debug**: Detailed processing information
- **Info**: General progress updates
- **Warn**: Non-critical issues
- **Error**: Critical errors

Configure via `[global]` section in TOML:

```toml
[global]
logLevel = "debug"
logTarget = "file"
logFile = "/tmp/sysbench_graph_creator.log"
```

### Error Handling

The application uses several strategies:

- **Immediate Exit**: Configuration errors, missing required parameters
- **Logged Errors**: Parsing issues, file access problems
- **Graceful Degradation**: Missing optional data, filter mismatches

Exit codes:
- `0`: Success
- `1`: General error (missing paths, validation failures)
- `2`: Configuration parsing error

---

## Troubleshooting

### Common Issues

**Issue**: "Source Path does not exist"
- **Solution**: Verify the `sourceDataPath` in your config file points to a valid directory containing sysbench output files.

**Issue**: No charts generated
- **Solution**: Check that:
  - Input files contain `META:` and `TEST SUMMARY:` sections
  - Files are not named with `_warmup_` (these are excluded)
  - Filters are not too restrictive

**Issue**: HTTP server doesn't start
- **Solution**: 
  - Check if port 8089 (or configured port) is already in use
  - Ensure `printData=true` is set
  - Check firewall settings

**Issue**: JPEG generation fails
- **Solution**: 
  - Ensure Chrome/Chromium is installed for headless rendering
  - Check system has sufficient memory
  - Review logs for chromedp errors

**Issue**: Incorrect statistics
- **Solution**: 
  - Verify all test runs have the same test names and dimensions
  - Check for outliers affecting calculations
  - Enable `filterOutliners=true` in config

### Debug Mode

Enable detailed logging:

```toml
[global]
logLevel = "debug"
logTarget = "file"
logFile = "/tmp/debug.log"
```

Then review the log file for detailed execution flow and data processing information.

---

## License

This project is licensed under the AGPL-3.0 License. See the LICENSE file for details.

Copyright (c) Marco Tusa 2021 - present

---

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/Tusamarco/sysbench-graph-creator/issues
- Pull Requests: https://github.com/Tusamarco/sysbench-graph-creator/pulls
