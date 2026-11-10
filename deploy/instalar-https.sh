#!/usr/bin/env bash
# Instala Caddy (HTTPS gratuito con Let's Encrypt) delante de Minimarket.
# Uso, dentro del servidor:  sudo bash instalar-https.sh
set -euo pipefail
cd "$(dirname "$0")"
if [ "$(id -u)" -ne 0 ]; then echo "Ejecuta con sudo: sudo bash instalar-https.sh"; exit 1; fi

CADDY_VERSION=2.10.2
if ! /usr/local/bin/caddy version 2>/dev/null | grep -q "v${CADDY_VERSION}"; then
  tmp=$(mktemp -d)
  base="https://github.com/caddyserver/caddy/releases/download/v${CADDY_VERSION}"
  tarball="caddy_${CADDY_VERSION}_linux_amd64.tar.gz"
  curl -fsSL -o "$tmp/$tarball" "$base/$tarball"
  curl -fsSL -o "$tmp/checksums.txt" "$base/caddy_${CADDY_VERSION}_checksums.txt"
  (cd "$tmp" && grep " ${tarball}\$" checksums.txt | sha512sum -c -)
  tar -xzf "$tmp/$tarball" -C "$tmp" caddy
  install -m 755 "$tmp/caddy" /usr/local/bin/caddy
  rm -rf "$tmp"
fi

id caddy >/dev/null 2>&1 || useradd --system --home-dir /var/lib/caddy --create-home --shell /sbin/nologin caddy
install -d -m 755 /etc/caddy
install -m 644 Caddyfile /etc/caddy/Caddyfile
install -m 644 caddy.service /etc/systemd/system/caddy.service

# La app pasa a escuchar solo en 127.0.0.1:8080; Caddy ocupa los puertos 80 y 443
install -m 644 minimarket.service /etc/systemd/system/minimarket.service
systemctl daemon-reload
systemctl restart minimarket
systemctl enable caddy
systemctl restart caddy
sleep 2
systemctl --no-pager status caddy | head -n 5
