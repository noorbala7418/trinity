package alert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/noorbala7418/trinity/internal/models"
	"github.com/sirupsen/logrus"
)

// GrabAlerts fetchs alerts from alertmanager api.
func GrabAlerts(alertmanagerAPI string) ([]models.Alert, error) {
	resp, err := http.Get(alertmanagerAPI)
	if err != nil {
		logrus.Error("function GetCurrentDay. Error in send GET request to calendar. err: ", err)
	}

	if resp != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Error("function GrabAlerts. Error in get body from alertmanager. err: ", err)
			return nil, fmt.Errorf("function GrabAlerts. Error in get body from alertmanager. err: %s", err)
		}

		var response []models.Alert
		if responseErr := json.Unmarshal(body, &response); responseErr != nil {
			logrus.Error("function GrabAlerts: Error in parsing json data. err: ", responseErr)
			return nil, fmt.Errorf("function GrabAlerts: Error in parsing json data. err: %s", responseErr)
		}
		resp.Body.Close()
		return response, nil
	}
	return nil, fmt.Errorf("function: GrabAlerts. No response to paging")
}

// CheckPingPacketLossAlert searchs ppl alert between received alerts. If critical alert appeared, then returns 2. warning is 1 and if there is no related alerts, returns 0.
func CheckPingPacketLossAlert(targetIP string, alerts []models.Alert) (int, string) {
	for i := len(alerts); i > 0; {
		i-- // instead of --i
		if alerts[i].Labels["service"] == "kwc-ping-exporter" &&
			alerts[i].Labels["type"] == "pingpacketlost" &&
			alerts[i].Labels["target"] == targetIP {

			if alerts[i].Labels["severity"] == "critical" {
				logrus.Info("function: CheckPingPacketLossAlert. Electricity Disconnected. Alert Verified. Alert status: ", alerts[i].Status.State)
				return 2, "Electricity Disconnected. Alert Verified."
			} else if alerts[i].Labels["severity"] == "warning" {
				logrus.Info("function: CheckPingPacketLossAlert. Electricity seems Disconnected. Still Waiting")
				return 1, "Electricity seems Disconnected. Still Waiting."
			}
		}
	}
	logrus.Debug("function: CheckPingPacketLossAlert. Electricity is OK. Sleeping")
	return 0, ""
}
