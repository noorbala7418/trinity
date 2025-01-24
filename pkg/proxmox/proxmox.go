package proxmox

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/noorbala7418/trinity/internal/models"
	"github.com/sirupsen/logrus"
)

// ShutdownProxmox will send The SHUTDOWN command to proxmox host.
func ShutdownProxmox(connection models.ProxmoxCredential) error {
	method := "POST"

	payload := strings.NewReader("command=shutdown")

	// Start checking site
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr,
	}

	req, err := http.NewRequest(method, connection.API, payload)

	if err != nil {
		logrus.Error("function ShutdownProxmox. Error in create http request. err: ", err)
		return fmt.Errorf("function ShutdownProxmox. Error in create http request. err: %s", err)
	}
	req.Header.Add("Authorization", connection.Token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		logrus.Error("function ShutdownProxmox. Error in send http request. err: ", err)
		return fmt.Errorf("function ShutdownProxmox. Error in send http request. err: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.Error("function ShutdownProxmox. Proxmox node does not reply correctly. Response Status is: ", res.StatusCode)
		return fmt.Errorf("function ShutdownProxmox. Proxmox node does not reply correctly. Response Status is: %d", res.StatusCode)
	}

	return nil
}
