package certificate

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	pkcs122 "software.sslmate.com/src/go-pkcs12"
)

const (
	icp = "ICP"
	ssl = "SSL"
)

type ICPCertificate struct {
	Name          string
	RsaPrivateKey interface{}
	Certificate   *x509.Certificate
}

func GetICPCertificate(name string, pfxBytes []byte, pass string) (ICPCertificate, error) {
	rsaPrivateKey, certificate, err := pkcs122.Decode(pfxBytes, pass)
	if err != nil {
		return ICPCertificate{}, err
	}

	iCPCertificate := new(ICPCertificate)
	iCPCertificate.Name = name
	iCPCertificate.RsaPrivateKey = rsaPrivateKey
	iCPCertificate.Certificate = certificate

	return *iCPCertificate, nil
}

type SSLCertificate struct {
	Name    string
	PemData []byte
}

func GetSSLCertificate(name string, pfxBytes []byte, pass string) (SSLCertificate, error) {
	pemData, err := localToPEM(pfxBytes, pass)
	if err != nil {
		return SSLCertificate{}, err
	}

	sslCertificate := new(SSLCertificate)
	sslCertificate.Name = name
	sslCertificate.PemData = pemData

	return *sslCertificate, nil
}

func localToPEM(pfxBytes []byte, pass string) ([]byte, error) {
	privateKey, certificate, caChain, err := pkcs122.DecodeChain(pfxBytes, pass)

	if err != nil {
		return nil, err
	}

	var pemBytes bytes.Buffer
	err = pem.Encode(&pemBytes, &pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey))})
	if err != nil {
		return nil, err
	}

	err = pem.Encode(&pemBytes, &pem.Block{Type: "CERTIFICATE", Bytes: certificate.Raw})
	if err != nil {
		return nil, err
	}

	for _, certChain := range caChain {
		if err := pem.Encode(&pemBytes, &pem.Block{Type: "CERTIFICATE", Bytes: certChain.Raw}); err != nil {
			return nil, err
		}
	}

	return pemBytes.Bytes(), nil
}
