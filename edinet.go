package jquants

import (
	"context"
	"net/url"
)

// edinetValues builds query parameters common to the EDINET endpoints, all of
// which accept optional edinet_code, code, and date filters plus pagination.
func edinetValues(edinetCode, code, date, paginationKey *string) (url.Values, error) {
	v := url.Values{}
	if edinetCode != nil {
		v.Add("edinet_code", *edinetCode)
	}
	if code != nil {
		v.Add("code", *code)
	}
	if date != nil {
		v.Add("date", *date)
	}
	if paginationKey != nil {
		v.Add("pagination_key", *paginationKey)
	}
	return v, nil
}

// EdinetRequest specifies the filter parameters shared by the EDINET endpoints.
// All parameters are optional. Note that supplying both EdinetCode and Code
// together is rejected by the API with a 400 response.
type EdinetRequest struct {
	// EdinetCode filters by EDINET code.
	EdinetCode *string
	// Code filters by security code.
	Code *string
	// Date filters by submission date.
	Date *string
}

// MajorShareholder represents a single major shareholder entry within a
// major-shareholders filing.
type MajorShareholder struct {
	// Rank is the shareholder's rank in the filing (JSON key "Rank").
	Rank int `json:"Rank"`
	// HolderName is the shareholder's name (JSON key "HldrName").
	HolderName string `json:"HldrName"`
	// HolderAddress is the shareholder's address (JSON key "HldrAddr").
	HolderAddress string `json:"HldrAddr"`
	// SharesHeld is the number of shares held (JSON key "ShsHeld").
	SharesHeld int64 `json:"ShsHeld"`
	// SharesRatio is the shareholding ratio as a decimal fraction (JSON key "ShsRatio").
	SharesRatio float64 `json:"ShsRatio"`
}

// MajorShareholders represents a major-shareholders filing extracted from EDINET.
type MajorShareholders struct {
	// DocumentID is the EDINET document ID (JSON key "DocId").
	DocumentID string `json:"DocId"`
	// Code is the security code (JSON key "Code").
	Code string `json:"Code"`
	// EdinetCode is the filer's EDINET code (JSON key "EdinetCode").
	EdinetCode string `json:"EdinetCode"`
	// FilerName is the filer's name in Japanese (JSON key "FilerName").
	FilerName string `json:"FilerName"`
	// FilerNameEnglish is the filer's name in English (JSON key "FilerNameEn").
	FilerNameEnglish string `json:"FilerNameEn"`
	// DocumentTypeCode is the EDINET document type code (JSON key "DocTypeCode").
	DocumentTypeCode string `json:"DocTypeCode"`
	// SubmissionDate is the filing submission date (JSON key "SubDate").
	SubmissionDate string `json:"SubDate"`
	// SubmissionTime is the filing submission time (JSON key "SubTime").
	SubmissionTime string `json:"SubTime"`
	// PeriodStart is the reporting period start date (JSON key "PerSt").
	PeriodStart string `json:"PerSt"`
	// PeriodEnd is the reporting period end date (JSON key "PerEn").
	PeriodEnd string `json:"PerEn"`
	// Holders is the list of major shareholders (JSON key "Hldrs").
	Holders []MajorShareholder `json:"Hldrs"`
}

type majorShareholdersParameters struct {
	EdinetRequest
	PaginationKey *string
}

func (p majorShareholdersParameters) values() (url.Values, error) {
	return edinetValues(p.EdinetCode, p.Code, p.Date, p.PaginationKey)
}

type majorShareholdersResponse struct {
	Data          []MajorShareholders `json:"data"`
	PaginationKey *string             `json:"pagination_key"`
}

func (r majorShareholdersResponse) Items() []MajorShareholders { return r.Data }
func (r majorShareholdersResponse) NextPageKey() *string       { return r.PaginationKey }

// MajorShareholders retrieves major-shareholders filings from the /edinet/major-shareholders endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/edinet-major-shareholders for API details.
func (c *Client) MajorShareholders(ctx context.Context, req EdinetRequest) ([]MajorShareholders, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (majorShareholdersResponse, error) {
		params := majorShareholdersParameters{EdinetRequest: req, PaginationKey: paginationKey}
		return getJSON[majorShareholdersResponse](ctx, c, "/edinet/major-shareholders", params)
	})
}

// CrossShareholdingIssue represents a single issuer held as a cross-shareholding,
// appearing in either the specified (Spec) or deemed (Deem) holding lists.
type CrossShareholdingIssue struct {
	// IssuerName is the issuer's name (JSON key "IsrName").
	IssuerName string `json:"IsrName"`
	// IssuerCode is the issuer's security code, or nil (JSON key "IsrCode").
	IssuerCode *string `json:"IsrCode"`
	// IssuerEdinetCode is the issuer's EDINET code, or nil (JSON key "IsrEdinetCode").
	IssuerEdinetCode *string `json:"IsrEdinetCode"`
	// CurrentShares is the number of shares held in the current period, or nil (JSON key "CurShs").
	CurrentShares *int64 `json:"CurShs"`
	// PriorShares is the number of shares held in the prior period, or nil (JSON key "PriShs").
	PriorShares *int64 `json:"PriShs"`
	// CurrentBookValue is the current-period book value, or nil (JSON key "CurBookVal").
	CurrentBookValue *int64 `json:"CurBookVal"`
	// PriorBookValue is the prior-period book value, or nil (JSON key "PriBookVal").
	PriorBookValue *int64 `json:"PriBookVal"`
	// CurrentSharesNotDisclosed is a note when current shares are not disclosed, or nil (JSON key "CurShsNotDisc").
	CurrentSharesNotDisclosed *string `json:"CurShsNotDisc"`
	// PriorSharesNotDisclosed is a note when prior shares are not disclosed, or nil (JSON key "PriShsNotDisc").
	PriorSharesNotDisclosed *string `json:"PriShsNotDisc"`
	// CurrentBookValueNotDisclosed is a note when the current book value is not disclosed, or nil (JSON key "CurBookValNotDisc").
	CurrentBookValueNotDisclosed *string `json:"CurBookValNotDisc"`
	// PriorBookValueNotDisclosed is a note when the prior book value is not disclosed, or nil (JSON key "PriBookValNotDisc").
	PriorBookValueNotDisclosed *string `json:"PriBookValNotDisc"`
	// HoldingRationale is the stated rationale for holding the shares (JSON key "HoldRat").
	HoldingRationale string `json:"HoldRat"`
	// IssuerHolds indicates whether the issuer holds the filer's shares (JSON key "IsrHolds").
	IssuerHolds string `json:"IsrHolds"`
	// IssuerHoldsCode is the coded form of IssuerHolds ("1", "0", or "2") (JSON key "IsrHoldsCode").
	IssuerHoldsCode string `json:"IsrHoldsCode"`
}

// CrossShareholdingEntry represents one of the Report, Largest, or SecondLargest
// cross-shareholding blocks in a filing.
type CrossShareholdingEntry struct {
	// HolderName is the holder's name (JSON key "HldrName").
	HolderName string `json:"HldrName"`
	// HolderCode is the holder's security code, or nil (JSON key "HldrCode").
	HolderCode *string `json:"HldrCode"`
	// HolderEdinetCode is the holder's EDINET code, or nil (JSON key "HldrEdinetCode").
	HolderEdinetCode *string `json:"HldrEdinetCode"`
	// ListedIssues is the number of listed cross-held issues, or nil (JSON key "ListedIss").
	ListedIssues *int64 `json:"ListedIss"`
	// ListedBookValue is the total book value of listed cross-holdings, or nil (JSON key "ListedBookVal").
	ListedBookValue *int64 `json:"ListedBookVal"`
	// ListedIncreaseIssues is the number of listed issues that increased, or nil (JSON key "ListedIncIss").
	ListedIncreaseIssues *int64 `json:"ListedIncIss"`
	// ListedIncreaseAcquisitionCost is the acquisition cost of listed increases, or nil (JSON key "ListedIncAcqCost").
	ListedIncreaseAcquisitionCost *int64 `json:"ListedIncAcqCost"`
	// ListedDecreaseIssues is the number of listed issues that decreased, or nil (JSON key "ListedDecIss").
	ListedDecreaseIssues *int64 `json:"ListedDecIss"`
	// ListedDecreaseSaleAmount is the sale amount of listed decreases, or nil (JSON key "ListedDecSaleAmt").
	ListedDecreaseSaleAmount *int64 `json:"ListedDecSaleAmt"`
	// ListedIncreaseReason is the reason for listed increases, or nil (JSON key "ListedIncRsn").
	ListedIncreaseReason *string `json:"ListedIncRsn"`
	// NonListedIssues is the number of non-listed cross-held issues, or nil (JSON key "NonListedIss").
	NonListedIssues *int64 `json:"NonListedIss"`
	// NonListedBookValue is the total book value of non-listed cross-holdings, or nil (JSON key "NonListedBookVal").
	NonListedBookValue *int64 `json:"NonListedBookVal"`
	// NonListedIncreaseIssues is the number of non-listed issues that increased, or nil (JSON key "NonListedIncIss").
	NonListedIncreaseIssues *int64 `json:"NonListedIncIss"`
	// NonListedIncreaseAcquisitionCost is the acquisition cost of non-listed increases, or nil (JSON key "NonListedIncAcqCost").
	NonListedIncreaseAcquisitionCost *int64 `json:"NonListedIncAcqCost"`
	// NonListedDecreaseIssues is the number of non-listed issues that decreased, or nil (JSON key "NonListedDecIss").
	NonListedDecreaseIssues *int64 `json:"NonListedDecIss"`
	// NonListedDecreaseSaleAmount is the sale amount of non-listed decreases, or nil (JSON key "NonListedDecSaleAmt").
	NonListedDecreaseSaleAmount *int64 `json:"NonListedDecSaleAmt"`
	// NonListedIncreaseReason is the reason for non-listed increases, or nil (JSON key "NonListedIncRsn").
	NonListedIncreaseReason *string `json:"NonListedIncRsn"`
	// Specified is the list of specified investment cross-holdings (JSON key "Spec").
	Specified []CrossShareholdingIssue `json:"Spec"`
	// Deemed is the list of deemed-held investment cross-holdings (JSON key "Deem").
	Deemed []CrossShareholdingIssue `json:"Deem"`
	// SpecifiedFootnote is a footnote for the specified list, or nil (JSON key "SpecFn").
	SpecifiedFootnote *string `json:"SpecFn"`
	// DeemedFootnote is a footnote for the deemed list, or nil (JSON key "DeemFn").
	DeemedFootnote *string `json:"DeemFn"`
}

// CrossShareholdings represents a cross-shareholdings filing extracted from EDINET.
type CrossShareholdings struct {
	// DocumentID is the EDINET document ID (JSON key "DocId").
	DocumentID string `json:"DocId"`
	// Code is the security code (JSON key "Code").
	Code string `json:"Code"`
	// EdinetCode is the filer's EDINET code (JSON key "EdinetCode").
	EdinetCode string `json:"EdinetCode"`
	// FilerName is the filer's name in Japanese (JSON key "FilerName").
	FilerName string `json:"FilerName"`
	// FilerNameEnglish is the filer's name in English (JSON key "FilerNameEn").
	FilerNameEnglish string `json:"FilerNameEn"`
	// DocumentTypeCode is the EDINET document type code (JSON key "DocTypeCode").
	DocumentTypeCode string `json:"DocTypeCode"`
	// SubmissionDate is the filing submission date (JSON key "SubDate").
	SubmissionDate string `json:"SubDate"`
	// SubmissionTime is the filing submission time (JSON key "SubTime").
	SubmissionTime string `json:"SubTime"`
	// PeriodStart is the reporting period start date (JSON key "PerSt").
	PeriodStart string `json:"PerSt"`
	// PeriodEnd is the reporting period end date (JSON key "PerEn").
	PeriodEnd string `json:"PerEn"`
	// Report is the reporting company's cross-shareholding block (JSON key "Report").
	Report CrossShareholdingEntry `json:"Report"`
	// Largest is the largest holding company's block, or nil (JSON key "Largest").
	Largest *CrossShareholdingEntry `json:"Largest"`
	// SecondLargest is the second-largest holding company's block, or nil (JSON key "SecondLargest").
	SecondLargest *CrossShareholdingEntry `json:"SecondLargest"`
}

type crossShareholdingsParameters struct {
	EdinetRequest
	PaginationKey *string
}

func (p crossShareholdingsParameters) values() (url.Values, error) {
	return edinetValues(p.EdinetCode, p.Code, p.Date, p.PaginationKey)
}

type crossShareholdingsResponse struct {
	Data          []CrossShareholdings `json:"data"`
	PaginationKey *string              `json:"pagination_key"`
}

func (r crossShareholdingsResponse) Items() []CrossShareholdings { return r.Data }
func (r crossShareholdingsResponse) NextPageKey() *string        { return r.PaginationKey }

// CrossShareholdings retrieves cross-shareholdings filings from the /edinet/cross-shareholdings endpoint.
// It automatically handles pagination to fetch all matching records.
// Some fields are extracted by an LLM and may be nil; callers should handle nil pointers.
// See https://jpx-jquants.com/en/spec/edinet-cross-shareholdings for API details.
func (c *Client) CrossShareholdings(ctx context.Context, req EdinetRequest) ([]CrossShareholdings, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (crossShareholdingsResponse, error) {
		params := crossShareholdingsParameters{EdinetRequest: req, PaginationKey: paginationKey}
		return getJSON[crossShareholdingsResponse](ctx, c, "/edinet/cross-shareholdings", params)
	})
}

// LargeVolumeAcquisitionDisposal represents a single acquisition or disposal
// record within a large-volume holding report.
type LargeVolumeAcquisitionDisposal struct {
	// Date is the transaction date (JSON key "Date").
	Date string `json:"Date"`
	// SecurityType is the type of security transacted (JSON key "SecType").
	SecurityType string `json:"SecType"`
	// Shares is the number of shares transacted (JSON key "Shs").
	Shares int64 `json:"Shs"`
	// Ratio is the shareholding ratio after the transaction (JSON key "Ratio").
	Ratio float64 `json:"Ratio"`
	// Market describes the market where the transaction occurred (JSON key "Mkt").
	Market string `json:"Mkt"`
	// MarketCode is the coded market ("1": on-market, "2": off-market) (JSON key "MktCode").
	MarketCode string `json:"MktCode"`
	// TransactionType describes the transaction type (JSON key "TxnType").
	TransactionType string `json:"TxnType"`
	// TransactionTypeCode is the coded transaction type ("1": acquisition, "2": disposal) (JSON key "TxnTypeCode").
	TransactionTypeCode string `json:"TxnTypeCode"`
	// Counterparty is the counterparty name, or nil (JSON key "Cptty").
	Counterparty *string `json:"Cptty"`
	// Price is the transaction price (JSON key "Price").
	Price int64 `json:"Price"`
	// PriceRaw is the raw price text as extracted, or nil (JSON key "PriceRaw").
	PriceRaw *string `json:"PriceRaw"`
}

// LargeVolumeBorrowing represents an entry in the list of counterparties from
// which shares are borrowed within a large-volume holding report.
type LargeVolumeBorrowing struct {
	// Name is the counterparty name, or nil (JSON key "Name").
	Name *string `json:"Name"`
	// Industry is the counterparty industry, or nil (JSON key "Ind").
	Industry *string `json:"Ind"`
	// Representative is the counterparty representative, or nil (JSON key "Rep").
	Representative *string `json:"Rep"`
	// Address is the counterparty address, or nil (JSON key "Addr").
	Address *string `json:"Addr"`
	// DisclosedBorrowingPurpose is coded disclosure of the borrowing purpose
	// ("1": not disclosed, "2": disclosed) (JSON key "DiscBrwPurp").
	DisclosedBorrowingPurpose string `json:"DiscBrwPurp"`
	// Amount is the borrowed amount (JSON key "Amt").
	Amount int64 `json:"Amt"`
}

// LargeVolumeCredit represents an entry in the list of credit counterparties
// within a large-volume holding report.
type LargeVolumeCredit struct {
	// Name is the counterparty name, or nil (JSON key "Name").
	Name *string `json:"Name"`
	// Representative is the counterparty representative, or nil (JSON key "Rep").
	Representative *string `json:"Rep"`
	// Address is the counterparty address, or nil (JSON key "Addr").
	Address *string `json:"Addr"`
}

// LargeVolumeHolder represents a single holder within a large-volume holding report.
type LargeVolumeHolder struct {
	// HolderName is the holder's name (JSON key "HldrName").
	HolderName string `json:"HldrName"`
	// HolderNameEnglish is the holder's name in English (JSON key "HldrNameEn").
	HolderNameEnglish string `json:"HldrNameEn"`
	// HolderEdinetCode is the holder's EDINET code (JSON key "HldrEdinetCode").
	HolderEdinetCode string `json:"HldrEdinetCode"`
	// HolderCode is the holder's security code, or nil (JSON key "HldrCode").
	HolderCode *string `json:"HldrCode"`
	// LargeHolderTypeCode is the coded holder type ("1": individual, "2": corporation, "0": unknown) (JSON key "LargeHldrTypeCode").
	LargeHolderTypeCode string `json:"LargeHldrTypeCode"`
	// LargeHolderTypeRaw is the raw holder type text (JSON key "LargeHldrTypeRaw").
	LargeHolderTypeRaw string `json:"LargeHldrTypeRaw"`
	// HoldingPurpose is the stated purpose of holding (JSON key "HldgPurp").
	HoldingPurpose string `json:"HldgPurp"`
	// ImportantProposal is a note on any important proposal, or nil (JSON key "ImpProp").
	ImportantProposal *string `json:"ImpProp"`
	// CollateralAgreement is a note on any collateral agreement, or nil (JSON key "ColAgr").
	CollateralAgreement *string `json:"ColAgr"`
	// SharesHeld is the number of shares held (JSON key "ShsHeld").
	SharesHeld int64 `json:"ShsHeld"`
	// SharesRatio is the shareholding ratio (JSON key "ShsRatio").
	SharesRatio float64 `json:"ShsRatio"`
	// SharesRatioLast is the shareholding ratio in the previous report, or nil (JSON key "ShsRatioLast").
	SharesRatioLast *float64 `json:"ShsRatioLast"`
	// OwnFund is the amount funded from the holder's own capital, or nil (JSON key "OwnFund").
	OwnFund *int64 `json:"OwnFund"`
	// TotalBorrowed is the total borrowed amount, or nil (JSON key "TotalBrw").
	TotalBorrowed *int64 `json:"TotalBrw"`
	// TotalOther is the total from other sources, or nil (JSON key "TotalOther").
	TotalOther *int64 `json:"TotalOther"`
	// OtherBreakdown is a breakdown note for other sources, or nil (JSON key "OtherBrk").
	OtherBreakdown *string `json:"OtherBrk"`
	// TotalFund is the total funding amount, or nil (JSON key "TotalFund").
	TotalFund *int64 `json:"TotalFund"`
	// AcquisitionsDisposals is the list of acquisition and disposal records (JSON key "AcqDisp").
	AcquisitionsDisposals []LargeVolumeAcquisitionDisposal `json:"AcqDisp"`
	// Borrowings is the list of borrowing counterparties (JSON key "BrwList").
	Borrowings []LargeVolumeBorrowing `json:"BrwList"`
	// Credits is the list of credit counterparties (JSON key "CredList").
	Credits []LargeVolumeCredit `json:"CredList"`
}

// LargeVolumeShareholders represents a large-volume holding report extracted from EDINET.
type LargeVolumeShareholders struct {
	// DocumentID is the EDINET document ID (JSON key "DocId").
	DocumentID string `json:"DocId"`
	// Code is the security code (JSON key "Code").
	Code string `json:"Code"`
	// EdinetCode is the filer's EDINET code (JSON key "EdinetCode").
	EdinetCode string `json:"EdinetCode"`
	// IssuerName is the issuer's name (JSON key "IsrName").
	IssuerName string `json:"IsrName"`
	// DocumentTypeCode is the EDINET document type code (JSON key "DocTypeCode").
	DocumentTypeCode string `json:"DocTypeCode"`
	// SubmissionDate is the filing submission date (JSON key "SubDate").
	SubmissionDate string `json:"SubDate"`
	// SubmissionTime is the filing submission time (JSON key "SubTime").
	SubmissionTime string `json:"SubTime"`
	// LargeHoldingTypeCode is the coded report type
	// ("1": report, "2": change, "3": short-term, "4": special, "5": special change, "0": unknown) (JSON key "LargeHldgTypeCode").
	LargeHoldingTypeCode string `json:"LargeHldgTypeCode"`
	// DocumentTitle is the document title (JSON key "DocTitle").
	DocumentTitle string `json:"DocTitle"`
	// ChangeReason is the reason for the change, or nil (JSON key "ChgRsn").
	ChangeReason *string `json:"ChgRsn"`
	// TotalSharesHeld is the total number of shares held (JSON key "TotalShsHeld").
	TotalSharesHeld int64 `json:"TotalShsHeld"`
	// TotalSharesRatio is the total shareholding ratio (JSON key "TotalShsRatio").
	TotalSharesRatio float64 `json:"TotalShsRatio"`
	// TotalSharesRatioLast is the total shareholding ratio in the previous report, or nil (JSON key "TotalShsRatioLast").
	TotalSharesRatioLast *float64 `json:"TotalShsRatioLast"`
	// TotalOutstandingStocks is the total number of outstanding shares (JSON key "TotalOutStks").
	TotalOutstandingStocks int64 `json:"TotalOutStks"`
	// Holders is the list of holders covered by the report (JSON key "Hldrs").
	Holders []LargeVolumeHolder `json:"Hldrs"`
}

type largeVolumeShareholdersParameters struct {
	EdinetRequest
	PaginationKey *string
}

func (p largeVolumeShareholdersParameters) values() (url.Values, error) {
	return edinetValues(p.EdinetCode, p.Code, p.Date, p.PaginationKey)
}

type largeVolumeShareholdersResponse struct {
	Data          []LargeVolumeShareholders `json:"data"`
	PaginationKey *string                   `json:"pagination_key"`
}

func (r largeVolumeShareholdersResponse) Items() []LargeVolumeShareholders { return r.Data }
func (r largeVolumeShareholdersResponse) NextPageKey() *string             { return r.PaginationKey }

// LargeVolumeShareholders retrieves large-volume holding reports from the /edinet/large-volume-shareholders endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/edinet-large-volume-shareholders for API details.
func (c *Client) LargeVolumeShareholders(ctx context.Context, req EdinetRequest) ([]LargeVolumeShareholders, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (largeVolumeShareholdersResponse, error) {
		params := largeVolumeShareholdersParameters{EdinetRequest: req, PaginationKey: paginationKey}
		return getJSON[largeVolumeShareholdersResponse](ctx, c, "/edinet/large-volume-shareholders", params)
	})
}
