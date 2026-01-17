package cli

import (
	"ca-server/config"
	"ca-server/internal/util"
	"path/filepath"

	"github.com/spf13/cobra"
)

var SetupCommand = &cobra.Command{
	Use:   "setup",
	Short: "Setup PKI Directory",
	Long:  "Setup PKI Directory with default Root CA or import existing RootCA",
	Run:   runSetup,
}

func runSetup(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg := config.New()
	pkiPath, _ := filepath.Abs(cfg.PKIPath)
	util.CreateDirectory(pkiPath)

	// create PKI root directory
	rootDir := pkiPath + "/roots"
	util.CreateDirectory(rootDir)

	// create PKI extra keys directory
	keyDir := pkiPath + "/keys"
	util.CreateDirectory(keyDir)

	// TODO: Setup Root CA
}
