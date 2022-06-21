package main

// func TestService(t *testing.T) {
// 	mockSQLiteRep := &mockSQLiteRepository{}
// 	mockExcludeRep := &mockFileRepository{}

// 	t.Run("GetAll", func(t *testing.T) {
// 		tests := map[string]struct {
// 			mockGetAllFunc  func() ([]Actor, error)
// 			mockExcludeFunc func() ([]string, error)
// 			want            []Actor
// 		}{
// 			"without_exclude": {
// 				mockGetAllFunc: func() ([]Actor, error) {
// 					return []Actor{
// 						{ID: 1, Name: "Watson", Age: 24},
// 						{ID: 2, Name: "Depp", Age: 54},
// 						{ID: 3, Name: "Portman", Age: 32},
// 					}, nil
// 				},
// 				mockExcludeFunc: func() ([]string, error) {
// 					return nil, nil
// 				},
// 				want: []Actor{
// 					{ID: 1, Name: "Watson", Age: 24},
// 					{ID: 2, Name: "Depp", Age: 54},
// 					{ID: 3, Name: "Portman", Age: 32},
// 				},
// 			},
// 			"with_exclude": {
// 				mockGetAllFunc: func() ([]Actor, error) {
// 					return []Actor{
// 						{ID: 1, Name: "Watson", Age: 24},
// 						{ID: 2, Name: "Depp", Age: 54},
// 						{ID: 3, Name: "Portman", Age: 32},
// 					}, nil
// 				},
// 				mockExcludeFunc: func() ([]string, error) {
// 					return []string{"Depp"}, nil
// 				},
// 				want: []Actor{
// 					{ID: 1, Name: "Watson", Age: 24},
// 					{ID: 3, Name: "Portman", Age: 32},
// 				},
// 			},
// 			"no_record": {
// 				mockGetAllFunc: func() ([]Actor, error) {
// 					return nil, nil
// 				},
// 				mockExcludeFunc: func() ([]string, error) {
// 					return nil, nil
// 				},
// 				want: nil,
// 			},
// 		}
// 		for name, tc := range tests {
// 			t.Run(name, func(t *testing.T) {
// 				mockSQLiteRep.mockGetAllFunc = tc.mockGetAllFunc
// 				mockExcludeRep.mockExcludedFunc = tc.mockExcludeFunc

// 				service := NewActorService(Config{}, mockSQLiteRep, mockExcludeRep)
// 				got, err := service.GetAll()
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				if diff := cmp.Diff(tc.want, got); diff != "" {
// 					t.Errorf("got: %v, want: %v\ndiff: %v", got, tc.want, diff)
// 				}
// 			})
// 		}
// 	})

// 	t.Run("Find", func(t *testing.T) {
// 		mockSQLiteRep.mockFindByIDFunc = func() ([]Actor, error) {
// 			return []Actor{
// 				{ID: 1, Name: "Watson", Age: 24},
// 			}, nil
// 		}
// 		mockSQLiteRep.mockFindByNameFunc = func() ([]Actor, error) {
// 			return []Actor{
// 				{ID: 2, Name: "Depp", Age: 54},
// 			}, nil
// 		}
// 		mockSQLiteRep.mockFindByAgeFunc = func() ([]Actor, error) {
// 			return []Actor{
// 				{ID: 3, Name: "Portman", Age: 32},
// 			}, nil
// 		}
// 		tests := map[string]struct {
// 			mockExcludeFunc func() ([]string, error)
// 			findCondition   Actor
// 			want            []Actor
// 			wantErr         string
// 		}{
// 			"findByID_without_exclude": {
// 				mockExcludeFunc: func() ([]string, error) {
// 					return nil, nil
// 				},
// 				findCondition: Actor{
// 					ID: 1,
// 				},
// 				want: []Actor{
// 					{ID: 1, Name: "Watson", Age: 24},
// 				},
// 			},
// 			"findByName_without_exclude": {
// 				mockExcludeFunc: func() ([]string, error) {
// 					return nil, nil
// 				},
// 				findCondition: Actor{
// 					Name: "Depp",
// 				},
// 				want: []Actor{
// 					{ID: 2, Name: "Depp", Age: 54},
// 				},
// 			},
// 			"findByAge_without_exclude": {
// 				mockExcludeFunc: func() ([]string, error) {
// 					return nil, nil
// 				},
// 				findCondition: Actor{
// 					Age: 32,
// 				},
// 				want: []Actor{
// 					{ID: 3, Name: "Portman", Age: 32},
// 				},
// 			},
// 			"exclude_different": {
// 				mockExcludeFunc: func() ([]string, error) {
// 					return []string{"Depp"}, nil
// 				},
// 				findCondition: Actor{
// 					ID: 1,
// 				},
// 				want: []Actor{
// 					{ID: 1, Name: "Watson", Age: 24},
// 				},
// 			},
// 			"exclude": {
// 				mockExcludeFunc: func() ([]string, error) {
// 					return []string{"Watson"}, nil
// 				},
// 				findCondition: Actor{
// 					ID: 1,
// 				},
// 				want: nil,
// 			},
// 			"no_condition": {
// 				mockExcludeFunc: func() ([]string, error) {
// 					return nil, nil
// 				},
// 				findCondition: Actor{},
// 				want:          []Actor{},
// 				wantErr:       `failed to findActors: invalid condition: main.Actor{ID:0, Name:"", Age:0}`,
// 			},
// 		}
// 		for name, tc := range tests {
// 			t.Run(name, func(t *testing.T) {
// 				mockExcludeRep.mockExcludedFunc = tc.mockExcludeFunc

// 				service := NewActorService(Config{}, mockSQLiteRep, mockExcludeRep)
// 				got, err := service.Find(tc.findCondition)
// 				if err != nil {
// 					if err.Error() != tc.wantErr {
// 						t.Fatalf("gotErr: %v, wantErr: %v", err.Error(), tc.wantErr)
// 					}
// 					return
// 				}
// 				if diff := cmp.Diff(tc.want, got); diff != "" {
// 					t.Errorf("got: %v, want: %v\ndiff: %v", got, tc.want, diff)
// 				}
// 			})
// 		}
// 	})

// 	t.Run("Update", func(t *testing.T) {
// 		tests := map[string]struct {
// 			mockFindByNameFunc func() ([]Actor, error)
// 			mockUpdateFunc     func() error
// 			updateTarget       Actor
// 			want               []Actor
// 			wantErr            string
// 		}{
// 			"update_new_record": {
// 				mockFindByNameFunc: func() ([]Actor, error) {
// 					return nil, nil // もともとはUpdate対象のデータはない
// 				},
// 				mockUpdateFunc: func() error {
// 					return nil
// 				},
// 				updateTarget: Actor{
// 					Name: "Portman",
// 				},
// 				want: []Actor{
// 					{ID: 3, Name: "Portman", Age: 32},
// 				},
// 			},
// 			"update_new_record_with_find_error": {
// 				mockFindByNameFunc: func() ([]Actor, error) {
// 					return nil, errors.New("error occurs at findByName")
// 				},
// 				mockUpdateFunc: func() error {
// 					return nil
// 				},
// 				updateTarget: Actor{
// 					Name: "Portman",
// 				},
// 				want: []Actor{
// 					{ID: 3, Name: "Portman", Age: 32},
// 				},
// 				wantErr: "failed to FindByName: error occurs at findByName",
// 			},
// 			"update_duplicated_record": {
// 				mockFindByNameFunc: func() ([]Actor, error) {
// 					return []Actor{
// 						{ID: 3, Name: "Portman", Age: 32},
// 					}, nil
// 				},
// 				mockUpdateFunc: func() error {
// 					return nil
// 				},
// 				updateTarget: Actor{
// 					Name: "Portman",
// 				},
// 				want: []Actor{
// 					{ID: 3, Name: "Portman", Age: 32},
// 				},
// 				wantErr: "actor Portman already exists",
// 			},
// 		}
// 		for name, tc := range tests {
// 			t.Run(name, func(t *testing.T) {
// 				mockSQLiteRep.mockFindByNameFunc = tc.mockFindByNameFunc
// 				mockSQLiteRep.mockUpdateFunc = tc.mockUpdateFunc
// 				mockExcludeRep.mockExcludedFunc = func() ([]string, error) { return nil, nil }

// 				service := NewActorService(Config{}, mockSQLiteRep, mockExcludeRep)
// 				if err := service.Update(tc.updateTarget); err != nil {
// 					if err.Error() != tc.wantErr {
// 						t.Errorf("gotErr: %v, wantErr: %v", err.Error(), tc.wantErr)
// 					}
// 				}
// 			})
// 		}
// 	})

// 	t.Run("Delete", func(t *testing.T) {
// 		tests := map[string]struct {
// 			mockFindByIDFunc   func(int) ([]Actor, error)
// 			mockDeleteByIDFunc func() error
// 			deleteID           int
// 			wantErr            string
// 		}{
// 			"delete_non_exist_record": {
// 				mockDeleteByIDFunc: func() error {
// 					return nil
// 				},
// 				deleteID: 999,
// 				wantErr:  "no row got affected",
// 			},
// 		}
// 		for name, tc := range tests {
// 			t.Run(name, func(t *testing.T) {
// 				mockSQLiteRep.mockFindByIDFunc,err = tc.mockFindByIDFunc(mock.Anything)
// 				mockSQLiteRep.mockDeleteByIDFunc = tc.mockDeleteByIDFunc
// 				mockExcludeRep.mockExcludedFunc = func() ([]string, error) { return nil, nil }

// 				service := NewActorService(Config{}, mockSQLiteRep, mockExcludeRep)
// 				a, err := service.DeleteByID(tc.deleteID)
// 				if err != nil {
// 					if err.Error() != tc.wantErr {
// 						t.Errorf("gotErr: %v, wantErr: %v", err.Error(), tc.wantErr)
// 					}
// 					return
// 				}

// 			})
// 		}
// 	})
// }
