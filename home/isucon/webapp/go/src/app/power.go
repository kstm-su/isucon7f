func (item *mItem) GetPower(count int) *big.Int {
	log.Println("prev=", item.GetPowerPrev(count))
	ret, acc := big.NewFloat(item.GetPowerFloat64(count)).Int(big.NewInt(0))
	if acc != -1 {
		//x := big.NewInt(int64(acc))
		//ans := big.NewInt(0)
		log.Println("ret= ", ret)
		//log.Println("x= ", x)
		//log.Println("acc=", acc)
		//ans.Add(ret, x)
		return ret
	} else {
		return item.GetPowerPrev(count)
	}
}

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

func (item *mItem) GetPowerFloat64(count int) float64 {
	// power(x):=(cx+1)*d^(ax+b)
	a := item.Power1
	b := item.Power2
	c := item.Power3
	d := item.Power4
	x := int64(count)

	return float64(c*x+1) * math.Pow(float64(d), float64(a*x+b))
	//	s := big.NewInt(c*x + 1)
	//	t := new(big.Int).Exp(big.NewInt(d), big.NewInt(a*x+b), nil)
	//	return new(big.Int).Mul(s, t)
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
}

