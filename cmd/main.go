package main

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/noorbala7418/trinity/internal/models"
	"github.com/noorbala7418/trinity/pkg/alert"
	"github.com/noorbala7418/trinity/pkg/email"
	"github.com/noorbala7418/trinity/pkg/proxmox"
	"github.com/sirupsen/logrus"
)

var alertmanagerAPI string
var targetIP string
var emailCreds models.EmailCredential
var proxmoxCreds models.ProxmoxCredential
var receiverMail string
var appMode string

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only logrus the warning severity or above.
	switch os.Getenv("APP_LOG_MODE") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	checkEnvs()
}

func main() {
	logrus.Info("Start Trinity")
	logrus.Info("App mode is ", appMode)
	sendHelloMail()
	logrus.Info("Check every 60s...")

	s, err := gocron.NewScheduler()
	if err != nil {
		logrus.Error("function main. Error in create scheduler. err: ", err)
	}

	// add a job to the scheduler
	job, err := s.NewJob(
		gocron.DurationJob(
			4*time.Minute,
		),
		gocron.NewTask(checkSystem),
	)
	if err != nil {
		logrus.Error("function main. Error in run cron task. err: ", err)
	}

	logrus.Info("function main. Cronjob scheduled. ID: ", job.ID())
	s.Start()

	// block until you are ready to shut down
	time.Sleep(time.Duration(math.MaxInt64))

}

// sendHelloMail Sends start email to inform service owner.
func sendHelloMail() {
	infoMail := models.Email{Sender: emailCreds.Username, Receiver: receiverMail, Subject: "Trinity Started", Body: "Trinity Started.\nNow we are observing.\n\nTime: " + time.Now().Local().String()}
	if mailErr := email.SendMail(infoMail, emailCreds); mailErr != nil {
		logrus.Error("function sendHelloMail. Send Email failed. err: ", mailErr)
	}
	logrus.Info("function sendHelloMail: Hello mail sent.")
}

// checkSystem will check alerts and send shutdown command to proxmox. This function is called using cronjob.
func checkSystem() {
	alerts, alertErr := alert.GrabAlerts(alertmanagerAPI)
	if alertErr != nil {
		logrus.Error("function checkSystem. Error in grab alerts. err: ", alertErr)
	}
	electricityCode, result := alert.CheckPingPacketLossAlert(targetIP, alerts)

	if electricityCode != 0 {
		logrus.Info("function checkSystem. Electricity Alert received. Code is: ", electricityCode, " Sending Email.")
		infoMail := models.Email{Sender: emailCreds.Username, Receiver: receiverMail, Subject: "Alert Detected!", Body: result}
		if mailErr := email.SendMail(infoMail, emailCreds); mailErr != nil {
			logrus.Error("function checkSystem. Send Email failed. err: ", mailErr)
		}

		if electricityCode == 2 {
			logrus.Info("function checkSystem. Alert Code is 2. State is CRITICAL. Going to SHUTDOWN the node.")

			if appMode == "action" {
				if shutErr := proxmox.ShutdownProxmox(proxmoxCreds); shutErr != nil {
					logrus.Error("function checkSystem. Shutdown Node Failed. err: ", shutErr)

					failEmail := models.Email{Sender: emailCreds.Username, Receiver: receiverMail, Subject: "Shutdown Failure", Body: "Node shutdown operation failed. Trintiy will try again. Error is: " + shutErr.Error()}

					if failMailErr := email.SendMail(failEmail, emailCreds); failMailErr != nil {
						logrus.Error("function checkSystem. Send failure email failed. err: ", failMailErr)
					}

					return
				}
			}

			logrus.Info("function checkSystem. Shutdown Node was success. Send Success Email.")

			successEmail := models.Email{
				Sender:   emailCreds.Username,
				Receiver: receiverMail,
				Subject:  "Shutdown Executed",
				Body:     "Node shutdown operation Succeded. Time: " + time.Now().Local().String(),
			}

			failMailErr := email.SendMail(successEmail, emailCreds)
			if failMailErr != nil {
				logrus.Error("function checkSystem. Success shutdown Email is sent. err: ", failMailErr)
			} else {
				logrus.Info("Job is done. Exit App.")
				os.Exit(0)
			}
		}
	}
}

// checkEnvs Checks environment variables and if one variable does not exist, Then it will Kill application.
func checkEnvs() {
	if os.Getenv("PROXMOX_API") == "" {
		logrus.Error("env variable $PROXMOX_API is not defined")
		os.Exit(1)
	}

	if os.Getenv("PROXMOX_TOKEN") == "" {
		logrus.Error("env variable $PROXMOX_TOKEN is not defined")
		os.Exit(1)
	}

	if os.Getenv("EMAIL_HOST") == "" {
		logrus.Error("env variable $EMAIL_HOST is not defined")
		os.Exit(1)
	}

	if os.Getenv("EMAIL_PORT") == "" {
		logrus.Error("env variable $EMAIL_PORT is not defined")
		os.Exit(1)
	}

	if os.Getenv("EMAIL_USERNAME") == "" {
		logrus.Error("env variable $EMAIL_USERNAME is not defined")
		os.Exit(1)
	}

	if os.Getenv("EMAIL_PASSWORD") == "" {
		logrus.Error("env variable $EMAIL_PASSWORD is not defined")
		os.Exit(1)
	}

	if os.Getenv("EMAIL_RECEIVER") == "" {
		logrus.Error("env variable $EMAIL_RECEIVER is not defined")
		os.Exit(1)
	}

	if os.Getenv("ALERTMANAGER_API") == "" {
		logrus.Error("env variable $ALERTMANAGER_API is not defined")
		os.Exit(1)
	}

	if os.Getenv("TARGET_IP") == "" {
		logrus.Error("env variable $TARGET_IP is not defined")
		os.Exit(1)
	}

	if os.Getenv("APP_MODE") == "" {
		logrus.Error("env variable $APP_MODE is not defined")
		os.Exit(1)
	}

	proxmoxCreds = models.ProxmoxCredential{
		API:   os.Getenv("PROXMOX_API"),
		Token: os.Getenv("PROXMOX_TOKEN"),
	}

	mailServerPort, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	emailCreds = models.EmailCredential{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     mailServerPort,
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}

	receiverMail = os.Getenv("EMAIL_RECEIVER")
	alertmanagerAPI = os.Getenv("ALERTMANAGER_API")
	targetIP = os.Getenv("TARGET_IP")
	appMode = os.Getenv("APP_MODE")
}
