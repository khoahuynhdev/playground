# Directory Structure

We use a standardized directory Structure like OpenSSL.

```cfg
#
# OpenSSL configuration for the Root Certification Authority.
#

#
# This definition doesn't work if HOME isn't defined.
CA_HOME                 = .

#
# Default Certification Authority
[ ca ]
default_ca              = root_ca

#
# Root Certification Authority
[ root_ca ]
dir                     = $ENV::CA_HOME
certs                   = $dir/certs
serial                  = $dir/ca.serial
database                = $dir/ca.index
new_certs_dir           = $dir/newcerts
certificate             = $dir/ca.cert
private_key             = $dir/private/ca.key.pem
default_days            = 1826 # 5 years
crl                     = $dir/ca.crl
crl_dir                 = $dir/crl
crlnumber               = $dir/ca.crlnum
```

We can represent a similar structure in Go

```go
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
```

This structure will be used to store the paths of the CA files. The paths will be set in the `PreparePKIDirectory` function and will be used throughout the application.
