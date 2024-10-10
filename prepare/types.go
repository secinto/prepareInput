package prepare

const VERSION = "0.2.2"

//const VERSION = "0.1.0" //Initial version

type Config struct {
	ProjectsPath string `yaml:"projects_path"`
	HostMapping  string `yaml:"host_mapping,omitempty"`
	PortsSimple  string `yaml:"ports_simple,omitempty"`
	ServicesFile string `yaml:"services_file,omitempty"`
}

type Project struct {
	Name string `json:"name"`
}

type Preparer struct {
	options       *Options
	preparedHosts []Host
}

type Host struct {
	IPv4      string
	Hostnames []string
	Ports     []string
}
