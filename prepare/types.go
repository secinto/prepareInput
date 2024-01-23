package prepare

const VERSION = "0.1.0"

//const VERSION = "0.1.0" //Initial version

type Config struct {
	S2SPath     string `yaml:"s2s_path"`
	HostMapping string `yaml:"host_mapping,omitempty"`
	PortsSinple string `yaml:"ports_simple,omitempty"`
}

type Project struct {
	Name string `json:"name"`
}

type Preparer struct {
	options *Options
}

type Host struct {
	IPv4      string
	Hostnames []string
	Ports     []string
}
