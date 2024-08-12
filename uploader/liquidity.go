package uploader

import (
	"fmt"

)

func (tx *IndexerTransactionExtractor) ExtractProvideLiquidityTx(u *Uploader) error {
	// pis, err := tx.GetProvideLiquidityInfo(u)
	// if err != nil {
	// 	u.Logger.Error(fmt.Sprintf("err provide liquidity information %s", err.Error()))
	// 	return err
	// }

	// for _, pi := range pis {
	// 	// add provide_liquidity activity
	// 	activity := models.Activities{
	// 		Address:        pi.Sender,
	// 		ActionType:     string(ProvideLiquidity),
	// 		BaseCurrency:   pi.Ticker.BaseCurrency,
	// 		BaseVolume:     pi.BaseVolume,
	// 		TargetCurrency: pi.Ticker.TargetCurrency,
	// 		TargetVolume:   pi.TargetVolume,
	// 		Timestamp:      tx.Time,
	// 	}

	// }

	return nil
}

func (tx *IndexerTransactionExtractor) GetProvideLiquidityInfo(u *Uploader) ([]ProvideLiquidityInformation, error) {
	// extract contract event to retrieve data
	wasmEvents, err := tx.GetContractEvents()
	if err != nil {
		return nil, fmt.Errorf("error get contract events")
	}

	var pis []ProvideLiquidityInformation
	for _, e := range wasmEvents {
		var sender string
		var baseToken string
		var targetToken string
		var baseVolume uint64
		var targetVolume uint64

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
			topic := string(topics[1].MustStr())

			switch topic {
			case "sender":
				sender, err = data.MustAddress().String()
				if err != nil {
					continue
				}
			case "token_a":
				baseTokenSorobanContract, err := data.MustAddress().String()
				if err != nil {
					continue
				}

				token, found := u.TokenList[baseTokenSorobanContract]
				if !found {
					u.Logger.Warnf("unknown soroban contract %s", baseTokenSorobanContract)
					continue
				}
				baseToken = token.Token
			case "token_a-amount":
				baseVolume = uint64(data.MustI128().Lo)
			case "token_b":
				targetTokenSorobanContract, err := data.MustAddress().String()
				if err != nil {
					continue
				}

				token, found := u.TokenList[targetTokenSorobanContract]
				if !found {
					u.Logger.Warnf("unknown soroban contract %s", targetTokenSorobanContract)
					continue
				}
				targetToken = token.Token
			case "token_b-amount":
				targetVolume = uint64(data.MustI128().Lo)
			default:
			}
		}

		if ticker.BaseCurrency != baseToken || ticker.TargetCurrency != targetToken {
			u.Logger.Warnf("unknown pool %s with pair %s - %s", ticker.PoolContract, baseToken, targetToken)
			continue
		}

		share, base, target, totalInUsd := u.GetPoolLiquidity(tx, ticker)

		pi := ProvideLiquidityInformation{
			Sender:              sender,
			Ticker:              ticker,
			BaseVolume:          baseVolume,
			TargetVolume:        targetVolume,
			ShareLiquidity:      share,
			BaseLiquidity:       base,
			TargetLiquidity:     target,
			TotalLiquidityInUsd: totalInUsd,
		}

		pis = append(pis, pi)
	}

	return pis, nil
}
