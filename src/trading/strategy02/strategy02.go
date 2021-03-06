package strategy02

import (
	"finantial/ema"
	"finantial/rsi"
	"fmt"
	"markets/bitstamp"
	"markets/generic"
	"time"
)

var UNDEF = int32(-1)
var TRUE = int32(1)
var FALSE = int32(0)

func Start(buycoin string, sellcoin string, invest float64, fee float64, period int, training_iters int, fast int, slow int, rsi_win_len int, rsi_buy_level float64, rsi_sell_level float64) {
	var ema_fast ema.TFinantial_EMA
	var ema_slow ema.TFinantial_EMA
	var rsi rsi.TFinantial_RSI
	var ema_vol ema.TFinantial_EMA

	ema_fast.Reset(fast)
	ema_slow.Reset(slow)
	ema_vol.Reset(10)
	rsi.Reset(rsi_win_len, rsi_buy_level, rsi_sell_level)

	var market generic.TMarket

	market.Reset(buycoin, sellcoin, invest, fee)

	var fast_gt_slow = int32(UNDEF)

	iter := 0
	pair := buycoin + sellcoin // btceur

	fmt.Println(pair)

	for {
		time.Sleep(time.Duration(period) * time.Second)

		price, err := bitstamp.DoGet(pair)
		if err != nil {
			fmt.Println("Error en el doget de bitstamp")
			continue
		}

		ema_fast.NewPrice(price)
		ema_slow.NewPrice(price)
		rsi.NewPrice(price)
		fmt.Println("price: ", fmt.Sprintf("%.2f", price),
			"\tema_fast: ", fmt.Sprintf("%.2f", ema_fast.Ema()),
			"\tema_slow: ", fmt.Sprintf("%.2f", ema_slow.Ema()),
			"\trsi: ", fmt.Sprintf("%.2f", rsi.RSI()),
			"\ttime: ", time.Now())

		if iter < training_iters {
			iter++
			continue
		}

		// End of training, start trading

		// Initialize fast_gt_slow only once after training
		if fast_gt_slow == UNDEF {
			if ema_fast.Ema() > ema_slow.Ema() {
				fast_gt_slow = TRUE
			} else {
				fast_gt_slow = FALSE
			}
			fmt.Println("Training ready. Starting trade now...")
			continue
		}

		// fast_gt_slow already defined

		/*
		   if (market.InsideMarket()) {
		       if (price < market.LastBuyPrice()) {
		           market.DoSell(price)
		           fmt.Println("********************************** Activated: CONTROL1")
		           fmt.Println("********************************** VENDE a: ", market.LastSellPrice())
		           fmt.Println("********************************** FIAT: ", market.Fiat())
		       } else {
		           fmt.Println("===> He comprado y esta subiendo, GOOD SIGNAL")
		       }
		   } */

		if fast_gt_slow == FALSE {
			if ema_fast.Ema() < ema_slow.Ema() {
				fmt.Println("ema_fast < ema_slow... Se mantiene la tendencia de bajada")
				// tendency is maintained (falling price)
				continue
			} else {
				fmt.Println("ema_fast > ema_slow... Cambio de tendencia, comprobemos si se puede comprar")
				if market.InsideMarket() == false {
					fmt.Println("InsideMarket = false")
					if rsi.Buy() {
						market.DoBuy(price)
						fmt.Println("********************************** Buy at: ", market.LastBuyPrice())
						fmt.Println("********************************** CRYPTO: ", market.Crypto())
					} else {
						fmt.Println("Improper RSI to buy: ", rsi.RSI(), "rsi.BuyLevel = ", rsi.BuyLevel())
						continue
					}
				} else {
					fmt.Println("===> Tocaba comprar pero ya estoy dentro")
				}
				fast_gt_slow = TRUE
			}
		} else {
			if ema_fast.Ema() > ema_slow.Ema() {
				fmt.Println("ema_fast > ema_slow... Se mantiene la tendencia de subida")
				// tendency is maintained (climbing price)
				continue
			} else {
				fmt.Println("ema_fast < ema_slow... Cambio de tendencia, comprobemos si se puede vender")
				if market.InsideMarket() == true {
					fmt.Println("InsideMarket = true")
					if rsi.Sell() {
						market.DoSell(price)
						fmt.Println("********************************** Sell at: ", market.LastSellPrice())
						fmt.Println("********************************** FIAT: ", market.Fiat())
					} else {
						fmt.Println("Improper RSI to sell: ", rsi.RSI(), "rsi.SellLevel = ", rsi.SellLevel())
						continue
					}

				} else {
					fmt.Println("===> Tocaba vender pero estoy fuera")

				}
				fast_gt_slow = FALSE
			}
		}
	}
}
