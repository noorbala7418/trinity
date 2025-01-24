# Trinity

Trinity is an watcher. This app checks triggered alerts in `alertmanger`every 4 minutes. when `PingPacketLoss` alert showed up, then the proxmox node will be shutdown.

## Environment Variables

- `PROXMOX_API="https://192.168.1.3:8006/api2/json/nodes/proxmox-0/status"`
- `PROXMOX_TOKEN="PVEAPIToken=trinity@pve!power=XXXXXXXXXX"`
- `APP_LOG_MODE=debug # debug or info. default is info.`
- `APP_MODE=action # action or notification_only`
- `EMAIL_HOST="mail.domain.com"`
- `EMAIL_PORT=465`
- `EMAIL_USERNAME="trinity@domain.com"`
- `EMAIL_PASSWORD="PASSWORD"`
- `EMAIL_RECEIVER="me@domain.com"`
- `ALERTMANAGER_API="http://alertmanager.DOMAIN.com/api/v2/alerts"`
- `TARGET_IP="192.168.1.2`

> If you want to change Timezone, Please Set `TZ` variable based on your region. Like: `TZ=America/Los_Angeles`
