# CNBC Quote SDK

## Usage
```go
import "github.com/ipanardian/cnbc-quote"

var res cnbc.CnbcData
cl := cnbc.NewCnbcQuote(CnbcURI)
headers := make(map[string]string)
err = cl.GetQuote(&res, headers, cnbc.CnbcQuoteRequest{
    Symbols:       symbol,
    RequestMethod: "itv",
    NoForm:        1,
    PartnerID:     2,
    Fund:          1,
    Output:        "json",
    Exthrs:        1,
    Events:        1,
})
```