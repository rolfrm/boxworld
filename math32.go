package main
import "math"
func Sin32(x float32) float32{
	return float32(math.Sin(float64(x)))
}

func Cos32(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sqrt32(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func Fabs32(x float32) float32 {
	if x < 0 {

		return -x
	}
	return x
}