package structure

import (
	"ca-server/internal/util"
	"fmt"
)

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

func (c *CertificateAuthorityPaths) SetupCAFileStructure(basePath string) *CertificateAuthorityPaths {
	// NOTE: the actual path to save the directory can be an SFTP server or blob storage
	// create root CA directory
	rootCAPath := basePath
	util.CreateDirectory(rootCAPath)

	// Create certificate requests (CSR) path
	rootCACertRequestsPath := rootCAPath + "/csrs"
	util.CreateDirectory(rootCACertRequestsPath)

	// Create certs path
	rootCACertsPath := rootCAPath + "/certs"
	util.CreateDirectory(rootCACertsPath)

	// create crls path
	rootCACertRevListPath := rootCAPath + "/crl"
	util.CreateDirectory(rootCACertRevListPath)

	// Answer from https://unix.stackexchange.com/questions/398280/openssl-basic-configuration-new-certs-dir-certs
	// new_certs_dir is used by the CA to output newly generated certs.
	rootCANewCertsPath := rootCAPath + "/newcerts"
	util.CreateDirectory(rootCANewCertsPath)

	// Create private path for CA keys
	rootCAKeysPath := rootCAPath + "/private"
	util.CreateDirectory(rootCAKeysPath)

	// Create private path for generated keys in the CA
	rootCACertKeysPath := rootCAPath + "/keys"
	util.CreateDirectory(rootCACertKeysPath)

	// Create intermediate CA path
	rootCAIntermediateCAPath := rootCAPath + "/intermed-ca"
	util.CreateDirectory(rootCAIntermediateCAPath)

	//  create index database file
	rootCACertIndexFilePath := rootCAPath + "/ca.index"
	indexExists, err := util.WriteFile(rootCACertIndexFilePath, "", 0600, false)
	if err != nil {
		fmt.Println("Error creating index file:", err)
	}
	if !indexExists {
		fmt.Println("Index file already exists")
	}

	// Create serial file
	rootCACertSerialFilePath := rootCAPath + "/ca.serial"
	serialExists, err := util.WriteFile(rootCACertSerialFilePath, "01", 0600, false)
	if err != nil {
		fmt.Println("Error creating serial file:", err)
	}
	if !serialExists {
		fmt.Println("Serial file already exists")
	}

	// Create certificate revocation number file
	rootCACrlnumFilePath := rootCAPath + "/ca.crlnum"
	crlNumExists, err := util.WriteFile(rootCACrlnumFilePath, "01", 0600, false)
	if err != nil {
		fmt.Println("Error creating crlnum file:", err)
	}
	if !crlNumExists {
		fmt.Println("Crlnum file already exists")
	}

	return &CertificateAuthorityPaths{
		rootCAPath,
		rootCACertRequestsPath,
		rootCANewCertsPath,
		rootCACertRevListPath,
		rootCANewCertsPath,
		rootCACertKeysPath,
		rootCAIntermediateCAPath,
		rootCACertIndexFilePath,
		rootCACertSerialFilePath,
		rootCACrlnumFilePath,
	}
}
