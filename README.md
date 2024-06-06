# QNAP LCD display manager

### Features
- Display local IPs
- Standby after 10 seconds
- Scroll through the local IPs with the up/down buttons

PRs to add more features are always appreciated!

### Build
```sh
go build main.go
```
### Build with Docker
```sh
sh build-with-docker.sh
```
### Install
Copy ``qnap-lcd-display-manager`` binary that you built and the content of ``deployment`` to your QNAP NAS (all in one folder).
Then install the service with:
```sh
sudo sh install.sh
```

Then enable the service:
```sh
sudo systemctl enable qnap-lcd-display-manager.service --now
```
