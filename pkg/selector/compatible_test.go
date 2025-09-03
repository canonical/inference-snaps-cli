package selector

import "testing"

type testValidInvalid struct {
	ValidMachines   []string
	InvalidMachines []string
}

var validInvalidSets = map[string]testValidInvalid{
	"ampere": {
		ValidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
		},
		InvalidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"ampere-altra": {
		ValidMachines: []string{
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"arm-neon": {
		ValidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
		},
		InvalidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"cpu-avx1": {
		ValidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
		},
	},

	"cpu-avx2": {
		ValidMachines: []string{
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-zbook-power-16-inch-g11",
			"i7-1165G7",
			"i7-10510U",
			"mustang",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-2600k+arc-a580",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
		},
	},

	"cpu-avx512": {
		ValidMachines: []string{
			"hp-pavilion-15-cs-3037nl",
			"i7-1165G7",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"example-memory": {
		ValidMachines: []string{
			"hp-zbook-power-16-inch-g11",
			"mustang",
			"system76-addw4",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"xps13-7390",
		},
	},

	"generic-cuda": {
		ValidMachines: []string{},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl", // Not enough vram
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"hp-zbook-power-16-inch-g11", // nvidia drivers not installed
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
			"xps13-9350",
		},
	},

	"intel-gpu": {
		ValidMachines: []string{
			"hp-zbook-power-16-inch-g11",
			"i7-2600k+arc-a580",
			"mustang",
			"system76-addw4",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l", // has intel gpu, but clinfo not working
			"hp-pavilion-15-cs-3037nl",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7", // 9a49 TigerLake-LP GT2 [Iris Xe Graphics]
			"i7-10510U",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"xps13-7390",
		},
	},

	"intel-npu": {
		ValidMachines: []string{
			"hp-zbook-power-16-inch-g11",
			"xps13-9350",
		},
		InvalidMachines: []string{
			"ampere-one-m-banshee-12",
			"ampere-one-siryn",
			"ampere-one-x-banshee-8",
			"asus-ux301l",
			"hp-pavilion-15-cs-3037nl",
			"hp-proliant-rl300-gen11-altra",
			"hp-proliant-rl300-gen11-altra-max",
			"i7-1165G7",
			"i7-2600k+arc-a580",
			"i7-10510U",
			"mustang",
			//"orange-pi-rv2",
			"raspberry-pi-5",
			"raspberry-pi-5+hailo-8",
			"system76-addw4",
			"xps13-7390",
		},
	},
}

func TestStack(t *testing.T) {
	for stackName, testSet := range validInvalidSets {
		for _, hwName := range testSet.ValidMachines {
			t.Run(stackName+" == "+hwName, func(t *testing.T) {
				testValidHw(t, stackName, hwName)
			})
		}

		for _, hwName := range testSet.InvalidMachines {
			t.Run(stackName+" != "+hwName, func(t *testing.T) {
				testInvalidHw(t, stackName, hwName)
			})
		}
	}
}
