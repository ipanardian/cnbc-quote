package cnbc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/google/go-querystring/query"
)

type Cnbc struct {
	BaseUrl string
	pool    sync.Pool
}

type CnbcRequest struct {
	Data        *string
	QueryParams *string
	Headers     map[string]string
}

type CnbcQuoteRequest struct {
	Symbols       string `json:"symbols" url:"symbols"`
	RequestMethod string `json:"requestMethod" url:"requestMethod"`
	NoForm        int8   `json:"noform" url:"noform"`
	PartnerID     int    `json:"partnerId" url:"partnerId"`
	Fund          int8   `json:"fund" url:"fund"`
	Exthrs        int8   `json:"exthrs" url:"exthrs"`
	Output        string `json:"output" url:"output"`
	Events        int8   `json:"events" url:"events"`
}

type CnbcData struct {
	FormattedQuoteResult struct {
		FormattedQuote []struct {
			Symbol                  string `json:"symbol"`
			SymbolType              string `json:"symbolType"`
			Code                    int    `json:"code"`
			Name                    string `json:"name"`
			ShortName               string `json:"shortName"`
			OnAirName               string `json:"onAirName"`
			AltName                 string `json:"altName"`
			Last                    string `json:"last"`
			LastTimedate            string `json:"last_timedate"`
			LastTime                string `json:"last_time"`
			Changetype              string `json:"changetype"`
			Type                    string `json:"type"`
			SubType                 string `json:"subType"`
			Exchange                string `json:"exchange"`
			Source                  string `json:"source"`
			Open                    string `json:"open"`
			High                    string `json:"high"`
			Low                     string `json:"low"`
			Change                  string `json:"change"`
			ChangePct               string `json:"change_pct"`
			Provider                string `json:"provider"`
			PreviousDayClosing      string `json:"previous_day_closing"`
			AltSymbol               string `json:"altSymbol"`
			RealTime                string `json:"realTime"`
			Curmktstatus            string `json:"curmktstatus"`
			Yrhiprice               string `json:"yrhiprice"`
			Yrhidate                string `json:"yrhidate"`
			Yrloprice               string `json:"yrloprice"`
			Yrlodate                string `json:"yrlodate"`
			Streamable              string `json:"streamable"`
			BondLastPrice           string `json:"bond_last_price"`
			BondChangePrice         string `json:"bond_change_price"`
			BondChangePctPrice      string `json:"bond_change_pct_price"`
			BondOpenPrice           string `json:"bond_open_price"`
			BondHighPrice           string `json:"bond_high_price"`
			BondLowPrice            string `json:"bond_low_price"`
			BondPrevDayClosingPrice string `json:"bond_prev_day_closing_price"`
			BondChangetype          string `json:"bond_changetype"`
			MaturityDate            string `json:"maturity_date"`
			Coupon                  string `json:"coupon"`
			IssueID                 string `json:"issue_id"`
			CountryCode             string `json:"countryCode"`
			TimeZone                string `json:"timeZone"`
			FeedSymbol              string `json:"feedSymbol"`
			Portfolioindicator      string `json:"portfolioindicator"`
			EventData               struct {
				Yrhiind  string `json:"yrhiind"`
				Yrloind  string `json:"yrloind"`
				IsHalted string `json:"is_halted"`
			} `json:"EventData"`
		} `json:"FormattedQuote"`
	} `json:"FormattedQuoteResult"`
}

func NewCnbcQuote(baseUrl string) *Cnbc {
	return &Cnbc{
		BaseUrl: baseUrl,
		pool: sync.Pool{
			New: func() interface{} {
				timeout := 5000 * time.Millisecond
				return httpclient.NewClient(
					httpclient.WithHTTPTimeout(timeout),
				)
			},
		},
	}
}

func (m *Cnbc) GetQuote(res any, headers map[string]string, data CnbcQuoteRequest) error {
	qs, e := query.Values(data)
	if e != nil {
		return e
	}
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	jsonStr := qs.Encode()
	return m.getAndUnmarshalJson(res, "/quote-html-webservice/restQuote/symbolType/symbol", CnbcRequest{
		QueryParams: &jsonStr,
		Headers:     headers,
	})
}

func (m *Cnbc) getAndUnmarshalJson(res any, path string, data CnbcRequest) error {
	client := m.pool.Get().(*httpclient.Client)
	url := fmt.Sprintf("%s%s", m.BaseUrl, path)

	if data.QueryParams != nil {
		url = url + "?" + *data.QueryParams
	}

	jB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	dataReader := strings.NewReader("")
	if data.Data != nil {
		dataReader = strings.NewReader(string(jB))
	}

	req, err := http.NewRequest(http.MethodGet, url, dataReader)
	if err != nil {
		return err
	}

	if data.Headers != nil {
		for key, value := range data.Headers {
			req.Header.Set(key, value)
		}
	}

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	return nil
}
