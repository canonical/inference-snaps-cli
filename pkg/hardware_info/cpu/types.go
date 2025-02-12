package cpu

type LsCpuContainer struct {
	LsCpu []LsCpuObject `json:"lscpu"`
}

type LsCpuObject struct {
	Field    string        `json:"field"`
	Data     string        `json:"data"`
	Children []LsCpuObject `json:"children"`
}
