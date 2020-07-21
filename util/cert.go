package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"

	dc "github.com/hugocarreira/go-decent-copy"
)

//ListCert Lista os Certificados necessários e chama o método que faz a cópia
func ListCert() error {

	list := []string{
		config.Get().CertBoletoPathCrt,
		config.Get().CertBoletoPathKey,
		config.Get().CertBoletoPathCa,
		config.Get().CertICP_PathPkey,
		config.Get().CertICP_PathChainCertificates,
	}

	var err error

	for _, v := range list {

		err = copyCert(v)

		if err != nil {
			return err
		}

	}

	return nil

}

func copyCert(certificateDestiny string) error {

	originPath := strings.Replace(certificateDestiny, "boleto_cert", "boleto_orig", 1)

	err := dc.Copy(originPath, certificateDestiny)
	if err != nil {
		logCopy("Copy Certificates Error:" + err.Error())
		return err
	}

	err = os.Chmod(certificateDestiny, 0777)
	if err != nil {
		logCopy("Copy Certificates Error:" + err.Error())

		return err
	}

	logCopy("Certificate Copy Sucessful: " + certificateDestiny)

	return nil

}

func logCopy(msg string) {
	fmt.Println(msg)
	log.Info("[{Application}: Certificates] - " + msg)
}
