package migrations

import (
	"bufio"
	"database/sql"
	"fmt"
	"my-pp/share/variables"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func MigrateDatabase() {
	currentVersion, err := readVersion()
	if err != nil {
		panic(err)
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		version, err := parseVersion(file)
		if err != nil {
			panic(err)
		}

		if version > currentVersion {
			if err := applyMigration(variables.DB, file); err != nil {
				panic(err)
			}
			currentVersion = version
			if err := writeVersion(currentVersion); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("All migrations applied successfully!")
}

func readVersion() (int, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)

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

func writeVersion(version int) error {
	file, err := os.Create("version.txt")
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
