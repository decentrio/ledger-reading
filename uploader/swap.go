package uploader

import (
	"fmt"
	"math"
	"time"

	"github.com/decentrio/ledger-reading/database/models"
)

func (tx *PhoenixTransactionExtractor) ExtractSwapTx(u *Uploader) error {
	tis, err := tx.GetTradeInfo(u)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("err get historical trade %s", err.Error()))
		return err
	}

	for i, ti := range tis {
		tradeId := (tx.Ops[0].ID() + time.Now().Unix() + int64(i)) % math.MaxInt32
		// create historical trade
		historicalTrade := models.HistoricalTrades{
			TradeId:        uint32(tradeId),
			Price:          ti.Price,
			TickerId:       ti.Ticker.TickerId,
			BaseVolume:     ti.BaseVolume,
			TargetVolume:   ti.TargetVolume,
			TradeTimestamp: tx.Time,
			TradeType:      ti.TradeType,
			TxHash:         tx.GetTransactionHash(),
			Maker:          tx.Tx.Envelope.SourceAccount().ToAccountId().Address(),
		}
		u.db.CreateHistoricalTrades(&historicalTrade)

		// create activities
		activity := models.Activities{
			Address:        tx.Tx.Envelope.SourceAccount().ToAccountId().Address(),
			ActionType:     ti.TradeType,
			BaseCurrency:   ti.Ticker.BaseCurrency,
			BaseVolume:     ti.BaseVolume,
			TargetCurrency: ti.Ticker.TargetCurrency,
			TargetVolume:   ti.TargetVolume,
			Timestamp:      tx.Time,
		}
		u.db.CreateActivities(&activity)

		// update price
		targetCurrencyPriceInUsd := u.getTokenPriceInUsd(ti.Ticker.TargetCurrency)
		if targetCurrencyPriceInUsd != float64(0) {
			price := targetCurrencyPriceInUsd * ti.Price
			u.db.SetTokenPrice(ti.Ticker.BaseCurrency, price)
		}

		// get liquidity data
		share, base, target, totalInUsd := u.GetPoolLiquidity(tx, ti.Ticker)
		// update data
		modelTicker := models.Tickers{
			TickerId:        ti.Ticker.TickerId,
			BaseCurrency:    ti.Ticker.BaseCurrency,
			TargetCurrency:  ti.Ticker.TargetCurrency,
			PoolId:          ti.Ticker.PoolContract,
			LastPrice:       ti.Price,
			ShareLiquidity:  share,
			BaseLiquidity:   base,
			TargetLiquidity: target,
			LiquidityInUsd:  totalInUsd,
		}
		u.db.SetTickers(&modelTicker)
	}

	return nil
}

func (tx *PhoenixTransactionExtractor) GetTradeInfo(u *Uploader) ([]TradeInformation, error) {
	//extract contract event to retrieve data
	wasmEvents, err := tx.GetContractEvents()
	if err != nil {
		return nil, fmt.Errorf("error get contract events")
	}

	var tis []TradeInformation
	for _, e := range wasmEvents {
		var tradeType string
		var askedToken string
		var offerToken string
		var offerAmount uint64
		var returnAmount uint64
		var baseVolume uint64
		var targetVolume uint64
		var price float64

		ticker, found := u.TickerListWithPoolKey[e.ContractId]
		if !found {
			u.Logger.Warnf("unknown contract pool %s", e.ContractId)
			continue
		}

		for _, eventBody := range e.EventBody {
			v0 := eventBody.MustV0()

			topics := v0.Topics
			data := v0.Data
			_ = data

			if len(topics) != 2 {
				continue
			}

			action := string(topics[1].MustStr())

			switch action {
			case "sell_token":
				offerTokenSorobanContract, err := data.MustAddress().String()
				if err != nil {
					continue
				}

				token, found := u.TokenList[offerTokenSorobanContract]
				if !found {
					u.Logger.Warnf("unknown soroban contract %s", offerTokenSorobanContract)
					continue
				}
				offerToken = token.Token
			case "buy_token":
				askedTokenSorobanContract, err := data.MustAddress().String()
				if err != nil {
					continue
				}

				token, found := u.TokenList[askedTokenSorobanContract]
				if !found {
					u.Logger.Warnf("unknown soroban contract %s", askedTokenSorobanContract)
					continue
				}
				askedToken = token.Token
			case "offer_amount":
				offerAmount = uint64(data.MustI128().Lo)
			case "return_amount":
				returnAmount = uint64(data.MustI128().Lo)
			default:
			}
		}

		if ticker.BaseCurrency == offerToken && ticker.TargetCurrency == askedToken {
			// sell
			tradeType = "sell"
			price = float64(returnAmount) / float64(offerAmount)
			baseVolume = offerAmount
			targetVolume = returnAmount
		} else if ticker.BaseCurrency == askedToken && ticker.TargetCurrency == offerToken {
			// buy
			tradeType = "buy"
			price = float64(offerAmount) / float64(returnAmount)
			baseVolume = returnAmount
			targetVolume = offerAmount
		} else {
			continue
		}

		ti := TradeInformation{
			TradeType:    tradeType,
			Ticker:       ticker,
			BaseVolume:   baseVolume,
			TargetVolume: targetVolume,
			Price:        price,
		}

		tis = append(tis, ti)
	}

	return tis, nil
}
