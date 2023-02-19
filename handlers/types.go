package handlers

type relocaliseInfo struct {
	Scene1Name string `json:"scene1name"`
	Scene1IP   string `json:"scene1ip"`
	Scene2Name string `json:"scene2name"`
	Scene2IP   string `json:"scene2ip"`
}

type ResourceInfo struct {
	GPUMemoryFree uint64    `json:"gpumemoryfree"`
	MemoryFree    uint64    `json:"memoryfree"`
	CpuUsage      []float64 `json:"cpuusage"`
}

type globalPose struct {
	Scene1Name string `json:"scene1name"`
	Scene1Ip   string `json:"scene1ip"`
	Scene2Name string `json:"scene2name"`
	Scene2Ip   string `json:"scene2ip"`
	Transform  Pose   `json:"transform"`
}

type MeshInfo struct {
	Scenes     map[string]bool `json:"scenes"`
	WorldScene string          `json:"worldscene"`
	FileName   string          `json:"filename"`
	Client     string          `json:"client"`
}

type MergeMeshInfo struct {
	File1Name  string          `json:"file1name"`
	Scenes1    map[string]bool `json:"scenes1"`
	File1Ip    string          `json:"file1ip"`
	File2Name  string          `json:"file2name"`
	Scenes2    map[string]bool `json:"scenes2"`
	File2Ip    string          `json:"file2ip"`
	PoseMatrix [4][4]float64   `json:"posematrix"`
}
