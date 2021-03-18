package flat

import (
	"reflect"
	"testing"
)

type AuthFields struct {
	ApiKey      string `json:"apiKey" hmac:",required"`
	Environment string `json:"environment" hmac:",required"`
	HmacKey     string `json:"HMACKey" hmac:"-"`
}

type GetTicketTypesRequest struct {
	*AuthFields
	UserId     int     `json:"userId" hmac:",omitempty"`
	Language   *string `json:"language" hmac:",omitempty"`
	TicketDate *string `json:"ticketDate" hmac:",omitempty"`
	HasTicket  bool    `json:"hasTicket" hmac:"has_ticket"`
	Ignore     string  `json:"-" hmac:"-"`
	Ignore2    string  `json:"ignore2"`
}

type GetAvailableUsersRequest struct {
	AuthFields
}

func String(s string) *string {
	return &s
}

func TestGenerateHmac(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want map[string]interface{}
	}{
		{name: "", args: GetTicketTypesRequest{
			AuthFields: &AuthFields{
				ApiKey:      "key abc",
				Environment: "test",
				HmacKey:     "hmac key",
			},
			UserId:     123,
			Language:   String("nl"),
			TicketDate: nil,
			HasTicket:  true,
			Ignore:     "nothing",
			Ignore2:    "something",
		}, want: map[string]interface{}{"UserId": 123, "Language": "nl", "has_ticket": true, "ApiKey": "key abc", "Environment": "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToMap(tt.args, String("hmac"))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() got = \n\t %v, want \n\t %v", got, tt.want)
			}
		})
	}
}
