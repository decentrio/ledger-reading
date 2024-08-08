package uploader

import (
	"fmt"
	"sync"

	db "github.com/decentrio/ledger-reading/database/handlers"
	"github.com/decentrio/ledger-reading/database/models"
	"github.com/decentrio/ledger-reading/importer"
	"github.com/decentrio/ledger-reading/lib/service"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
)

type Uploader struct {
	service.BaseService

	// networkPassPhrase
	networkPassPhrase string
	// ledgerQueue channel for trigger new ledger
	ledgerReadChan <-chan xdr.LedgerCloseMeta

	TickerList            map[importer.TokenPair]importer.Ticker
	TickerListWithPoolKey map[string]importer.Ticker
	TokenList             map[string]importer.Token

	db *db.DBHandler

	wg sync.WaitGroup
}

// ExporterOption sets an optional parameter on the State.
type UploaderOption func(*Uploader)

func NewUploader(
	// cfg *config.ExporterConfig,
	baseLogger *log.Entry,
	readChan <-chan xdr.LedgerCloseMeta,
	networkPassPhrase string,
	options ...UploaderOption,
) *Uploader {
	u := &Uploader{
		ledgerReadChan:        readChan,
		TickerListWithPoolKey: make(map[string]importer.Ticker),
		TickerList:            make(map[importer.TokenPair]importer.Ticker),
		TokenList:             make(map[string]importer.Token),
		networkPassPhrase:     networkPassPhrase,
	}

	u.BaseService = *service.NewBaseService("uploader", u)
	for _, opt := range options {
		opt(u)
	}

	logger := baseLogger.WithField("module", "uploader")
	logger.SetLevel(log.DebugLevel)
	u.BaseService.SetLogger(logger)

	u.db = db.NewDBHandler()

	return u
}

func (u *Uploader) OnStart() error {
	u.Logger.Info("start services")
	u.wg.Add(1)

	go u.ledgerProcessing()
	return nil
}

func (u *Uploader) OnStop() error {
	u.Logger.Info("stop services")
	u.wg.Wait()

	return nil
}

func (u *Uploader) UploadNewToken(t importer.Token) {
	isUpload := false
	token, found := u.TokenList[t.SorobanContract]
	if found {
		if token.Symbol != t.Symbol ||
			token.SorobanContract != t.SorobanContract ||
			token.Decimals != t.Decimals {
			isUpload = true
			u.Logger.Warn(fmt.Sprintf("ticker unmatch %v - %v", t, token))

			// Update Token list for now
			u.TokenList[t.SorobanContract] = t
		}
	} else {
		// check on db
		tk, err := u.db.GetToken(t.Token)
		if err != nil {
			isUpload = true
		}

		if tk.Symbol != t.Symbol ||
			tk.SorobanContract != t.SorobanContract ||
			tk.Decimal != int(t.Decimals) {
			isUpload = true
			u.Logger.Warn(fmt.Sprintf("ticker unmatch %v - %v", t, tk))
		}

		u.Logger.Info(fmt.Sprintf("new token uploaded %v", t))
		u.TokenList[t.SorobanContract] = t
	}

	if isUpload {
		token := models.Tokens{
			Symbol:          t.Symbol,
			TokenName:       t.Token,
			SorobanContract: t.SorobanContract,
			Decimal:         int(t.Decimals),
			PriceInUsd:      0,
		}
		u.db.SetTokens(&token)
	}
}

func (u *Uploader) UploadNewTicker(t importer.Ticker) {
	isUpload := false
	pair := importer.GetTokenPair(t.BaseCurrency, t.TargetCurrency)
	ticker, found := u.TickerList[pair]
	if found {
		if ticker.BaseCurrency != t.BaseCurrency ||
			ticker.TargetCurrency != t.TargetCurrency ||
			ticker.PoolContract != t.PoolContract {
			isUpload = true
			u.Logger.Warn(fmt.Sprintf("ticker unmatch %v - %v", t, ticker))

			// Update Ticker list for now
			u.TickerList[pair] = t
			u.TickerListWithPoolKey[t.PoolContract] = t
		}
	} else {
		// check on db
		tk, err := u.db.GetTicker(t.TickerId)
		if err != nil {
			isUpload = true
		}

		if tk.BaseCurrency != t.BaseCurrency ||
			tk.TargetCurrency != t.TargetCurrency ||
			tk.PoolId != t.PoolContract {
			isUpload = true
			u.Logger.Warn(fmt.Sprintf("ticker unmatch %v - %v", t, tk))
		}

		u.Logger.Info(fmt.Sprintf("new ticker uploaded %v", t))
		u.TickerList[pair] = t
		u.TickerListWithPoolKey[t.PoolContract] = t
	}

	if isUpload {
		ticker := models.Tickers{
			TickerId:        t.TickerId,
			BaseCurrency:    t.BaseCurrency,
			TargetCurrency:  t.TargetCurrency,
			PoolId:          t.PoolContract,
			LastPrice:       0,
			BaseVolume:      0,
			TargetVolume:    0,
			High:            0,
			Low:             0,
			ShareLiquidity:  0,
			BaseLiquidity:   0,
			TargetLiquidity: 0,
			LiquidityInUsd:  0,
		}
		u.db.SetTickers(&ticker)
	}
}
