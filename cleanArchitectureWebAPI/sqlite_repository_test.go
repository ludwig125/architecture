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
	cleanup, err := SetupTestDB(t, dbName, tableName)
	if err != nil {
		t.Fatalf("failed to SetupTestDB: %v", err)
	}
	// 以下のcleanup関数を呼ばなければローカルにテスト用DBが残る
	defer t.Cleanup(cleanup)

	r, err := NewSQLiteActorRepository(dbName)
	if err != nil {
		t.Fatalf("failed to NewSQLiteActorRepository: %v", err)
	}

	// テストケースをmapにすることで順序依存がないことを確認
	tests := map[string]struct {
		actors               []Actor
		wantGetAll           []Actor
		actorID              int
		wantFindByID         []Actor
		actorName            string
		wantFindByName       []Actor
		actorAge             int
		wantFindByAge        []Actor
		wantGetAllDuplicated []Actor
		deleteID             int
		wantDeleteErr        bool
		wantAfterDelete      []Actor
	}{
		"one_actor": {
			actors: []Actor{
				{Name: "Depp", Age: 54},
			},
			wantGetAll: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
			},
			actorID: 1,
			wantFindByID: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
			},
			actorName: "Depp",
			wantFindByName: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
			},
			actorAge: 54,
			wantFindByAge: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
			},
			wantGetAllDuplicated: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
				{ID: 2, Name: "Depp", Age: 54},
			},
			deleteID: 2,
			wantAfterDelete: []Actor{
				{ID: 1, Name: "Depp", Age: 54},
			},
		},
		"three_actoress": {
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
			wantFindByID: []Actor{
				{ID: 2, Name: "Knightley", Age: 36},
			},
			actorName: "Portman",
			wantFindByName: []Actor{
				{ID: 1, Name: "Portman", Age: 32},
			},
			actorAge: 24,
			wantFindByAge: []Actor{
				{ID: 3, Name: "Watson", Age: 24},
			},
			wantGetAllDuplicated: []Actor{
				{ID: 1, Name: "Portman", Age: 32},
				{ID: 2, Name: "Knightley", Age: 36},
				{ID: 3, Name: "Watson", Age: 24},
				{ID: 4, Name: "Portman", Age: 32},
				{ID: 5, Name: "Knightley", Age: 36},
				{ID: 6, Name: "Watson", Age: 24},
			},
			deleteID: 3,
			wantAfterDelete: []Actor{
				{ID: 1, Name: "Portman", Age: 32},
				{ID: 2, Name: "Knightley", Age: 36},
				{ID: 4, Name: "Portman", Age: 32},
				{ID: 5, Name: "Knightley", Age: 36},
				{ID: 6, Name: "Watson", Age: 24},
			},
		},
		"no_data": {
			actors:               []Actor{},
			wantGetAll:           nil,
			actorID:              1,
			wantFindByID:         nil,
			actorName:            "Depp",
			wantFindByName:       nil,
			actorAge:             54,
			wantFindByAge:        nil,
			wantGetAllDuplicated: nil,
			deleteID:             2,
			wantDeleteErr:        true, // 何も消すものがないと'no row got affected' エラーが出る
			wantAfterDelete:      nil,
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
			t.Run("FindByID", func(t *testing.T) {
				got, err := r.FindByID(tc.actorID)
				if err != nil {
					t.Errorf("failed to FindByID: %v", err)
				}
				if diff := cmp.Diff(tc.wantFindByID, got); diff != "" {
					t.Errorf("got: %v, wantFindByID: %v\ndiff: %v", got, tc.wantFindByID, diff)
				}
			})
			t.Run("FindByName", func(t *testing.T) {
				got, err := r.FindByName(tc.actorName)
				if err != nil {
					t.Errorf("failed to FindByName: %v", err)
				}
				if diff := cmp.Diff(tc.wantFindByName, got); diff != "" {
					t.Errorf("got: %v, wantFindByName: %v\ndiff: %v", got, tc.wantFindByName, diff)
				}
			})
			t.Run("FindByAge", func(t *testing.T) {
				got, err := r.FindByAge(tc.actorAge)
				if err != nil {
					t.Errorf("failed to FindByAge: %v", err)
				}
				if diff := cmp.Diff(tc.wantFindByAge, got); diff != "" {
					t.Errorf("got: %v, wantFindByAge: %v\ndiff: %v", got, tc.wantFindByAge, diff)
				}
			})
			t.Run("Update_same", func(t *testing.T) {
				// 同じデータを入れると別のIDで入ってしまうことを確認
				// 重複の判定はrepositoryではしない
				for _, a := range tc.actors {
					if err := r.Update(a); err != nil {
						t.Fatalf("failed to Update: %v", err)
					}
				}
				got, err := r.GetAll()
				if err != nil {
					t.Errorf("failed to GetAll: %v", err)
				}
				if diff := cmp.Diff(tc.wantGetAllDuplicated, got); diff != "" {
					t.Errorf("got: %v, wantGetAllDuplicated: %v\ndiff: %v", got, tc.wantGetAllDuplicated, diff)
				}
			})
			t.Run("DeleteByID", func(t *testing.T) {
				// 一件消す
				if err := r.DeleteByID(tc.deleteID); err != nil {
					if !tc.wantDeleteErr {
						t.Log(err)
						t.Errorf("failed to DeleteByID: %v", err)
					}
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
		if err := cleanupTestDB(t, dbName); err != nil {
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
func cleanupTestDB(t *testing.T, dbName string) error {
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
