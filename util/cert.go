package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/mundipagg/boleto-api/config"

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

	return err

}

func copyCert(d string) error {
	o := strings.Replace(d, "boleto_orig", "boleto_cert", 1)

	err := dc.Copy(o, d)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return err
	}

	err = os.Chmod(d, 0777)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err
	}

	fmt.Println("Cert Copy Sucessful: ", d)

	return err
}
