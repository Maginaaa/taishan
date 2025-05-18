package libs

import (
	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/shopspring/decimal"
	"reflect"
)

// DiffCheck 比较大小前检查参数是否正确
func DiffCheck(value, target interface{}) (valueOf, targetOf reflect.Value) {
	valueOf = reflect.ValueOf(value)
	targetOf = reflect.ValueOf(target)

	// 确保两个值都是可比较的
	if !valueOf.IsValid() || !targetOf.IsValid() || !valueOf.Type().Comparable() || !targetOf.Type().Comparable() {
		panic("incomparable types")
	}

	// 确保两个值的类型相同
	if valueOf.Type() != targetOf.Type() {
		panic("different types")
	}
	return
}

// Max 取两个数中较大的那个值
func Max(value, target interface{}) interface{} {
	valueOf, targetOf := DiffCheck(value, target)

	switch valueOf.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valueOf.Int() > targetOf.Int() {
			return value
		}
		return target
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if valueOf.Uint() > targetOf.Uint() {
			return value
		}
		return target
	case reflect.Float32, reflect.Float64:
		if valueOf.Float() > targetOf.Float() {
			return target
		}
		return target
	case reflect.String:
		if valueOf.String() > targetOf.String() {
			return value
		}
		return target
	default:
		return value
	}
}

// CalcRps 计算RPS(每秒请求数)
func CalcRps(requestNum int64, startTime int64, endTime int64, isTotal bool) (rps float64) {
	// 计算持续时间，如果持续时间(结束时间 - 开始时间)大于5000ms, 则持续时间设定为5000ms
	var duration = endTime - startTime
	if isTotal == false {
		duration = mathutil.Min(5000, duration)
	}
	//如果分子或分母为0, 则直接返回0
	if requestNum == 0 || float64(duration) == 0 {
		return 0
	}
	// RPS = 请求数 / 持续时间 * 1000
	rps, _ = decimal.NewFromFloat(float64(requestNum) / float64(duration) * 1000).Round(2).Float64()
	return
}

// CalcRpsNew 计算RPS(每秒请求数) 时间单位毫秒
func CalcRpsNew(requestNum int64, duration int64) (rps float64) {
	// 计算持续时间，如果持续时间(结束时间 - 开始时间)大于5000ms, 则持续时间设定为5000ms
	//如果分子或分母为0, 则直接返回0
	if requestNum == 0 || float64(duration) == 0 {
		return 0
	}
	// RPS = 请求数 / 持续时间 * 1000
	rps, _ = decimal.NewFromFloat(float64(requestNum) / float64(duration) * 1000).Round(2).Float64()
	return
}

// CalcDiv 计算除数
func CalcDiv(numerator int64, denominator int64) (avg float64) {
	if numerator == 0 || denominator == 0 {
		return 0
	}
	avg, _ = decimal.NewFromFloat(float64(numerator) / float64(denominator)).Round(2).Float64()
	return
}

// CalcRate 计算比率
func CalcRate(numerator int64, denominator int64) (rate float64) {
	if numerator == 0 || denominator == 0 {
		return 0
	}
	rate, _ = decimal.NewFromFloat(float64(numerator) / float64(denominator) * 100).RoundFloor(2).Float64()
	return
}
