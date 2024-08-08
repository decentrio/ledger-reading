package importer

import (
	"time"
)

func (i *Importer) FetchTokenList() {
	defer i.wg.Done()
	for {
		select {
		// Terminate process
		case <-i.BaseService.Terminate():
			return
		default:
			var newTokenList []Token
			err := getJSON(TokenListUrl, &newTokenList)
			if err != nil {
				i.Logger.Error(err)
			}

			for _, newToken := range newTokenList {
				i.tokenListUpdateCb(newToken)
			}
		}

		time.Sleep(time.Minute * 5)
	}

}

func (i *Importer) FetchTickerList() {
	defer i.wg.Done()
	for {
		select {
		// Terminate process
		case <-i.BaseService.Terminate():
			return
		default:
			var newTickerList []Ticker
			err := getJSON(TickerListUrl, &newTickerList)
			if err != nil {
				i.Logger.Error(err)
			}

			for _, newTicker := range newTickerList {
				i.tickerListUpdateCb(newTicker)
			}
		}

		time.Sleep(time.Minute * 5)
	}
}
