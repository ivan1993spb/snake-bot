package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-bot/internal/core"
	"github.com/ivan1993spb/snake-bot/internal/http/handlers"
	"github.com/ivan1993spb/snake-bot/internal/http/handlers/handlersfakes"
)

func Test_SetStateHandler_XWWWFormURLEncoded(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}
	app.SetOneReturns(map[int]int{
		1: 1,
		2: 2,
	}, nil)
	expectBody := "games:\n- game: 1\n  bots: 1\n- game: 2\n  bots: 2\n"

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	form := url.Values{}
	form.Add("game", "1")
	form.Add("bots", "1")

	// Content-Type: application/x-www-form-urlencoded
	resp, err := server.Client().PostForm(server.URL, form)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, 201, resp.StatusCode)
	require.Equal(t, "text/yaml", resp.Header.Get("Content-Type"))

	buffer := bytes.NewBuffer(nil)
	_, err = buffer.ReadFrom(resp.Body)
	require.NoError(t, err)

	require.Equal(t, expectBody, buffer.String())
	require.Equal(t, 1, app.SetOneCallCount())
}

func Test_SetStateHandler_Json(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}
	app.SetStateReturns(map[int]int{
		1: 1,
		2: 21,
	}, nil)

	expectBody := `{"games":[{"game":1,"bots":1},{"game":2,"bots":21}]}` + "\n"

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	data := []byte(`{"games":[{"game":1,"bots":1}]}`)

	buffer := bytes.NewBuffer(data)

	resp, err := server.Client().Post(server.URL, "application/json", buffer)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, 201, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	buffer = bytes.NewBuffer(nil)
	_, err = buffer.ReadFrom(resp.Body)
	require.NoError(t, err)

	require.Equal(t, expectBody, buffer.String())
	require.Equal(t, 1, app.SetStateCallCount())
}

func Test_SetStateHandler_Yaml(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}
	app.SetStateReturns(map[int]int{
		1:  51,
		2:  2,
		15: 25,
	}, nil)

	expectBody := "games:\n- game: 1\n  bots: 51\n- game: 2\n  bots: 2\n- game: 15\n  bots: 25\n"

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	//data := []byte(`{"games":[{"game":1,"bots":1}]}`)
	data := []byte("games:\n  - game: 1\n    bots: 1")

	buffer := bytes.NewBuffer(data)

	resp, err := server.Client().Post(server.URL, "text/yaml", buffer)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, 201, resp.StatusCode)
	require.Equal(t, "text/yaml", resp.Header.Get("Content-Type"))

	buffer = bytes.NewBuffer(nil)
	_, err = buffer.ReadFrom(resp.Body)
	require.NoError(t, err)

	require.Equal(t, expectBody, buffer.String())
	require.Equal(t, 1, app.SetStateCallCount())
}

func Test_SetStateHandler_MediaDeadbeef(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	resp, err := server.Client().Post(server.URL, "dead/beef", nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}

func Test_SetStateHandler_MediaYaml_AcceptJson(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}
	app.SetStateReturns(map[int]int{
		16: 1,
		2:  21,
		31: 8,
	}, nil)

	expectBody := `{"games":[{"game":2,"bots":21},{"game":16,"bots":1},{"game":31,"bots":8}]}` + "\n"

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	data := []byte("games:\n  - game: 1\n    bots: 1")
	buffer := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodPost, server.URL, buffer)
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "text/yaml")

	resp, err := server.Client().Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, 201, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	buffer = bytes.NewBuffer(nil)
	_, err = buffer.ReadFrom(resp.Body)
	require.NoError(t, err)

	require.Equal(t, expectBody, buffer.String())
	require.Equal(t, 1, app.SetStateCallCount())
}

func Test_SetStateHandler_AppError(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}
	app.SetStateReturns(nil, errors.New("app error"))

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	data := []byte("games:\n  - game: 1\n    bots: 1")
	buffer := bytes.NewBuffer(data)

	resp, err := server.Client().Post(server.URL, "text/yaml", buffer)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Equal(t, 1, app.SetStateCallCount())
}

func Test_SetStateHandler_TooManyBots(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}
	app.SetOneReturns(nil, core.ErrRequestedTooManyBots)

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	form := url.Values{}
	form.Add("game", "1")
	form.Add("bots", "1")

	// Content-Type: application/x-www-form-urlencoded
	resp, err := server.Client().PostForm(server.URL, form)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, 1, app.SetOneCallCount())
}

func Test_SetStateHandler_ServiceUnavailable(t *testing.T) {
	app := &handlersfakes.FakeAppSetState{}
	app.SetOneReturns(nil, context.DeadlineExceeded)

	server := httptest.NewServer(handlers.NewSetStateHandler(app))
	defer server.Close()

	form := url.Values{}
	form.Add("game", "1")
	form.Add("bots", "1")

	// Content-Type: application/x-www-form-urlencoded
	resp, err := server.Client().PostForm(server.URL, form)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	require.Equal(t, 1, app.SetOneCallCount())
}
