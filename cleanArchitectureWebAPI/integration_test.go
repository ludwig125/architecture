package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestIntegrationErrorCase(t *testing.T) {
	dbName := "integrationtest_error_sqlite_db" // 他のunit testと並行でテストが動くと名前が衝突するので
	tableName := "actor"
	dbCleanup, err := SetupTestDB(t, dbName, tableName)
	if err != nil {
		t.Fatalf("failed to SetupTestDB: %v", err)
	}
	// 以下のcleanup関数を呼ばなければローカルにテスト用DBが残る
	defer t.Cleanup(dbCleanup)

	fileName := "integrationtest_error_file"
	fCleanup, err := SetupTestFile(t, fileName)
	if err != nil {
		t.Fatalf("failed to SetupTestFile: %v", err)
	}
	// 以下のcleanup関数を呼ばなければローカルにテスト用ファイルが残る
	defer t.Cleanup(fCleanup)

	tests := map[string]struct {
		envs    map[string]string
		wantErr string
	}{
		"invalid_db_type": {
			envs: map[string]string{
				"DB_TYPE":            "invalid_db_Type",
				"DB_NAME":            "",
				"EXCLUDE_ACTOR_FILE": "",
			},
			wantErr: "invalid dbType: invalid_db_Type",
		},
		"invalid_db_name": {
			envs: map[string]string{
				"DB_TYPE":            "sqlite",
				"DB_NAME":            "invalid_db_name",
				"EXCLUDE_ACTOR_FILE": "",
			},
			wantErr: "failed to NewMySQLActorRepository: no such db file: invalid_db_name",
		},
		"invalid_exclude_file": {
			envs: map[string]string{
				"DB_TYPE":            "sqlite",
				"DB_NAME":            "integrationtest_error_sqlite_db",
				"EXCLUDE_ACTOR_FILE": "invalid_exclude_file",
			},
			wantErr: "failed to NewExcludeRepository: failed to open file:open invalid_exclude_file: no such file or directory",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// 環境変数の設定
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			// 次のテストのために環境変数の削除
			defer func() {
				for k := range tc.envs {
					os.Unsetenv(k)
				}
			}()
			err := runActorAPI()
			if err == nil {
				t.Errorf("should be error")
				return
			}
			if err.Error() != tc.wantErr {
				t.Errorf("got error: %s, want error: %s", err.Error(), tc.wantErr)
			}
		})
	}

}

func TestIntegration(t *testing.T) {
	dbName := "integrationtest_sqlite_db" // 他のunit testと並行でテストが動くと名前が衝突するので
	tableName := "actor"
	dbCleanup, err := SetupTestDB(t, dbName, tableName)
	if err != nil {
		t.Fatalf("failed to SetupTestDB: %v", err)
	}
	// 以下のcleanup関数を呼ばなければローカルにテスト用DBが残る
	defer t.Cleanup(dbCleanup)

	fileName := "integrationtest_file"
	fCleanup, err := SetupTestFile(t, fileName)
	if err != nil {
		t.Fatalf("failed to SetupTestFile: %v", err)
	}
	// 以下のcleanup関数を呼ばなければローカルにテスト用ファイルが残る
	defer t.Cleanup(fCleanup)

	// 環境変数を事前に設定
	os.Setenv("DB_TYPE", "sqlite")
	os.Setenv("DB_NAME", dbName)

	go func(t *testing.T) {
		// 別goroutineでサーバを起動
		if err := runActorAPI(); err != nil {
			t.Errorf("failed to runActorAPI: %v", err)
		}
	}(t)

	// 上のgoroutineのサーバが起動するまで1秒待つ
	time.Sleep(100 * time.Millisecond)

	// status check
	t.Run("status", func(t *testing.T) {
		got, err := testRequest("GET", "http://localhost:8080/status", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		want := "ok"
		if got != want {
			t.Errorf("got: %s, want %s", got, want)
		}
	})
	t.Run("update", func(t *testing.T) {
		t.Run("update_Watson", func(t *testing.T) {
			req := decodeActor(t, &Actor{Name: "Watson", Age: 24})
			got, err := testRequest("POST", "http://localhost:8080/update", "application/json", bytes.NewBuffer(req))
			if err != nil {
				t.Fatal(err)
			}
			want := "Updated Actor name: 'Watson', age: '24'"
			if got != want {
				t.Errorf("got: %s, want %s", got, want)
			}
		})
		t.Run("update_Portman", func(t *testing.T) {
			req := decodeActor(t, &Actor{Name: "Portman", Age: 32})
			got, err := testRequest("POST", "http://localhost:8080/update", "application/json", bytes.NewBuffer(req))
			if err != nil {
				t.Fatal(err)
			}
			want := "Updated Actor name: 'Portman', age: '32'"
			if got != want {
				t.Errorf("got: %s, want %s", got, want)
			}
		})
		t.Run("update_Knightley", func(t *testing.T) {
			req := decodeActor(t, &Actor{Name: "Knightley", Age: 36})
			got, err := testRequest("POST", "http://localhost:8080/update", "application/json", bytes.NewBuffer(req))
			if err != nil {
				t.Fatal(err)
			}
			want := "Updated Actor name: 'Knightley', age: '36'"
			if got != want {
				t.Errorf("got: %s, want %s", got, want)
			}
		})
		// 上と同じデータを入れても重複データが入らないことを確認
		t.Run("update_Knightley", func(t *testing.T) {
			req := decodeActor(t, &Actor{Name: "Knightley", Age: 36})
			got, err := testRequest("POST", "http://localhost:8080/update", "application/json", bytes.NewBuffer(req))
			if err != nil {
				t.Fatal(err)
			}
			want := "Internal Server Error\n"
			if got != want {
				t.Errorf("got: %s, want %s", got, want)
			}
		})
	})
	t.Run("getall", func(t *testing.T) {
		got, err := testRequest("GET", "http://localhost:8080/getall", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		want := `[{"id":1,"name":"Watson","age":24},{"id":2,"name":"Portman","age":32},{"id":3,"name":"Knightley","age":36}]
`
		if got != want {
			t.Errorf("got: %s, want %s", got, want)
		}
	})
	t.Run("delete", func(t *testing.T) {
		req := decodeActor(t, &Actor{Name: "Knightley", Age: 36})
		got, err := testRequest("POST", "http://localhost:8080/delete", "application/json", bytes.NewBuffer(req))
		if err != nil {
			t.Fatal(err)
		}
		want := "Actor name: 'Knightley', age: '36'"
		if got != want {
			t.Errorf("got: %s, want %s", got, want)
		}
	})
}

func testRequest(method, url, contentType string, body io.Reader) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.TODO(), method, url, body)
	if err != nil {
		return "", fmt.Errorf("failed to NewRequestWithContext: %v", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to ReadAll: %v", err)
	}
	return string(res), nil
}
