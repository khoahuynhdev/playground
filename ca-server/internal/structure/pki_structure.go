package structure

type CertificateAuthorityPaths struct {
	RootCAPath               string
	RootCACertRequestsPath   string
	RootCACertsPath          string
	RootCACertRevListPath    string
	RootCANewCertsPath       string
	RootCACertKeysPath       string
	RootCAIntermediateCAPath string
	RootCACertIndexFilePath  string
	RootCACertSerialFilePath string
	RootCACrlnumFilePath     string
}

func NewCertificateAuthorityPaths(rootCAPath string) *CertificateAuthorityPaths {
	return &CertificateAuthorityPaths{
		RootCAPath:               rootCAPath,
		RootCACertRequestsPath:   rootCAPath + "/certreqs",
		RootCACertsPath:          rootCAPath + "/certs",
		RootCACertRevListPath:    rootCAPath + "/crl",
		RootCANewCertsPath:       rootCAPath + "/newcerts",
		RootCACertKeysPath:       rootCAPath + "/private",
		RootCAIntermediateCAPath: rootCAPath + "/intermediate",
		RootCACertIndexFilePath:  rootCAPath + "/index.txt",
		RootCACertSerialFilePath: rootCAPath + "/serial",
		RootCACrlnumFilePath:     rootCAPath + "/crlnumber",
	}
}

func (c *CertificateAuthorityPaths) SetupCAFileStructure(basePath string) *CertificateAuthorityPaths {
	rootCAPath := basePath
	CreateDirectory(rootCAPath)
	return nil
}
