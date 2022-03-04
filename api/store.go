package api

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HostSandBox    = "https://api.storekit-sandbox.itunes.apple.com"
	HostProduction = "https://api.storekit.itunes.apple.com"

	PathLookUp             = "/inApps/v1/lookup/{orderId}"
	PathTransactionHistory = "/inApps/v1/history/{originalTransactionId}"
	PathRefundHistory      = "/inApps/v1/refund/lookup/"
)

type StoreClient struct {
	ProjectID string
	Token     *Token
	cert      *Cert
}

// NewStoreClient create appstore server api client
func NewStoreClient(token *Token) *StoreClient {
	client := &StoreClient{
		Token: token,
		cert:  &Cert{},
	}
	return client
}

// LookupOrderID
func (a *StoreClient) LookupOrderID(invoiceOrderId string) (rsp *OrderLookupResponse, err error) {
	URL := HostProduction + PathLookUp
	URL = strings.Replace(URL, "{orderId}", invoiceOrderId, -1)
	statusCode, body, err := a.Do(http.MethodGet, URL, nil)
	if err != nil {
		return
	}

	if statusCode != http.StatusOK {
		err = fmt.Errorf("LookupOrderID inApps/v1/lookup api return status code %v", statusCode)
		return
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		fmt.Errorf("LookupOrderID unmarshal err %+v", err)
		return nil, err
	}

	return
}

// GetTransactionHistory https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (a *StoreClient) GetTransactionHistory(originalTransactionId string) (responses []*HistoryResponse, err error) {
	URL := HostProduction + PathTransactionHistory
	URL = strings.Replace(URL, "{originalTransactionId}", originalTransactionId, -1)
	rsp := HistoryResponse{}

	for ;; {
		data := url.Values{}
		if rsp.HasMore && rsp.Revision != "" {
			data.Set("revision", rsp.Revision)
		}
		statusCode, body, errOmit := a.Do(http.MethodGet, URL+"?"+data.Encode(), nil)
		if errOmit != nil {
			return nil, errOmit
		}

		if statusCode != http.StatusOK {
			return nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
		}

		err = json.Unmarshal(body, &rsp)
		if err != nil {
			fmt.Errorf("GetTransactionHistory unmarshal err %+v", err)
			return
		}

		responses = append(responses, &rsp)

		if !rsp.HasMore {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	return
}

// GetRefundHistory https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (a *StoreClient) GetRefundHistory(originalTransactionID string) (rsp *RefundLookupResponse, statusCode int, err error) {
	URL := HostProduction + PathRefundHistory + originalTransactionID

	statusCode, body, err := a.Do(http.MethodGet, URL, nil)
	if err != nil {
		fmt.Errorf("GetRefundHistory Get err %+v", err)
		return nil, statusCode, err
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return nil, statusCode, fmt.Errorf("GetRefundHistory Unmarshal err %w", err)
	}

	return
}

func (a *StoreClient) ParseSignedTransactions(transactions []string) ([]*JWSTransaction, error)  {
	result := make([]*JWSTransaction, 0)
	for _, v := range transactions {
		trans, err := a.parseSignedTransaction(v)
		if err == nil && trans != nil {
			result = append(result, trans)
		}
	}

	return result, nil
}

func (a *StoreClient) parseSignedTransaction(transaction string) (*JWSTransaction, error) {
	tran := &JWSTransaction{}

	rootCertStr, err := a.cert.extractCertByIndex(transaction, 2)
	if err != nil {
		return nil, err
	}
	intermediaCertStr, err := a.cert.extractCertByIndex(transaction, 1)
	if err != nil {
		return nil, err
	}
	if err = a.cert.verifyCert(rootCertStr, intermediaCertStr); err != nil {
		return nil, err
	}

	_, err = jwt.ParseWithClaims(transaction, tran, func(token *jwt.Token) (interface{}, error) {
		return a.cert.extractPublicKeyFromToken(transaction)
	})
	if err != nil {
		return nil, err
	}

	return tran, nil
}

// Per doc: https://developer.apple.com/documentation/appstoreserverapi#topics
func (a *StoreClient) Do(method string, url string, body io.Reader) (int, []byte, error) {
	authToken, err := a.Token.GenerateIfExpired()
	if err != nil {
		return 0, nil, fmt.Errorf("appstore generate token err %w", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, fmt.Errorf("appstore new http request err %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("User-Agent", "App Store Client")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("appstore http client do err %w", err)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("appstore read http body err %w", err)
	}

	return resp.StatusCode, bytes, err
}

