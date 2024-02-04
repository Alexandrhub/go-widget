package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/getlantern/systray"
	"github.com/robfig/cron/v3"
)

const (
	btc = "BTC"
	eth = "ETH"
	thr = "USDT"
	bnb = "BNB"
)

type state struct {
	Price            string
	Cron             *cron.Cron
	SelectedCurrency string
	CurrencyNames    map[string]string
	MenuItems        map[string]*systray.MenuItem
}

func main() {
	s := &state{
		SelectedCurrency: btc,
		CurrencyNames: map[string]string{
			btc: "Bitcoin",
			eth: "Ethereum",
			thr: "Tether",
			bnb: "Binance Coin",
		},
		MenuItems: map[string]*systray.MenuItem{},
	}
	systray.Run(s.onReady, s.onExit)
}

func (s *state) onReady() {
	s.updatePrice()

	s.Cron = cron.New()
	s.Cron.AddFunc("@every 30s", func() { s.updatePrice() })
	s.Cron.Start()

	for currency := range s.CurrencyNames {
		s.MenuItems[currency] = systray.AddMenuItem(currency, "")
	}

	for {
		select {
		case <-s.MenuItems[btc].ClickedCh:
			s.SelectedCurrency = btc
		case <-s.MenuItems[eth].ClickedCh:
			s.SelectedCurrency = eth
		case <-s.MenuItems[thr].ClickedCh:
			s.SelectedCurrency = thr
		case <-s.MenuItems[bnb].ClickedCh:
			s.SelectedCurrency = bnb
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *state) onExit() {
	s.Cron.Stop()
}

func (s *state) updatePrice() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := "https://coinmarketcap.com/currencies/" + s.CurrencyNames[s.SelectedCurrency]
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	price := doc.Find(".sc-f70bb44c-0.jxpCgO.base-text").Text()

	currency := fmt.Sprintf("%s: %s", s.SelectedCurrency, price)

	systray.SetTitle(currency)
}

// https://coinmarketcap.com/currencies/bitcoin/
//<span class="sc-f70bb44c-0 jxpCgO base-text">$42,851.89</span>
