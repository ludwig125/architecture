package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteRepository struct {
	db *sql.DB
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ActorRepository = (*sqliteRepository)(nil)

func NewSQLiteActorRepository(dbName string) (ActorRepository, error) {
	db, err := connSQLite(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to connection db: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}

	return &sqliteRepository{db: db}, nil
}

func connSQLite(dbName string) (*sql.DB, error) {
	// DNS: root:password@tcp(ipaddress:port)/dbname
	// https://github.com/go-sql-driver/mysql#examples
	// パスワードなしで、localhostに対して、デフォルトの3306 portに接続する場合は以下でいい
	return sql.Open("sqlite3", dbName)
}

func (r *sqliteRepository) GetAll() ([]Actor, error) {
	rows, err := r.db.Query("SELECT * FROM actor")
	if err != nil {
		return nil, fmt.Errorf("failed to select all actors, err: %v", err)
	}
	defer rows.Close()

	var actors []Actor
	for rows.Next() {
		var a Actor
		err := rows.Scan(&a.ID, &a.Name, &a.Age)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %v", err)
		}
		actors = append(actors, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %v", err)
	}
	return actors, nil
}

func (r *sqliteRepository) SearchByID(id int) ([]Actor, error) {
	rows, err := r.db.Query("SELECT * FROM actor WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return scanActors(rows)
}

func (r *sqliteRepository) SearchByName(name string) ([]Actor, error) {
	rows, err := r.db.Query("SELECT * FROM actor WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return scanActors(rows)
}

func (r *sqliteRepository) SearchByAge(age int) ([]Actor, error) {
	rows, err := r.db.Query("SELECT * FROM actor WHERE age = ?", age)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return scanActors(rows)
}

func scanActors(rows *sql.Rows) ([]Actor, error) {
	var actors []Actor
	defer rows.Close()
	for rows.Next() {
		var a Actor
		err := rows.Scan(&a.ID, &a.Name, &a.Age)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %v", err)
		}
		actors = append(actors, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %v", err)
	}
	return actors, nil
}

func (r *sqliteRepository) Update(a Actor) error {
	q := "INSERT OR REPLACE INTO actor(name, age) VALUES($1, $2);"
	if _, err := r.db.Exec(q, a.Name, a.Age); err != nil {
		return fmt.Errorf("failed to update db: %v", err)
	}
	return nil
}

func (r *sqliteRepository) DeleteByID(id int) error {
	q := "DELETE FROM actor WHERE id = $1;"
	if _, err := r.db.Exec(q, id); err != nil {
		return fmt.Errorf("failed to delete from db: %v", err)
	}
	return nil
}
