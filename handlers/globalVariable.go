package handlers

import "sync"

var ClientId string = "1"
var CUDA_DEVICE string
var CenterServerAddr string = "http://127.0.0.1:23333"
var HostAddr string

var requestFile map[string]*sync.Mutex
var requestFileLock sync.RWMutex

// var RenderFinish chan string
// var RelocaliseFinish chan relocaliseInfo //two scene
// var MergeMeshFinish chan MeshInfo
