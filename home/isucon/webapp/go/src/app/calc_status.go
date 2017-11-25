package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

func calcStatus(currentTime int64, mItems map[int]mItem, addings []Adding, buyings []Buying) (*GameStatus, error) {
	var (
		// 1ミリ秒に生産できる椅子の単位をミリ椅子とする
		totalMilliIsu = big.NewInt(0)
		totalPower    = big.NewInt(0)

		itemPower    = map[int]*big.Int{}    // ItemID => Power
		itemPrice    = map[int]*big.Int{}    // ItemID => Price
		itemOnSale   = map[int]int64{}       // ItemID => OnSale
		itemBuilt    = map[int]int{}         // ItemID => BuiltCount
		itemBought   = map[int]int{}         // ItemID => CountBought
		itemBuilding = map[int][]Building{}  // ItemID => Buildings
		itemPower0   = map[int]Exponential{} // ItemID => currentTime における Power
		itemBuilt0   = map[int]int{}         // ItemID => currentTime における BuiltCount

		addingAt = map[int64]Adding{}   // Time => currentTime より先の Adding
		buyingAt = map[int64][]Buying{} // Time => currentTime より先の Buying
	)

	for itemID := range mItems {
		itemPower[itemID] = big.NewInt(0)
		itemBuilding[itemID] = []Building{}
	}

	for _, a := range addings {
		// adding は adding.time に isu を増加させる
		if a.Time <= currentTime {
			totalMilliIsu.Add(totalMilliIsu, new(big.Int).Mul(str2big(a.Isu), big.NewInt(1000)))
		} else {
			addingAt[a.Time] = a
		}
	}

	for _, b := range buyings {
		// buying は 即座に isu を消費し buying.time からアイテムの効果を発揮する
		itemBought[b.ItemID]++
		m := mItems[b.ItemID]
		totalMilliIsu.Sub(totalMilliIsu, new(big.Int).Mul(m.GetPrice(b.Ordinal), big.NewInt(1000)))

		if b.Time <= currentTime {
			itemBuilt[b.ItemID]++
			power := m.GetPower(itemBought[b.ItemID])
			totalMilliIsu.Add(totalMilliIsu, new(big.Int).Mul(power, big.NewInt(currentTime-b.Time)))
			totalPower.Add(totalPower, power)
			itemPower[b.ItemID].Add(itemPower[b.ItemID], power)
		} else {
			buyingAt[b.Time] = append(buyingAt[b.Time], b)
		}
	}

	for _, m := range mItems {
		itemPower0[m.ItemID] = big2exp(itemPower[m.ItemID])
		itemBuilt0[m.ItemID] = itemBuilt[m.ItemID]
		price := m.GetPrice(itemBought[m.ItemID] + 1)
		itemPrice[m.ItemID] = price
		if 0 <= totalMilliIsu.Cmp(new(big.Int).Mul(price, big.NewInt(1000))) {
			itemOnSale[m.ItemID] = 0 // 0 は 時刻 currentTime で購入可能であることを表す
		}
	}

	schedule := []Schedule{
		Schedule{
			Time:       currentTime,
			MilliIsu:   big2exp(totalMilliIsu),
			TotalPower: big2exp(totalPower),
		},
	}

	// currentTime から 1000 ミリ秒先までシミュレーションする
	for t := currentTime + 1; t <= currentTime+1000; t++ {
		totalMilliIsu.Add(totalMilliIsu, totalPower)
		updated := false

		// 時刻 t で発生する adding を計算する
		if a, ok := addingAt[t]; ok {
			updated = true
			totalMilliIsu.Add(totalMilliIsu, new(big.Int).Mul(str2big(a.Isu), big.NewInt(1000)))
		}

		// 時刻 t で発生する buying を計算する
		if _, ok := buyingAt[t]; ok {
			updated = true
			updatedID := map[int]bool{}
			for _, b := range buyingAt[t] {
				m := mItems[b.ItemID]
				updatedID[b.ItemID] = true
				itemBuilt[b.ItemID]++
				power := m.GetPower(b.Ordinal)
				itemPower[b.ItemID].Add(itemPower[b.ItemID], power)
				totalPower.Add(totalPower, power)
			}
			for id := range updatedID {
				itemBuilding[id] = append(itemBuilding[id], Building{
					Time:       t,
					CountBuilt: itemBuilt[id],
					Power:      big2exp(itemPower[id]),
				})
			}
		}

		if updated {
			schedule = append(schedule, Schedule{
				Time:       t,
				MilliIsu:   big2exp(totalMilliIsu),
				TotalPower: big2exp(totalPower),
			})
		}

		// 時刻 t で購入可能になったアイテムを記録する
		for itemID := range mItems {
			if _, ok := itemOnSale[itemID]; ok {
				continue
			}
			if 0 <= totalMilliIsu.Cmp(new(big.Int).Mul(itemPrice[itemID], big.NewInt(1000))) {
				itemOnSale[itemID] = t
			}
		}
	}

	gsAdding := []Adding{}
	for _, a := range addingAt {
		gsAdding = append(gsAdding, a)
	}

	gsItems := []Item{}
	for itemID, _ := range mItems {
		gsItems = append(gsItems, Item{
			ItemID:      itemID,
			CountBought: itemBought[itemID],
			CountBuilt:  itemBuilt0[itemID],
			NextPrice:   big2exp(itemPrice[itemID]),
			Power:       itemPower0[itemID],
			Building:    itemBuilding[itemID],
		})
	}

	gsOnSale := []OnSale{}
	for itemID, t := range itemOnSale {
		gsOnSale = append(gsOnSale, OnSale{
			ItemID: itemID,
			Time:   t,
		})
	}

	return &GameStatus{
		Adding:   gsAdding,
		Schedule: schedule,
		Items:    gsItems,
		OnSale:   gsOnSale,
	}, nil
}

