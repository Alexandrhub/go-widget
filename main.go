package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/getlantern/systray"
	"github.com/robfig/cron/v3"
)

type state struct {
	Price string
	Cron  *cron.Cron
}

func main() {
	s := &state{}
	systray.Run(s.onReady, s.onExit)
}

func (s *state) onReady() {
	s.updatePrice()
	s.Cron = cron.New()
	s.Cron.AddFunc("@every 1m", func() { s.updatePrice() })
	s.Cron.Start()
}

func (s *state) onExit() {
	s.Cron.Stop()
}

func (s *state) updatePrice() {
	url := "https://coinmarketcap.com/currencies/bitcoin/"
	resp, err := http.Get(url)
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

	systray.SetTitle("BTC $: " + price)
}

// https://coinmarketcap.com/currencies/bitcoin/
//<span class="sc-f70bb44c-0 jxpCgO base-text">$42,851.89</span>
