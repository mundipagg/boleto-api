package util

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mundipagg/boleto-api/certificate"

	s "github.com/fullsailor/pkcs7"
	"github.com/mundipagg/boleto-api/config"
)

var defaultDialer = &net.Dialer{Timeout: 16 * time.Second, KeepAlive: 16 * time.Second}

var (
	client            *http.Client
	onceDefaultClient = &sync.Once{}
	onceTransport     = &sync.Once{}
	icpCert           certificate.ICPCertificate
	transport         *http.Transport
)

// DefaultHTTPClient retorna um cliente http configurado para dar um skip na validação do certificado digital
func DefaultHTTPClient() *http.Client {
	onceDefaultClient.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				Dial:                defaultDialer.Dial,
				TLSHandshakeTimeout: 16 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	})
	return client
}

//Post faz um requisição POST para uma URL e retorna o response, status e erro
func Post(url, body, timeout string, header map[string]string) (string, int, error) {
	return doRequest("POST", url, body, timeout, header)
}

//Get faz um requisição GET para uma URL e retorna o response, status e erro
func Get(url, body, timeout string, header map[string]string) (string, int, error) {
	return doRequest("GET", url, body, timeout, header)
}

func doRequest(method, url, body, timeout string, header map[string]string) (string, int, error) {
	t := GetDurationTimeoutRequest(timeout) * time.Second

	ctx, cls := context.WithTimeout(context.Background(), t)
	defer cls()

	client := DefaultHTTPClient()

	message := strings.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, method, url, message)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return "", 0, errResp
	}
	defer resp.Body.Close()
	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		return "", resp.StatusCode, errResponse
	}
	sData := string(data)
	return sData, resp.StatusCode, nil
}

// BuildTLSTransport creates a TLS Client Transport from crt, ca and key files
func BuildTLSTransport() (*http.Transport, error) {
	var errF error
	onceTransport.Do(func() {

		ssl, err := certificate.GetCertificateFromStore(config.Get().CertificateSSLName)
		if err != nil {
			errF = err
			return
		}

		cert, err := tls.X509KeyPair(ssl.(certificate.SSLCertificate).PemData, ssl.(certificate.SSLCertificate).PemData)
		if err != nil {
			errF = err
			return
		}

		transport = &http.Transport{
			Dial:                defaultDialer.Dial,
			TLSHandshakeTimeout: 16 * time.Second,
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		}
		return
	})
	return transport, errF
}

//Sign request
func SignRequest(request string) (string, error) {

	if icpCert == (certificate.ICPCertificate{}) {
		icp, err := certificate.GetCertificateFromStore(config.Get().CertificateICPName)
		if err != nil {
			return "", err
		}
		icpCert = icp.(certificate.ICPCertificate)
	}

	signedData, err := s.NewSignedData([]byte(request))
	if err != nil {
		return "", err
	}

	if err := signedData.AddSigner(icpCert.Certificate, icpCert.RsaPrivateKey, s.SignerInfoConfig{}); err != nil {
		return "", err
	}

	detachedSignature, err := signedData.Finish()
	if err != nil {
		return "", err
	}

	signedRequest := base64.StdEncoding.EncodeToString(detachedSignature)

	return signedRequest, nil
}

//Read privatekey and parse to PKCS#1
func parsePrivateKey() (crypto.PrivateKey, error) {

	pkeyBytes, err := ioutil.ReadFile(config.Get().CertICP_PathPkey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pkeyBytes)
	if block == nil {
		return nil, errors.New("Key Not Found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return rsa, nil
	default:
		return nil, fmt.Errorf("SSH: Unsupported key type %q", block.Type)
	}

}

///Read chainCertificates and adapter to x509.Certificate
func parseChainCertificates() (*x509.Certificate, error) {

	chainCertsBytes, err := ioutil.ReadFile(config.Get().CertICP_PathChainCertificates)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(chainCertsBytes)
	if block == nil {
		return nil, errors.New("Key Not Found")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func doRequestTLS(method, url, body, timeout string, header map[string]string, transport *http.Transport) (string, int, error) {
	tlsClient := &http.Client{}
	tlsClient.Transport = transport
	tlsClient.Timeout = GetDurationTimeoutRequest(timeout) * time.Second
	b := strings.NewReader(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return "", 0, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, err := tlsClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	sData := string(data)
	return sData, resp.StatusCode, nil
}

func PostTLS(url, body, timeout string, header map[string]string, transport *http.Transport) (string, int, error) {
	return doRequestTLS("POST", url, body, timeout, header, transport)
}

//HeaderToMap converte um http Header para um dicionário string -> string
func HeaderToMap(h http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range h {
		m[k] = v[0]
	}
	return m
}
