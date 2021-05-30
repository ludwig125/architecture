package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSQLite(t *testing.T) {
	dbName := "test_sqlite_db"
	tableName := "actor"
	_, err := SetupTestDB(t, dbName, tableName)
	if err != nil {
		t.Fatalf("failed to SetupTestDB: %v", err)
	}
	// 以下のcleanup関数を呼ばなければローカルにテスト用DBが残る
	// defer t.Cleanup(cleanup)

	r, err := NewSQLiteActorRepository(dbName)
	if err != nil {
		t.Fatalf("failed to NewSQLiteActorRepository: %v", err)
	}

	// テストケースをmapにすることで順序依存がないことを確認
	tests := map[string]struct {
		actors           []Actor
		wantGetAll       []Actor
		actorID          int
		wantSearchByID   []Actor
		actorName        string
		wantSearchByName []Actor
		actorAge         int
		wantSearchByAge  []Actor
		deleteID         int
		wantAfterDelete  []Actor
	}{
		"actors": {
			actors: []Actor{
				{Name: "Depp", Age: 54},
				{Name: "Jackson", Age: 63},
				{Name: "Hopkins", Age: 74},
			},
			wantGetAll: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
				{ID: 2, Name: "Jackson", Age: 63},
				{ID: 3, Name: "Hopkins", Age: 74},
			},
			actorID: 2,
			wantSearchByID: []Actor{
				{ID: 2, Name: "Jackson", Age: 63},
			},
			actorName: "Hopkins",
			wantSearchByName: []Actor{
				{ID: 3, Name: "Hopkins", Age: 74},
			},
			actorAge: 54,
			wantSearchByAge: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
			},
			deleteID: 3,
			wantAfterDelete: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
				{ID: 2, Name: "Jackson", Age: 63},
			},
		},
		"actoress": {
			actors: []Actor{
				{Name: "Portman", Age: 32},
				{Name: "Knightley", Age: 36},
				{Name: "Watson", Age: 24},
			},
			wantGetAll: []Actor{
				{ID: 1, Name: "Portman", Age: 32},
				{ID: 2, Name: "Knightley", Age: 36},
				{ID: 3, Name: "Watson", Age: 24},
			},
			actorID: 2,
			wantSearchByID: []Actor{
				{ID: 2, Name: "Knightley", Age: 36},
			},
			actorName: "Portman",
			wantSearchByName: []Actor{
				{ID: 1, Name: "Portman", Age: 32},
			},
			actorAge: 24,
			wantSearchByAge: []Actor{
				{ID: 3, Name: "Watson", Age: 24},
			},
			deleteID: 3,
			wantAfterDelete: []Actor{
				{ID: 1, Name: "Portman", Age: 32},
				{ID: 2, Name: "Knightley", Age: 36},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run("Update", func(t *testing.T) {
				for _, a := range tc.actors {
					if err := r.Update(a); err != nil {
						t.Fatalf("failed to Update: %v", err)
					}
				}
			})
			t.Run("GetAll", func(t *testing.T) {
				got, err := r.GetAll()
				if err != nil {
					t.Errorf("failed to GetAll: %v", err)
				}
				if diff := cmp.Diff(tc.wantGetAll, got); diff != "" {
					t.Errorf("got: %v, wantGetAll: %v\ndiff: %v", got, tc.wantGetAll, diff)
				}
			})
			t.Run("SearchByID", func(t *testing.T) {
				got, err := r.SearchByID(tc.actorID)
				if err != nil {
					t.Errorf("failed to SearchByID: %v", err)
				}
				if diff := cmp.Diff(tc.wantSearchByID, got); diff != "" {
					t.Errorf("got: %v, wantSearchByID: %v\ndiff: %v", got, tc.wantSearchByID, diff)
				}
			})
			t.Run("SearchByName", func(t *testing.T) {
				got, err := r.SearchByName(tc.actorName)
				if err != nil {
					t.Errorf("failed to SearchByName: %v", err)
				}
				if diff := cmp.Diff(tc.wantSearchByName, got); diff != "" {
					t.Errorf("got: %v, wantSearchByName: %v\ndiff: %v", got, tc.wantSearchByName, diff)
				}
			})
			t.Run("SearchByAge", func(t *testing.T) {
				got, err := r.SearchByAge(tc.actorAge)
				if err != nil {
					t.Errorf("failed to SearchByAge: %v", err)
				}
				if diff := cmp.Diff(tc.wantSearchByAge, got); diff != "" {
					t.Errorf("got: %v, wantSearchByAge: %v\ndiff: %v", got, tc.wantSearchByAge, diff)
				}
			})
			t.Run("SearchByAge", func(t *testing.T) {
				got, err := r.SearchByAge(tc.actorAge)
				if err != nil {
					t.Errorf("failed to SearchByAge: %v", err)
				}
				if diff := cmp.Diff(tc.wantSearchByAge, got); diff != "" {
					t.Errorf("got: %v, wantSearchByAge: %v\ndiff: %v", got, tc.wantSearchByAge, diff)
				}
			})
			t.Run("DeleteByID", func(t *testing.T) {
				// 一件消す
				if err := r.DeleteByID(tc.deleteID); err != nil {
					t.Errorf("failed to DeleteByID: %v", err)
				}
				got, err := r.GetAll()
				if err != nil {
					t.Errorf("failed to GetAll: %v", err)
				}
				if diff := cmp.Diff(tc.wantAfterDelete, got); diff != "" {
					t.Errorf("got: %v, wantAfterDelete: %v\ndiff: %v", got, tc.wantAfterDelete, diff)
				}

				// 残りも消しておく
				for _, a := range got {
					if err := r.DeleteByID(a.ID); err != nil {
						t.Errorf("failed to DeleteByID: %v", err)
					}
				}
				got, err = r.GetAll()
				if err != nil {
					t.Errorf("failed to GetAll: %v", err)
				}
				if len(got) > 0 {
					t.Errorf("got: %v, remains after delete all", got)
				}
			})
		})
	}
}

// SetupTestDB creates test Database and table
// cleanupTestDB関数を返すので、呼び出し元は"defer SetupTestDB(t)()"
// とするだけで、test用DatabaseとTableの作成と、テスト終了時の削除を担保できる
func SetupTestDB(t *testing.T, dbName, tableName string) (func(), error) {
	t.Helper()
	if dbName == "" {
		return nil, errors.New("dbName is not set")
	}
	if tableName == "" {
		return nil, errors.New("tableName is not set")
	}
	t.Logf("test database name: '%s', table name: '%s'", dbName, tableName)

	// Database の作成
	f, err := createTestDB(t, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to createTestDB: %v", err)
	}
	defer f.Close()
	t.Logf("created test table '%s' successfully", tableName)

	db, err := connSQLite(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to connSQLite: %v", err)
	}

	// test用tableの作成
	if err := createTestTable(t, db, tableName); err != nil {
		return nil, fmt.Errorf("failed to createTestTable: %v", err)
	}

	// cleanupTestDB関数を返す
	return func() {
		if err := cleanupTestDB(t, db, dbName); err != nil {
			t.Fatalf("failed to cleanupTestDB: %v", err)
		}
	}, nil
}

// test用Database作成
func createTestDB(t *testing.T, dbName string) (*os.File, error) {
	t.Helper()
	return os.Create(dbName)
}

// test用Databaseの削除
func cleanupTestDB(t *testing.T, db *sql.DB, dbName string) error {
	t.Helper()
	return os.Remove(dbName)
}

// test用tableの作成
func createTestTable(t *testing.T, db *sql.DB, tableName string) error {
	t.Helper()
	q := fmt.Sprintf("CREATE TABLE %s(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);", tableName)
	if _, err := db.Exec(q); err != nil {
		return err
	}
	return nil
}
