package types

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
	BogoMips      float64  `json:"bogo_mips"`
	Flags         []string `json:"flags"`
}
