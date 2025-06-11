package types

type CpuInfo struct {
	Architecture string `json:"architecture"`
	CpuInfoAmd64
	CpuInfoArm64
}

type CpuInfoAmd64 struct {
	ManufacturerId string   `json:"manufacturer_id"`
	Flags          []string `json:"flags"`
}

type CpuInfoArm64 struct {
	ImplementerId HexInt `json:"implementer_id"`
	PartNumber    HexInt `json:"part_number"`
}
