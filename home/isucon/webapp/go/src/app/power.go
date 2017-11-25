package main

import (
	"math/big"

	_ "github.com/go-sql-driver/mysql"
)

func Mul(s, t *big.Int) *big.Int {
	sBit := s.BitLen()
	tBit := t.BitLen()
	sMove := uint(0)
	tMove := uint(0)
	if sBit > prec {
		sMove = uint(sBit - prec)
	}
	if tBit > prec {
		tMove = uint(tBit - prec)
	}

	s.Rsh(s, sMove)
	t.Rsh(t, tMove)

	ret := new(big.Int).Mul(s, t)
	return ret.Lsh(ret, sMove+tMove)
}

func (item *mItem) GetPower(count int) *big.Int {
	return item.GetPowerNext(count)
}

/*
func (item *mItem) GetPowerPrev(count int) *big.Int {
	// power(x):=(cx+1)*d^(ax+b)
	a := item.Power1
	b := item.Power2
	c := item.Power3
	d := item.Power4
	x := int64(count)

	s := big.NewInt(c*x + 1)
	t := new(big.Int).Exp(big.NewInt(d), big.NewInt(a*x+b), nil)
	return new(big.Int).Mul(s, t)
}
*/
func (item *mItem) GetPowerNext(count int) *big.Int {
	// power(x):=(cx+1)*d^(ax+b)
	a := item.Power1
	b := item.Power2
	c := item.Power3
	d := item.Power4
	x := int64(count)

	s := big.NewInt(c*x + 1)
	t := new(big.Int).Exp(big.NewInt(d), big.NewInt(a*x+b), nil)
	//return new(big.Int).Mul(s, t)
	return Mul(s, t)
	/*
		sBit := s.BitLen()
		tBit := t.BitLen()
		sMove := 0
		tMove := 0
		if sBit > prec {
			sMove = sBit - prec
		}
		if tBit > prec {
			tMove = tBit - prec
		}

		s.Rsh(s, uint(sMove))
		t.Rsh(t, uint(tMove))

		ret := new(big.Int).Mul(s, t)
		return ret.Lsh(ret, uint(sMove+tMove))
	*/
}
func (item *mItem) GetPrice(count int) *big.Int {
	// price(x):=(cx+1)*d^(ax+b)
	a := item.Price1
	b := item.Price2
	c := item.Price3
	d := item.Price4
	x := int64(count)

	s := big.NewInt(c*x + 1)
	t := new(big.Int).Exp(big.NewInt(d), big.NewInt(a*x+b), nil)
	return new(big.Int).Mul(s, t)
	/*
		//return Mul(s, t)
		sBit := s.BitLen()
		tBit := t.BitLen()
		sMove := 0
		tMove := 0
		if sBit > prec {
			sMove = sBit - prec
		}
		if tBit > prec {
			tMove = tBit - prec
		}

		s.Rsh(s, uint(sMove))
		t.Rsh(t, uint(tMove))

		ret := new(big.Int).Mul(s, t)
		return ret.Lsh(ret, uint(sMove+tMove))
	*/
}

/*
func (item *mItem) GetPrice(count int) *big.Int {
	// price(x):=(cx+1)*d^(ax+b)
	a := item.Price1
	b := item.Price2
	c := item.Price3
	d := item.Price4
	x := int64(count)

	s := big.NewInt(c*x + 1)
	t := new(big.Int).Exp(big.NewInt(d), big.NewInt(a*x+b), nil)
	return new(big.Int).Mul(s, t)
}

*/
