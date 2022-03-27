package computation

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lnikon/glcs/common"

	_ "github.com/lib/pq"
)

type ComputationServiceDbConnector struct {
	db *sql.DB
}

func NewComputationServiceDbConnector() (*ComputationServiceDbConnector, error) {
	dbUser := os.Getenv(string(common.DBUser))
	dbPassword := os.Getenv(string(common.DBUser))
	dbName := os.Getenv(string(common.DBName))
	dbHost := os.Getenv(string(common.DBHost))
	dbPort, err := strconv.Atoi(os.Getenv(string(common.DBPort)))
	if err != nil {
		return nil, err
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	return &ComputationServiceDbConnector{
		db: db,
	}, nil
}

func (cc *ComputationServiceDbConnector) WriteNewComputationIntoDb(computation *Computation) error {
	stmt, err := cc.db.Prepare("insert into public.\"Computations\"(Name, Algorithm, VertexCount, Density, Replicas, StartTime, Status) values($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	desc := computation.Description()
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

func (cc *ComputationServiceDbConnector) ReadComputationFromDb(name string) (*Computation, error) {
	stmt, err := cc.db.Prepare("select Name, Algorithm, VertexCount, Replicas, Result from public.\"Computations\" c where c.Name=$1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	description := &ComputationDescription{}
	computation := &Computation{description: description}
	result := ""
	err = stmt.QueryRow(name).Scan(&description.Name, &description.Algorithm, &description.VertexCount, &description.Replicas, &result)
	if err != nil {
		return nil, err
	}

	computation.result = bytes.NewBufferString(result)

	return computation, nil
}
