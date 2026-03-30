package config

type GKIConfig struct {
	Version            string   `yaml:"version"`
	BaseDefconfig      string   `yaml:"base_defconfig"`
	Image              string   `yaml:"image"`
	BootHeaderVersion  int      `yaml:"boot_header_version"`
	RamdiskCompression string   `yaml:"ramdisk_compression"`
	HasInitBoot        bool     `yaml:"has_init_boot"`
	VendorDLKMPath     string   `yaml:"vendor_dlkm_path"`
	VendorBootPath     string   `yaml:"vendor_boot_path"`
	MakeTargets        []string `yaml:"make_targets"`
}

type ToolchainProfile struct {
	ExtraCFlags   []string `yaml:"extra_cflags"`
	KbuildLDFlags []string `yaml:"kbuild_ldflags"`
}

type ToolchainConfig struct {
	Clang     string                      `yaml:"clang"`
	Prebuilts string                      `yaml:"prebuilts"`
	Jobs      string                      `yaml:"jobs"`
	Flags     map[string]string           `yaml:"flags"`
	Profiles  map[string]ToolchainProfile `yaml:"profiles"`
	Profile   string                      `yaml:"profile"`
}

type ExtModules struct {
	Root    string   `yaml:"root"`
	Sources []string `yaml:"sources"`
}

type DefconfigConfig struct {
	Base      string   `yaml:"base"`
	Fragments []string `yaml:"fragments"`
	Custom    []string `yaml:"custom"`
}

type DTBConfig struct {
	Wildcard     string   `yaml:"wildcard"`
	DTBOWildcard string   `yaml:"dtbo_wildcard"`
	TechpackDirs []string `yaml:"techpack_dirs"`
}

type ModulesConfig struct {
	SecondStage string `yaml:"second_stage"`
	VendorDLKM  string `yaml:"vendor_dlkm"`
}

type OutputConfig struct {
	Dist string `yaml:"dist"`
	Zip  string `yaml:"zip"`
}

type KernelConfig struct {
	Version         string            `yaml:"version"`
	Localversion    string            `yaml:"localversion"`
	Defconfig       DefconfigConfig   `yaml:"defconfig"`
	ConfigOverrides map[string]string `yaml:"config_overrides"`
}

type DeviceInfo struct {
	Name     string `yaml:"name"`
	Platform string `yaml:"platform"`
	SOC      string `yaml:"soc"`
}

type BuildConfig struct {
	Extends    string          `yaml:"extends"`
	GKI        GKIConfig       `yaml:"gki"`
	Platform   string          `yaml:"platform"`
	KernelSrc  string          `yaml:"kernel_source"`
	Devicetrees string         `yaml:"devicetrees"`
	Device     DeviceInfo      `yaml:"device"`
	Kernel     KernelConfig    `yaml:"kernel"`
	Toolchain  ToolchainConfig `yaml:"toolchain"`
	ExtModules ExtModules      `yaml:"ext_modules"`
	DTB        DTBConfig       `yaml:"dtb"`
	Modules    ModulesConfig   `yaml:"modules"`
	Output     OutputConfig    `yaml:"output"`
}

type BuildMode int

const (
	BuildModeDevice BuildMode = iota
	BuildModeGKI
)
