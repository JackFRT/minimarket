#!/usr/bin/env bash
# Instala o actualiza Minimarket en EC2 (Amazon Linux 2023 / Ubuntu).
# Uso, dentro del servidor y en la carpeta con los archivos subidos:  sudo bash instalar.sh
set -euo pipefail
cd "$(dirname "$0")"

if [ "$(id -u)" -ne 0 ]; then echo "Ejecuta con sudo: sudo bash instalar.sh"; exit 1; fi

id minimarket >/dev/null 2>&1 || useradd --system --no-create-home --shell /sbin/nologin minimarket

install -d -m 755 /opt/minimarket
install -d -m 750 -o minimarket -g minimarket /var/lib/minimarket

systemctl stop minimarket 2>/dev/null || true
install -m 755 minimarket /opt/minimarket/minimarket

# Solo copia data.json la primera vez, para no borrar lo que ya se guardó en el servidor
if [ -f data.json ] && [ ! -f /var/lib/minimarket/data.json ]; then
  install -m 600 -o minimarket -g minimarket data.json /var/lib/minimarket/data.json
fi

install -m 644 minimarket.service /etc/systemd/system/minimarket.service
systemctl daemon-reload
systemctl enable --now minimarket
sleep 1
systemctl --no-pager status minimarket | head -n 5

if curl -fsS -o /dev/null http://127.0.0.1:8080/api/productos; then
  echo "OK: Minimarket responde (puerto interno 8080; Caddy lo publica por HTTPS)."
else
  echo "Aviso: el servicio no respondió. Revisa: sudo journalctl -u minimarket -n 50"
fi
