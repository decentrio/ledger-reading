package handlers

import (
	"time"

	"github.com/decentrio/ledger-reading/database/models"
)

func (h *DBHandler) GetTokenPriceInUsd(name string) (float64, error) {
	var token models.Tokens
	if err := h.db.Table("tokens").Where("token_name = ?", name).Scan(&token).Error; err != nil {
		return 0, err
	}

	return token.PriceInUsd, nil
}

func (h *DBHandler) GetToken(name string) (models.Tokens, error) {
	var token models.Tokens
	if err := h.db.Table("tokens").Where("token_name = ?", name).Scan(&token).Error; err != nil {
		return models.Tokens{}, err
	}

	return token, nil
}

func (h *DBHandler) SetTokens(data *models.Tokens) (string, error) {
	if err := h.db.Table("tokens").Where("token_name = ?", data.TokenName).Save(data).Error; err != nil {
		if err := h.db.Create(data).Error; err != nil {
			return "", err
		}
	}

	return data.TokenName, nil
}

func (h *DBHandler) SetTokenPrice(name string, price float64) (string, error) {
	var token models.Tokens
	h.db.Table("tokens").Where("token_name = ?", name).Scan(&token)
	token.PriceInUsd = price

	if err := h.db.Table("tokens").Where("token_name = ?", name).Save(token).Error; err != nil {
		return "", err
	}

	return name, nil
}

func (h *DBHandler) GetTicker(id string) (models.Tickers, error) {
	var ticker models.Tickers
	if err := h.db.Table("tickers").Where("ticker_id = ?", id).Scan(&ticker).Error; err != nil {
		return models.Tickers{}, err
	}

	return ticker, nil
}

func (h *DBHandler) SetTickers(data *models.Tickers) (string, error) {
	var baseVolume uint64
	h.db.Table("historical_trades").
		Where("ticker_id = ?", data.TickerId).
		Where("trade_timestamp >= ?", time.Now().Unix()-86400).
		Select("sum(base_volume) as total").Scan(&baseVolume)

	var targetVolume uint64
	h.db.Table("historical_trades").
		Where("ticker_id = ?", data.TickerId).
		Where("trade_timestamp >= ?", time.Now().Unix()-86400).
		Select("sum(target_volume) as total").Scan(&targetVolume)

	var high float64
	if err := h.db.Table("historical_trades").
		Where("ticker_id = ?", data.TickerId).
		Where("trade_timestamp >= ?", time.Now().Unix()-86400).
		Select("max(price)").Scan(&high).Error; err != nil {
		data.High = data.LastPrice
	} else {
		if data.LastPrice > high {
			data.High = data.LastPrice
		} else {
			data.High = high
		}
	}

	var low float64
	if err := h.db.Table("historical_trades").
		Where("ticker_id = ?", data.TickerId).
		Where("trade_timestamp >= ?", time.Now().Unix()-86400).
		Select("min(price)").Scan(&low).Error; err != nil {
		data.Low = data.LastPrice
	} else {
		if data.LastPrice < low {
			data.Low = data.LastPrice
		} else {
			data.Low = low
		}
	}

	data.BaseVolume = baseVolume
	data.TargetVolume = targetVolume

	// update share liquidity
	if data.ShareLiquidity == 0 {
		var ticker models.Tickers
		err := h.db.Table("tickers").
			Where("ticker_id = ?", data.TickerId).
			Scan(&ticker).Error
		if err == nil {
			data.ShareLiquidity = ticker.ShareLiquidity
		}
	}

	if err := h.db.Table("tickers").Where("ticker_id = ?", data.TickerId).Save(data).Error; err != nil {
		if err := h.db.Create(data).Error; err != nil {
			return "", err
		}
	}

	return data.TickerId, nil
}

func (h *DBHandler) UpdateTickerLiquidity(tickerId string, shareLiquidity, baseLiquidity, targetLiquidity, liquidityInUsd uint64) (string, error) {
	var ticker models.Tickers
	h.db.Table("tickers").
		Where("ticker_id = ?", tickerId).
		Scan(&ticker)

	ticker.ShareLiquidity = shareLiquidity
	ticker.BaseLiquidity = baseLiquidity
	ticker.TargetLiquidity = targetLiquidity
	ticker.LiquidityInUsd = liquidityInUsd

	if err := h.db.Table("tickers").
		Where("ticker_id = ?", tickerId).
		Save(ticker).Error; err != nil {
		return "ERROR: update old contract data entry", err
	}

	return tickerId, nil
}

func (h *DBHandler) CreateActivities(data *models.Activities) error {
	if err := h.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (h *DBHandler) CreateHistoricalTrades(data *models.HistoricalTrades) (uint32, error) {
	if err := h.db.Create(data).Error; err != nil {
		return 0, err
	}

	return data.TradeId, nil
}
