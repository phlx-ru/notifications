package data

import (
	"database/sql"
	"fmt"
	"regexp"

	entDialectSQL "entgo.io/ent/dialect/sql"
)

const (
	databaseNameRegex    = `dbname=([a-zA-Z_][a-zA-Z0-9_]*)`
	databaseNamePostgres = `dbname=postgres`
)

//goland:noinspection SqlDialectInspection
func createDatabaseIfNotExists(driver, source string) error {
	if driver != "postgres" {
		return nil // TODO Support other databases
	}
	original := extractDatabaseNameFromSource(source)
	if original == "" {
		return nil
	}
	db, err := openDefaultDatabase(driver, source)
	if err != nil {
		return err
	}
	rows, err := db.Query(`select true as exists from pg_database where datname = $1`, original)
	if err != nil {
		return err
	}
	if rows.Next() {
		return nil
	}
	_, err = db.Exec(fmt.Sprintf(`create database %s`, original))
	return err
}

func extractDatabaseNameFromSource(source string) string {
	regex := regexp.MustCompile(databaseNameRegex)
	submatch := regex.FindAllStringSubmatch(source, 1)
	if len(submatch) < 1 {
		return ""
	}
	match := submatch[0]
	if len(match) < 2 {
		return ""
	}
	databaseName := match[1]
	if databaseName == "" {
		return ""
	}
	if databaseName == `postrges` {
		return ""
	}
	return databaseName
}

func openDefaultDatabase(driver, source string) (*sql.DB, error) {
	regex := regexp.MustCompile(databaseNameRegex)
	baseSource := regex.ReplaceAllString(source, databaseNamePostgres)
	db, err := entDialectSQL.Open(driver, baseSource)
	if err != nil {
		return nil, err
	}
	return db.DB(), nil
}
