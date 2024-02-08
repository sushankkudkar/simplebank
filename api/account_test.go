package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	mockdb "github.com/sushankkudkar/simplebank/db/mock"
	db "github.com/sushankkudkar/simplebank/db/sqlc"
	"github.com/sushankkudkar/simplebank/util"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	fmt.Println(account)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Set up an expectation for GetAccount
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	// Generate the URL with the correct account ID
	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Serve the HTTP request
	server.router.ServeHTTP(recorder, request)

	// Check the recorded response status
	require.Equal(t, http.StatusOK, recorder.Code)

	// Print the recorder details for debugging
	fmt.Println(recorder)

	// Validate the mock call
	gomock.InOrder(
		store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil),
	)

	// Ensure the recorder contains the expected response
	expectedResponse := fmt.Sprintf(`{"id":%d,"owner":"%s","balance":%d,"currency":"%s"}`, account.ID, account.Owner, account.Balance, account.Currency)
	require.JSONEq(t, expectedResponse, recorder.Body.String())
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
