package certificate

import (
	"io/ioutil"

	"github.com/mundipagg/boleto-api/config"
)

const fileServerEnv = "fileServer"
const formatCertificate = ".pfx"

func InstanceStoreCertificatesFromFileServer(certificatesName ...string) (err error) {

	for _, certificateName := range certificatesName {
		pfxCertificate, err := ioutil.ReadFile(config.Get().CertificatesPath + certificateName + formatCertificate)
		if err != nil {
			return err
		}

		err = loadCertificatesOnStorage(fileServerEnv, certificateName, pfxCertificate)
		if err != nil {
			return err
		}
	}

	return nil
}
