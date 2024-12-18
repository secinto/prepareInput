package prepare

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var (
	log           = NewLogger()
	appConfig     Config
	project       Project
	sslPorts      = []string{"22", "465", "995", "3306"}
	startTLSPorts = map[string]string{"21": "ftp", "25": "smtp", "110": "pop3", "143": "imap", "587": "smtp", "993": "imap"}
)

/*
--------------------------------------------------------------------------------

	Initialization functions for the application

-------------------------------------------------------------------------------
*/
func (p *Preparer) initialize(configLocation string) {
	appConfig = loadConfigFrom(configLocation)
	if !strings.HasSuffix(appConfig.ProjectsPath, "/") {
		appConfig.ProjectsPath = appConfig.ProjectsPath + "/"
	}
	p.options.BaseFolder = appConfig.ProjectsPath + p.options.Project
	if !strings.HasSuffix(p.options.BaseFolder, "/") {
		p.options.BaseFolder = p.options.BaseFolder + "/"
	}
	//appConfig.DpuxHostToIP = strings.Replace(appConfig.DpuxHostToIP, "{project_name}", p.options.Project, -1)
	//appConfig.UniqueOpenPorts = strings.Replace(appConfig.UniqueOpenPorts, "{project_name}", p.options.Project, -1)

	project = Project{
		Name: p.options.Project,
	}
}

func loadConfigFrom(location string) Config {
	var config Config
	var yamlFile []byte
	var err error

	yamlFile, err = os.ReadFile(location)
	if err != nil {
		yamlFile, err = os.ReadFile(defaultSettingsLocation)
		if err != nil {
			log.Fatalf("yamlFile.Get err   #%v ", err)
		}
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if &config == nil {
		config = Config{
			ProjectsPath: "/checkfix/projects",
			HostMapping:  "dpux_host_to_ip.json",
			ServicesFile: "services.json",
		}
	}
	return config
}

func NewPreparer(options *Options) (*Preparer, error) {
	finder := &Preparer{options: options}
	finder.initialize(options.SettingsFile)
	return finder, nil
}

func (p *Preparer) Prepare() error {
	if p.options.Project != "" {
		if p.options.HostMapping {
			log.Infof("Create host mapping for project %s", p.options.Project)
			p.createHostMapping()
		}
		if p.options.TLSCheckFull {
			log.Infof("Preparing input for tls_check for all avaiable ports for project %s", p.options.Project)
			p.createHostMapping()
			p.prepareForTestSSL()
		}
	} else {
		log.Info("No project specified. Exiting application")
	}
	return nil
}

/*
--------------------------------------------------------------------------------
 Public functions of the application
-------------------------------------------------------------------------------
*/

func (p *Preparer) createHostMapping() {
	dnsToIPInfo := GetDocumentFromFile(p.options.BaseFolder+"recon/"+appConfig.HostMapping, true)
	ipToPortInfo := GetDocumentFromFile(p.options.BaseFolder+"findings/"+appConfig.ServicesFile, false)

	cleanedIPs := GetValuesForKey(dnsToIPInfo, "ip")

	for _, ip := range cleanedIPs {
		hostnames := GetValueForQueryKey(dnsToIPInfo, "host", "ip", []string{ip})
		ports := GetValueForQueryKey(ipToPortInfo, "port", "ip", []string{ip})
		p.preparedHosts = append(p.preparedHosts, Host{
			IPv4:      ip,
			Hostnames: hostnames,
			Ports:     ports,
		})
	}
	data, _ := json.MarshalIndent(p.preparedHosts, "", " ")
	hostMappingFile := p.options.BaseFolder + "findings/hostMapping.json"
	WriteToTextFileInProject(hostMappingFile, string(data))
	log.Infof("Processed %d host entries and their associated names and ports, created host mapping file %s", len(p.preparedHosts), hostMappingFile)
}

func (p *Preparer) prepareForTestSSL() {

	var tlsCheckEntries []string

	for _, host := range p.preparedHosts {
		for _, port := range host.Ports {
			if value, ok := startTLSPorts[port]; ok {
				log.Debugf("Port %s will be used for tls_check with StartTLS", port)
				tlsCheckEntries = append(tlsCheckEntries, "-t "+value+" "+host.IPv4+":"+port)
			} else if ExistsInArray(sslPorts, port) {
				log.Debugf("Port %s will be used for tls_check with SSL", port)
				tlsCheckEntries = append(tlsCheckEntries, host.IPv4+":"+port)
			} else {
				log.Debugf("Port %s will not be used for tls_check", port)
			}

		}
	}
	additionalTLSCheckFile := p.options.BaseFolder + "recon/tls_check_additional.txt"
	WriteToTextFileInProject(additionalTLSCheckFile, strings.Join(tlsCheckEntries[:], "\n")+"\n")
	log.Infof("Created %d entries for additional TLS check file %s", len(tlsCheckEntries), additionalTLSCheckFile)

	log.Info("Finished")
}
