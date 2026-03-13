package vpn

import (
	"fmt"
	"os"
	"os/exec"
)

// CertificateAuthority represents a local Root CA
type CertificateAuthority struct {
	Name         string `json:"name"`
	Organization string `json:"org"`
	CommonName   string `json:"cn"`
	KeyLength    int    `json:"key_len"`
	LifetimeDays int    `json:"lifetime"`
}

// GenerateCA initializes a new CA using openssl
func (ca CertificateAuthority) GenerateCA(basePath string) error {
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return err
	}

	keyFile := fmt.Sprintf("%s/ca.key", basePath)
	certFile := fmt.Sprintf("%s/ca.crt", basePath)

	// 1. Generate Private Key
	cmd := exec.Command("openssl", "genrsa", "-out", keyFile, fmt.Sprintf("%d", ca.KeyLength))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ca key gen failed: %s", string(out))
	}

	// 2. Generate Self-Signed Cert
	subj := fmt.Sprintf("/O=%s/CN=%s", ca.Organization, ca.CommonName)
	cmd = exec.Command("openssl", "req", "-new", "-x509", "-days", fmt.Sprintf("%d", ca.LifetimeDays), 
		"-key", keyFile, "-out", certFile, "-subj", subj)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ca cert gen failed: %s", string(out))
	}

	return nil
}

// IssueCertificate signs a new client/server cert
func IssueCertificate(caName string, certName string, commonName string, basePath string, keyLen int, days int) error {
	caKey := fmt.Sprintf("%s/ca.key", basePath)
	caCert := fmt.Sprintf("%s/ca.crt", basePath)
	
	keyFile := fmt.Sprintf("%s/%s.key", basePath, certName)
	csrFile := fmt.Sprintf("%s/%s.csr", basePath, certName)
	crtFile := fmt.Sprintf("%s/%s.crt", basePath, certName)

	if keyLen == 0 {
		keyLen = 2048
	}
	if days == 0 {
		days = 365
	}

	// 1. Generate Key
	cmd := exec.Command("openssl", "genrsa", "-out", keyFile, fmt.Sprintf("%d", keyLen))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("key generation failed: %s", string(out))
	}

	// 2. Generate CSR
	subj := fmt.Sprintf("/CN=%s", commonName)
	cmd = exec.Command("openssl", "req", "-new", "-key", keyFile, "-out", csrFile, "-subj", subj)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("csr generation failed: %s", string(out))
	}

	// 3. Sign Cert
	cmd = exec.Command("openssl", "x509", "-req", "-in", csrFile, "-CA", caCert, "-CAkey", caKey, 
		"-CAcreateserial", "-out", crtFile, "-days", fmt.Sprintf("%d", days))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("signing failed: %s", string(out))
	}

	return nil
}
