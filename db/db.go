package database

import (
	"database/sql"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type CommandType uint8

const (
	Up CommandType = iota + 1
	Down
)

const (
	dbFile = "urlshortnerDB.db"
)

// connect to sqllite
func ConnectToSQLite() (*sql.DB, error) {
	err := createDatabaseIfExist()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	err = executeTableCmd(db, Up)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db, nil
}

// create a new database if it does not exist
func createDatabaseIfExist() error {
	var (
		err  error
		file *os.File
	)

	exist := doesFileExist()
	if exist {
		return nil
	}

	_, err = os.Create(dbFile)
	if err != nil {
		return err
	}

	defer file.Close()
	return nil
}

// function to check if file exists
func doesFileExist() bool {
	_, error := os.Stat(dbFile)

	// check if error is "file not exists"
	return os.IsNotExist(error)
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

	path := filepath.Join(currentWorkingDirectory, "db")
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	switch cmdType {
	case Up:
		query, err = getFileContent(files, path, "createtable.sql")
		if err != nil {
			return err
		}
	case Down:
		query, err = getFileContent(files, path, "droptable.sql")
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
