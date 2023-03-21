/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package wrapper_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/api"
	"github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/wrapper"
)

type mockServer struct {
	t              *testing.T
	headersToCheck *api.Headers
}

func (m *mockServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if m.headersToCheck != nil {
		for _, headerToCheck := range m.headersToCheck.GetAll() {
			// Note: for these tests, we're assuming that there aren't multiple values under a single name/key.
			value := request.Header.Get(headerToCheck.Name)
			require.Equal(m.t, headerToCheck.Value, value)
		}
	}
}

func TestHTTPClient_Do(t *testing.T) {
	server := &mockServer{}
	testServer := httptest.NewServer(server)

	defer testServer.Close()

	request, err := http.NewRequest(http.MethodGet, testServer.URL, http.NoBody)
	require.NoError(t, err)

	httpClient := wrapper.NewHTTPClient()

	t.Run("Default client settings", func(t *testing.T) {
		response, err := httpClient.Do(request)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NoError(t, response.Body.Close())
	})
	t.Run("TLS verification disabled", func(t *testing.T) {
		httpClient.DisableTLSVerification = true

		response, err := httpClient.Do(request)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NoError(t, response.Body.Close())
	})
	t.Run("With additional headers", func(t *testing.T) {
		additionalHeaders := api.NewHeaders()

		additionalHeaders.Add(api.NewHeader("header-name-1", "header-value-1"))
		additionalHeaders.Add(api.NewHeader("header-name-2", "header-value-2"))

		httpClient.AddHeaders(additionalHeaders)

		response, err := httpClient.Do(request)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NoError(t, response.Body.Close())
	})
}