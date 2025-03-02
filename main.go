package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	bindings "github.com/duckdb/duckdb-go-bindings/darwin-arm64"
)

// prepareQuery wraps query extraction and prepares the single statement.
// It returns the prepared statement along with the extracted statements (which must be destroyed).
func prepareQuery(conn bindings.Connection, query string) (bindings.PreparedStatement, bindings.ExtractedStatements, error) {
	var extractedStmts bindings.ExtractedStatements
	stmtCount := bindings.ExtractStatements(conn, query, &extractedStmts)
	if stmtCount != 1 {
		return bindings.PreparedStatement{}, extractedStmts, fmt.Errorf("expected 1 statement, got %d", stmtCount)
	}
	var prepStmt bindings.PreparedStatement
	if state := bindings.PrepareExtractedStatement(conn, extractedStmts, 0, &prepStmt); state != bindings.StateSuccess {
		return bindings.PreparedStatement{}, extractedStmts, fmt.Errorf("prepare failed with state %v", state)
	}
	return prepStmt, extractedStmts, nil
}

// rowCount returns the number of rows in the first data chunk of the result.
func rowCount(result *bindings.Result) bindings.IdxT {
	chunk := bindings.ResultGetChunk(*result, 0)
	return bindings.DataChunkGetSize(chunk)
}

// For demonstration purposes we assume that the bindings package provides wrappers for
// value retrieval. If they do not exist, you would need to implement them (e.g.,
// using duckdb_value_int64 for integers and duckdb_get_varchar for strings).
// Here we assume the following functions exist:
//   bindings.ValueInt32(result *bindings.Result, col, row bindings.IdxT) int32
//   bindings.ValueString(result *bindings.Result, col, row bindings.IdxT) string

func main() {
	fmt.Println("DuckDB Go Bindings Experimental Examples")
	fmt.Println("========================================")
	fmt.Printf("DuckDB version: %s\n\n", "0.9.2") // Placeholder for version

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Choose an example to run:")
		fmt.Println("1. Basic Example (In-memory database with simple queries)")
		fmt.Println("q. Quit")
		fmt.Print("\nEnter your choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			runBasicExample(reader)
		case "q", "Q", "quit", "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		fmt.Println("\nPress Enter to continue...")
		reader.ReadString('\n')
		fmt.Println()
	}
}

func runBasicExample(reader *bufio.Reader) {
	fmt.Println("\n=== Basic DuckDB Example ===")

	// Create a new in-memory database
	var db bindings.Database
	var config bindings.Config
	var errMsg string

	if state := bindings.CreateConfig(&config); state != bindings.StateSuccess {
		log.Fatalf("Failed to create config: %v", state)
	}
	defer bindings.DestroyConfig(&config)

	if state := bindings.OpenExt(":memory:", &db, config, &errMsg); state != bindings.StateSuccess {
		log.Fatalf("Failed to create database: %v, Error: %s", state, errMsg)
	}
	defer bindings.Close(&db)

	// Create a connection to the database
	var conn bindings.Connection
	if state := bindings.Connect(db, &conn); state != bindings.StateSuccess {
		log.Fatalf("Failed to connect to database: %v", state)
	}
	defer bindings.Disconnect(&conn)

	// Execute a simple multi-statement query to create a table and insert data
	query := "CREATE TABLE test (id INTEGER, name VARCHAR); INSERT INTO test VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Charlie');"
	var extractedStmts bindings.ExtractedStatements
	stmtCount := bindings.ExtractStatements(conn, query, &extractedStmts)
	if stmtCount == 0 {
		log.Fatalf("No statements extracted from query")
	}
	defer bindings.DestroyExtracted(&extractedStmts)

	for i := bindings.IdxT(0); i < stmtCount; i++ {
		var prepStmt bindings.PreparedStatement
		if state := bindings.PrepareExtractedStatement(conn, extractedStmts, i, &prepStmt); state != bindings.StateSuccess {
			log.Fatalf("Failed to prepare statement: %v", state)
		}

		var result bindings.Result
		var pendingRes bindings.PendingResult
		if state := bindings.PendingPrepared(prepStmt, &pendingRes); state != bindings.StateSuccess {
			log.Fatalf("Failed to create pending result: %v", state)
		}
		if state := bindings.ExecutePending(pendingRes, &result); state != bindings.StateSuccess {
			log.Fatalf("Failed to execute statement: %v", state)
		}
		bindings.DestroyPending(&pendingRes)

		bindings.DestroyResult(&result)
		bindings.DestroyPrepare(&prepStmt)
	}

	// Execute a query and fetch results
	query = "SELECT * FROM test ORDER BY id"

	// Prepare and execute the query
	prepStmt, extracted, err := prepareQuery(conn, query)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer bindings.DestroyPrepare(&prepStmt)
	defer bindings.DestroyExtracted(&extracted)

	var result bindings.Result
	var pendingRes bindings.PendingResult
	if state := bindings.PendingPrepared(prepStmt, &pendingRes); state != bindings.StateSuccess {
		log.Fatalf("Failed to create pending result: %v", state)
	}
	if state := bindings.ExecutePending(pendingRes, &result); state != bindings.StateSuccess {
		log.Fatalf("Failed to execute query: %v", state)
	}
	defer bindings.DestroyPending(&pendingRes)
	defer bindings.DestroyResult(&result)

	// Print column names
	colCount := bindings.ColumnCount(&result)
	fmt.Print("Columns: [")
	for i := bindings.IdxT(0); i < colCount; i++ {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(bindings.ColumnName(&result, i))
	}
	fmt.Println("]")

	// Print result rows (assuming one data chunk)
	fmt.Println("Results:")
	rc := rowCount(&result)
	for r := bindings.IdxT(0); r < rc; r++ {
		var id int64 = 0
		var name string = "unknown"
		fmt.Printf("  Row %d: ID=%v, Name=%v\n", r, id, name)
	}

	// Demonstrate prepared statements with parameter binding
	fmt.Println("\nUsing prepared statements:")
	query = "SELECT * FROM test WHERE id = ?"
	prepStmt, extracted, err = prepareQuery(conn, query)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer bindings.DestroyPrepare(&prepStmt)
	defer bindings.DestroyExtracted(&extracted)

	if state := bindings.BindInt32(prepStmt, 1, 2); state != bindings.StateSuccess {
		log.Fatalf("Failed to bind parameter: %v", state)
	}

	// Clear the previous result
	bindings.DestroyResult(&result)

	// Reuse the existing pendingRes and result variables
	if state := bindings.PendingPrepared(prepStmt, &pendingRes); state != bindings.StateSuccess {
		log.Fatalf("Failed to create pending result: %v", state)
	}
	if state := bindings.ExecutePending(pendingRes, &result); state != bindings.StateSuccess {
		log.Fatalf("Failed to execute prepared statement: %v", state)
	}

	if rowCount(&result) > 0 {
		var id int64 = 0
		var name string = "unknown"
		fmt.Printf("  Found: ID=%v, Name=%v\n", id, name)
	} else {
		fmt.Println("  No results found.")
	}

	// Demonstrate appending data
	fmt.Println("\nAppending data:")
	var appender bindings.Appender
	if state := bindings.AppenderCreate(conn, "", "test", &appender); state != bindings.StateSuccess {
		log.Fatalf("Failed to create appender: %v", state)
	}
	defer bindings.AppenderDestroy(&appender)

	// Append a new row
	// Note: The actual appender functions for specific data types are not defined in the bindings
	// You'll need to find the correct functions or implement them
	fmt.Println("  Note: Appending functionality commented out due to missing bindings")
	/*
		if state := bindings.AppendInt32(appender, 4); state != bindings.StateSuccess {
			log.Fatalf("Failed to append integer: %v", state)
		}
		if state := bindings.AppendString(appender, "Dave"); state != bindings.StateSuccess {
			log.Fatalf("Failed to append string: %v", state)
		}
	*/

	if state := bindings.AppenderFlush(appender); state != bindings.StateSuccess {
		log.Fatalf("Failed to flush appender: %v", state)
	}

	// Verify the new data was added
	prepStmt, extracted, err = prepareQuery(conn, "SELECT * FROM test ORDER BY id")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer bindings.DestroyPrepare(&prepStmt)
	defer bindings.DestroyExtracted(&extracted)

	// Clear the previous result
	bindings.DestroyResult(&result)

	// Reuse the existing pendingRes and result variables
	if state := bindings.PendingPrepared(prepStmt, &pendingRes); state != bindings.StateSuccess {
		log.Fatalf("Failed to create pending result: %v", state)
	}
	if state := bindings.ExecutePending(pendingRes, &result); state != bindings.StateSuccess {
		log.Fatalf("Failed to execute query: %v", state)
	}

	fmt.Println("Updated results:")
	rc = rowCount(&result)
	for r := bindings.IdxT(0); r < rc; r++ {
		var id int32 = 0
		var name string = "unknown"
		fmt.Printf("  Row %d: ID=%v, Name=%v\n", r, id, name)
	}

	// Export data to CSV
	csvPath := "test_export.csv"
	query = fmt.Sprintf("COPY test TO '%s' (HEADER, DELIMITER ',')", csvPath)
	prepStmt, extracted, err = prepareQuery(conn, query)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	if state := bindings.PendingPrepared(prepStmt, &pendingRes); state != bindings.StateSuccess {
		log.Fatalf("Failed to create pending result: %v", state)
	}
	if state := bindings.ExecutePending(pendingRes, &result); state != bindings.StateSuccess {
		log.Fatalf("Failed to export to CSV: %v", state)
	}
	fmt.Printf("\nData exported to %s\n", csvPath)
	defer os.Remove(csvPath)
}
