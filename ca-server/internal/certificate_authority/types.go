package ca

import "net"

/*
CertificateConfiguration is a struct to pass Certificate Config Information into the setup functions

`Subject` is a CertificateConfigurationSubject object

`ExpirationDate` is expressed as a slice of 3 ints [ years, months, days ] in the future

`RSAPrivateKey` is optional - this is used to sign a certificate request with an external key instead of one generated in the PKI

`RSAPrivateKeyPassphrase` is optional - this is used to secure the key if generated via PKI

`SANData` is a SANData object

`CertificateType` is a string representing what type of certificate is being requested or generated and is used in validation checks.  Options: server|client|authority|authority-no-subs
*/
type CertificateConfiguration struct {
	Subject                 CertificateConfigurationSubject `json:"subject"`
	ExpirationDate          []int                           `json:"expiration_date,omitempty"`
	RSAPrivateKey           string                          `json:"rsa_private_key,omitempty"`
	RSAPrivateKeyPassphrase string                          `json:"rsa_private_key_passphrase,omitempty"`
	SerialNumber            string                          `json:"serial_number,omitempty"`
	SANData                 SANData                         `json:"san_data,omitempty"`
	CertificateType         string                          `json:"certificate_type,omitempty"`
}

// SANData provides a collection of SANData for a certificate
type SANData struct {
	IPAddresses    []net.IP `json:"ip_addresses,omitempty"`
	EmailAddresses []string `json:"email_addresses,omitempty"`
	DNSNames       []string `json:"dns_names,omitempty"`
	URIs           []string `json:"uris,omitempty"`
	// URIs           []*url.URL `json:"uris,omitempty"`
}

// CertificateConfigurationSubject is simply a redefinition of pkix.Name
type CertificateConfigurationSubject struct {
	CommonName         string   `json:"common_name"`
	Organization       []string `json:"organization,omitempty"`
	OrganizationalUnit []string `json:"organizational_unit,omitempty"`
	Country            []string `json:"country,omitempty"`
	Province           []string `json:"province,omitempty"`
	Locality           []string `json:"locality,omitempty"`
	StreetAddress      []string `json:"street_address,omitempty"`
	PostalCode         []string `json:"postal_code,omitempty"`
}
