package cpu

import (
	"fmt"
	"log"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

var procCpuInfoPath = "/proc/cpuinfo"

func GetInfo() {
	cpuInfo, err := linuxproc.ReadCPUInfo(procCpuInfoPath)
	if err != nil {
		log.Fatal("stat read fail")
	}

	fmt.Printf("Number of Sockets: %d\n", cpuInfo.NumPhysicalCPU())
	fmt.Printf("Number of Cores: %d\n", cpuInfo.NumCore())
	fmt.Printf("Number of CPUs: %d\n", cpuInfo.NumCPU())

	type CpuModel struct {
		Vendor string
		//Family string // Family is only reported for Intel
		Model string
	}

	type ModelCount struct {
		SocketCoreThreads map[string]int
	}

	coreInfo := map[CpuModel]ModelCount{}

	for _, core := range cpuInfo.Processors {
		//log.Printf("%+v", core)
		pid := core.PhysicalId
		cid := core.CoreId

		// to avoid fmt import
		key := fmt.Sprintf("s%dc%d", pid, cid)
		//log.Println(key)

		model := CpuModel{
			Vendor: core.VendorId,
			Model:  core.ModelName,
		}
		if _, ok := coreInfo[model]; !ok {
			coreInfo[model] = ModelCount{SocketCoreThreads: make(map[string]int)}
		}
		//if _, ok := coreInfo[model].SocketCoreThreads[int(pid)]; !ok {
		//	coreInfo[model].SocketCoreThreads[int(pid)] = 0
		//}
		coreInfo[model].SocketCoreThreads[key]++
	}

	for vendorModel, cpu := range coreInfo {
		cores := len(cpu.SocketCoreThreads)
		threads := 0
		for _, coreThreads := range cpu.SocketCoreThreads {
			threads += coreThreads
		}

		fmt.Printf("%s %s - cores %d - threads %d\n", vendorModel.Vendor, vendorModel.Model, cores, threads)
	}
}
