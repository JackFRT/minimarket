# Subir Minimarket a AWS EC2

Resultado: la página queda en `http://<IP-pública>` funcionando las 24 horas, y los datos se guardan en el disco del servidor.

Costo: una instancia `t3.micro` cuesta unos 7–9 US$/mes (incluye la IP pública), y se descuenta de los créditos de AWS.

---

## 0. Antes de empezar (una sola vez)

1. **Alerta de presupuesto:** Facturación → *Presupuestos* → *Crear presupuesto* → plantilla *Presupuesto de gasto mensual* de **5 US$** con tu correo.
2. **Región:** arriba a la derecha elige una región, por ejemplo **EE. UU. Este (Norte de Virginia) `us-east-1`**, y usa siempre la misma.

## 1. Compilar en tu PC

En PowerShell, dentro de la carpeta del proyecto:

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\compilar-linux.ps1
```

Ejecuta las pruebas y crea la carpeta `dist\` con:
`minimarket` (programa para Linux), `minimarket.service`, `instalar.sh` y `data.json`.

## 2. Crear la instancia EC2

Consola → **EC2** → **Lanzar instancias**:

| Campo | Valor |
|---|---|
| Nombre | `minimarket` |
| Imagen (AMI) | **Amazon Linux 2023** (64 bits **x86**) |
| Tipo de instancia | **t3.micro** |
| Par de claves | *Crear un nuevo par de claves* → nombre `minimarket-key`, tipo RSA, formato **.pem**. Se descarga el archivo: **guárdalo**, porque no se puede volver a descargar. |
| Configuración de red → Firewall | *Crear grupo de seguridad* y marca: ✅ **Permitir tráfico SSH desde → Mi IP** · ✅ **Permitir tráfico HTTP desde Internet** |
| Almacenamiento | 8 GiB gp3 (el valor por defecto) |

Pulsa **Lanzar instancia**. Cuando el estado sea *En ejecución*, copia la **Dirección IPv4 pública**.

> Opcional pero recomendado: **Direcciones IP elásticas** → *Asignar* → *Asociar* a la instancia `minimarket`. Así la IP no cambia si apagas y enciendes el servidor.

## 3. Proteger la llave `.pem` (Windows)

Mueve el archivo a `C:\Users\<tu-usuario>\.ssh\minimarket-key.pem` y, en PowerShell:

```powershell
icacls "$env:USERPROFILE\.ssh\minimarket-key.pem" /inheritance:r /grant:r "$($env:USERNAME):(R)"
```

(Sin esto, `ssh` se niega a usar la llave porque "los permisos son demasiado abiertos".)

## 4. Subir los archivos

Reemplaza `IP` por la IP pública:

```powershell
scp -i "$env:USERPROFILE\.ssh\minimarket-key.pem" dist\* ec2-user@IP:~/
```

La primera vez pregunta si confías en el servidor: escribe `yes`.

## 5. Instalar en el servidor

```powershell
ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@IP "sudo bash ~/instalar.sh"
```

Debe terminar con `OK: Minimarket responde en el puerto 80.`

Abre en el navegador: **`http://IP`** 🎉

---

## HTTPS con dominio gratuito (ya configurado)

- Dominio: **https://minimarket-grupo2.duckdns.org** (DuckDNS → IP elástica `100.63.140.101`).
- **Caddy** atiende los puertos 80/443, obtiene y renueva solo el certificado de Let's Encrypt y redirige `http://` a `https://`. La app escucha solo en `127.0.0.1:8080`.
- Para instalarlo o reinstalarlo: sube `dist\*` y ejecuta `sudo bash ~/instalar-https.sh`. Para cambiar el dominio, edita `deploy\Caddyfile`.
- Si cambias la IP elástica, actualiza también la IP en duckdns.org.

## Actualizar después de hacer cambios

Repite los pasos **1, 4 y 5**. El instalador reemplaza el programa, pero **no borra** el `data.json` del servidor.

## Comandos útiles (dentro del servidor, con `ssh ... ec2-user@IP`)

| Para | Comando |
|---|---|
| Ver si está funcionando | `sudo systemctl status minimarket` |
| Ver los registros | `sudo journalctl -u minimarket -n 50` |
| Reiniciar | `sudo systemctl restart minimarket` |
| Descargar una copia de los datos (desde tu PC) | `ssh -i ...pem ec2-user@IP "sudo cat /var/lib/minimarket/data.json" > respaldo.json` |

## Para no gastar créditos

- **Detener** la instancia (EC2 → Instancias → *Estado de la instancia* → *Detener*) deja de cobrar el procesador; el disco y la IP elástica cobran centavos.
- **Terminar** la instancia la borra junto con los datos. Si asignaste una IP elástica, **libérala** también, porque una IP elástica sin usar sí cobra.

## ⚠️ Importante

- La aplicación **no tiene usuario ni contraseña**: cualquiera con la URL puede crear o borrar registros. Úsala para demostraciones y apágala cuando no la necesites.
- Funciona por `http://`, sin candado. Para tener `https://` hace falta un dominio, que se puede añadir después con Caddy o un balanceador de carga.
