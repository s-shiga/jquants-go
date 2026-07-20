# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `jquants-go`, a Go client library for the J-Quants API, which provides access to Japanese stock market data from the Tokyo Stock Exchange (TSE).

## Build and Test Commands

```bash
# Run all tests (requires J_QUANTS_API_KEY environment variable)
go test ./...

# Run a specific test
go test -run TestClient_IssueInformation

# Build the package
go build ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

**Note:** Tests make real API calls and require the `J_QUANTS_API_KEY` environment variable to be set. There are no mocked tests.

## Architecture

### Client Pattern

The library uses a single `Client` struct (`client.go`) that holds HTTP client, base URL, API key, and retry/timeout settings. All API methods are methods on this `Client`.

The constructor `NewClient(baseURL, apiKey string, opts ...Option)` returns `*Client` (no error). It uses a functional options pattern with `WithHTTPClient`, `WithRetryInterval`, and `WithLoopTimeout`.

### API Method Structure

Each API endpoint follows a consistent pattern:

1. **Response struct** - Go struct with custom `UnmarshalJSON` to handle J-Quants API quirks (e.g., numeric strings, floats that should be ints, abbreviated JSON keys like `"O"` → `Open`)
2. **Request struct** - Public struct with optional filter parameters (uses `*string` for optional fields)
3. **Parameters struct** - Internal struct embedding the request, adding `PaginationKey`
4. **`values()` method** - Implements the `parameters` interface to convert to URL query params
5. **`send*Request` method** - Internal method to make a single paginated request
6. **Public method** - Loops through pagination, handles 500 error retries, returns complete data

### Pagination Handling

APIs that return large datasets use pagination. The client automatically fetches all pages in a loop until `pagination_key` is nil. Some methods also offer `*WithChannel` variants for streaming results (`StockPriceWithChannel`, `IndexOptionPriceWithChannel`, `OptionPriceWithChannel`, `FuturesPriceWithChannel`, `BreakdownTradingWithChannel`).

### Error Types

Custom error types in `client.go` wrap HTTP status codes: `NoContent` (210), `BadRequest`, `Unauthorized`, `Forbidden`, `PayloadTooLarge`, `TooManyRequests`, `InternalServerError`, `BadGateway`, `ServiceUnavailable`, `GatewayTimeout`. The client auto-retries on 429/500/502/503/504, honoring `Retry-After` for 429.

### Module Organization

- `client.go` - Client initialization, HTTP request handling, error types, pagination helpers (`fetchAllPages`, `fetchAllPagesWithChannel`)
- `generics.go` - Generic `Request` and `Response` interfaces
- `equity.go` - Stock-related APIs:
  - Issue information (`/equities/master`)
  - Stock prices (`/equities/bars/daily`)
  - Morning session stock prices (`/equities/bars/daily/am`, Premium)
  - Earnings calendar (`/equities/earnings-calendar`)
  - Investor type trading (`/equities/investor-types`)
- `markets.go` - Market data APIs:
  - Margin trading outstanding (`/markets/margin-interest`)
  - Short selling value (`/markets/short-ratio`)
  - Trading calendar (`/markets/calendar`)
  - Outstanding short positions (`/markets/short-sale-report`, Standard)
  - Margin alert (`/markets/margin-alert`, Standard)
  - Breakdown trading (`/markets/breakdown`, Premium)
- `indices.go` - Index APIs:
  - Index prices (`/indices/bars/daily`)
  - TOPIX prices (`/indices/bars/daily/topix`)
- `option.go` - Options APIs:
  - Index option prices (`/derivatives/bars/daily/options/225`)
  - Option prices, all underlyings (`/derivatives/bars/daily/options`, Premium)
- `future.go` - Futures APIs:
  - Futures prices (`/derivatives/bars/daily/futures`, Premium)
- `fins.go` - Financial data APIs:
  - Financial summary (`/fins/summary`)
  - Financial statement details (`/fins/details`, Premium)
  - Cash dividend data (`/fins/dividend`, Premium)
- `edinet.go` - EDINET filing APIs (Standard):
  - Major shareholders (`/edinet/major-shareholders`)
  - Cross-shareholdings (`/edinet/cross-shareholdings`)
  - Large volume shareholders (`/edinet/large-volume-shareholders`)
- `codes/codes.go` - Constants for market sections, 33-sector codes, and index codes
- `testutil.go` - Test helper that reads `J_QUANTS_API_KEY` from env and creates a client

### JSON Unmarshaling

The J-Quants API returns some numeric fields as strings and uses abbreviated JSON keys (e.g., `"O"`, `"H"`, `"L"`, `"C"` for OHLC prices, `"CoName"` for company name). Custom `UnmarshalJSON` methods translate these to proper Go types with descriptive field names. Price fields use `*json.Number` and volume fields may be `nil` when no trading occurred.
