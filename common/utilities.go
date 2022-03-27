package common

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Described environment variables necessary for proper operation of UPCXX
type UPCXXEnvVars string

const (
	UpcxxInstall     UPCXXEnvVars = "UPCXX_INSTALL"
	PGASGraphInstall              = "PGASGRAPH_INSTALL"
)

type UPCXXBinaries string

const (
	UPCXXMeta    UPCXXBinaries = "upcxx-meta"
	UPCXX                      = "upcxx"
	UPCXXRun                   = "upcxx-run"
	PGASGraphRun               = "pgas-graph-cli-smp"
)

type DBEnvVar string

// Move into .env.
const (
	DBHost     DBEnvVar = "DB_HOST"
	DBName     DBEnvVar = "DB_NAME"
	DBUser     DBEnvVar = "DB_USER"
	DBPassword DBEnvVar = "DB_PASSWORD"
)

func CheckUpcxxEnvVars() error {
	upcxxInstall := os.Getenv(string(UpcxxInstall))
	if len(upcxxInstall) == 0 {
		return errors.New("$UPCXX_INSTALL env var is not set")
	}

	pgasgraphInstall := os.Getenv(string(PGASGraphInstall))
	if len(pgasgraphInstall) == 0 {
		return errors.New("$PGASGRAPH_INSTALL env var is not set")
	}

	return nil
}

func CheckDbEnvVars() error {
	value := os.Getenv(string(DBName))
	if len(value) == 0 {
		return fmt.Errorf("$%s env var is not set", string(DBName))
	}

	value = os.Getenv(string(DBUser))
	if len(value) == 0 {
		return fmt.Errorf("$%s env var is not set", string(DBUser))
	}

	value = os.Getenv(string(DBPassword))
	if len(value) == 0 {
		return fmt.Errorf("$%s env var is not set", string(DBPassword))
	}

	return nil
}

func CheckUpcxxBinaries() error {
	_, err := exec.LookPath(string(UPCXXMeta))
	if err != nil {
		return fmt.Errorf("Can not find %s in $PATH", UPCXXMeta)
	}

	_, err = exec.LookPath(string(UPCXX))
	if err != nil {
		return fmt.Errorf("Can not find %s in $PATH", UPCXX)
	}

	_, err = exec.LookPath(string(UPCXXRun))
	if err != nil {
		return fmt.Errorf("Can not find %s in $PATH", UPCXXRun)
	}

	return nil
}
