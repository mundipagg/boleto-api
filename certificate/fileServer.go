package certificate

import (
	"fmt"
	"io/ioutil"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
)

const fileServerEnv = "fileServer"
const formatCertificate = ".pfx"

func InstanceStoreCertificatesFromFileServer(certificatesName ...string) (err error) {
	l := log.CreateLog()

	for _, certificateName := range certificatesName {
		pfxCertificate, err := ioutil.ReadFile(config.Get().CertificatesPath + certificateName + formatCertificate)
		if err != nil {
			return err
		}

		err = loadCertificatesOnStorage(fileServerEnv, certificateName, pfxCertificate)
		if err != nil {
			return err
		}

		l.Info(fmt.Sprintf("Success in load certificate [%s] from azureVault", certificateName))
	}

	return nil
}
