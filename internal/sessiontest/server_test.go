package sessiontest

import (
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/privacybydesign/irmago/internal/test"
	"github.com/privacybydesign/irmago/server"
	"github.com/privacybydesign/irmago/server/irmaserver"
)

var (
	logger   = logrus.New()
	testdata = test.FindTestdataFolder(nil)
)

func init() {
	logger.Level = logrus.WarnLevel
	logger.Formatter = &logrus.TextFormatter{}
}

func StartIrmaServer(configuration *irmaserver.Configuration) {
	go func() {
		err := irmaserver.Start(configuration)
		if err != nil {
			panic("Starting server failed: " + err.Error())
		}
	}()
	time.Sleep(100 * time.Millisecond) // Give server time to start
}

func StopIrmaServer() {
	irmaserver.Stop()
}

var IrmaServerConfiguration = &irmaserver.Configuration{
	Configuration: &server.Configuration{
		Logger:                logger,
		IrmaConfigurationPath: filepath.Join(testdata, "irma_configuration"),
		IssuerPrivateKeysPath: filepath.Join(testdata, "privatekeys"),
	},
	DisableRequestorAuthentication: true,
	Port: 48682,
}

var JwtServerConfiguration = &irmaserver.Configuration{
	Configuration: &server.Configuration{
		Logger:                logger,
		IrmaConfigurationPath: filepath.Join(testdata, "irma_configuration"),
		IssuerPrivateKeysPath: filepath.Join(testdata, "privatekeys"),
	},
	Port: 48682,
	DisableRequestorAuthentication: false,
	GlobalPermissions: irmaserver.Permissions{
		Disclosing: []string{"*"},
		Signing:    []string{"*"},
		Issuing:    []string{"*"},
	},
	Requestors: map[string]irmaserver.Requestor{
		"requestor1": irmaserver.Requestor{
			AuthenticationMethod: irmaserver.AuthenticationMethodPublicKey,
			AuthenticationKey:    filepath.Join(testdata, "jwtkeys", "requestor1.pem"),
		},
		"requestor2": {
			AuthenticationMethod: irmaserver.AuthenticationMethodToken,
			AuthenticationKey:    "xa6=*&9?8jeUu5>.f-%rVg`f63pHim",
		},
	},
	JwtPrivateKey: filepath.Join(testdata, "jwtkeys", "sk.pem"),
}
