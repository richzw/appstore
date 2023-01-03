
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

### Look Up Order ID

```go
import(
    "github.com/richzw/appstore/api"
)

// ACCOUNTPRIVATEKEY is the key file generated from previous step
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `

func main() {
    c := &StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    invoiceOrderId := "FAKEORDERID"

    a := NewStoreClient(c)
    rsp, err := a.LookupOrderID(invoiceOrderId)

    orders, err := a.ParseSignedTransactions(context.TODO(), rsp.SignedTransactions)
}
```

### Get Transaction History

```go
import(
    "github.com/richzw/appstore/api"
)

// ACCOUNTPRIVATEKEY is the key file generated from previous step
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `

func main() {
    c := &StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    originalTransactionId := "FAKEORDERID"
    a := NewStoreClient(c)
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
    "github.com/richzw/appstore/api"
)

// ACCOUNTPRIVATEKEY is the key file generated from previous step
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `

func main() {
    c := &StoreConfig{
        KeyContent: []byte(ACCOUNTPRIVATEKEY),
        KeyID:      "FAKEKEYID",
        BundleID:   "fake.bundle.id",
        Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",
        Sandbox:    false,
    }
    originalTransactionId := "FAKEORDERID"
    a := NewStoreClient(c)
    gotRsp, err := a.GetRefundHistory(context.TODO(), originalTransactionId)

    for _, rsp := range gotRsp {
       trans, err := a.ParseSignedTransactions(rsp.SignedTransactions)
    }
}
```

# Support

App Store Server API [1.6](https://developer.apple.com/documentation/appstoreserverapi)

# License

appstore is licensed under the MIT.

