#!/bin/bash

# Check for root access
if [ $(whoami) != "root" ]; then
	echo "Error: permission denied"
	exit 1
fi

# Copy program and config to /var/www/web-msg-handler
mkdir -p /var/www/web-msg-handler
cp web-msg-handler /var/www/web-msg-handler
cp examples/config.json /var/www/web-msg-handler/config.json.example
chown -R www-data:www-data /var/www/web-msg-handler

# Copy systemd unit
cp configs/systemd/web-msg-handler.service /lib/systemd/system
chown root:root /lib/systemd/system/web-msg-handler.service
chmod 0644 /lib/systemd/system/web-msg-handler.service
systemctl enable web-msg-handler.service

# Copy nginx config
cp config/nginx/web-msg-handler.conf /etc/nginx/sites
chown root:root /etc/nginx/sites/web-msg-handler.conf
chmod 0644 /etc/nginx/sites/web-msg-handler.conf

# Notify
echo "- Please review the configurations -"
echo "CONFIG:  /var/www/web-msg-handler/config.json.example"
echo "SYSTEMD: /lib/systemd/system/web-msg-handler.conf"
echo "NGINX:   /etc/nginx/sites/web-msg-handler.conf"

exit 0
