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
func IssueCertificate(caName string, certName string, commonName string, basePath string) error {
	caKey := fmt.Sprintf("%s/ca.key", basePath)
	caCert := fmt.Sprintf("%s/ca.crt", basePath)
	
	keyFile := fmt.Sprintf("%s/%s.key", basePath, certName)
	csrFile := fmt.Sprintf("%s/%s.csr", basePath, certName)
	crtFile := fmt.Sprintf("%s/%s.crt", basePath, certName)

	// Generate Key
	exec.Command("openssl", "genrsa", "-out", keyFile, "2048").Run()

	// Generate CSR
	subj := fmt.Sprintf("/CN=%s", commonName)
	exec.Command("openssl", "req", "-new", "-key", keyFile, "-out", csrFile, "-subj", subj).Run()

	// Sign Cert
	cmd := exec.Command("openssl", "x509", "-req", "-in", csrFile, "-CA", caCert, "-CAkey", caKey, 
		"-CAcreateserial", "-out", crtFile, "-days", "365")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("signing failed: %s", string(out))
	}

	return nil
}
