# Minimarket Cloud SaaS — Sistema de Gestión Multiempresa (Go + React)

Sistema SaaS multiempresa para la gestión integral de minimarkets, desarrollado con **Go 1.22+** en el backend (aplicando **Arquitectura Cebolla** y los patrones de diseño **Adapter** y **Prototype**) y **React 18 (Vite)** en el frontend, empaquetado en un único ejecutable listo para desplegar en **AWS EC2**.

---

## Índice

1. [Módulos y Funcionalidades](#1-módulos-y-funcionalidades)
2. [Credenciales de Acceso Rápido (Demo SaaS)](#2-credenciales-de-acceso-rápido-demo-saas)
3. [Paso a Paso: Ejecutar en tu PC (Local)](#3-paso-a-paso-ejecutar-en-tu-pc-local)
4. [Paso a Paso: Desplegar en AWS EC2 desde Cero](#4-paso-a-paso-desplegar-en-aws-ec2-desde-cero)
5. [Paso a Paso: Configurar HTTPS Gratis (Dominio + Caddy)](#5-paso-a-paso-configurar-https-gratis-dominio--caddy)
6. [Cómo Actualizar el Servidor después de Hacer Cambios](#6-cómo-actualizar-el-servidor-después-de-hacer-cambios)
7. [Arquitectura Cebolla y Patrones de Diseño](#7-arquitectura-cebolla-y-patrones-de-diseño)
8. [Rutas de la API REST](#8-rutas-de-la-api-rest)

---

## 1. Módulos y Funcionalidades

- **Multiempresa (SaaS Multi-tenant):** Aislamiento de datos por empresa (`empresa_id` / cabecera `X-Empresa-ID`), selector de empresa activa en la barra lateral y creación (onboarding) de nuevas empresas con planes (`Starter`, `Pro`, `Enterprise`).
- **Usuarios, Autenticación y Roles:** Inicio de sesión y gestión de usuarios por empresa con roles (`admin`, `gerente`, `cajero`, `almacenero`, `rrhh`).
- **Productos y Categorías (Patrón Prototype):** Catálogo con SKU, costo de compra, precio de venta, stock actual, stock mínimo, unidad, categoría, almacén y proveedor. Incluye botón **`⧉ Clonar (Prototype)`** para duplicar productos y categorías con nueva identidad.
- **Control de Inventario (Kardex), Almacenes y Proveedores:** Registro auditable de movimientos (`entrada`, `salida`, `ajuste`) que actualiza automáticamente el stock del producto, gestión de múltiples almacenes/sucursales y directorio de proveedores.
- **Punto de Venta (POS) e Historial de Ventas:** Terminal POS interactiva con buscador, filtro por categoría, carrito de compras con cálculo automático de Subtotal e IGV (18%), métodos de pago (`efectivo`, `yape_plin`, `tarjeta`, `transferencia`), emisión de comprobante (`Boleta`, `Factura`, `Ticket`) y descuento automático de stock.
- **Recursos Humanos (RRHH) y Asistencia:** Directorio de trabajadores (DNI, cargo, departamento, almacén asignado, turno horario, salario mensual, estado laboral), botón **`⧉ Clonar Puesto`** (Patrón Prototype) y control diario de asistencia (`puntual`, `tardanza`, `justificado`, `falta`).
- **Dashboard Ejecutivo y Reportes:** KPIs de ingresos, ticket promedio, valorización de inventario, alertas de stock bajo, planilla mensual, gráficos de ventas por categoría, recaudación por método de pago y top productos más vendidos.

---

## 2. Credenciales de Acceso Rápido (Demo SaaS)

Puedes cambiar de usuario o empresa con un solo clic desde el botón **🔐 ADMIN / Cambiar** en la aplicación, o ingresar manualmente con:

| Empresa (Tenant SaaS) | Rol | Correo | Contraseña |
|---|---|---|---|
| **Minimarket El Sol - Sede Central** (`emp-default`) | Administrador | `admin@minimarket.pe` | `admin123` |
| **Minimarket El Sol - Sede Central** (`emp-default`) | Cajera POS | `caja@minimarket.pe` | `caja123` |
| **Minimarket El Sol - Sede Central** (`emp-default`) | Jefa de RRHH | `rrhh@minimarket.pe` | `rrhh123` |
| **Market Express Los Olivos** (`emp-norte`) | Administrador | `admin@marketnorte.pe` | `admin123` |

---

## 3. Paso a Paso: Ejecutar en tu PC (Local)

### Opción 1: Ejecutar todo con un solo comando (Recomendado)
Abre **PowerShell** en la carpeta del proyecto (`c:\Proyectos EN\minimarket-visual\minimarket`) y ejecuta:

```powershell
./iniciar.ps1
```
*(O con Go instalado: `go run ./cmd/api`)*.

Luego abre en tu navegador:
👉 **http://127.0.0.1:8080**

Detén el servidor cuando quieras con `Ctrl + C`. Todos los datos se guardan automáticamente en `data.json`.

### Opción 2: Si modificas el código del Frontend React (`frontend/`)
Si editas archivos dentro de `frontend/src/`:

```powershell
# 1. Instalar dependencias de React (solo la primera vez)
npm --prefix frontend install

# 2. Compilar el frontend React dentro de internal/presentation/web
npm --prefix frontend run build

# 3. Iniciar el servidor Go con los cambios incrustados
./iniciar.ps1
```

### Ejecutar las pruebas automatizadas
```powershell
go test -v ./...
```

---

## 4. Paso a Paso: Desplegar en AWS EC2 desde Cero

### Paso 1. Compilar el paquete para Linux en tu PC
En PowerShell, dentro de la carpeta del proyecto, ejecuta:

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\compilar-linux.ps1
```
Esto compila el frontend React, ejecuta las pruebas de Go y genera la carpeta `dist\` con todo listo para subir a Linux (`minimarket`, `data.json`, `instalar.sh`, `instalar-https.sh`, `minimarket.service`, `caddy.service` y `Caddyfile`).

---

### Paso 2. Lanzar la instancia en AWS EC2
1. Entra a la consola de AWS → **EC2** → **Lanzar instancias**.
2. Completa los campos:
   - **Nombre:** `minimarket`
   - **Imagen (AMI):** **Amazon Linux 2023** (64 bits **x86**)
   - **Tipo de instancia:** **t3.micro**
   - **Par de claves (inicio de sesión):**
     - Selecciona *Crear un nuevo par de claves* → nombre: `minimarket-key`, tipo: **RSA**, formato: **`.pem`**.
     - Se descargará `minimarket-key.pem`. Muévelo a la carpeta `C:\Users\<tu-usuario>\.ssh\minimarket-key.pem`.
   - **Configuración de red → Firewall (Grupo de seguridad):**
     - Marca ✅ **Permitir el tráfico de SSH**
     - Marca ✅ **Permitir el tráfico de HTTPS desde Internet**
     - Marca ✅ **Permitir el tráfico de HTTP desde Internet**
   - **Almacenamiento:** `8 GiB gp3` (valor por defecto).
3. Pulsa el botón naranja **Lanzar instancia**.
4. Ve a **EC2 → Instancias**, espera a que el estado sea **En ejecución** (*Running*) y copia la **Dirección IPv4 pública** (por ejemplo `54.123.45.67`).

---

### Paso 3. Ajustar permisos de la llave `.pem` en Windows (Solo la primera vez)
Asegúrate de que exista la carpeta `.ssh` y restringe los permisos de la llave para que `ssh` y `scp` la acepten:

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.ssh" | Out-Null
icacls "$env:USERPROFILE\.ssh\minimarket-key.pem" /inheritance:r /grant:r "$($env:USERNAME):(R)"
```

---

### Paso 4. Subir los archivos al servidor EC2
En PowerShell (dentro de la carpeta del proyecto), reemplaza `TU_IP` por la **IPv4 pública** de tu instancia:

```powershell
scp -i "$env:USERPROFILE\.ssh\minimarket-key.pem" dist\* ec2-user@TU_IP:~/
```
*(La primera vez que te conectes te preguntará `Are you sure you want to continue connecting (yes/no)?`: escribe `yes` y presiona Enter).*

---

### Paso 5. Instalar y arrancar el sistema en el servidor
Ejecuta en PowerShell (reemplazando `TU_IP` por tu IPv4 pública):

```powershell
ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo bash ~/instalar.sh"
```

Cuando termine de ejecutar verás el mensaje de confirmación. Ahora abre en cualquier navegador:
👉 **`http://TU_IP`** 🎉

---

## 5. Paso a Paso: Configurar HTTPS Gratis (Dominio + Caddy)

Si deseas que el sistema tenga candado de seguridad **`https://`** con certificado automático de Let's Encrypt:

1. Entra a [https://www.duckdns.org](https://www.duckdns.org) y apunta tu subdominio (por ejemplo `minimarket-grupo2.duckdns.org`) a la **IPv4 pública** de tu instancia EC2.
2. Si usas un nombre de dominio distinto, edita la primera línea de `deploy\Caddyfile` y vuelve a ejecutar `.\deploy\compilar-linux.ps1` y el `scp` del Paso 4.
3. Ejecuta en PowerShell:
   ```powershell
   ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo bash ~/instalar-https.sh"
   ```
4. Abre en tu navegador:
   👉 **`https://minimarket-grupo2.duckdns.org`**

---

## 6. Cómo Actualizar el Servidor después de Hacer Cambios

Cada vez que hagas cambios en el código (Go o React) y quieras subirlos a tu EC2:

```powershell
# 1. Recompilar todo el paquete para Linux
powershell -ExecutionPolicy Bypass -File .\deploy\compilar-linux.ps1

# 2. Subir los nuevos archivos al servidor
scp -i "$env:USERPROFILE\.ssh\minimarket-key.pem" dist\* ec2-user@TU_IP:~/

# 3. Reiniciar el servicio (conserva los datos ya guardados en el servidor)
ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo bash ~/instalar.sh"
```

### Comandos útiles de administración remota

| Qué necesitas hacer | Comando en PowerShell |
|---|---|
| **Ver si el servicio está activo** | `ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo systemctl status minimarket"` |
| **Ver los últimos registros (logs)** | `ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo journalctl -u minimarket -n 50 --no-pager"` |
| **Reiniciar el servidor** | `ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo systemctl restart minimarket"` |
| **Reemplazar la base de datos del servidor con el `data.json` local** | `ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo install -m 600 -o minimarket -g minimarket ~/data.json /var/lib/minimarket/data.json && sudo systemctl restart minimarket"` |
| **Descargar una copia de seguridad (`respaldo.json`) a tu PC** | `ssh -i "$env:USERPROFILE\.ssh\minimarket-key.pem" ec2-user@TU_IP "sudo cat /var/lib/minimarket/data.json" > respaldo.json` |

---

## 7. Arquitectura Cebolla y Patrones de Diseño

```text
cmd/api/principal.go                           Composición de dependencias y arranque HTTP
internal/domain/entidades.go                   Núcleo de dominio, validaciones y Patrón Prototype
internal/application/servicio.go               Casos de uso de negocio y puerto Repository
internal/infrastructure/repositorio_json.go    Patrón Adapter para almacenamiento JSON (data.json)
internal/infrastructure/repositorio_memoria.go Patrón Adapter en memoria efímera
internal/infrastructure/semilla.go             Datos iniciales multiempresa de demostración
internal/presentation/interfaz_http.go         API RESTful + servidor embebido React SPA
frontend/                                      Código fuente de la interfaz React (Vite)
deploy/                                        Scripts de compilación Linux y despliegue AWS EC2
```

- **Patrón Prototype (`internal/domain/entidades.go`):** `Product.Clone`, `Category.Clone` y `Employee.Clone` permiten clonar entidades complejas asignándoles un nuevo identificador sin alterar el registro original.
- **Patrón Adapter (`internal/infrastructure/`):** `JSONRepository` y `MemoryAdapter` adaptan la interfaz `application.Repository` al sistema de archivos JSON o memoria sin acoplar la lógica de negocio.

---

## 8. Rutas de la API REST

Todas las rutas aceptan la cabecera `X-Empresa-ID` (o el parámetro `?empresa_id=`) para filtrar y operar sobre la empresa activa:

| Método | Ruta | Acción |
|---|---|---|
| `GET` / `POST` | `/api/empresas` | Listar / crear empresa SaaS (Tenant) |
| `PUT` | `/api/empresas/{id}` | Actualizar empresa |
| `POST` | `/api/auth/login` | Iniciar sesión (`email`, `password`) |
| `GET` / `POST` | `/api/usuarios` | Listar / crear usuarios por empresa |
| `PUT` / `DELETE` | `/api/usuarios/{id}` | Actualizar / eliminar usuario |
| `GET` | `/api/dashboard` | Obtener métricas, gráficos y alertas de la empresa |
| `GET` / `POST` | `/api/categorias` | Listar / crear categorías de productos |
| `PUT` / `DELETE` | `/api/categorias/{id}` | Actualizar / eliminar categoría |
| `POST` | `/api/categorias/{id}/clonar` | Clonar categoría con **Prototype** |
| `GET` / `POST` | `/api/productos` | Listar / crear productos |
| `GET` / `PUT` / `DELETE` | `/api/productos/{id}` | Consultar / actualizar / eliminar producto |
| `POST` | `/api/productos/{id}/clonar` | Clonar producto con **Prototype** |
| `GET` / `POST` | `/api/almacenes` | Listar / crear almacenes |
| `PUT` / `DELETE` | `/api/almacenes/{id}` | Actualizar / eliminar almacén |
| `GET` / `POST` | `/api/proveedores` | Listar / crear proveedores |
| `GET` / `PUT` / `DELETE` | `/api/proveedores/{id}` | Consultar / actualizar / eliminar proveedor |
| `GET` / `POST` | `/api/inventario/movimientos` | Listar Kardex / registrar entrada, salida o ajuste |
| `GET` / `POST` | `/api/ventas` | Listar ventas / emitir venta POS descontando stock |
| `GET` / `POST` | `/api/rrhh/empleados` | Listar / registrar trabajadores en RRHH |
| `PUT` / `DELETE` | `/api/rrhh/empleados/{id}` | Actualizar / retirar trabajador |
| `POST` | `/api/rrhh/empleados/{id}/clonar` | Clonar puesto de trabajador con **Prototype** |
| `GET` / `POST` | `/api/rrhh/asistencias` | Listar / registrar asistencia diaria del personal |
