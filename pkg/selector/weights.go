package selector

const (
	WeightPci       = 1000
	WeightPciVendor = 200

	WeightGpu                  = 100
	WeightGpuVendor            = 20
	WeightGpuVRam              = 10
	WeightGpuComputeCapability = 10

	WeightCpu       = 10
	WeightCpuVendor = 3
	WeightCpuModel  = 2
	WeightCpuFlag   = 1
)
