package handlers

import (
	"sync"
	"time"
)

var ClientId string = "1"
var CUDA_DEVICE string
var CenterServerAddr string = "http://172.24.109.142:23333"

// var CenterServerAddr string = "http://172.26.43.12:23333"
// var CenterServerAddr string = "http://127.0.0.1:23333"
var HostAddr string

var requestFile map[string]*sync.Mutex
var requestFileLock sync.RWMutex

// measurement
var TimeCost map[string]time.Duration
var TimeCostLock sync.RWMutex
var MemoryCost map[string]float64
var MemoryCostLock sync.RWMutex
var CpuUsage map[string]float64
var CpuUsageLock sync.RWMutex

// var RenderFinish chan string
// var RelocaliseFinish chan relocaliseInfo //two scene
// var MergeMeshFinish chan MeshInfo
