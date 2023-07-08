package api

import (
	"testing"

	"github.com/KHarshit1203/simple-bank/service/db/mocks"
	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/stretchr/testify/suite"
)

type ApiServerSuite struct {
	suite.Suite
	mockConfig util.Config
	mockStore  *mocks.Store
	server     Server
}

func TestApiServerSuite(t *testing.T) {
	suite.Run(t, &ApiServerSuite{})
}

func (at *ApiServerSuite) SetupSubTest() {
	mockConfig := util.Config{TokenSymmetricKey: util.RandomString(32)}
	mockStore := mocks.NewStore(at.T())
	testServer, err := NewServer(mockConfig, mockStore)
	at.NoError(err)

	at.mockConfig = mockConfig
	at.mockStore = mockStore
	at.server = testServer
}
