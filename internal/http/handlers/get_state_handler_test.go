package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/ivan1993spb/snake-bot/internal/http/handlers"
	"github.com/ivan1993spb/snake-bot/internal/http/handlers/handlersfakes"
	"github.com/ivan1993spb/snake-bot/internal/models"
)

func Test_GetStateHandler_AcceptNotSpecified(t *testing.T) {
	expectedState := map[int]int{
		1: 1,
		2: 2,
		5: 12,
	}
	app := &handlersfakes.FakeAppGetState{}
	app.GetStateReturns(expectedState)

	server := httptest.NewServer(handlers.NewGetStateHandler(app))
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "text/yaml", resp.Header.Get("Content-Type"))

	var games *models.Games
	err = yaml.NewDecoder(resp.Body).Decode(&games)
	require.NoError(t, err)
	require.Equal(t, expectedState, games.ToMapState())
}

func Test_GetStateHandler_AcceptJson(t *testing.T) {
	expectedState := map[int]int{
		2:  2,
		7:  7,
		9:  9,
		45: 100,
	}
	app := &handlersfakes.FakeAppGetState{}
	app.GetStateReturns(expectedState)

	server := httptest.NewServer(handlers.NewGetStateHandler(app))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	resp, err := server.Client().Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var games *models.Games
	err = json.NewDecoder(resp.Body).Decode(&games)
	require.NoError(t, err)
	require.Equal(t, expectedState, games.ToMapState())
}

func Test_GetStateHandler_AcceptYaml(t *testing.T) {
	expectedState := map[int]int{
		2:  2,
		7:  7,
		9:  9,
		45: 100,
	}
	app := &handlersfakes.FakeAppGetState{}
	app.GetStateReturns(expectedState)

	server := httptest.NewServer(handlers.NewGetStateHandler(app))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "text/yaml")
	resp, err := server.Client().Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "text/yaml", resp.Header.Get("Content-Type"))

	var games *models.Games
	err = yaml.NewDecoder(resp.Body).Decode(&games)
	require.NoError(t, err)
	require.Equal(t, expectedState, games.ToMapState())
}

func Test_GetStateHandler_AcceptDeadbeef(t *testing.T) {
	app := &handlersfakes.FakeAppGetState{}
	app.GetStateReturns(map[int]int{})

	server := httptest.NewServer(handlers.NewGetStateHandler(app))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "Deadbeef")
	resp, err := server.Client().Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
}
