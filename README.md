
[App Store Server API](https://developer.apple.com/documentation/appstoreserverapi) Golang Client
================

The App Store Server API is a REST API that you call from your server to request and provide information about your customers' in-app purchases. 

The App Store Server API is independent of the app’s installation status on the customers’ devices. The App Store server returns information based on a customer’s in-app purchase history regardless of whether the customer installs, removes, or reinstalls the app on their devices.

# Quick Start

### Installation

```shell
go get github.com/richzw/appstore
```

### [Generate a Private Key](https://developer.apple.com/documentation/appstoreserverapi/creating_api_keys_to_use_with_the_app_store_server_api)

> Log in to [App Store Connect](https://appstoreconnect.apple.com/login) and complete the following steps:
> - Select Users and Access, and then select the Keys tab.
> - Select In-App Purchase under the Key Type.
> - Click Generate API Key or the Add (+) button.
> - Enter a name for the key. The name is for your reference only and isn’t part of the key itself. Click Generate.
> - Click Download API Key next to the new API key. And store your private key in a secure place.

### Get Transaction Info

```go
import(
	"github.com/richzw/appstore"
)

// ACCOUNTPRIVATEKEY is the key file generated from previous step
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `
func main() {
    c := &appstore.StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    transactionId := "FAKETRANSACTIONID"
    a := appstore.NewStoreClient(c)
    response, err := a.GetTransactionInfo(context.TODO(), transactionId)

    transantions, err := a.ParseSignedTransactions([]{response.SignedTransactionInfo})

    if transactions[0].TransactionId == transactionId {
        // the transaction is valid
    }
}
```

- Validate the receipt
  - One option could be to validate the receipt with the App Store server through `GetTransactionInfo` API, and then check the `transactionId` in the response matches the one you are looking for.

- Error handling
    - handler error per [apple store server api error](https://developer.apple.com/documentation/appstoreserverapi/error_codes) document
    - [error definition](./error.go)

### Look Up Order ID

```go
import(
    "github.com/richzw/appstore"
)

// ACCOUNTPRIVATEKEY is the key file generated from previous step
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `

func main() {
    c := &appstore.StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    invoiceOrderId := "FAKEORDERID"

    a := appstore.NewStoreClient(c)
    rsp, err := a.LookupOrderID(context.TODO(), invoiceOrderId)

    orders, err := a.ParseSignedTransactions(rsp.SignedTransactions)
}
```

### Get Transaction History

```go
import(
    "github.com/richzw/appstore"
)

// ACCOUNTPRIVATEKEY is the key file generated from previous step
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `

func main() {
    c := &appstore.StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    originalTransactionId := "FAKEORDERID"
    a := appstore.NewStoreClient(c)
    query := &url.Values{}
    query.Set("productType", "AUTO_RENEWABLE")
    query.Set("productType", "NON_CONSUMABLE")
    gotRsp, err := a.GetTransactionHistory(context.TODO(), originalTransactionId, query)

    for _, rsp := range gotRsp {
       trans, err := a.ParseSignedTransactions(rsp.SignedTransactions)
    }
}
```

### Get Refund History

```go
import(
    "github.com/richzw/appstore"
)

// ACCOUNTPRIVATEKEY is the key file generated from previous step
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `

func main() {
    c := &appstore.StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    originalTransactionId := "FAKEORDERID"
    a := appstore.NewStoreClient(c)
    gotRsp, err := a.GetRefundHistory(context.TODO(), originalTransactionId)

    for _, rsp := range gotRsp {
       trans, err := a.ParseSignedTransactions(rsp.SignedTransactions)
    }
}
```

### Parse Notification from App Store

```go
import (
    "github.com/richzw/appstore"
    "github.com/golang-jwt/jwt/v4"
)

func main() {
    c := &appstore.StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    tokenStr := "SignedRenewalInfo Encode String" // or SignedTransactionInfo string
    a := appstore.NewStoreClient(c)
    token, err := a.ParseNotificationV2(tokenStr)

    claims, ok := token.Claims.(jwt.MapClaims)
    for key, val := range claims {
        fmt.Printf("Key: %v, value: %v\n", key, val) // key value of TransactionInfo
    }
}
```

```go
import (
    "github.com/richzw/appstore"
    "github.com/golang-jwt/jwt/v4"
)

func main() {
    c := &appstore.StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    tokenStr := "JWSTransactionDecodedPayload Encode String"
    a := appstore.NewStoreClient(c)

    jws, err := a.ParseNotificationV2WithClaim(tokenStr)
    // access the fields of JWSTransactionDecodedPayload from jws directly
}
```

# Support

App Store Server API [1.8](https://developer.apple.com/documentation/appstoreserverapi)

# License

appstore is licensed under the MIT.

