package main

import (
	_ "github.com/mattn/go-sqlite3"
)

type mockService struct {
	mockGetAllFunc     func() ([]Actor, error)
	mockFindFunc       func() ([]Actor, error)
	mockUpdateFunc     func() error
	mockDeleteByIDFunc func() (Actor, error)
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ActorService = (*mockService)(nil)

// func NewMockSQLiteActorRepository() ActorRepository {
// 	// func NewMockSQLiteActorRepository(dbName string) (ActorRepository, error) {
// 	// db, err := connSQLite(dbName)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to connection db: %v", err)
// 	// }
// 	// if err := db.Ping(); err != nil {
// 	// 	return nil, fmt.Errorf("failed to ping db: %v", err)
// 	// }

// 	return &mockService{}
// 	// return &mockService{db: db}, nil
// }

// func connSQLite(dbName string) (*sql.DB, error) {
// 	// DNS: root:password@tcp(ipaddress:port)/dbname
// 	// https://github.com/go-sql-driver/mysql#examples
// 	// パスワードなしで、localhostに対して、デフォルトの3306 portに接続する場合は以下でいい
// 	return sql.Open("sqlite3", dbName)
// }

func (r *mockService) GetAll() ([]Actor, error) {
	// rows, err := r.db.Query("SELECT * FROM actor")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to select all actors, err: %v", err)
	// }
	// defer rows.Close()

	// var actors []Actor
	// for rows.Next() {
	// 	var a Actor
	// 	err := rows.Scan(&a.ID, &a.Name, &a.Age)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to scan: %v", err)
	// 	}
	// 	actors = append(actors, a)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, fmt.Errorf("row error: %v", err)
	// }
	// return actors, nil
	return r.mockGetAllFunc()
}

func (r *mockService) Find(cond Actor) ([]Actor, error) {
	// rows, err := r.db.Query("SELECT * FROM actor WHERE id = ?", id)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to query: %v", err)
	// }
	// return scanActors(rows)
	return r.mockFindFunc()
}

// func scanActors(rows *sql.Rows) ([]Actor, error) {
// 	var actors []Actor
// 	defer rows.Close()
// 	for rows.Next() {
// 		var a Actor
// 		err := rows.Scan(&a.ID, &a.Name, &a.Age)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan: %v", err)
// 		}
// 		actors = append(actors, a)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("row error: %v", err)
// 	}
// 	return actors, nil
// }

func (r *mockService) Update(a Actor) error {

	// 	q := "INSERT OR REPLACE INTO actor(name, age) VALUES($1, $2);"
	// 	res, err := r.db.Exec(q, a.Name, a.Age)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to update db: %v", err)
	// 	}
	// 	// コマンドで影響を受けた件数が０ならエラーとする
	// 	row, err := res.RowsAffected()
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get RowsAffected: %v", err)
	// 	}
	// 	if row == 0 {
	// 		return errors.New("no row got affected")
	// 	}
	return r.mockUpdateFunc()
}

func (r *mockService) DeleteByID(id int) (Actor, error) {
	// 	q := "DELETE FROM actor WHERE id = $1;"
	// 	res, err := r.db.Exec(q, id)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to delete from db: %v", err)
	// 	}

	// 	// コマンドで影響を受けた件数が０ならエラーとする
	// 	row, err := res.RowsAffected()
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get RowsAffected: %v", err)
	// 	}
	// 	if row == 0 {
	// 		return errors.New("no row got affected")
	// 	}
	return r.mockDeleteByIDFunc()
}

// func (r *mockService) SetMockGetAllFunc(f func() ([]Actor, error)) {
// 	r.mockGetAllFunc = f
// }
