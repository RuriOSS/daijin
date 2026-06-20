#!/bin/sh
# Linux-specific container fixup script
# This runs inside the container on first boot

# Fix PATH
PATH=$PATH:/bin:/sbin:/usr/bin:/usr/sbin

# Fix DNS
rm -f /etc/resolv.conf >/dev/null 2>&1
echo "nameserver 8.8.8.8" >>/etc/resolv.conf
echo "nameserver 1.1.1.1" >>/etc/resolv.conf

# Create mountpoints if needed
[ -e /dev ] || mkdir /dev
[ -e /proc ] || mkdir /proc
[ -e /sys ] || mkdir /sys

echo "Linux container initialized"
