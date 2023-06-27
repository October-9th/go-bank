package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/October-9th/simple-bank/database/mock"
	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/October-9th/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserAPIParams struct {
	arg      sqlc.CreateUserParams
	password string
}

func (e eqCreateUserAPIParams) Matches(x interface{}) bool {
	arg, ok := x.(sqlc.CreateUserParams)

	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserAPIParams) String() string {
	return fmt.Sprintf("Matches arg %v and password %v", e.arg, e.password)
}

// Add function to return an instance of eqCreateUserAPIParams
func EqCreateUserAPIParams(arg sqlc.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserAPIParams{arg, password}
}
func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCase := []struct {
		name          string
		body          gin.H
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"Username": user.Username,
				"Password": password,
				"Fullname": user.FullName,
				"Email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := sqlc.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.
					EXPECT().
					CreateUser(gomock.Any(), EqCreateUserAPIParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"Username": "!@#$%",
				"Password": password,
				"Fullname": user.FullName,
				"Email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"Fullname": user.FullName,
				"Email":    "123",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username": user.Username,
				"password": "123",
				"Fullname": user.FullName,
				"Email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"Fullname": user.FullName,
				"Email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(sqlc.User{}, sql.ErrConnDone)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStub(store)
			server := NewServer(store)

			r := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/users"

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(r, req)

			tc.checkResponse(r)
		})
	}
}

func randomUser(t *testing.T) (sqlc.User, string) {
	password := util.RandomString(10)
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		log.Fatal("Can't hash password", err)
	}
	return sqlc.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}, password
}
