# DuckDB Go Bindings Example

This project demonstrates how to use the DuckDB Go bindings to interact with DuckDB databases from Go applications. It provides a simple interactive CLI that showcases basic DuckDB operations including creating tables, executing queries, using prepared statements, and exporting data.

## Overview

DuckDB is an in-process SQL OLAP database management system designed for analytical queries. This example application demonstrates how to:

- Create an in-memory DuckDB database
- Execute SQL statements
- Use prepared statements with parameters
- Fetch and display query results
- Append data to tables
- Export data to CSV files

## Requirements

- Go 1.18 or higher
- macOS with Apple Silicon (M1/M2/M3) for this specific example

## Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/yourusername/duck-duck-go.git
   cd duck-duck-go
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

Run the application:

```bash
go run main.go
```

The application presents a menu with the following options:

1. Basic Example - Demonstrates in-memory database with simple queries
2. Quit - Exit the application

## Implementation Details

The application uses the platform-specific DuckDB Go bindings for macOS ARM64:

```go
import (
    bindings "github.com/duckdb/duckdb-go-bindings/darwin-arm64"
)
```

Key components:

- **Database Creation**: Creates an in-memory DuckDB database
- **Query Execution**: Demonstrates how to execute SQL queries
- **Prepared Statements**: Shows how to use prepared statements with parameters
- **Result Processing**: Retrieves and displays query results
- **Data Export**: Exports table data to CSV format

## Limitations

This is an experimental implementation with some limitations:

- The value retrieval functions are not fully implemented, so actual data values are displayed as placeholders
- The advanced example option is not yet implemented
- Error handling could be more robust
- Platform-specific (currently only works on macOS with ARM64 architecture)

## Future Improvements

- Implement value retrieval functions to display actual data
- Add the advanced example functionality
- Improve error handling
- Add support for additional platforms
- Implement more DuckDB features (transactions, user-defined functions, etc.)

## License

MIT

## Acknowledgments

- [DuckDB Team](https://github.com/duckdb/duckdb) for creating DuckDB
- [DuckDB Go Bindings](https://github.com/duckdb/duckdb-go-bindings) for providing the Go bindings
