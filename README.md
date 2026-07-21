# jquants-go

Go client library for the [J-Quants API](https://jpx-jquants.com/), providing access to Japanese stock market data from the Tokyo Stock Exchange (TSE).

## Installation

```bash
go get github.com/s-shiga/jquants-go/v2
```

## Quick Start

1. Get your API key from [J-Quants](https://jpx-jquants.com/)
2. Use the client:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/s-shiga/jquants-go/v2"
)

func main() {
    client := jquants.NewClient(jquants.BaseURL, os.Getenv("J_QUANTS_API_KEY"))

    ctx := context.Background()

    // Get issue information for all listed securities
    issues, err := client.IssueInformation(ctx, jquants.IssueInformationRequest{})
    if err != nil {
        log.Fatal(err)
    }

    for _, issue := range issues {
        fmt.Printf("%s: %s\n", issue.Code, issue.CompanyName)
    }
}
```

## Client Options

`NewClient` accepts functional options to customize behavior:

```go
client := jquants.NewClient(
    jquants.BaseURL,
    os.Getenv("J_QUANTS_API_KEY"),
    jquants.WithHTTPClient(customHTTPClient),       // custom *http.Client (default: http.DefaultClient)
    jquants.WithRetryInterval(10 * time.Second),    // retry interval for retryable errors (default: 5s)
    jquants.WithLoopTimeout(60 * time.Second),      // timeout for paginated requests (default: 20s)
)
```

## Available APIs

### Equities

#### Issue Information

Retrieves master data for listed securities from the `/equities/master` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/eq-master) for details.

```go
// Get all securities
issues, err := client.IssueInformation(ctx, jquants.IssueInformationRequest{})

// Filter by code
code := "7203"
issues, err := client.IssueInformation(ctx, jquants.IssueInformationRequest{
    Code: &code,
})

// Filter by date
date := "2024-01-15"
issues, err := client.IssueInformation(ctx, jquants.IssueInformationRequest{
    Date: &date,
})
```

#### Minute Stock Prices

Retrieves 1-minute OHLCV bars from the `/equities/bars/minute` endpoint (minute-bars add-on plan).
See [API reference](https://jpx-jquants.com/en/spec/eq-bars-minute) for details.

```go
code := "8697"
from, to := "2024-01-15", "2024-01-15"
bars, err := client.MinuteStockPrice(ctx, jquants.MinuteStockPriceRequest{
    Code: &code,
    From: &from,
    To:   &to,
})
```

A `MinuteStockPriceWithChannel` streaming variant is also available.

#### Morning Session Stock Prices

Retrieves the current day's morning session OHLCV data from the `/equities/bars/daily/am` endpoint (Premium plan).
See [API reference](https://jpx-jquants.com/en/spec/eq-bars-daily-am) for details.

Outside the morning-session publication window, the API responds with HTTP 210 and the method returns a `NoContent` error.

```go
code := "7203"
prices, err := client.MorningSessionStockPrice(ctx, jquants.MorningSessionStockPriceRequest{
    Code: &code,
})
if err != nil {
    var noContent jquants.NoContent
    if errors.As(err, &noContent) {
        // outside the morning session window
    }
}
```

#### Stock Prices

Retrieves daily OHLCV data for stocks from the `/equities/bars/daily` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/eq-bars-daily) for details.

Price fields (`Open`, `High`, `Low`, `Close`, etc.) use `*json.Number` because the API returns numeric strings. Volume fields may be `nil` when no trading occurred.

```go
// Get prices for a specific stock
code := "7203"
prices, err := client.StockPrice(ctx, jquants.StockPriceRequest{
    Code: &code,
})

// With date range
from, to := "2024-01-01", "2024-01-31"
prices, err := client.StockPrice(ctx, jquants.StockPriceRequest{
    Code: &code,
    From: &from,
    To:   &to,
})

// Get all stocks for a specific date
date := "2024-01-15"
prices, err := client.StockPrice(ctx, jquants.StockPriceRequest{
    Date: &date,
})

// Stream results via channel
ch := make(chan jquants.StockPrice)
go func() {
    err := client.StockPriceWithChannel(ctx, jquants.StockPriceRequest{Code: &code}, ch)
    if err != nil {
        log.Println(err)
    }
}()
for price := range ch {
    fmt.Printf("%s: %v\n", price.Date, price.Close)
}
```

#### Investor Type Trading

Retrieves weekly trading data by investor category from the `/equities/investor-types` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/eq-investor-types) for details.

```go
// Get all investor type data for a date range
from, to := "2024-01-01", "2024-03-31"
data, err := client.InvestorType(ctx, jquants.InvestorTypeRequest{
    From: &from,
    To:   &to,
})

// Filter by market section
section := codes.SectionPrime
data, err := client.InvestorType(ctx, jquants.InvestorTypeRequest{
    Section: &section,
    From:    &from,
    To:      &to,
})
```

#### Earnings Calendar

Retrieves upcoming earnings announcement dates from the `/equities/earnings-calendar` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/eq-earnings-cal) for details.

```go
calendar, err := client.EarningsCalendar(ctx, jquants.EarningsCalendarRequest{})
```

### Markets

#### Margin Trading Outstanding

Retrieves margin trading balance data from the `/markets/margin-interest` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/mkt-margin-int) for details.

```go
code := "7203"
data, err := client.MarginTradingOutstanding(ctx, jquants.MarginTradingOutstandingRequest{
    Code: &code,
})
```

#### Short Selling Value

Retrieves short selling turnover data by sector from the `/markets/short-ratio` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/mkt-short-ratio) for details.

```go
date := "2024-01-15"
data, err := client.ShortSellingValue(ctx, jquants.ShortSellingValueRequest{
    Date: &date,
})
```

#### Trading Calendar

Retrieves the TSE trading calendar from the `/markets/calendar` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/mkt-calendar) for details.

```go
from, to := "2024-01-01", "2024-12-31"
calendar, err := client.TradingCalendar(ctx, jquants.TradingCalendarRequest{
    From: &from,
    To:   &to,
})
```

#### Outstanding Short Positions

Retrieves outstanding short selling position reports from the `/markets/short-sale-report` endpoint (Standard plan or above).
See [API reference](https://jpx-jquants.com/en/spec/mkt-short-sale) for details.

```go
calcDate := "2024-03-01"
positions, err := client.OutstandingShortPosition(ctx, jquants.OutstandingShortPositionRequest{
    CalculationDate: &calcDate,
})
```

#### Margin Alert

Retrieves issues under margin trading restrictions from the `/markets/margin-alert` endpoint (Standard plan or above).
See [API reference](https://jpx-jquants.com/en/spec/mkt-margin-alert) for details.

Change and ratio fields may be `nil` when the API reports them as unavailable (`"-"`) or not applicable for ETFs (`"*"`).

```go
date := "2024-02-08"
alerts, err := client.MarginAlert(ctx, jquants.MarginAlertRequest{
    Date: &date,
})
```

#### Breakdown Trading

Retrieves the breakdown of trading value and volume (long selling/buying, short selling, margin transactions) from the `/markets/breakdown` endpoint (Premium plan).
See [API reference](https://jpx-jquants.com/en/spec/mkt-breakdown) for details.

```go
code := "7203"
from, to := "2024-01-01", "2024-01-31"
data, err := client.BreakdownTrading(ctx, jquants.BreakdownTradingRequest{
    Code: &code,
    From: &from,
    To:   &to,
})
```

A `BreakdownTradingWithChannel` streaming variant is also available.

### Indices

#### Index Prices

Retrieves daily OHLC data for market indices from the `/indices/bars/daily` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/idx-bars-daily) for details.

```go
code := "0000" // TOPIX
prices, err := client.IndexPrice(ctx, jquants.IndexPriceRequest{
    Code: &code,
})
```

#### TOPIX Prices

Retrieves TOPIX index prices directly from the `/indices/bars/daily/topix` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/idx-bars-daily-topix) for details.

```go
from, to := "2024-01-01", "2024-01-31"
prices, err := client.TopixPrices(ctx, jquants.TopixPriceRequest{
    From: &from,
    To:   &to,
})
```

### Derivatives

#### Index Option Prices

Retrieves Nikkei 225 index option prices from the `/derivatives/bars/daily/options/225` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/deriv-bars-daily-options-225) for details.

```go
data, err := client.IndexOptionPrice(ctx, jquants.IndexOptionPriceRequest{
    Date: "2024-01-15",
})

// Stream results via channel
ch := make(chan jquants.IndexOptionPrice)
go func() {
    err := client.IndexOptionPriceWithChannel(ctx, jquants.IndexOptionPriceRequest{Date: "2024-01-15"}, ch)
    if err != nil {
        log.Println(err)
    }
}()
for option := range ch {
    fmt.Printf("%s: Strike=%d, Close=%v\n", option.Code, option.StrikePrice, option.WholeDayClose)
}
```

#### Futures Prices

Retrieves daily futures prices for all products (TOPIX futures, Nikkei 225 futures, etc.) from the `/derivatives/bars/daily/futures` endpoint (Premium plan).
See [API reference](https://jpx-jquants.com/en/spec/drv-bars-daily-fut) for details.

Price fields use `*json.Number` (prices can be fractional and may be `nil` when no trading occurred). Whole-day, morning, night, and day session OHLC are all provided.

```go
category := "TOPIXF"
futures, err := client.FuturesPrice(ctx, jquants.FuturesPriceRequest{
    Date:     "2024-07-23",
    Category: &category,
})
```

A `FuturesPriceWithChannel` streaming variant is also available.

#### Option Prices (All Underlyings)

Retrieves daily option prices for all underlyings (TOPIX options, Nikkei 225 options, securities options, etc.) from the `/derivatives/bars/daily/options` endpoint (Premium plan).
See [API reference](https://jpx-jquants.com/en/spec/drv-bars-daily-opt) for details.

```go
category := "NK225E"
options, err := client.OptionPrice(ctx, jquants.OptionPriceRequest{
    Date:     "2024-07-23",
    Category: &category,
})
```

For securities options, pass the underlying security code via `Code`. A `OptionPriceWithChannel` streaming variant is also available.

### Financials

#### Financial Summary

Retrieves financial data summaries (results, forecasts, dividends, per-share metrics) from the `/fins/summary` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/fin-summary) for details.

The API returns numeric values as strings; this library parses them into `*float64` fields, which are `nil` when no data is available.

```go
code := "7203"
summaries, err := client.FinancialSummary(ctx, jquants.FinancialSummaryRequest{
    Code: &code,
})
```

#### Financial Statement Details

Retrieves detailed financial statement line items (balance sheet, income statement, etc.) from the `/fins/details` endpoint (Premium plan).
See [API reference](https://jpx-jquants.com/en/spec/fin-details) for details.

Line items are returned in the `FinancialStatement` field as a `map[string]string` keyed by verbose English XBRL labels.

```go
code := "7203"
details, err := client.FinancialDetails(ctx, jquants.FinancialDetailsRequest{
    Code: &code,
})
```

#### Cash Dividend Data

Retrieves cash dividend announcements (record dates, ex-dates, dividend rates, etc.) from the `/fins/dividend` endpoint (Premium plan).
See [API reference](https://jpx-jquants.com/en/spec/fin-dividend) for details.

Numeric fields use `*float64` and are `nil` when the value is undetermined (`"-"`) or not applicable (`""`).

```go
code := "7203"
dividends, err := client.Dividend(ctx, jquants.DividendRequest{
    Code: &code,
})
```

### EDINET Filings

The EDINET endpoints (Standard plan or above) share a common `EdinetRequest` with `EdinetCode`, `Code`, and `Date` filters. Specifying both `EdinetCode` and `Code` at once is rejected by the API.

#### Major Shareholders

Retrieves major shareholder data from annual securities reports via the `/edinet/major-shareholders` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/edinet-major-shareholders) for details.

```go
code := "7203"
reports, err := client.MajorShareholders(ctx, jquants.EdinetRequest{Code: &code})
```

#### Cross-Shareholdings

Retrieves cross-shareholding (policy holdings) data from the `/edinet/cross-shareholdings` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/edinet-cross-shareholdings) for details.

```go
code := "7203"
reports, err := client.CrossShareholdings(ctx, jquants.EdinetRequest{Code: &code})
```

#### Large Volume Shareholders

Retrieves large volume holding reports from the `/edinet/large-volume-shareholders` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/edinet-large-volume-shareholders) for details.

```go
code := "7203"
reports, err := client.LargeVolumeShareholders(ctx, jquants.EdinetRequest{Code: &code})
```

### TDnet Timely Disclosure (Add-on)

The TDnet endpoints require the TimelyDisclosure add-on plan.

#### Disclosure Index List

Retrieves the TDnet disclosure index from the `/td/list` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/td-list) for details.

```go
date := "2024-04-01"
disclosures, err := client.TimelyDisclosure(ctx, jquants.TimelyDisclosureRequest{
    Date: &date,
})
```

`DisclosureStatus` is `nil` for new disclosures, `"revision"` for corrections, and `"delete"` for deletions. A `TimelyDisclosureWithChannel` streaming variant is also available.

#### Disclosure Files

Retrieves signed download URLs for a disclosure's documents (full PDF, summary PDF, XBRL) from the `/td/files` endpoint. URLs expire in 15 minutes.
See [API reference](https://jpx-jquants.com/en/spec/td-files) for details.

```go
files, err := client.TimelyDisclosureFiles(ctx, jquants.TimelyDisclosureFilesRequest{
    DisclosureNumber: "20250401130100",
})
if files.Files.PDF != nil {
    // download *files.Files.PDF
}
```

#### Disclosure CSV Download

Retrieves a signed URL for a gzip CSV covering five years of disclosures from the `/td/bulk` endpoint.
See [API reference](https://jpx-jquants.com/en/spec/td-bulk) for details.

```go
bulk, err := client.TimelyDisclosureBulk(ctx)
// bulk.URL is a short-lived download link; bulk.LastUpdated is ISO 8601
```

### Bulk Download

Lists and fetches whole-dataset gzip CSV files via the `/bulk/list` and `/bulk/get` endpoints. Tick-level stock trades (Stock Prices (Tick) add-on) have no REST endpoint and are delivered exclusively this way, under keys prefixed `equities/trades/`.
See [API reference](https://jpx-jquants.com/en/spec/bulk-list) for details.

```go
// List files for a date (or by endpoint name)
date := "2024-01-15"
files, err := client.BulkList(ctx, jquants.BulkListRequest{Date: &date})

// Get a signed download URL (valid ~5 minutes, single-use)
url, err := client.BulkGet(ctx, jquants.BulkGetRequest{Key: &files[0].Key})
```

### Channel API (Streaming)

Methods with a `WithChannel` suffix (`StockPriceWithChannel`, `MinuteStockPriceWithChannel`, `IndexOptionPriceWithChannel`, `OptionPriceWithChannel`, `FuturesPriceWithChannel`, `BreakdownTradingWithChannel`, `TimelyDisclosureWithChannel`) stream results through a channel instead of returning a slice. This is useful when processing large datasets incrementally.

- The caller must create the channel and pass it in.
- The channel is **automatically closed** when all pages have been sent or when an error occurs.
- The method respects context cancellation via the `loopTimeout` setting.
- Errors are returned from the goroutine; use a separate goroutine to call the method and check the error after the channel is drained.

## Codes Package

The `codes` package provides constants for market sections, sector codes, and index codes.

```go
import "github.com/s-shiga/jquants-go/v2/codes"

// Market sections
section := codes.SectionPrime

// Sector codes (33-sector classification)
sector := codes.Sector33Banks

// Index codes
index := codes.IndexTOPIX
```

Available constants:

- **Sections**: `SectionPrime`, `SectionStandard`, `SectionGrowth`, `SectionTokyoNagoya` (current market segments)
- **Legacy sections**: `SectionTSE1st`, `SectionTSE2nd`, `SectionMothers`, `SectionJASDAQ` (pre-2022 market restructuring)
- **Sector33 codes**: `Sector33Banks`, `Sector33Chemicals`, `Sector33Construction`, etc. (all 33 TSE sector classifications)
- **Index codes**: `IndexTOPIX`, `IndexTOPIXCore30`, `IndexTOPIX500`, `IndexREIT`, TOPIX-17 sector indices, etc.

Convenience slices:

- `codes.Sections` — current market sections (Prime, Standard, Growth, TokyoNagoya)
- `codes.SectionsAll` — all market sections including legacy ones
- `codes.Sector33Codes` — all 33-sector classification codes

Note: The 17-sector classification (`Sector17Code` in `IssueInformation`) uses integer codes returned by the API directly. The TOPIX-17 index codes (e.g., `IndexTOPIX17FOODS`, `IndexTOPIX17Banks`) are available in the codes package.

## Error Handling

The client returns typed errors for different HTTP status codes:

```go
prices, err := client.StockPrice(ctx, req)
if err != nil {
    var badReq jquants.BadRequest
    var unauthorized jquants.Unauthorized
    var forbidden jquants.Forbidden

    switch {
    case errors.As(err, &badReq):
        log.Println("Bad request:", err)
    case errors.As(err, &unauthorized):
        log.Println("Invalid API key:", err)
    case errors.As(err, &forbidden):
        log.Println("Access forbidden:", err)
    default:
        log.Println("Error:", err)
    }
}
```

A `NoContent` error (HTTP 210) is returned by endpoints that have no data for the requested window, such as morning session prices outside publication hours; it is not retried.

The client automatically retries on HTTP 429, 500, 502, 503, and 504 errors with a configurable interval (`TooManyRequests`, `InternalServerError`, `BadGateway`, `ServiceUnavailable`, `GatewayTimeout`). For 429 responses, a `Retry-After` header is honored when present.

Transient transport-level failures are retried under the same policy and surface as `TransientTransportError`: a response body truncated mid-stream (unexpected EOF while decoding a very large page), a connection reset, or a network error from the HTTP round trip itself. Only the failing page is re-requested (with the same pagination key). Retries are bounded by `WithLoopTimeout`; caller cancellation (`context.Canceled` / `context.DeadlineExceeded`) is never retried and remains fatal.

## License

MIT License - see [LICENSE](LICENSE) for details.
