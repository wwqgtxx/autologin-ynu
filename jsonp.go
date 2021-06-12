package main

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"time"
)

type Jsonp struct {
	CallbackString string
}

func randIntRange(max, min *big.Int) *big.Int {
	n := big.NewInt(0)
	tmp, err := rand.Int(rand.Reader, n.Sub(max, min))
	if err != nil {
		return n.Set(min)
	}
	return n.Add(tmp, min)
}

func randNumString(long int) string {
	min := big.NewInt(1)
	ten := big.NewInt(10)
	for i := 1; i < long; i++ {
		min.Mul(min, ten)
	}
	max := new(big.Int)
	max.Set(min)
	max.Mul(max, ten)
	return randIntRange(max, min).String()
}

func timestampString() string {
	return time.Now().Format("20060102150405")
}

func NewJsonp() *Jsonp {
	return &Jsonp{
		CallbackString: "jQuery" + randNumString(22) + "_" + timestampString(),
	}
}

func (j *Jsonp) RemoveJsonP(b []byte) []byte {
	b = bytes.TrimLeft(b, j.CallbackString+"(")
	b = bytes.TrimRight(b, ")")
	return b
}
