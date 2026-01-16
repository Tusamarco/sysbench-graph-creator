# Sysbench Graph Creator

Sysbench Graph Creator is a Go-based utility that generates graphs from the output of [sysbench](https://github.com/akopytov/sysbench), a popular benchmark tool for evaluating database systems' performance.

## Features
- Parses sysbench output to extract performance metrics.
- Generates visual graphs to represent benchmarking results.
- Supports customization via configuration files.
- Efficient and lightweight, designed for quick analysis of benchmark data.

## Requirements
- [Go](https://go.dev/) 1.16 or higher.
- Sysbench installed and configured on your system.

## Installation

### Clone the Repository
```sh
git clone https://github.com/Tusamarco/sysbench-graph-creator.git
cd sysbench-graph-creator
```

### Build the Project
```sh
go build .
```

## Usage

1. Run sysbench to generate benchmark data.
2. Parse the output using Sysbench Graph Creator.

```sh
./sysbench_graph_creator -input sysbench_output.txt -output graph.png
```

Replace `sysbench_output.txt` with the path to your sysbench output file and `graph.png` with your desired graph output file.

### Options
- `-input`: Path to the sysbench output file.
- `-output`: Desired output graph file name.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## License
This project is licensed under the AGPL-3.0 License. See the [LICENSE](./LICENSE) file for details.

---

For more detailed usage and examples, check the documentation or [raise an issue](https://github.com/Tusamarco/sysbench-graph-creator/issues) if you encounter any problems.