package appstore

import (
	"crypto/x509"
	"embed"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

//go:embed certs/*.cer
var certs embed.FS

const srcUrl = "https://www.apple.com/certificateauthority/"
const outDir = "certs/"

var certLinkPattern = regexp.MustCompile(`<a [^>]*href="([^"]+\.cer)"`)

type CertPool struct {
	pool     *x509.CertPool
	poolOnce sync.Once
}

func NewCertPool() (*CertPool, error) {
	cp := &CertPool{}
	err := cp.Init()
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func (cp *CertPool) Init() error {
	var err error
	cp.poolOnce.Do(func() {
		cp.pool = x509.NewCertPool()
		err = cp.downloadCerts()
		err = cp.loadCerts()
	})
	return err
}

func (cp *CertPool) downloadCerts() error {
	resp, err := http.Get(srcUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(outDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	matches := certLinkPattern.FindAllSubmatch(content, -1)
	for _, match := range matches {
		certUrl, err := cp.constructCertUrl(string(match[1]))
		if err != nil {
			return err
		}

		if err := cp.downloadAndSaveCert(certUrl); err != nil {
			return err
		}
	}
	return nil
}

func (cp *CertPool) constructCertUrl(certPath string) (string, error) {
	if certPath[0] == '/' {
		baseUrl, err := url.Parse(srcUrl)
		if err != nil {
			return "", err
		}
		baseUrl.Path = certPath
		return baseUrl.String(), nil
	} else if strings.HasPrefix(certPath, "https://www.apple.com/") || strings.HasPrefix(certPath, "https://developer.apple.com/") {
		return certPath, nil
	} else {
		return url.JoinPath(srcUrl, certPath)
	}
}

func (cp *CertPool) downloadAndSaveCert(certUrl string) error {
	resp, err := http.Get(certUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	fileName := path.Base(certUrl)
	filePath := filepath.Join(outDir, fileName)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func (cp *CertPool) loadCerts() error {
	entries, err := certs.ReadDir("certs")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() && entry.Type().IsRegular() {
			cert, err := certs.ReadFile("certs/" + entry.Name())
			if err != nil {
				continue
			}
			if ok := cp.pool.AppendCertsFromPEM(cert); ok {
				continue
			}
			if cer, err := x509.ParseCertificate(cert); err == nil {
				cp.pool.AddCert(cer)
			}
		}
	}
	return nil
}

func (cp *CertPool) GetCertPool() *x509.CertPool {
	return cp.pool
}
