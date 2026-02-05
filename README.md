# jquants-go

Go client library for the [J-Quants API](https://jpx-jquants.com/), providing access to Japanese stock market data from the Tokyo Stock Exchange (TSE).

## Installation

```bash
go get github.com/S-Shiga/jquants-go/v2
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

    "github.com/S-Shiga/jquants-go/v2"
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
    jquants.WithRetryInterval(10 * time.Second),    // retry interval for 500 errors (default: 5s)
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

### Not Yet Implemented

The following J-Quants API endpoints are not yet implemented in this library:

- Morning Session Stock Prices
- Outstanding Short Selling Positions Reported
- Margin Trading Outstanding (Breakdown)
- Breakdown Trading

### Channel API (Streaming)

Methods with a `WithChannel` suffix (`StockPriceWithChannel`, `IndexOptionPriceWithChannel`) stream results through a channel instead of returning a slice. This is useful when processing large datasets incrementally.

- The caller must create the channel and pass it in.
- The channel is **automatically closed** when all pages have been sent or when an error occurs.
- The method respects context cancellation via the `loopTimeout` setting.
- Errors are returned from the goroutine; use a separate goroutine to call the method and check the error after the channel is drained.

## Codes Package

The `codes` package provides constants for market sections, sector codes, and index codes.

```go
import "github.com/S-Shiga/jquants-go/v2/codes"

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

The client automatically retries on HTTP 500 errors with a configurable interval.

## License

MIT License - see [LICENSE](LICENSE) for details.
