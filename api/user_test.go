package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/ProstoyVadila/simple_bank/db/mock"
	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// Implementation of gomock.Matcher interface
type eqCreateUserParamsMathcher struct {
	args     db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMathcher) Matches(x interface{}) bool {
	params, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := utils.CheckPassword(e.password, params.HashedPassword)
	if err != nil {
		return false
	}
	e.args.HashedPassword = params.HashedPassword
	return reflect.DeepEqual(e.args, params)
}

func (e eqCreateUserParamsMathcher) String() string {
	return fmt.Sprintf("matches args: %v, password %v", e.args, e.password)
}

// EqCreateUserParams is a custom implementation of gomock.Eq
func EqCreateUserParams(args db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMathcher{args: args, password: password}
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	body, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(body)
}

func createRandomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(14)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username:          utils.RandomOwner(),
		FullName:          utils.RandomString(12),
		Email:             utils.RandomEmail(),
		HashedPassword:    hashedPassword,
		CreatedAt:         time.Now(),
		PasswordChangedAt: time.Now(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var resp defaultUserResponse
	err = json.Unmarshal(data, &resp)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, user.Email, resp.Email)
	require.Equal(t, user.Username, resp.Username)
	require.Equal(t, user.FullName, resp.FullName)
	require.WithinDuration(t, user.CreatedAt, resp.CreatedAt, time.Minute)
}

func TestCreateUserApi(t *testing.T) {
	user, password := createRandomUser(t)
	tests := []struct {
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, recorder *httptest.ResponseRecorder)
		userReq    createUserRequest
		name       string
	}{
		{
			name: "OK",
			userReq: createUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(
					gomock.Any(),
					EqCreateUserParams(args, password),
				).Times(1).Return(user, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/users", jsonBody(t, tt.userReq))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tt.checkResp(t, recorder)
		})
	}
}

// func TestCreateUserApi(t *testing.T) {
// 	type fields struct {
// 		store  db.Store
// 		router *gin.Engine
// 	}
// 	type args struct {
// 		ctx *gin.Context
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Server{
// 				store:  tt.fields.store,
// 				router: tt.fields.router,
// 			}
// 			s.createUser(tt.args.ctx)
// 		})
// 	}
// }
