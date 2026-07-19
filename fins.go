package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// parseFloatPtr parses a numeric string into a *float64.
// Empty strings and the placeholder values "-" and "*" (used by the J-Quants
// API to signal "no data") are converted to a nil pointer.
func parseFloatPtr(s string) (*float64, error) {
	switch s {
	case "", "-", "*":
		return nil, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// floatAccumulator accumulates the first error encountered while converting a
// series of numeric strings to *float64, allowing cleaner unmarshaling code.
type floatAccumulator struct {
	err error
}

func (a *floatAccumulator) f(s string) *float64 {
	if a.err != nil {
		return nil
	}
	result, err := parseFloatPtr(s)
	a.err = err
	return result
}

// unmarshalFloatFromAny converts a value decoded from JSON (which may be a
// float64, or one of the placeholder strings ""/"-"/"*") into a *float64.
// Placeholder strings and nil become a nil pointer.
func unmarshalFloatFromAny(value any) (*float64, error) {
	switch v := value.(type) {
	case float64:
		return &v, nil
	case string:
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unmarshalFloatFromAny: unknown type %T", v)
	}
}

// FinancialSummary represents a summary of financial statement data disclosed
// by a listed company. All numeric values are reported by the API as strings;
// numeric fields are parsed into *float64 with a nil pointer indicating no data.
type FinancialSummary struct {
	// DisclosureDate is the disclosure date (JSON key "DiscDate").
	DisclosureDate string
	// DisclosureTime is the disclosure time (JSON key "DiscTime").
	DisclosureTime string
	// Code is the security code (JSON key "Code").
	Code string
	// DisclosureNumber is the disclosure number (JSON key "DiscNo").
	DisclosureNumber string
	// DocumentType is the type of disclosure document (JSON key "DocType").
	DocumentType string
	// CurrentPeriodType is the type of the current fiscal period, e.g. "1Q" (JSON key "CurPerType").
	CurrentPeriodType string
	// CurrentPeriodStart is the start date of the current period (JSON key "CurPerSt").
	CurrentPeriodStart string
	// CurrentPeriodEnd is the end date of the current period (JSON key "CurPerEn").
	CurrentPeriodEnd string
	// CurrentFiscalYearStart is the start date of the current fiscal year (JSON key "CurFYSt").
	CurrentFiscalYearStart string
	// CurrentFiscalYearEnd is the end date of the current fiscal year (JSON key "CurFYEn").
	CurrentFiscalYearEnd string
	// NextFiscalYearStart is the start date of the next fiscal year (JSON key "NxtFYSt").
	NextFiscalYearStart string
	// NextFiscalYearEnd is the end date of the next fiscal year (JSON key "NxtFYEn").
	NextFiscalYearEnd string

	// NetSales is net sales (JSON key "Sales").
	NetSales *float64
	// OperatingProfit is operating profit (JSON key "OP").
	OperatingProfit *float64
	// OrdinaryProfit is ordinary profit (JSON key "OdP").
	OrdinaryProfit *float64
	// NetProfit is profit attributable to owners of parent (JSON key "NP").
	NetProfit *float64
	// EarningsPerShare is earnings per share (JSON key "EPS").
	EarningsPerShare *float64
	// DilutedEarningsPerShare is diluted earnings per share (JSON key "DEPS").
	DilutedEarningsPerShare *float64
	// TotalAssets is total assets (JSON key "TA").
	TotalAssets *float64
	// Equity is equity (JSON key "Eq").
	Equity *float64
	// EquityToAssetRatio is the equity-to-asset ratio (JSON key "EqAR").
	EquityToAssetRatio *float64
	// BookValuePerShare is book value per share (JSON key "BPS").
	BookValuePerShare *float64
	// CashFlowsFromOperating is cash flows from operating activities (JSON key "CFO").
	CashFlowsFromOperating *float64
	// CashFlowsFromInvesting is cash flows from investing activities (JSON key "CFI").
	CashFlowsFromInvesting *float64
	// CashFlowsFromFinancing is cash flows from financing activities (JSON key "CFF").
	CashFlowsFromFinancing *float64
	// CashAndEquivalents is cash and cash equivalents at end of period (JSON key "CashEq").
	CashAndEquivalents *float64

	// Dividend1Q is the first-quarter dividend per share (JSON key "Div1Q").
	Dividend1Q *float64
	// Dividend2Q is the second-quarter dividend per share (JSON key "Div2Q").
	Dividend2Q *float64
	// Dividend3Q is the third-quarter dividend per share (JSON key "Div3Q").
	Dividend3Q *float64
	// DividendFY is the fiscal-year-end dividend per share (JSON key "DivFY").
	DividendFY *float64
	// DividendAnnual is the annual dividend per share (JSON key "DivAnn").
	DividendAnnual *float64
	// DividendUnit is the per-unit dividend (JSON key "DivUnit").
	DividendUnit *float64
	// DividendTotalAnnual is the total annual dividend paid (JSON key "DivTotalAnn").
	DividendTotalAnnual *float64
	// PayoutRatioAnnual is the annual payout ratio (JSON key "PayoutRatioAnn").
	PayoutRatioAnnual *float64

	// ForecastDividend1Q is the forecast first-quarter dividend per share (JSON key "FDiv1Q").
	ForecastDividend1Q *float64
	// ForecastDividend2Q is the forecast second-quarter dividend per share (JSON key "FDiv2Q").
	ForecastDividend2Q *float64
	// ForecastDividend3Q is the forecast third-quarter dividend per share (JSON key "FDiv3Q").
	ForecastDividend3Q *float64
	// ForecastDividendFY is the forecast fiscal-year-end dividend per share (JSON key "FDivFY").
	ForecastDividendFY *float64
	// ForecastDividendAnnual is the forecast annual dividend per share (JSON key "FDivAnn").
	ForecastDividendAnnual *float64
	// ForecastDividendUnit is the forecast per-unit dividend (JSON key "FDivUnit").
	ForecastDividendUnit *float64
	// ForecastDividendTotalAnnual is the forecast total annual dividend (JSON key "FDivTotalAnn").
	ForecastDividendTotalAnnual *float64
	// ForecastPayoutRatioAnnual is the forecast annual payout ratio (JSON key "FPayoutRatioAnn").
	ForecastPayoutRatioAnnual *float64

	// NextForecastDividend1Q is the next-year forecast first-quarter dividend per share (JSON key "NxFDiv1Q").
	NextForecastDividend1Q *float64
	// NextForecastDividend2Q is the next-year forecast second-quarter dividend per share (JSON key "NxFDiv2Q").
	NextForecastDividend2Q *float64
	// NextForecastDividend3Q is the next-year forecast third-quarter dividend per share (JSON key "NxFDiv3Q").
	NextForecastDividend3Q *float64
	// NextForecastDividendFY is the next-year forecast fiscal-year-end dividend per share (JSON key "NxFDivFY").
	NextForecastDividendFY *float64
	// NextForecastDividendAnnual is the next-year forecast annual dividend per share (JSON key "NxFDivAnn").
	NextForecastDividendAnnual *float64
	// NextForecastDividendUnit is the next-year forecast per-unit dividend (JSON key "NxFDivUnit").
	NextForecastDividendUnit *float64
	// NextForecastPayoutRatioAnnual is the next-year forecast annual payout ratio (JSON key "NxFPayoutRatioAnn").
	NextForecastPayoutRatioAnnual *float64

	// ForecastNetSales2Q is the forecast half-year net sales (JSON key "FSales2Q").
	ForecastNetSales2Q *float64
	// ForecastOperatingProfit2Q is the forecast half-year operating profit (JSON key "FOP2Q").
	ForecastOperatingProfit2Q *float64
	// ForecastOrdinaryProfit2Q is the forecast half-year ordinary profit (JSON key "FOdP2Q").
	ForecastOrdinaryProfit2Q *float64
	// ForecastNetProfit2Q is the forecast half-year net profit (JSON key "FNP2Q").
	ForecastNetProfit2Q *float64
	// ForecastEarningsPerShare2Q is the forecast half-year EPS (JSON key "FEPS2Q").
	ForecastEarningsPerShare2Q *float64
	// NextForecastNetSales2Q is the next-year forecast half-year net sales (JSON key "NxFSales2Q").
	NextForecastNetSales2Q *float64
	// NextForecastOperatingProfit2Q is the next-year forecast half-year operating profit (JSON key "NxFOP2Q").
	NextForecastOperatingProfit2Q *float64
	// NextForecastOrdinaryProfit2Q is the next-year forecast half-year ordinary profit (JSON key "NxFOdP2Q").
	NextForecastOrdinaryProfit2Q *float64
	// NextForecastNetProfit2Q is the next-year forecast half-year net profit (JSON key "NxFNp2Q").
	NextForecastNetProfit2Q *float64
	// NextForecastEarningsPerShare2Q is the next-year forecast half-year EPS (JSON key "NxFEPS2Q").
	NextForecastEarningsPerShare2Q *float64

	// ForecastNetSales is the forecast full-year net sales (JSON key "FSales").
	ForecastNetSales *float64
	// ForecastOperatingProfit is the forecast full-year operating profit (JSON key "FOP").
	ForecastOperatingProfit *float64
	// ForecastOrdinaryProfit is the forecast full-year ordinary profit (JSON key "FOdP").
	ForecastOrdinaryProfit *float64
	// ForecastNetProfit is the forecast full-year net profit (JSON key "FNP").
	ForecastNetProfit *float64
	// ForecastEarningsPerShare is the forecast full-year EPS (JSON key "FEPS").
	ForecastEarningsPerShare *float64
	// NextForecastNetSales is the next-year forecast full-year net sales (JSON key "NxFSales").
	NextForecastNetSales *float64
	// NextForecastOperatingProfit is the next-year forecast full-year operating profit (JSON key "NxFOP").
	NextForecastOperatingProfit *float64
	// NextForecastOrdinaryProfit is the next-year forecast full-year ordinary profit (JSON key "NxFOdP").
	NextForecastOrdinaryProfit *float64
	// NextForecastNetProfit is the next-year forecast full-year net profit (JSON key "NxFNp").
	NextForecastNetProfit *float64
	// NextForecastEarningsPerShare is the next-year forecast full-year EPS (JSON key "NxFEPS").
	NextForecastEarningsPerShare *float64

	// MaterialChangesInSubsidiaries indicates material changes in the scope of consolidation (JSON key "MatChgSub").
	MaterialChangesInSubsidiaries string
	// SignificantChangesInConsolidation indicates significant changes in the consolidation scope (JSON key "SigChgInC").
	SignificantChangesInConsolidation string
	// ChangesByAccountingStandardsRevisions indicates changes due to revisions of accounting standards (JSON key "ChgByASRev").
	ChangesByAccountingStandardsRevisions string
	// ChangesOtherThanAccountingStandardsRevisions indicates changes other than by revisions of accounting standards (JSON key "ChgNoASRev").
	ChangesOtherThanAccountingStandardsRevisions string
	// ChangesInAccountingEstimates indicates changes in accounting estimates (JSON key "ChgAcEst").
	ChangesInAccountingEstimates string
	// RetrospectiveRestatement indicates a retrospective restatement (JSON key "RetroRst").
	RetrospectiveRestatement string

	// SharesOutstandingFY is the number of issued shares at fiscal year end (JSON key "ShOutFY").
	SharesOutstandingFY *float64
	// TreasurySharesFY is the number of treasury shares at fiscal year end (JSON key "TrShFY").
	TreasurySharesFY *float64
	// AverageShares is the average number of shares during the period (JSON key "AvgSh").
	AverageShares *float64

	// NonConsolidatedNetSales is non-consolidated net sales (JSON key "NCSales").
	NonConsolidatedNetSales *float64
	// NonConsolidatedOperatingProfit is non-consolidated operating profit (JSON key "NCOP").
	NonConsolidatedOperatingProfit *float64
	// NonConsolidatedOrdinaryProfit is non-consolidated ordinary profit (JSON key "NCOdP").
	NonConsolidatedOrdinaryProfit *float64
	// NonConsolidatedNetProfit is non-consolidated net profit (JSON key "NCNP").
	NonConsolidatedNetProfit *float64
	// NonConsolidatedEarningsPerShare is non-consolidated EPS (JSON key "NCEPS").
	NonConsolidatedEarningsPerShare *float64
	// NonConsolidatedTotalAssets is non-consolidated total assets (JSON key "NCTA").
	NonConsolidatedTotalAssets *float64
	// NonConsolidatedEquity is non-consolidated equity (JSON key "NCEq").
	NonConsolidatedEquity *float64
	// NonConsolidatedEquityToAssetRatio is the non-consolidated equity-to-asset ratio (JSON key "NCEqAR").
	NonConsolidatedEquityToAssetRatio *float64
	// NonConsolidatedBookValuePerShare is non-consolidated book value per share (JSON key "NCBPS").
	NonConsolidatedBookValuePerShare *float64

	// ForecastNonConsolidatedNetSales2Q is the forecast half-year non-consolidated net sales (JSON key "FNCSales2Q").
	ForecastNonConsolidatedNetSales2Q *float64
	// ForecastNonConsolidatedOperatingProfit2Q is the forecast half-year non-consolidated operating profit (JSON key "FNCOP2Q").
	ForecastNonConsolidatedOperatingProfit2Q *float64
	// ForecastNonConsolidatedOrdinaryProfit2Q is the forecast half-year non-consolidated ordinary profit (JSON key "FNCOdP2Q").
	ForecastNonConsolidatedOrdinaryProfit2Q *float64
	// ForecastNonConsolidatedNetProfit2Q is the forecast half-year non-consolidated net profit (JSON key "FNCNP2Q").
	ForecastNonConsolidatedNetProfit2Q *float64
	// ForecastNonConsolidatedEarningsPerShare2Q is the forecast half-year non-consolidated EPS (JSON key "FNCEPS2Q").
	ForecastNonConsolidatedEarningsPerShare2Q *float64
	// NextForecastNonConsolidatedNetSales2Q is the next-year forecast half-year non-consolidated net sales (JSON key "NxFNCSales2Q").
	NextForecastNonConsolidatedNetSales2Q *float64
	// NextForecastNonConsolidatedOperatingProfit2Q is the next-year forecast half-year non-consolidated operating profit (JSON key "NxFNCOP2Q").
	NextForecastNonConsolidatedOperatingProfit2Q *float64
	// NextForecastNonConsolidatedOrdinaryProfit2Q is the next-year forecast half-year non-consolidated ordinary profit (JSON key "NxFNCOdP2Q").
	NextForecastNonConsolidatedOrdinaryProfit2Q *float64
	// NextForecastNonConsolidatedNetProfit2Q is the next-year forecast half-year non-consolidated net profit (JSON key "NxFNCNP2Q").
	NextForecastNonConsolidatedNetProfit2Q *float64
	// NextForecastNonConsolidatedEarningsPerShare2Q is the next-year forecast half-year non-consolidated EPS (JSON key "NxFNCEPS2Q").
	NextForecastNonConsolidatedEarningsPerShare2Q *float64

	// ForecastNonConsolidatedNetSales is the forecast full-year non-consolidated net sales (JSON key "FNCSales").
	ForecastNonConsolidatedNetSales *float64
	// ForecastNonConsolidatedOperatingProfit is the forecast full-year non-consolidated operating profit (JSON key "FNCOP").
	ForecastNonConsolidatedOperatingProfit *float64
	// ForecastNonConsolidatedOrdinaryProfit is the forecast full-year non-consolidated ordinary profit (JSON key "FNCOdP").
	ForecastNonConsolidatedOrdinaryProfit *float64
	// ForecastNonConsolidatedNetProfit is the forecast full-year non-consolidated net profit (JSON key "FNCNP").
	ForecastNonConsolidatedNetProfit *float64
	// ForecastNonConsolidatedEarningsPerShare is the forecast full-year non-consolidated EPS (JSON key "FNCEPS").
	ForecastNonConsolidatedEarningsPerShare *float64
	// NextForecastNonConsolidatedNetSales is the next-year forecast full-year non-consolidated net sales (JSON key "NxFNCSales").
	NextForecastNonConsolidatedNetSales *float64
	// NextForecastNonConsolidatedOperatingProfit is the next-year forecast full-year non-consolidated operating profit (JSON key "NxFNCOP").
	NextForecastNonConsolidatedOperatingProfit *float64
	// NextForecastNonConsolidatedOrdinaryProfit is the next-year forecast full-year non-consolidated ordinary profit (JSON key "NxFNCOdP").
	NextForecastNonConsolidatedOrdinaryProfit *float64
	// NextForecastNonConsolidatedNetProfit is the next-year forecast full-year non-consolidated net profit (JSON key "NxFNCNP").
	NextForecastNonConsolidatedNetProfit *float64
	// NextForecastNonConsolidatedEarningsPerShare is the next-year forecast full-year non-consolidated EPS (JSON key "NxFNCEPS").
	NextForecastNonConsolidatedEarningsPerShare *float64
}

func (fs *FinancialSummary) UnmarshalJSON(b []byte) error {
	var raw struct {
		DiscDate    string `json:"DiscDate"`
		DiscTime    string `json:"DiscTime"`
		Code        string `json:"Code"`
		DiscNo      string `json:"DiscNo"`
		DocType     string `json:"DocType"`
		CurPerType  string `json:"CurPerType"`
		CurPerSt    string `json:"CurPerSt"`
		CurPerEn    string `json:"CurPerEn"`
		CurFYSt     string `json:"CurFYSt"`
		CurFYEn     string `json:"CurFYEn"`
		NxtFYSt     string `json:"NxtFYSt"`
		NxtFYEn     string `json:"NxtFYEn"`
		Sales       string `json:"Sales"`
		OP          string `json:"OP"`
		OdP         string `json:"OdP"`
		NP          string `json:"NP"`
		EPS         string `json:"EPS"`
		DEPS        string `json:"DEPS"`
		TA          string `json:"TA"`
		Eq          string `json:"Eq"`
		EqAR        string `json:"EqAR"`
		BPS         string `json:"BPS"`
		CFO         string `json:"CFO"`
		CFI         string `json:"CFI"`
		CFF         string `json:"CFF"`
		CashEq      string `json:"CashEq"`
		Div1Q       string `json:"Div1Q"`
		Div2Q       string `json:"Div2Q"`
		Div3Q       string `json:"Div3Q"`
		DivFY       string `json:"DivFY"`
		DivAnn      string `json:"DivAnn"`
		DivUnit     string `json:"DivUnit"`
		DivTotalAnn string `json:"DivTotalAnn"`

		PayoutRatioAnn string `json:"PayoutRatioAnn"`

		FDiv1Q       string `json:"FDiv1Q"`
		FDiv2Q       string `json:"FDiv2Q"`
		FDiv3Q       string `json:"FDiv3Q"`
		FDivFY       string `json:"FDivFY"`
		FDivAnn      string `json:"FDivAnn"`
		FDivUnit     string `json:"FDivUnit"`
		FDivTotalAnn string `json:"FDivTotalAnn"`

		FPayoutRatioAnn string `json:"FPayoutRatioAnn"`

		NxFDiv1Q  string `json:"NxFDiv1Q"`
		NxFDiv2Q  string `json:"NxFDiv2Q"`
		NxFDiv3Q  string `json:"NxFDiv3Q"`
		NxFDivFY  string `json:"NxFDivFY"`
		NxFDivAnn string `json:"NxFDivAnn"`

		NxFDivUnit string `json:"NxFDivUnit"`

		NxFPayoutRatioAnn string `json:"NxFPayoutRatioAnn"`

		FSales2Q   string `json:"FSales2Q"`
		FOP2Q      string `json:"FOP2Q"`
		FOdP2Q     string `json:"FOdP2Q"`
		FNP2Q      string `json:"FNP2Q"`
		FEPS2Q     string `json:"FEPS2Q"`
		NxFSales2Q string `json:"NxFSales2Q"`
		NxFOP2Q    string `json:"NxFOP2Q"`
		NxFOdP2Q   string `json:"NxFOdP2Q"`
		NxFNp2Q    string `json:"NxFNp2Q"`
		NxFEPS2Q   string `json:"NxFEPS2Q"`

		FSales   string `json:"FSales"`
		FOP      string `json:"FOP"`
		FOdP     string `json:"FOdP"`
		FNP      string `json:"FNP"`
		FEPS     string `json:"FEPS"`
		NxFSales string `json:"NxFSales"`
		NxFOP    string `json:"NxFOP"`
		NxFOdP   string `json:"NxFOdP"`
		NxFNp    string `json:"NxFNp"`
		NxFEPS   string `json:"NxFEPS"`

		MatChgSub  string `json:"MatChgSub"`
		SigChgInC  string `json:"SigChgInC"`
		ChgByASRev string `json:"ChgByASRev"`
		ChgNoASRev string `json:"ChgNoASRev"`
		ChgAcEst   string `json:"ChgAcEst"`
		RetroRst   string `json:"RetroRst"`

		ShOutFY string `json:"ShOutFY"`
		TrShFY  string `json:"TrShFY"`
		AvgSh   string `json:"AvgSh"`

		NCSales string `json:"NCSales"`
		NCOP    string `json:"NCOP"`
		NCOdP   string `json:"NCOdP"`
		NCNP    string `json:"NCNP"`
		NCEPS   string `json:"NCEPS"`
		NCTA    string `json:"NCTA"`
		NCEq    string `json:"NCEq"`
		NCEqAR  string `json:"NCEqAR"`
		NCBPS   string `json:"NCBPS"`

		FNCSales2Q   string `json:"FNCSales2Q"`
		FNCOP2Q      string `json:"FNCOP2Q"`
		FNCOdP2Q     string `json:"FNCOdP2Q"`
		FNCNP2Q      string `json:"FNCNP2Q"`
		FNCEPS2Q     string `json:"FNCEPS2Q"`
		NxFNCSales2Q string `json:"NxFNCSales2Q"`
		NxFNCOP2Q    string `json:"NxFNCOP2Q"`
		NxFNCOdP2Q   string `json:"NxFNCOdP2Q"`
		NxFNCNP2Q    string `json:"NxFNCNP2Q"`
		NxFNCEPS2Q   string `json:"NxFNCEPS2Q"`

		FNCSales   string `json:"FNCSales"`
		FNCOP      string `json:"FNCOP"`
		FNCOdP     string `json:"FNCOdP"`
		FNCNP      string `json:"FNCNP"`
		FNCEPS     string `json:"FNCEPS"`
		NxFNCSales string `json:"NxFNCSales"`
		NxFNCOP    string `json:"NxFNCOP"`
		NxFNCOdP   string `json:"NxFNCOdP"`
		NxFNCNP    string `json:"NxFNCNP"`
		NxFNCEPS   string `json:"NxFNCEPS"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal financial summary: %w", err)
	}

	a := &floatAccumulator{}

	fs.DisclosureDate = raw.DiscDate
	fs.DisclosureTime = raw.DiscTime
	fs.Code = raw.Code
	fs.DisclosureNumber = raw.DiscNo
	fs.DocumentType = raw.DocType
	fs.CurrentPeriodType = raw.CurPerType
	fs.CurrentPeriodStart = raw.CurPerSt
	fs.CurrentPeriodEnd = raw.CurPerEn
	fs.CurrentFiscalYearStart = raw.CurFYSt
	fs.CurrentFiscalYearEnd = raw.CurFYEn
	fs.NextFiscalYearStart = raw.NxtFYSt
	fs.NextFiscalYearEnd = raw.NxtFYEn

	fs.NetSales = a.f(raw.Sales)
	fs.OperatingProfit = a.f(raw.OP)
	fs.OrdinaryProfit = a.f(raw.OdP)
	fs.NetProfit = a.f(raw.NP)
	fs.EarningsPerShare = a.f(raw.EPS)
	fs.DilutedEarningsPerShare = a.f(raw.DEPS)
	fs.TotalAssets = a.f(raw.TA)
	fs.Equity = a.f(raw.Eq)
	fs.EquityToAssetRatio = a.f(raw.EqAR)
	fs.BookValuePerShare = a.f(raw.BPS)
	fs.CashFlowsFromOperating = a.f(raw.CFO)
	fs.CashFlowsFromInvesting = a.f(raw.CFI)
	fs.CashFlowsFromFinancing = a.f(raw.CFF)
	fs.CashAndEquivalents = a.f(raw.CashEq)

	fs.Dividend1Q = a.f(raw.Div1Q)
	fs.Dividend2Q = a.f(raw.Div2Q)
	fs.Dividend3Q = a.f(raw.Div3Q)
	fs.DividendFY = a.f(raw.DivFY)
	fs.DividendAnnual = a.f(raw.DivAnn)
	fs.DividendUnit = a.f(raw.DivUnit)
	fs.DividendTotalAnnual = a.f(raw.DivTotalAnn)
	fs.PayoutRatioAnnual = a.f(raw.PayoutRatioAnn)

	fs.ForecastDividend1Q = a.f(raw.FDiv1Q)
	fs.ForecastDividend2Q = a.f(raw.FDiv2Q)
	fs.ForecastDividend3Q = a.f(raw.FDiv3Q)
	fs.ForecastDividendFY = a.f(raw.FDivFY)
	fs.ForecastDividendAnnual = a.f(raw.FDivAnn)
	fs.ForecastDividendUnit = a.f(raw.FDivUnit)
	fs.ForecastDividendTotalAnnual = a.f(raw.FDivTotalAnn)
	fs.ForecastPayoutRatioAnnual = a.f(raw.FPayoutRatioAnn)

	fs.NextForecastDividend1Q = a.f(raw.NxFDiv1Q)
	fs.NextForecastDividend2Q = a.f(raw.NxFDiv2Q)
	fs.NextForecastDividend3Q = a.f(raw.NxFDiv3Q)
	fs.NextForecastDividendFY = a.f(raw.NxFDivFY)
	fs.NextForecastDividendAnnual = a.f(raw.NxFDivAnn)
	fs.NextForecastDividendUnit = a.f(raw.NxFDivUnit)
	fs.NextForecastPayoutRatioAnnual = a.f(raw.NxFPayoutRatioAnn)

	fs.ForecastNetSales2Q = a.f(raw.FSales2Q)
	fs.ForecastOperatingProfit2Q = a.f(raw.FOP2Q)
	fs.ForecastOrdinaryProfit2Q = a.f(raw.FOdP2Q)
	fs.ForecastNetProfit2Q = a.f(raw.FNP2Q)
	fs.ForecastEarningsPerShare2Q = a.f(raw.FEPS2Q)
	fs.NextForecastNetSales2Q = a.f(raw.NxFSales2Q)
	fs.NextForecastOperatingProfit2Q = a.f(raw.NxFOP2Q)
	fs.NextForecastOrdinaryProfit2Q = a.f(raw.NxFOdP2Q)
	fs.NextForecastNetProfit2Q = a.f(raw.NxFNp2Q)
	fs.NextForecastEarningsPerShare2Q = a.f(raw.NxFEPS2Q)

	fs.ForecastNetSales = a.f(raw.FSales)
	fs.ForecastOperatingProfit = a.f(raw.FOP)
	fs.ForecastOrdinaryProfit = a.f(raw.FOdP)
	fs.ForecastNetProfit = a.f(raw.FNP)
	fs.ForecastEarningsPerShare = a.f(raw.FEPS)
	fs.NextForecastNetSales = a.f(raw.NxFSales)
	fs.NextForecastOperatingProfit = a.f(raw.NxFOP)
	fs.NextForecastOrdinaryProfit = a.f(raw.NxFOdP)
	fs.NextForecastNetProfit = a.f(raw.NxFNp)
	fs.NextForecastEarningsPerShare = a.f(raw.NxFEPS)

	fs.MaterialChangesInSubsidiaries = raw.MatChgSub
	fs.SignificantChangesInConsolidation = raw.SigChgInC
	fs.ChangesByAccountingStandardsRevisions = raw.ChgByASRev
	fs.ChangesOtherThanAccountingStandardsRevisions = raw.ChgNoASRev
	fs.ChangesInAccountingEstimates = raw.ChgAcEst
	fs.RetrospectiveRestatement = raw.RetroRst

	fs.SharesOutstandingFY = a.f(raw.ShOutFY)
	fs.TreasurySharesFY = a.f(raw.TrShFY)
	fs.AverageShares = a.f(raw.AvgSh)

	fs.NonConsolidatedNetSales = a.f(raw.NCSales)
	fs.NonConsolidatedOperatingProfit = a.f(raw.NCOP)
	fs.NonConsolidatedOrdinaryProfit = a.f(raw.NCOdP)
	fs.NonConsolidatedNetProfit = a.f(raw.NCNP)
	fs.NonConsolidatedEarningsPerShare = a.f(raw.NCEPS)
	fs.NonConsolidatedTotalAssets = a.f(raw.NCTA)
	fs.NonConsolidatedEquity = a.f(raw.NCEq)
	fs.NonConsolidatedEquityToAssetRatio = a.f(raw.NCEqAR)
	fs.NonConsolidatedBookValuePerShare = a.f(raw.NCBPS)

	fs.ForecastNonConsolidatedNetSales2Q = a.f(raw.FNCSales2Q)
	fs.ForecastNonConsolidatedOperatingProfit2Q = a.f(raw.FNCOP2Q)
	fs.ForecastNonConsolidatedOrdinaryProfit2Q = a.f(raw.FNCOdP2Q)
	fs.ForecastNonConsolidatedNetProfit2Q = a.f(raw.FNCNP2Q)
	fs.ForecastNonConsolidatedEarningsPerShare2Q = a.f(raw.FNCEPS2Q)
	fs.NextForecastNonConsolidatedNetSales2Q = a.f(raw.NxFNCSales2Q)
	fs.NextForecastNonConsolidatedOperatingProfit2Q = a.f(raw.NxFNCOP2Q)
	fs.NextForecastNonConsolidatedOrdinaryProfit2Q = a.f(raw.NxFNCOdP2Q)
	fs.NextForecastNonConsolidatedNetProfit2Q = a.f(raw.NxFNCNP2Q)
	fs.NextForecastNonConsolidatedEarningsPerShare2Q = a.f(raw.NxFNCEPS2Q)

	fs.ForecastNonConsolidatedNetSales = a.f(raw.FNCSales)
	fs.ForecastNonConsolidatedOperatingProfit = a.f(raw.FNCOP)
	fs.ForecastNonConsolidatedOrdinaryProfit = a.f(raw.FNCOdP)
	fs.ForecastNonConsolidatedNetProfit = a.f(raw.FNCNP)
	fs.ForecastNonConsolidatedEarningsPerShare = a.f(raw.FNCEPS)
	fs.NextForecastNonConsolidatedNetSales = a.f(raw.NxFNCSales)
	fs.NextForecastNonConsolidatedOperatingProfit = a.f(raw.NxFNCOP)
	fs.NextForecastNonConsolidatedOrdinaryProfit = a.f(raw.NxFNCOdP)
	fs.NextForecastNonConsolidatedNetProfit = a.f(raw.NxFNCNP)
	fs.NextForecastNonConsolidatedEarningsPerShare = a.f(raw.NxFNCEPS)

	return a.err
}

// FinancialSummaryRequest specifies filter parameters for the FinancialSummary API.
// At least one of Code or Date should be provided.
type FinancialSummaryRequest struct {
	// Code filters by security code.
	Code *string
	// Date filters by disclosure date in YYYYMMDD or YYYY-MM-DD format.
	Date *string
}

type financialSummaryParameters struct {
	FinancialSummaryRequest
	PaginationKey *string
}

func (p financialSummaryParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.Date != nil {
		v.Add("date", *p.Date)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type financialSummaryResponse struct {
	Data          []FinancialSummary `json:"data"`
	PaginationKey *string            `json:"pagination_key"`
}

func (r financialSummaryResponse) Items() []FinancialSummary { return r.Data }
func (r financialSummaryResponse) NextPageKey() *string      { return r.PaginationKey }

// FinancialSummary retrieves financial statement summary data from the /fins/summary endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/fin-summary for API details.
func (c *Client) FinancialSummary(ctx context.Context, req FinancialSummaryRequest) ([]FinancialSummary, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (financialSummaryResponse, error) {
		params := financialSummaryParameters{FinancialSummaryRequest: req, PaginationKey: paginationKey}
		return getJSON[financialSummaryResponse](ctx, c, "/fins/summary", params)
	})
}
