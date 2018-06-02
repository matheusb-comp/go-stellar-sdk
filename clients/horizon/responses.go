package horizon

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Problem struct {
	Type     string                     `json:"type"`
	Title    string                     `json:"title"`
	Status   int                        `json:"status"`
	Detail   string                     `json:"detail,omitempty"`
	Instance string                     `json:"instance,omitempty"`
	Extras   map[string]json.RawMessage `json:"extras,omitempty"`
}

type Root struct {
	Links struct {
		Account             Link `json:"account"`
		AccountTransactions Link `json:"account_transactions"`
		Friendbot           Link `json:"friendbot"`
		Metrics             Link `json:"metrics"`
		OrderBook           Link `json:"order_book"`
		Self                Link `json:"self"`
		Transaction         Link `json:"transaction"`
		Transactions        Link `json:"transactions"`
	} `json:"_links"`

	HorizonVersion       string `json:"horizon_version"`
	StellarCoreVersion   string `json:"core_version"`
	HorizonSequence      int32  `json:"history_latest_ledger"`
	HistoryElderSequence int32  `json:"history_elder_ledger"`
	CoreSequence         int32  `json:"core_latest_ledger"`
	CoreElderSequence    int32  `json:"core_elder_ledger"`
	NetworkPassphrase    string `json:"network_passphrase"`
	ProtocolVersion      int32  `json:"protocol_version"`
}

type Account struct {
	Links struct {
		Self         Link `json:"self"`
		Transactions Link `json:"transactions"`
		Operations   Link `json:"operations"`
		Payments     Link `json:"payments"`
		Effects      Link `json:"effects"`
		Offers       Link `json:"offers"`
	} `json:"_links"`

	HistoryAccount
	Sequence             string            `json:"sequence"`
	SubentryCount        int32             `json:"subentry_count"`
	InflationDestination string            `json:"inflation_destination,omitempty"`
	HomeDomain           string            `json:"home_domain,omitempty"`
	Thresholds           AccountThresholds `json:"thresholds"`
	Flags                AccountFlags      `json:"flags"`
	Balances             []Balance         `json:"balances"`
	Signers              []Signer          `json:"signers"`
	Data                 map[string]string `json:"data"`
}

func (a Account) GetNativeBalance() string {
	for _, balance := range a.Balances {
		if balance.Asset.Type == "native" {
			return balance.Balance
		}
	}

	return "0"
}

func (a Account) GetCreditBalance(code, issuer string) string {
	for _, balance := range a.Balances {
		if balance.Asset.Code == code && balance.Asset.Issuer == issuer {
			return balance.Balance
		}
	}

	return "0"
}

// MustGetData returns decoded value for a given key. If the key does
// not exist, empty slice will be returned. If there is an error
// decoding a value, it will panic.
func (this *Account) MustGetData(key string) []byte {
	bytes, err := this.GetData(key)
	if err != nil {
		panic(err)
	}
	return bytes
}

// GetData returns decoded value for a given key. If the key does
// not exist, empty slice will be returned.
func (this *Account) GetData(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(this.Data[key])
}

type AccountFlags struct {
	AuthRequired  bool `json:"auth_required"`
	AuthRevocable bool `json:"auth_revocable"`
}

type AccountThresholds struct {
	LowThreshold  byte `json:"low_threshold"`
	MedThreshold  byte `json:"med_threshold"`
	HighThreshold byte `json:"high_threshold"`
}

type Asset struct {
	Type   string `json:"asset_type"`
	Code   string `json:"asset_code,omitempty"`
	Issuer string `json:"asset_issuer,omitempty"`
}

type Balance struct {
	Balance string `json:"balance"`
	Limit   string `json:"limit,omitempty"`
	Asset
}

type HistoryAccount struct {
	ID        string `json:"id"`
	PT        string `json:"paging_token"`
	AccountID string `json:"account_id"`
}

type Effect struct {
	Links struct {
		Operation Link `json:"operation"`
		Succeeds  Link `json:"succeeds"`
		Precedes  Link `json:"precedes"`
	} `json:"_links"`
	ID              string `json:"id"`
	PT              string `json:"paging_token"`
	Account         string `json:"account"`
	Amount          string `json:"amount"`
	Type            string `json:"type"`
	TypeI           int32  `json:"type_i"`
	StartingBalance string `json:"starting_balance"`
	Balance
	Signer
}

type Operation struct {
	Links struct {
		Transaction Link `json:"transaction"`
		Effects     Link `json:"effects"`
		Succeeds    Link `json:"succeeds"`
		Precedes    Link `json:"precedes"`
	} `json:"_links"`
	ID              string `json:"id"`
	PT              string `json:"paging_token"`
	Account         string `json:"account"`
	SourceAccount   string `json:"source_account"`
	Type            string `json:"type"`
	TypeI           int32  `json:"type_i"`
	StartingBalance string `json:"starting_balance"`
	CreatedAt       string `json:"created_at"`
	Funder          string `json:"funder"`
	TransactionHash string `json:"transaction_hash"`
}

type Ledger struct {
	Links struct {
		Self         Link `json:"self"`
		Transactions Link `json:"transactions"`
		Operations   Link `json:"operations"`
		Payments     Link `json:"payments"`
		Effects      Link `json:"effects"`
	} `json:"_links"`
	ID               string    `json:"id"`
	PT               string    `json:"paging_token"`
	Hash             string    `json:"hash"`
	PrevHash         string    `json:"prev_hash,omitempty"`
	Sequence         int32     `json:"sequence"`
	TransactionCount int32     `json:"transaction_count"`
	OperationCount   int32     `json:"operation_count"`
	ClosedAt         time.Time `json:"closed_at"`
	TotalCoins       string    `json:"total_coins"`
	FeePool          string    `json:"fee_pool"`
	BaseFee          int32     `json:"base_fee_in_stroops"`
	BaseReserve      int32     `json:"base_reserve_in_stroops"`
	MaxTxSetSize     int32     `json:"max_tx_set_size"`
	ProtocolVersion  int32     `json:"protocol_version"`
}

type Link struct {
	Href      string `json:"href"`
	Templated bool   `json:"templated,omitempty"`
}

type Offer struct {
	Links struct {
		Self       Link `json:"self"`
		OfferMaker Link `json:"offer_maker"`
	} `json:"_links"`

	ID      int64  `json:"id"`
	PT      string `json:"paging_token"`
	Seller  string `json:"seller"`
	Selling Asset  `json:"selling"`
	Buying  Asset  `json:"buying"`
	Amount  string `json:"amount"`
	PriceR  Price  `json:"price_r"`
	Price   string `json:"price"`
}

// TradeAggregationsPage returns a list of aggregated trade records, aggregated by resolution
type TradeAggregationsPage struct {
	Links struct {
		Self Link `json:"self"`
		Next Link `json:"next"`
		Prev Link `json:"prev"`
	} `json:"_links"`
	Embedded struct {
		Records []TradeAggregation `json:"records"`
	} `json:"_embedded"`
}

// TradeAggregation represents trade data aggregation over a period of time
type TradeAggregation struct {
	Timestamp     int64  `json:"timestamp"`
	TradeCount    int64  `json:"trade_count"`
	BaseVolume    string `json:"base_volume"`
	CounterVolume string `json:"counter_volume"`
	Average       string `json:"avg"`
	High          string `json:"high"`
	HighR         Price  `json:"high_r"`
	Low           string `json:"low"`
	LowR          Price  `json:"low_r"`
	Open          string `json:"open"`
	OpenR         Price  `json:"open_r"`
	Close         string `json:"close"`
	CloseR        Price  `json:"close_r"`
}

// TradesPage returns a list of trade records
type TradesPage struct {
	Links struct {
		Self Link `json:"self"`
		Next Link `json:"next"`
		Prev Link `json:"prev"`
	} `json:"_links"`
	Embedded struct {
		Records []Trade `json:"records"`
	} `json:"_embedded"`
}

// Trade represents a horizon digested trade
type Trade struct {
	Links struct {
		Self      Link `json:"self"`
		Base      Link `json:"base"`
		Counter   Link `json:"counter"`
		Operation Link `json:"operation"`
	} `json:"_links"`

	ID                 string    `json:"id"`
	PT                 string    `json:"paging_token"`
	LedgerCloseTime    time.Time `json:"ledger_close_time"`
	OfferID            string    `json:"offer_id"`
	BaseAccount        string    `json:"base_account"`
	BaseAmount         string    `json:"base_amount"`
	BaseAssetType      string    `json:"base_asset_type"`
	BaseAssetCode      string    `json:"base_asset_code,omitempty"`
	BaseAssetIssuer    string    `json:"base_asset_issuer,omitempty"`
	CounterAccount     string    `json:"counter_account"`
	CounterAmount      string    `json:"counter_amount"`
	CounterAssetType   string    `json:"counter_asset_type"`
	CounterAssetCode   string    `json:"counter_asset_code,omitempty"`
	CounterAssetIssuer string    `json:"counter_asset_issuer,omitempty"`
	BaseIsSeller       bool      `json:"base_is_seller"`
	Price              *Price    `json:"price"`
}

type OrderBookSummary struct {
	Bids    []PriceLevel `json:"bids"`
	Asks    []PriceLevel `json:"asks"`
	Selling Asset        `json:"base"`
	Buying  Asset        `json:"counter"`
}

type TransactionSuccess struct {
	Links struct {
		Transaction Link `json:"transaction"`
	} `json:"_links"`
	Hash   string `json:"hash"`
	Ledger int32  `json:"ledger"`
	Env    string `json:"envelope_xdr"`
	Result string `json:"result_xdr"`
	Meta   string `json:"result_meta_xdr"`
}

// TransactionResultCodes represent a summary of result codes returned from
// a single xdr TransactionResult
type TransactionResultCodes struct {
	TransactionCode string   `json:"transaction"`
	OperationCodes  []string `json:"operations,omitempty"`
}

type Signer struct {
	PublicKey string `json:"public_key"`
	Weight    int32  `json:"weight"`
	Key       string `json:"key"`
	Type      string `json:"type"`
}

type OffersPage struct {
	Links struct {
		Self Link `json:"self"`
		Next Link `json:"next"`
		Prev Link `json:"prev"`
	} `json:"_links"`
	Embedded struct {
		Records []Offer `json:"records"`
	} `json:"_embedded"`
}

type Payment struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	PagingToken string `json:"paging_token"`

	Links struct {
		Transaction struct {
			Href string `json:"href"`
		} `json:"transaction"`
	} `json:"_links"`

	TransactionHash string `json:"transaction_hash"`
	SourceAccount   string `json:"source_account"`
	CreatedAt       string `json:"created_at"`

	// create_account and account_merge field
	Account string `json:"account"`

	// create_account fields
	Funder          string `json:"funder"`
	StartingBalance string `json:"starting_balance"`

	// account_merge fields
	Into string `json:into"`

	// payment/path_payment fields
	From        string `json:"from"`
	To          string `json:"to"`
	AssetType   string `json:"asset_type"`
	AssetCode   string `json:"asset_code"`
	AssetIssuer string `json:"asset_issuer"`
	Amount      string `json:"amount"`

	// transaction fields
	Memo struct {
		Type  string `json:"memo_type"`
		Value string `json:"memo"`
	}
}

type Price struct {
	N int32 `json:"n"`
	D int32 `json:"d"`
}

type PriceLevel struct {
	PriceR Price  `json:"price_r"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

type Transaction struct {
	ID              string    `json:"id"`
	PagingToken     string    `json:"paging_token"`
	Hash            string    `json:"hash"`
	Ledger          int32     `json:"ledger"`
	LedgerCloseTime time.Time `json:"created_at"`
	Account         string    `json:"source_account"`
	AccountSequence string    `json:"source_account_sequence"`
	FeePaid         int32     `json:"fee_paid"`
	OperationCount  int32     `json:"operation_count"`
	EnvelopeXdr     string    `json:"envelope_xdr"`
	ResultXdr       string    `json:"result_xdr"`
	ResultMetaXdr   string    `json:"result_meta_xdr"`
	FeeMetaXdr      string    `json:"fee_meta_xdr"`
	MemoType        string    `json:"memo_type"`
	Memo            string    `json:"memo,omitempty"`
	Signatures      []string  `json:"signatures"`
	ValidAfter      string    `json:"valid_after,omitempty"`
	ValidBefore     string    `json:"valid_before,omitempty"`
}
