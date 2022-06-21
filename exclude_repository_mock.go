package main

type mockFileRepository struct {
	// File     *os.File
	// FileName string
	mockExcludedFunc func() ([]string, error)
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ExcludeRepository = (*mockFileRepository)(nil)

// func NewExcludeRepository(fileName string) (ExcludeRepository, error) {
// 	f, err := openExcludedFile(fileName)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open file:%v", err)
// 	}
// 	return &mockFileRepository{File: f, FileName: fileName}, nil
// }

// func openExcludedFile(fileName string) (*os.File, error) {
// 	return os.Open(fileName)
// }

func (e *mockFileRepository) Excluded() ([]string, error) {
	// var excluded []string
	// scanner := bufio.NewScanner(e.File)
	// for scanner.Scan() {
	// 	excluded = append(excluded, string(scanner.Text()))
	// }
	// if err := scanner.Err(); err != nil {
	// 	return nil, fmt.Errorf("failed to scan file %s: %v", e.FileName, err)
	// }
	// return excluded, nil
	return e.mockExcludedFunc()
}
