#!/bin/bash

INSTALLATION_PATH="/opt/web-msg-handler"
SETTINGS_PATH="/etc/opt/web-msg-handler"
PLUGINS_PATH="$SETTINGS_PATH/plugins"
SITES_PATH="$SETTINGS_PATH/sites"
SYSTEMD_SERVICE_PATH="/lib/systemd/system/web-msg-handler.service"
NGINX_SITE_PATH="/etc/nginx/sites/web-msg-handler.conf"

# Check for root access
if [ $(whoami) != "root" ]; then
	echo "Error: permission denied"
	exit 1
fi

# Change owner of all files to root
chown -R root:root *

# Copy program to installation path
mkdir -p $INSTALLATION_PATH
cp web-msg-handler $INSTALLATION_PATH

# Copy configs and plugins
mkdir -p $PLUGINS_PATH $SITES_PATH
cp examples/config.toml $SETTINGS_PATH/config.toml.example
cp examples/sites/mail.toml $SITES_PATH/mail.toml.example
cp examples/sites/telegram.toml $SITES_PATH/telegram.toml.example
cp plugins/* $PLUGINS_PATH

# Copy systemd unit
cp configs/systemd/web-msg-handler.service $SYSTEMD_SERVICE_PATH
chmod 0644 $SYSTEMD_SERVICE_PATH
systemctl enable web-msg-handler.service

# Copy nginx config
cp config/nginx/web-msg-handler.conf $NGINX_SITE_PATH
chmod 0644 $NGINX_SITE_PATH

# Notify
echo "- Please review the configurations -"
echo "CONFIG:  $SETTINGS_PATH"
echo "SYSTEMD: $SYSTEMD_SERVICE_PATH"
echo "NGINX:   $NGINX_SITE_PATH"

exit 0
