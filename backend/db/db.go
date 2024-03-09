package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	TargetFolderPath = "backend"
)

type CommandType uint8

const (
	Up CommandType = iota + 1
	Down
)

// connect to sqllite
func ConnectToSQLite(filepath string) (*sql.DB, error) {
	err := createDatabaseIfExist(filepath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	err = executeTableCmd(db, Up)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// create a new database if it does not exist
func createDatabaseIfExist(path string) error {
	var (
		err  error
		file *os.File
	)

	fileExist := doesFileExist(path)
	if fileExist {
		return nil
	}

	_, err = os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	return nil
}

// function to check if file exists
func doesFileExist(path string) bool {
	_, err := os.Stat(path)

	// check if error is "file not exists"
	return !os.IsNotExist(err)
}

func executeTableCmd(db *sql.DB, cmdType CommandType) error {
	var (
		query string
		err   error
	)
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	index := strings.Index(currentWorkingDirectory, TargetFolderPath)
	if index == -1 {
		return errors.New("app Root Folder Path not found")
	}

	path := filepath.Join(currentWorkingDirectory[:index], TargetFolderPath, "db")
	fmt.Println("newpath :", path)
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	switch cmdType {
	case Up:
		query, err = getFileContent(files, path, "create_table.sql")
		if err != nil {
			return err
		}
	case Down:
		query, err = getFileContent(files, path, "drop_table.sql")
		if err != nil {
			return err
		}
	default:
		return errors.New("no command type available")
	}

	//TODO:Change this to ExecContext later
	_, err = db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func getFileContent(files []fs.DirEntry, basePath string, lookupFile string) (string, error) {
	for _, file := range files {
		if file.Name() == lookupFile {
			path := filepath.Join(basePath, lookupFile)
			content, err := os.ReadFile(path)
			if err != nil {
				return "", err
			}
			return string(content), nil
		}
	}

	return "", errors.New("file not found")
}
