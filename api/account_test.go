package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/October-9th/simple-bank/database/mock"
	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/October-9th/simple-bank/token"
	"github.com/October-9th/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	testCase := []struct {
		name string
		// Declare the request body
		body          gin.H
		setupAuth     func(t *testing.T, r *http.Request, token token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{

			name: "OK",
			body: gin.H{
				"Owner":    account.Owner,
				"Currency": account.Currency,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
				addAuthorization(t, r, token, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := sqlc.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"Owner":    account.Owner,
				"Currency": "Failed Currency",
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
				addAuthorization(t, r, token, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"Owner":    account.Owner,
				"Currency": account.Currency,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
				addAuthorization(t, r, token, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := sqlc.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sqlc.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"Owner":    account.Owner,
				"Currency": account.Currency,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		// run each case as a separate sub-test of this unit test
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// Build stubs for api testing
			tc.buildStubs(store)
			// Start test server and send request
			server := newTestSever(t, store)

			// Don't need to start the real http server, instead can just use the Recorder feature
			// of the httptest package to record the response of the API request
			recorder := httptest.NewRecorder()

			// Marshall body data to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			// Path to the url path of the API
			url := "/api/v1/accounts"

			// Initialize new HTTP request using method POST to that url
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			// Send API request through the server router and record the response
			server.router.ServeHTTP(recorder, request)

			// Check the response of the recorder
			tc.checkResponse(t, recorder)

		})
	}
}
func TestGetAcocuntAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	// Declare test cases
	testCases := []struct {
		name      string
		accountID int64
		// GetAccountStub for each scenario will be built differently suit the purpose of each test case
		setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		// checkResponse check the output of the API request
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			}, buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},

		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			}, buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			}, buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					// Remove the authentication then the request will be aborted by the middleware, so we can expect that the GetAccount request will not be called
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// TODO: add more cases
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sqlc.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sqlc.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// TODO: add more cases
		{
			name:      "BadRequest",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()). // pass invalid ID to uri
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {

		// Declare variables to store data of the current test case
		tc := testCases[i]
		// run each case as a separate sub-test of this unit test
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// Build stubs for api testing
			tc.buildStubs(store)
			// Start test server and send request
			server := newTestSever(t, store)

			// Don't need to start the real http server, instead can just use the Recorder feature
			// of the httptest package to record the response of the API request
			recorder := httptest.NewRecorder()

			// Declare the url path of the API need to be called
			url := fmt.Sprintf("/api/v1/accounts/%d", tc.accountID)

			// Create a new HTTP request with method GET to that url, since it the GET request
			// use nil for the request body

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			// Send the API request through the server router
			// and record its response in the recorder
			server.router.ServeHTTP(recorder, request)

			// Check that response of the recorder
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetListAccount(t *testing.T) {
	n := 5
	user, _ := randomUser(t)
	accounts := make([]sqlc.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}
	type Query struct {
		PageID   int
		PageSize int
	}
	testCase := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, r *http.Request, token token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
				addAuthorization(t, r, token, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := sqlc.ListAccountsParams{
					Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				PageID:   -1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
				addAuthorization(t, r, token, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageNumber",
			query: Query{
				PageID:   1,
				PageSize: 1000000000000,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
				addAuthorization(t, r, token, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
				addAuthorization(t, r, token, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]sqlc.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, r *http.Request, token token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			}},
	}
	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestSever(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/v1/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			// Add query parameter to request url
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount(owner string) sqlc.Account {
	return sqlc.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account sqlc.Account) {
	// Read all data from the response body
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	getAccount := sqlc.Account{}
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)
	require.Equal(t, account, getAccount)
}
