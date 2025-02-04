package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"my-pp/modules/databases"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func main() {
	db, err := databases.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)

	currentVersion, err := readVersion(dir)
	if err != nil {
		panic(err)
	}

	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		version, err := parseVersion(file)
		if err != nil {
			panic(err)
		}

		if version > currentVersion {
			if err := applyMigration(db, file); err != nil {
				panic(err)
			}
			currentVersion = version
			if err := writeVersion(currentVersion, dir); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("All migrations applied successfully!")
}

func readVersion(dir string) (int, error) {
	file, err := os.Open(filepath.Join(dir, "version.txt"))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return strconv.Atoi(scanner.Text())
	}
	return 0, scanner.Err()
}

func writeVersion(version int, dir string) error {
	file, err := os.Create(filepath.Join(dir, "version.txt"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(version))
	return err
}

func parseVersion(file string) (int, error) {
	base := filepath.Base(file)
	return strconv.Atoi(base[:3])
}

func applyMigration(db *sql.DB, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to apply migration %s: %w", file, err)
	}

	fmt.Printf("Applied migration: %s\n", file)
	return nil
}
