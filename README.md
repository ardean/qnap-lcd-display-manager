# QNAP LCD display manager

## Build
```sh
go build main.go
```
## Build with Docker
```sh
sh build-with-docker.sh
```
## Install
Copy qnap-lcd-display-manager binary and the content of ``deployment`` to your QNAP NAS (all in one folder) and install it with
```sh
sudo sh install.sh
```

Then enable the service with
```sh
sudo systemctl enable qnap-lcd-display-manager.service --now
```
