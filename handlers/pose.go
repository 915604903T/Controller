package handlers

import (
	"log"
	"os"
	"strconv"
)

type Pose struct {
	Matrix [4][4]float64 `json:"matrix"`
	HasM   bool          `json:"hasM"`
	//w, x, y, z
	Q [4]float64 `json:"r"`
	P [4]float64 `json:"d"`
}

func NewPoseDq(dq [2][4]float64) Pose {
	return Pose{
		Q:    dq[0],
		P:    dq[1],
		HasM: false,
	}
}
func NewPoseMatrix(matrix [4][4]float64) Pose {
	return Pose{
		Matrix: matrix,
		HasM:   true,
	}
}

func (p *Pose) GetM() [4][4]float64 {
	if p.HasM {
		return p.Matrix
	}
	p.Matrix = p.DualQuat2Matrix()
	p.HasM = true
	return p.Matrix
}

func (p *Pose) DualQuat2Matrix() [4][4]float64 {
	mat := [4][4]float64{}

	mat[0][0] = p.Q[0]*p.Q[0] + p.Q[1]*p.Q[1] - p.Q[2]*p.Q[2] - p.Q[3]*p.Q[3]
	mat[0][1] = 2.0 * (p.Q[1]*p.Q[2] - p.Q[0]*p.Q[3])
	mat[0][2] = 2.0 * (p.Q[1]*p.Q[3] + p.Q[0]*p.Q[2])
	mat[0][3] = 2.0 * (p.Q[0]*p.P[1] - p.Q[1]*p.P[0] + p.Q[2]*p.P[3] - p.Q[3]*p.P[2])

	mat[1][0] = 2.0 * (p.Q[1]*p.Q[2] + p.Q[0]*p.Q[3])
	mat[1][1] = p.Q[0]*p.Q[0] - p.Q[1]*p.Q[1] + p.Q[2]*p.Q[2] - p.Q[3]*p.Q[3]
	mat[1][2] = 2.0 * (p.Q[2]*p.Q[3] - p.Q[0]*p.Q[1])
	mat[1][3] = 2.0 * (p.Q[0]*p.P[2] - p.Q[1]*p.P[3] - p.Q[2]*p.P[0] + p.Q[3]*p.P[1])

	mat[2][0] = 2.0 * (p.Q[1]*p.Q[3] - p.Q[0]*p.Q[2])
	mat[2][1] = 2.0 * (p.Q[2]*p.Q[3] + p.Q[0]*p.Q[1])
	mat[2][2] = p.Q[0]*p.Q[0] - p.Q[1]*p.Q[1] - p.Q[2]*p.Q[2] + p.Q[3]*p.Q[3]
	mat[2][3] = 2.0 * (p.Q[0]*p.P[3] + p.Q[1]*p.P[2] - p.Q[2]*p.P[1] - p.Q[3]*p.P[0])

	mat[3][0] = 0.0
	mat[3][1] = 0.0
	mat[3][2] = 0.0
	mat[3][3] = 1.0
	return mat
}

func writeTmpPoseFile(fileName string, poseM [4][4]float64) {
	fp, err := os.Create(fileName)
	defer fp.Close()
	if err != nil {
		log.Println("create ", fileName, " err: ", err)
		panic(err)
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			strNum := strconv.FormatFloat(poseM[i][j], 'f', -1, 64)
			_, err := fp.Write([]byte(strNum))
			if err != nil {
				log.Println("write ", strNum, " to file ", fileName, " err:", err)
				panic(err)
			}
			if j != 3 {
				_, err := fp.Write([]byte(" "))
				if err != nil {
					log.Println("write to file ", fileName, " err:", err)
					panic(err)
				}
			}
		}
		if i != 3 {
			_, err := fp.Write([]byte("\n"))
			if err != nil {
				log.Println("write enter to file ", fileName, " err:", err)
				panic(err)
			}
		}
	}
}
