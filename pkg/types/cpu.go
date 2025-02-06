package types

type LsCpuContainer struct {
	LsCpu []LsCpuObject `json:"lscpu"`
}

type LsCpuObject struct {
	Field    string        `json:"field"`
	Data     string        `json:"data"`
	Children []LsCpuObject `json:"children"`
}

type CpuInfo struct {
	Architecture  string   `json:"architecture"`
	VendorId      string   `json:"vendor_id"`
	FamilyId      *int     `json:"family_id"`
	ModelId       int      `json:"model_id"`
	ModelName     string   `json:"model_name"`
	PhysicalCores int      `json:"physical_cores"`
	LogicalCores  int      `json:"logical_cores"`
	MaxFrequency  float64  `json:"max_frequency"`
	MinFrequency  float64  `json:"min_frequency"`
	Flags         []string `json:"flags"`
}
