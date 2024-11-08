package main

import "hardware-info/pkg/cpu"

func main() {
	//memoryInfo, err := memory.GetInfo()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v", memoryInfo)

	cpu.GetInfo()
}
