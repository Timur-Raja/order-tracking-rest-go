package testing

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// todo improbement: make tests independent of each other and add failing tests

type endpointTest struct {
	name         string
	method       string
	path         string
	body         string
	wantStatus   int
	wantContains string
}

var today = time.Now().Format("2006-01-02")

var tests = []endpointTest{
	{
		name:   "Make an Order",
		method: http.MethodPost,
		path:   "/orders",
		body: `{
		        "shippingAddress": "test address 123",
				"orderItems": [
				{
				"productID": 1,
				"quantity": 3
				},
				{
				"productID": 3,
				"quantity": 1
				}
				 ]
		}`,
		wantStatus:   201,
		wantContains: `"status":"created"`,
	},
	{
		name:         "Get order",
		method:       http.MethodGet,
		path:         "/orders/1",
		wantStatus:   200,
		wantContains: `"id":1,"userID":1,"status":"created"`,
	},
	{
		name:         "Search orders by text",
		method:       http.MethodGet,
		path:         "/orders?q=test",
		wantStatus:   200,
		wantContains: `"shippingAddress":"test address 123"`,
	},
	{
		name:         "Filter orders by date",
		method:       http.MethodGet,
		path:         fmt.Sprintf("/orders?start_date=%s&end_date=%s", today, today),
		wantStatus:   200,
		wantContains: `"shippingAddress":"test address 123"`,
	},
	{
		name:   "update order shipping address",
		method: http.MethodPatch,
		path:   "/orders/1",
		body: `{
				"shippingAddress": "New address 123"
				}`,
		wantStatus:   200,
		wantContains: `"shippingAddress":"New address 123"`,
	},
	{
		name:   "cancel order",
		method: http.MethodPatch,
		path:   "/orders/1",
		body: `{
				"status": "cancelled"
				}`,
		wantStatus:   200,
		wantContains: `"status":"cancelled"`,
	},
	{
		name:         "Create user",
		method:       http.MethodPost,
		path:         "/signup",
		body:         `{"email":"email@example2.com","name":"John Doe", "password":"password"}`,
		wantStatus:   201,
		wantContains: `"id":2`, // second user after the seeded one
	},
	{
		name:         "Sign in with user",
		method:       http.MethodPost,
		path:         "/signin",
		body:         `{"email":"email@example2.com","password":"password"}`,
		wantStatus:   201,
		wantContains: `"id":2,"`,
	},
}

func TestEndpoints(t *testing.T) {
	for _, tc := range tests {
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			var bodyReader io.Reader
			if tc.body != "" {
				bodyReader = strings.NewReader(tc.body)
			}
			req, err := http.NewRequest(tc.method, baseURL+tc.path, bodyReader)
			c := &http.Cookie{
				Name:  "session_token",
				Value: "test123",
			}
			req.AddCookie(c)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantStatus, resp.StatusCode, "status code")
			assert.Contains(t, string(data), tc.wantContains)
		})
	}
}
