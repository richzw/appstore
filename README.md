
[App Store Server API](https://developer.apple.com/documentation/appstoreserverapi) Golang Client
================

# Quick Start

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

    orders, err := a.ParseSignedTransactions(rsp.SignedTransactions)
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
    gotRsp, err := a.GetTransactionHistory(originalTransactionId, query)

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
    gotRsp, err := a.GetRefundHistory(originalTransactionId)

    for _, rsp := range gotRsp {
       trans, err := a.ParseSignedTransactions(rsp.SignedTransactions)
    }
}
```



# License

-------------------

appstore is licensed under the MIT.

