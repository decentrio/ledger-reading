package uploader

import (
	"fmt"
)

func (tx *IndexerTransactionExtractor) GetTradeInfo(u *Uploader) ([]TradeInformation, error) {
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
