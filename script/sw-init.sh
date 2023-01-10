#!/usr/bin/env bash

set -e

pacman -Syy
pacman -S --noconfirm archlinux-keyring
pacman --noconfirm -Syyu

pacman -S --noconfirm vim wget openresolv wireguard-tools wgcf

wgcf register
wgcf generate

sed -i 's|engage.cloudflareclient.com:2408|[2606:4700:d0::a29f:c001]:2408|' wgcf-profile.conf
sed -i 's|AllowedIPs = ::/0||' wgcf-profile.conf
sed -i 's|1.1.1.1|2001:4860:4860::8888,2001:4860:4860::8844,8.8.8.8,8.8.4.4|' wgcf-profile.conf

cp wgcf-profile.conf /etc/wireguard/wgcf.conf

reboot

# wg-quick up wgcf
# wg-quick down wgcf

# systemctl enable --now wg-quick@wgcf
