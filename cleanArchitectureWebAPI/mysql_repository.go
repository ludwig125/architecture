package main

import (
	_ "github.com/go-sql-driver/mysql"
)

/*
type mySQLRepository struct {
	db *sql.DB
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ActorRepository = (*mySQLRepository)(nil)

func NewMySQLActorRepository(dbName string) (ActorRepository, error) {
	db, err := ConnMySQL(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to connection db: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db. try 'sudo service mysql start', err: %v", err)
	}

	return &mySQLRepository{db: db}, nil
}

// 事前に'sudo service mysql start'でmysqlを立ち上げておく必要がある
// WSLのUbuntuではserviceコマンドしかデフォルトでは使えなかったが、
// systemdが使えるなら'systemctl start mysql'でもいい
func ConnMySQL(dbName string) (*sql.DB, error) {
	// DNS: root:password@tcp(ipaddress:port)/dbname
	// https://github.com/go-sql-driver/mysql#examples
	// パスワードなしで、localhostに対して、デフォルトの3306 portに接続する場合は以下でいい
	return sql.Open("mysql", fmt.Sprintf("root:@tcp/%s", dbName))
}

func (r *mySQLRepository) GetAll() ([]Actor, error) {
	rows, err := r.db.Query("SELECT * FROM actor")
	if err != nil {
		return nil, fmt.Errorf("failed to select all actors, err: %v", err)
	}
	defer rows.Close()

	var actors []Actor
	for rows.Next() {
		var p Actor
		err := rows.Scan(&p.ID, &p.Name, &p.Age)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %v", err)
		}
		actors = append(actors, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %v", err)
	}
	return actors, nil
}
*/
