# jquants-go

Go client library for the [J-Quants API](https://jpx-jquants.com/), providing access to Japanese stock market data from the Tokyo Stock Exchange (TSE).

## Installation

```bash
go get github.com/S-Shiga/jquants-go/v2
```

## Quick Start

1. Get your API key from [J-Quants](https://jpx-jquants.com/)
2. Set the environment variable:

```bash
export J_QUANTS_API_KEY=your_api_key
```

3. Use the client:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/S-Shiga/jquants-go/v2"
)

func main() {
    client, err := jquants.NewClient(http.DefaultClient)
    if err != nil {
        log.Fatal(err)
    }

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

## Available APIs

### Equities

#### Issue Information

Retrieves master data for listed securities.

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

Retrieves daily OHLCV data for stocks.

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

Retrieves weekly trading data by investor category.

```go
from, to := "2024-01-01", "2024-03-31"
data, err := client.InvestorType(ctx, jquants.InvestorTypeRequest{
    From: &from,
    To:   &to,
})
```

### Markets

#### Margin Trading Outstanding

Retrieves margin trading balance data.

```go
code := "7203"
data, err := client.MarginTradingOutstanding(ctx, jquants.MarginTradingOutstandingRequest{
    Code: &code,
})
```

#### Short Selling Value

Retrieves short selling turnover data by sector.

```go
date := "2024-01-15"
data, err := client.ShortSellingValue(ctx, jquants.ShortSellingValueRequest{
    Date: &date,
})
```

#### Trading Calendar

Retrieves the TSE trading calendar.

```go
from, to := "2024-01-01", "2024-12-31"
calendar, err := client.TradingCalendar(ctx, jquants.TradingCalendarRequest{
    From: &from,
    To:   &to,
})
```

### Indices

#### Index Prices

Retrieves daily OHLC data for market indices.

```go
code := "0000" // TOPIX
prices, err := client.IndexPrice(ctx, jquants.IndexPriceRequest{
    Code: &code,
})
```

#### TOPIX Prices

Retrieves TOPIX index prices directly.

```go
from, to := "2024-01-01", "2024-01-31"
prices, err := client.TopixPrices(ctx, jquants.TopixPriceRequest{
    From: &from,
    To:   &to,
})
```

### Derivatives

#### Index Option Prices

Retrieves Nikkei 225 index option prices.

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

- **Sections**: `SectionPrime`, `SectionStandard`, `SectionGrowth`, `SectionTSE1st`, `SectionTSE2nd`, etc.
- **Sector33 codes**: `Sector33Banks`, `Sector33Chemicals`, `Sector33Construction`, etc.
- **Index codes**: `IndexTOPIX`, `IndexTOPIXCore30`, `IndexTOPIX500`, `IndexREIT`, etc.

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

## Configuration

The `Client` struct has configurable fields:

```go
client, _ := jquants.NewClient(http.DefaultClient)

// Customize retry interval for 500 errors (default: 5s)
client.RetryInterval = 10 * time.Second

// Customize timeout for paginated requests (default: 20s)
client.LoopTimeout = 60 * time.Second
```

## License

MIT License - see [LICENSE](LICENSE) for details.
