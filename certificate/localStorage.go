package certificate

import (
	"errors"

	"github.com/mundipagg/boleto-api/config"
)

var localCertificateStorage = map[string]interface{}{}

func SetCertificateOnStore(key string, value interface{}) {
	localCertificateStorage[key] = value
}

func GetCertificateFromStore(key string) (interface{}, error) {
	if value, ok := localCertificateStorage[key]; ok {
		return value, nil
	}
	return nil, errors.New("Certificate not found.")
}

func getPassWordToCertificate(from string, certType string) string {
	if from == azureVaultEnv {
		return ""
	} else if from == fileServerEnv && certType == icp {
		return config.Get().PswCertificateICP
	} else if from == fileServerEnv && certType == ssl {
		return config.Get().PswCertificateSSL
	}
	return ""
}

func loadCertificatesOnStorage(from string, certificateName string, pfxBytes []byte) error {
	switch certificateName {
	case config.Get().CertificateICPName:
		var certificateICP, err = GetICPCertificate(certificateName, pfxBytes, getPassWordToCertificate(from, icp))
		if err != nil {
			return err
		}
		SetCertificateOnStore(certificateName, certificateICP)

	case config.Get().CertificateSSLName:
		var certificateSSL, err = GetSSLCertificate(certificateName, pfxBytes, getPassWordToCertificate(from, ssl))
		if err != nil {
			return err
		}
		SetCertificateOnStore(certificateName, certificateSSL)

	}
	return nil
}
