package handlers

import (
	"math"
)

type Matrix struct {
	//矩阵结构
	N, M int //m是列数,n是⾏数
	Data [][]float64
}

func Mul(a [4][4]float64, b [4][4]float64) [4][4]float64 {
	res := [4][4]float64{}
	for i := 0; i < 4; i++ {
		t := [4]float64{}
		for j := 0; j < 4; j++ {
			r := float64(0)
			for k := 0; k < 4; k++ {
				r += a[i][k] * b[k][j]
			}
			t[j] = r
		}
		res[i] = t
	}
	return res
}

func Det(Matrix [4][4]float64, N int) float64 {
	var T0, T1, T2, Cha int
	var Num float64
	var B [4][4]float64

	if N > 0 {
		Cha = 0
		for i := 0; i < N; i++ {
			var tmpArr [4]float64
			for j := 0; j < N; j++ {
				tmpArr[j] = 0
			}
			B[i] = tmpArr
		}
		Num = 0
		for T0 = 0; T0 <= N; T0++ { //T0循环
			for T1 = 1; T1 <= N; T1++ { //T1循环
				for T2 = 0; T2 <= N-1; T2++ { //T2循环
					if T2 == T0 {
						Cha = 1
					}
					B[T1-1][T2] = Matrix[T1][T2+Cha]
				} //T2循环
				Cha = 0
			} //T1循环
			Num = Num + Matrix[0][T0]*Det(B, N-1)*math.Pow(-1, float64(T0))
		} //T0循环
		return Num
	} else if N == 0 {
		return Matrix[0][0]
	}
	return 0
}

func Inverse(S1 [4][4]float64) (MatrixC [4][4]float64) {
	N := 4 - 1
	Matrix := S1
	var T0, T1, T2, T3 int
	var B [4][4]float64
	for i := 0; i < N; i++ {
		var tmpArr [4]float64
		for j := 0; j < N; j++ {
			tmpArr[j] = 0
		}
		B[i] = tmpArr
	}

	for i := 0; i < N+1; i++ {
		var tmpArr [4]float64
		for j := 0; j < N+1; j++ {
			tmpArr[j] = 0
		}
		MatrixC[i] = tmpArr
	}

	Chay := 0
	Chax := 0
	add := 1 / Det(Matrix, N)
	for T0 = 0; T0 <= N; T0++ {
		for T3 = 0; T3 <= N; T3++ {
			for T1 = 0; T1 <= N-1; T1++ {
				if T1 < T0 {
					Chax = 0
				} else {
					Chax = 1
				}
				for T2 = 0; T2 <= N-1; T2++ {
					if T2 < T3 {
						Chay = 0
					} else {
						Chay = 1
					}
					B[T1][T2] = Matrix[T1+Chax][T2+Chay]
				} //T2循环
			} //T1循环
			Det(B, N-1)
			MatrixC[T3][T0] = Det(B, N-1) * add * (math.Pow(-1, float64(T0+T3)))
		}
	}
	return MatrixC
}
