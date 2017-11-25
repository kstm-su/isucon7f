package main

import (
	"log"
	"math/big"
)

func buyItem(roomName string, itemID int, countBought int, reqTime int64) bool {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return false
	}

	_, ok := updateRoomTime(tx, roomName, reqTime)
	if !ok {
		tx.Rollback()
		return false
	}

	var countBuying int
	err = tx.Get(&countBuying, "SELECT COUNT(*) FROM buying WHERE room_name = ? AND item_id = ?", roomName, itemID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}
	if countBuying != countBought {
		tx.Rollback()
		log.Println(roomName, itemID, countBought+1, " is already bought")
		return false
	}

	totalMilliIsu := new(big.Int)
	var addings []Adding
	err = tx.Select(&addings, "SELECT isu FROM adding WHERE room_name = ? AND time <= ?", roomName, reqTime)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}

	for _, a := range addings {
		totalMilliIsu.Add(totalMilliIsu, new(big.Int).Mul(str2big(a.Isu), big.NewInt(1000)))
		//totalMilliIsu.Add(totalMilliIsu, Mul(str2big(a.Isu), big.NewInt(1000)))
	}

	var buyings []Buying
	err = tx.Select(&buyings, "SELECT item_id, ordinal, time FROM buying WHERE room_name = ?", roomName)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}
	for _, b := range buyings {
		//var item mItem
		//tx.Get(&item, "SELECT * FROM m_item WHERE item_id = ?", b.ItemID)
		item := MItems[b.ItemID]
		cost := new(big.Int).Mul(item.GetPrice(b.Ordinal), big.NewInt(1000))
		//cost := Mul(item.GetPrice(b.Ordinal), big.NewInt(1000))
		totalMilliIsu.Sub(totalMilliIsu, cost)
		if b.Time <= reqTime {
			gain := new(big.Int).Mul(item.GetPower(b.Ordinal), big.NewInt(reqTime-b.Time))
			//gain := Mul(item.GetPower(b.Ordinal), big.NewInt(reqTime-b.Time))
			totalMilliIsu.Add(totalMilliIsu, gain)
		}
	}

	//var item mItem
	//tx.Get(&item, "SELECT * FROM m_item WHERE item_id = ?", itemID)
	item := MItems[itemID]
	need := new(big.Int).Mul(item.GetPrice(countBought+1), big.NewInt(1000))
	//need := Mul(item.GetPrice(countBought+1), big.NewInt(1000))
	if totalMilliIsu.Cmp(need) < 0 {
		log.Println("not enough")
		tx.Rollback()
		return false
	}

	_, err = tx.Exec("INSERT INTO buying(room_name, item_id, ordinal, time) VALUES(?, ?, ?, ?)", roomName, itemID, countBought+1, reqTime)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return false
	}

	return true
}
