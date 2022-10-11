
[App Store Server API](https://developer.apple.com/documentation/appstoreserverapi) Golang Client
================

# Quick Start

### Look Up Order ID

```go
import(
    "github.com/richzw/appstore/api"
)

const ACCOUNTKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `

func main() {
    c := &StoreConfig{
        KeyContent: []byte(ACCOUNTKEY),
        KeyID:      "SKEYID",
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


# License

-------------------

appstore is licensed under the MIT.

