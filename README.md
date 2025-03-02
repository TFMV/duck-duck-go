# DuckDB Go Bindings Example

🚧 Warning: The DuckDB Go bindings are in early development. 🚧

Expect potential API changes and missing functionality.

This project demonstrates how to use the DuckDB Go bindings to interact with DuckDB databases from Go applications.
It provides a simple interactive CLI showcasing:

✅ Creating an in-memory DuckDB database
✅ Executing SQL queries
✅ Using prepared statements
✅ Retrieving query results
✅ Appending data to tables
✅ Exporting data to CSV files

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

- Database Creation → Creates an in-memory DuckDB instance
- Query Execution → Executes SQL queries directly
- Prepared Statements → Uses parameterized SQL for safe queries
- Result Processing → Retrieves query results
- Data Export → Saves tables as CSV

## Limitations

🚨 This is an early-stage implementation with limitations:

❌ Limited data retrieval – Values are placeholders, as full data retrieval functions aren’t implemented
❌ Basic error handling – Some edge cases may not be handled gracefully
❌ No Appending - Appends have not been implemented.

## Future Improvements

- Implement value retrieval functions to display actual data
- Appender Functionality
- Add the advanced example functionality
- Improve error handling
- Add support for additional platforms
- Implement more DuckDB features (transactions, user-defined functions, etc.)

## License

MIT

## Acknowledgments

- [DuckDB Team](https://github.com/duckdb/duckdb) for creating DuckDB
- [DuckDB Go Bindings](https://github.com/duckdb/duckdb-go-bindings) for providing the Go bindings
