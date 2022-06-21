package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFile(t *testing.T) {
	fileName := "unittest_file"
	cleanup, err := SetupTestFile(t, fileName)
	if err != nil {
		t.Fatalf("failed to SetupTestFile: %v", err)
	}
	// 以下のcleanup関数を呼ばなければローカルにテスト用ファイルが残る
	defer t.Cleanup(cleanup)

	t.Run("invalid_file", func(t *testing.T) {
		if _, err := NewExcludeRepository("invalid_file"); err != nil {
			wantErr := "failed to open file:open invalid_file: no such file or directory"
			if err.Error() != wantErr {
				t.Fatalf("gotErr: %v, wantErr: %v", err.Error(), wantErr)
			}
		}
	})

	e, err := NewExcludeRepository(fileName)
	if err != nil {
		t.Fatalf("failed to NewExcludeRepository: %v", err)
	}

	_, err = e.Excluded()
	if err != nil {
		t.Fatal(err)
	}

	// テストケースをmapにすることで順序依存がないことを確認
	tests := map[string]struct {
		testData string
		wants    []string
	}{
		"actors": {
			testData: "Watson\nDepp\n",
			wants: []string{
				"Watson",
				"Depp",
			},
		},
		"no_data": {
			testData: "",
			wants:    nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := writeTestDataToFile(t, fileName, tc.testData); err != nil {
				t.Fatalf("failed to write testData: %v", err)
			}
			got, err := e.Excluded()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.wants, got); diff != "" {
				t.Errorf("got: %v, wants: %v\ndiff: %v", got, tc.wants, diff)
			}
		})
	}
}

// SetupTestFile creates test Database and table
// cleanupTestFile関数を返すので、呼び出し元は"defer SetupTestFile(t)()"
// とするだけで、test用DatabaseとTableの作成と、テスト終了時の削除を担保できる
func SetupTestFile(t *testing.T, fileName string) (func(), error) {
	t.Helper()
	if fileName == "" {
		return nil, errors.New("fileName is not set")
	}
	t.Logf("test file name: '%s'", fileName)

	// file の作成
	f, err := createTestFile(t, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to createTestFile: %v", err)
	}
	defer f.Close()
	t.Logf("created test file '%s' successfully", fileName)

	// cleanupTestFile関数を返す
	return func() {
		if err := cleanupTestFile(t, fileName); err != nil {
			t.Fatalf("failed to cleanupTestFile: %v", err)
		}
	}, nil
}

// test用file作成
func createTestFile(t *testing.T, fileName string) (*os.File, error) {
	t.Helper()
	return os.Create(fileName)
}

// test用Fileの削除
func cleanupTestFile(t *testing.T, fileName string) error {
	t.Helper()
	return os.Remove(fileName)
}

func writeTestDataToFile(t *testing.T, fileName, testData string) error {
	t.Helper()
	f, err := os.OpenFile(fileName, os.O_RDWR, 0664)
	if err != nil {
		return fmt.Errorf("failed to Open: %v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(testData); err != nil {
		return fmt.Errorf("failed to WriteString: %v", err)
	}
	return nil
}
