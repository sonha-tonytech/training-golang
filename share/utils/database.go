package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Development struct {
		Dialect    string `yaml:"dialect"`
		Datasource string `yaml:"datasource"`
		LogQueries bool   `yaml:"log_queries"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		Database   string `yaml:"database"`
	} `yaml:"development"`

	Production struct {
		Dialect    string `yaml:"dialect"`
		Datasource string `yaml:"datasource"`
		LogQueries bool   `yaml:"log_queries"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		Database   string `yaml:"database"`
	} `yaml:"production"`
}

func LoadConfig() error {
	file, err := os.Open("dbconfig.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	var config DatabaseConfig
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}
	return nil
}

func OpenDatabase() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var (
		dbDriver = os.Getenv("DB_DRIVER")
		dbUser   = os.Getenv("DB_USER")
		dbPass   = os.Getenv("DB_PASS")
		dbName   = os.Getenv("DB_NAME")
	)

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db, nil
}

func RemoveOldBackups(dir string, days int) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	expiry := time.Now().AddDate(0, 0, -days)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(expiry) {
			err := os.Remove(fmt.Sprintf("%s/%s", dir, file.Name()))
			if err != nil {
				return err
			}
			fmt.Printf("Removed: %s\n", file.Name())
		}
	}

	return nil
}

func BackupDatabase() {
	backupDir := "/Users/Tony/Documents/my-golang-app/backup"
	date := time.Now().Format("2006-01-02")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.sql", date))

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			fmt.Printf("Fail to create backup directory: %v\n", err)
			return
		}
	}

	cmd := exec.Command("mysqldump", "-u", "root", "-pletmein", "golang_crud")
	outputFile, err := os.Create(backupFile)
	if err != nil {
		fmt.Printf("Fail to create backup file: %v\n", err)
		return
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Fail to backup database: %v\n", err)
		return
	}

	fmt.Printf("Backup successful: %s\n", backupFile)

	err = RemoveOldBackups(backupDir, 7)
	if err != nil {
		fmt.Printf("Cannot remove old backups: %v\n", err)
	}
}
