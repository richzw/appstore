package api

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

const ACCOUNTKEY = `
-----BEGIN PRIVATE KEY-----
FAKEACCOUNTKEYBASE64FORMAT
-----END PRIVATE KEY-----
`

func TestStoreClient_LookupOrderID(t *testing.T) {
	type args struct {
		invoiceOrderId string
	}
	tests := []struct {
		name    string
		args    args
		wantRsp *OrderLookupResponse
		wantErr bool
	}{
		{
			name: "Lookup api test",
			args: args{invoiceOrderId: "FAKEINVOICEID"},
			wantRsp: &OrderLookupResponse{
				Status: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreConfig{
				KeyContent: []byte(ACCOUNTKEY),
				KeyID:      "SKEYID",
				BundleID:   "fake.bundle.id",
				Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
				Sandbox:    false,
			}

			a := NewStoreClient(c)
			gotRsp, err := a.LookupOrderID(tt.args.invoiceOrderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupOrderID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRsp.Status, tt.wantRsp.Status) {
				t.Errorf("LookupOrderID() gotRsp = %v, want %v", gotRsp, tt.wantRsp)
			}

			orders, err := a.ParseSignedTransactions(gotRsp.SignedTransactions)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupOrderID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, o := range orders {
				t.Log(o)
			}
		})
	}
}

func TestStoreClient_GetTransactionHistory(t *testing.T) {
	type args struct {
		originalTransactionId string
		query                 *url.Values
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetTransactionHistory api test",
			args: args{originalTransactionId: "123321",
				query: &url.Values{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreConfig{
				KeyContent: []byte(ACCOUNTKEY),
				KeyID:      "SKEYID",
				BundleID:   "fake.bundle.id",
				Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
				Sandbox:    false,
			}

			a := NewStoreClient(c)
			tt.args.query.Set("productType", "AUTO_RENEWABLE")
			tt.args.query.Set("productType", "NON_CONSUMABLE")
			tt.args.query.Set("productType", "CONSUMABLE")
			gotRsp, err := a.GetTransactionHistory(tt.args.originalTransactionId, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, rsp := range gotRsp {
				trans, err := a.ParseSignedTransactions(rsp.SignedTransactions)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetTransactionHistory() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				for _, tran := range trans {
					t.Logf("%+v", tran)
				}
			}

		})
	}
}

func TestStoreClient_GetRefundHistory(t *testing.T) {
	type args struct {
		originalTransactionId string
		query                 *url.Values
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetRefundHistory api test",
			args: args{originalTransactionId: "123321",
				query: &url.Values{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreConfig{
				KeyContent: []byte(ACCOUNTKEY),
				KeyID:      "SKEYID",
				BundleID:   "fake.bundle.id",
				Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
				Sandbox:    false,
			}

			a := NewStoreClient(c)
			gotRsp, err := a.GetRefundHistory(tt.args.originalTransactionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRefundHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, rsp := range gotRsp {
				t.Logf("%+v", rsp.SignedTransactions)
				trans, err := a.ParseSignedTransactions(rsp.SignedTransactions)
				t.Logf("%+v", trans)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetTransactionHistory() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				for _, tran := range trans {
					t.Logf("%+v", tran)
				}
			}
		})
	}
}

func TestStoreClient_GetNotificationHistory(t *testing.T) {
	type args struct {
		body NotificationHistoryRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetNotificationHistory api test",
			args: args{
				body: NotificationHistoryRequest{
					StartDate: time.Now().Add(-time.Hour * time.Duration(480)).UnixMilli(),
					EndDate:   time.Now().Add(-time.Hour * time.Duration(24)).UnixMilli(),
					//OriginalTransactionId: "123321",
					NotificationType: NotificationTypeV2Refund,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreConfig{
				KeyContent: []byte(ACCOUNTKEY),
				KeyID:      "SKEYID",
				BundleID:   "fake.bundle.id",
				Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",q
			}

			a := NewStoreClient(c)
			gotRsp, err := a.GetNotificationHistory(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRefundHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, rsp := range gotRsp {
				t.Logf("%+v", rsp)
			}

		})
	}
}
