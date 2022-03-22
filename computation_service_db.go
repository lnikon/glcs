package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBEnvVar string

// Move into .env.
const (
	DBName     DBEnvVar = "DB_NAME"
	DBUser     DBEnvVar = "DB_USER"
	DBPassword DBEnvVar = "DB_PASSWORD"
)

type ComputationServiceDbConnector struct {
	db *sql.DB
}

func NewComputationServiceDbConnector() (*ComputationServiceDbConnector, error) {
	dbUser := os.Getenv(string(DBUser))
	dbPassword := os.Getenv(string(DBUser))
	dbName := os.Getenv(string(DBName))
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	return &ComputationServiceDbConnector{
		db: db,
	}, nil
}

func (cc *ComputationServiceDbConnector) WriteNewComputationIntoDb(desc *ComputationDescription) error {
	stmt, err := cc.db.Prepare("insert into public.\"Computations\"(Name, Algorithm, VertexCount, Density, Replicas, StartTime, Status) values($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(desc.Name, desc.Algorithm, desc.VertexCount, desc.Density, desc.Replicas, time.Now(), Starting)
	if err != nil {
		return err
	}

	return nil
}

func (cc *ComputationServiceDbConnector) UpdateComputationStatusInDb(name string, status string, result string) error {
	var stmt *sql.Stmt
	var err error

	if len(result) == 0 {
		stmt, err = cc.db.Prepare("update public.\"Computations\" set EndTime=$1, Status=$2 where Name=$3")
	} else {
		stmt, err = cc.db.Prepare("update public.\"Computations\" set EndTime=$1, Status=$2, Result=$3 where Name=$4")
	}
	defer stmt.Close()

	if err != nil {
		return err
	}

	if len(result) == 0 {
		_, err = stmt.Exec(time.Now(), status, name)
	} else {
		_, err = stmt.Exec(time.Now(), status, result, name)
	}

	if err != nil {
		return err
	}

	return nil
}

func (cc *ComputationServiceDbConnector) ReadComputationFromDb(name string) (string, error) {
	stmt, err := cc.db.Prepare("select result from public.\"Computations\" c where c.Name=$1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var result string
	err = stmt.QueryRow(name).Scan(&result)
	if err != nil {
		return "", err
	}

	return result, nil
}
