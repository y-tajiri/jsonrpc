// +build !go1.7

package jsonrpc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestHandler(t *testing.T) {

	c := context.Background()
	rec := httptest.NewRecorder()
	r, err := http.NewRequest("", "", nil)
	require.NoError(t, err)

	Handler(c, rec, r)

	res := Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	Handler(c, rec, r)
	res = Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	require.NoError(t, RegisterMethod("hello", func(c context.Context, params *json.RawMessage) (interface{}, *Error) {
		return "hello", nil
	}))

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	Handler(c, rec, r)
	res = Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.Nil(t, res.Error)
	assert.Equal(t, "hello", res.Result)
}