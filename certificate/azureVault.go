package certificate

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"

	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

const azureVaultEnv = "vault"

type azureKeyVaultCertificate struct {
	Ctx           context.Context
	VaultName     string
	Client        keyvault.BaseClient
	authenticated bool
	vaultBaseURL  string
}

func InstanceStoreCertificatesFromAzureVault(vaultName string, certificatesName ...string) (err error) {
	ctx := context.Background()
	configCertificate := azureKeyVaultCertificate{
		Ctx:       ctx,
		VaultName: vaultName,
	}

	if err = configCertificate.getKeyVaultClient(); err != nil {
		return err
	}

	for _, certificateName := range certificatesName {
		pfxCertificate, err := configCertificate.requestCertificatesPFX(certificateName)
		if err != nil {
			return err
		}

		err = loadCertificatesOnStorage(azureVaultEnv, certificateName, pfxCertificate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (akv *azureKeyVaultCertificate) getKeyVaultClient() (err error) {

	akv.Client = keyvault.New()
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	akv.Client.Authorizer = authorizer
	akv.authenticated = true

	akv.vaultBaseURL = fmt.Sprintf("https://%s.%s", akv.VaultName, azure.PublicCloud.KeyVaultDNSSuffix)

	return nil
}

func (akv *azureKeyVaultCertificate) requestCertificateVersion(certificateName string) (version string, err error) {

	list, err := akv.Client.GetCertificateVersionsComplete(akv.Ctx, akv.vaultBaseURL, certificateName, nil)
	if err != nil {
		return "", err
	}

	var lastItemDate time.Time
	var lastItemVersion string
	for list.NotDone() {

		item := list.Value()

		if *item.Attributes.Enabled {

			updatedTime := time.Time(*item.Attributes.Updated)
			if lastItemDate.IsZero() || updatedTime.After(lastItemDate) {
				lastItemDate = updatedTime

				parts := strings.Split(*item.ID, "/")
				lastItemVersion = parts[len(parts)-1]
			}
		}
		list.Next()
	}

	return lastItemVersion, nil
}

func (akv *azureKeyVaultCertificate) requestCertificatesPFX(certificateName string) ([]byte, error) {

	if !akv.authenticated {
		return nil, errors.New("Need to invoke GetKeyVaultClient() first")
	}

	certificateVersion, err := akv.requestCertificateVersion(certificateName)
	if err != nil {
		return nil, err
	}

	pfx, err := akv.Client.GetSecret(akv.Ctx, akv.vaultBaseURL, certificateName, certificateVersion)
	if err != nil {
		return nil, err
	}

	pfxBytes, err := base64.StdEncoding.DecodeString(*pfx.Value)
	if err != nil {
		return nil, err
	}

	return pfxBytes, nil
}
