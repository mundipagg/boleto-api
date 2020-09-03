package app

import (
	"github.com/mundipagg/boleto-api/api"
	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/usermanagement"
)

//Params this struct contains all execution parameters to run application
type Params struct {
	DevMode    bool
	MockMode   bool
	DisableLog bool
}

//NewParams returns new Empty pointer to ExecutionParameters
func NewParams() *Params {
	return new(Params)
}

//Run starts boleto api Application
func Run(params *Params) {
	env.Config(params.DevMode, params.MockMode, params.DisableLog)

	if config.Get().MockMode {
		go mock.Run("9091")
	}
	log.Install()

	installCertificates()

	usermanagement.LoadUserCredentials()

	api.InstallRestAPI()

}

func installCertificates() {
	l := log.CreateLog()

	if config.Get().MockMode == false && config.Get().EnableFileServerCertificate == false {
		err := certificate.InstanceStoreCertificatesFromAzureVault(config.Get().VaultName, config.Get().CertificateICPName, config.Get().CertificateSSLName)
		if err == nil {
			l.Info("Success in load certificates from azureVault")
		} else {
			l.Error(err.Error(), "Error in load certificates from azureVault")
		}
	}

	if config.Get().MockMode == false && config.Get().EnableFileServerCertificate == true {
		err := certificate.InstanceStoreCertificatesFromFileServer(config.Get().CertificateICPName, config.Get().CertificateSSLName)
		if err == nil {
			l.Info("Success in load certificates from fileserver")
		} else {
			l.Error(err.Error(), "Error in load certificates from fileServer")
		}
	}
}
