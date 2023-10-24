package helper

import (
	"math/big"
	"strconv"
)

// 加
func Add(num1 string, num2 string) string {
	result := new(big.Float)

	num1F64, _ := strconv.ParseFloat(num1, 64)
	num2F64, _ := strconv.ParseFloat(num2, 64)
	result.Add(big.NewFloat(num1F64), big.NewFloat(num2F64))

	return result.String()
}

// 減
func Sub(num1 string, num2 string) string {
	result := new(big.Float)

	num1F64, _ := strconv.ParseFloat(num1, 64)
	num2F64, _ := strconv.ParseFloat(num2, 64)
	result.Sub(big.NewFloat(num1F64), big.NewFloat(num2F64))

	return result.String()
}

// 乘
func Mul(num1 string, num2 string) string {
	result := new(big.Float)

	num1F64, _ := strconv.ParseFloat(num1, 64)
	num2F64, _ := strconv.ParseFloat(num2, 64)
	result.Mul(big.NewFloat(num1F64), big.NewFloat(num2F64))

	return result.String()
}

// 除
func Div(num1 string, num2 string) string {
	result := new(big.Float)

	num1F64, _ := strconv.ParseFloat(num1, 64)
	num2F64, _ := strconv.ParseFloat(num2, 64)
	result.Quo(big.NewFloat(num1F64), big.NewFloat(num2F64))

	return result.String()
}
