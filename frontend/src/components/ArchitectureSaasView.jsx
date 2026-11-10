import React, { useState } from 'react';
import { apiRequest } from '../api';

export default function ArchitectureSaasView({
  companies,
  activeCompanyId,
  onSelectCompany,
  users,
  onRefresh,
  notify,
}) {
  const [compForm, setCompForm] = useState({
    nombre: '',
    ruc: '',
    plan: 'Enterprise',
    direccion: '',
    telefono: '',
  });

  const [userForm, setUserForm] = useState({
    id: '',
    nombre: '',
    email: '',
    password: '',
    rol: 'cajero',
  });

  async function handleCreateCompany(e) {
    e.preventDefault();
    try {
      const created = await apiRequest('empresas', {
        method: 'POST',
        body: compForm,
      });
      setCompForm({ nombre: '', ruc: '', plan: 'Enterprise', direccion: '', telefono: '' });
      await onRefresh();
      onSelectCompany(created.id);
      notify(`Nueva empresa SaaS "${created.nombre}" creada y activada con almacén y catálogo base.`);
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleSaveUser(e) {
    e.preventDefault();
    try {
      if (userForm.id) {
        await apiRequest(`usuarios/${userForm.id}`, { method: 'PUT', body: userForm });
        notify('Usuario y rol actualizados.');
      } else {
        await apiRequest('usuarios', { method: 'POST', body: userForm });
        notify('Usuario creado y vinculado a la empresa actual.');
      }
      setUserForm({ id: '', nombre: '', email: '', password: '', rol: 'cajero' });
      await onRefresh();
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleDeleteUser(id) {
    if (!window.confirm('¿Eliminar este usuario del sistema?')) return;
    try {
      await apiRequest(`usuarios/${id}`, { method: 'DELETE' });
      await onRefresh();
      notify('Usuario eliminado.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  const roleLabels = {
    admin: 'Administrador General',
    gerente: 'Gerente de Tienda',
    cajero: 'Cajero POS',
    almacenero: 'Jefe de Almacén',
    rrhh: 'Recursos Humanos',
  };

  return (
    <div>
      {/* GESTIÓN DE EMPRESAS TENANTS SAAS */}
      <div className="grid-equal">
        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Empresas Registradas (Tenants SaaS Multiempresa)</h2>
              <p>Aislamiento lógico de catálogo, almacenes, ventas y RRHH por <code>empresa_id</code></p>
            </div>
          </div>

          <form className="form-grid" onSubmit={handleCreateCompany}>
            <div className="field">
              <label>Razón Social / Marca *</label>
              <input
                required
                placeholder="Ej. Minimarket Surco Express"
                value={compForm.nombre}
                onChange={(e) => setCompForm({ ...compForm, nombre: e.target.value })}
              />
            </div>
            <div className="field">
              <label>RUC</label>
              <input
                placeholder="20601122334"
                value={compForm.ruc}
                onChange={(e) => setCompForm({ ...compForm, ruc: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Plan de Suscripción SaaS</label>
              <select
                value={compForm.plan}
                onChange={(e) => setCompForm({ ...compForm, plan: e.target.value })}
              >
                <option value="Starter">Plan Starter</option>
                <option value="Pro">Plan Pro</option>
                <option value="Enterprise">Plan Enterprise</option>
              </select>
            </div>
            <div className="field">
              <label>Dirección Fiscal</label>
              <input
                placeholder="Av. Caminos del Inca 450"
                value={compForm.direccion}
                onChange={(e) => setCompForm({ ...compForm, direccion: e.target.value })}
              />
            </div>
            <div>
              <button type="submit" className="btn btn-primary">
                + Crear Empresa SaaS
              </button>
            </div>
          </form>

          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>ID Tenant</th>
                  <th>Empresa</th>
                  <th>Plan</th>
                  <th>Estado</th>
                  <th>Acción</th>
                </tr>
              </thead>
              <tbody>
                {companies.map((c) => (
                  <tr key={c.id}>
                    <td><code>{c.id}</code></td>
                    <td>
                      <strong>{c.nombre}</strong>
                      <div style={{ fontSize: 11, color: '#68776e' }}>
                        RUC: {c.ruc || 'S/R'} · {c.direccion}
                      </div>
                    </td>
                    <td><span className="badge badge-purple">{c.plan}</span></td>
                    <td>
                      {c.id === activeCompanyId ? (
                        <span className="badge badge-green">● En uso activo</span>
                      ) : (
                        <span className="badge badge-blue">Disponible</span>
                      )}
                    </td>
                    <td>
                      {c.id !== activeCompanyId && (
                        <button
                          type="button"
                          className="btn btn-secondary btn-sm"
                          onClick={() => onSelectCompany(c.id)}
                        >
                          Cambiar a esta empresa
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* GESTIÓN DE USUARIOS Y ROLES */}
        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Usuarios, Autenticación y Roles</h2>
              <p>Control de accesos por rol dentro de la empresa seleccionada</p>
            </div>
          </div>

          <form className="form-grid" onSubmit={handleSaveUser}>
            <div className="field">
              <label>Nombre Completo *</label>
              <input
                required
                placeholder="Ej. Diego Torres"
                value={userForm.nombre}
                onChange={(e) => setUserForm({ ...userForm, nombre: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Correo Electrónico *</label>
              <input
                required
                type="email"
                placeholder="diego@minimarket.pe"
                value={userForm.email}
                onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Contraseña</label>
              <input
                type="password"
                placeholder={userForm.id ? 'Dejar vacío para mantener' : 'Por defecto: 123456'}
                value={userForm.password}
                onChange={(e) => setUserForm({ ...userForm, password: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Rol en el Sistema</label>
              <select
                value={userForm.rol}
                onChange={(e) => setUserForm({ ...userForm, rol: e.target.value })}
              >
                <option value="admin">Administrador (Control Total)</option>
                <option value="gerente">Gerente de Tienda</option>
                <option value="cajero">Cajero POS (Ventas)</option>
                <option value="almacenero">Almacenero (Inventario / Kardex)</option>
                <option value="rrhh">Especialista RRHH</option>
              </select>
            </div>
            <div>
              <button type="submit" className="btn btn-primary">
                {userForm.id ? 'Actualizar Usuario' : 'Registrar Usuario'}
              </button>
            </div>
          </form>

          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Usuario</th>
                  <th>Correo</th>
                  <th>Rol Asignado</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {users.map((u) => (
                  <tr key={u.id}>
                    <td><strong>{u.nombre}</strong></td>
                    <td><code>{u.email}</code></td>
                    <td>
                      <span className="badge badge-green">
                        {roleLabels[u.rol] || u.rol}
                      </span>
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: 6 }}>
                        <button
                          type="button"
                          className="btn btn-secondary btn-sm"
                          onClick={() =>
                            setUserForm({
                              id: u.id,
                              nombre: u.nombre,
                              email: u.email,
                              password: '',
                              rol: u.rol,
                            })
                          }
                        >
                          Editar
                        </button>
                        <button
                          type="button"
                          className="btn btn-danger btn-sm"
                          onClick={() => handleDeleteUser(u.id)}
                        >
                          Eliminar
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* ARQUITECTURA CEBOLLA, PATRONES Y AWS */}
      <div className="panel">
        <div className="panel-header">
          <div className="panel-title">
            <h2>Arquitectura Cebolla (Onion Architecture), Patrones de Diseño y Despliegue AWS EC2</h2>
            <p>Estructura técnica del backend en Go y su integración con el frontend en React</p>
          </div>
          <span className="badge badge-green">Go 1.22+ · React 18 · AWS EC2 Ready</span>
        </div>

        <div className="grid-equal">
          <div className="onion-diagram">
            <div className="onion-layer">
              <strong>1. Capa de Dominio (Núcleo Interno) — <code>internal/domain/entidades.go</code></strong>
              <p style={{ margin: '4px 0 0', fontSize: 12.5, color: '#33473c' }}>
                Entidades puras sin dependencias externas: <code>Company</code>, <code>Product</code>, <code>Category</code>,{' '}
                <code>Warehouse</code>, <code>InventoryMovement</code>, <code>Sale</code>, <code>User</code>, <code>Employee</code> y{' '}
                <code>Attendance</code>. Implementa el <strong>Patrón Prototype</strong> mediante <code>Clone(id)</code>.
              </p>
            </div>
            <div className="onion-layer">
              <strong>2. Capa de Aplicación (Casos de Uso) — <code>internal/application/servicio.go</code></strong>
              <p style={{ margin: '4px 0 0', fontSize: 12.5, color: '#33473c' }}>
                Define el puerto <code>application.Repository</code> y coordina la lógica de negocio: descuento automático de stock en Ventas POS, trazabilidad Kardex, cálculo de IGV (18%), onboarding multiempresa y consolidación del Dashboard.
              </p>
            </div>
            <div className="onion-layer">
              <strong>3. Capa de Infraestructura (Patrón Adapter) — <code>internal/infrastructure/</code></strong>
              <p style={{ margin: '4px 0 0', fontSize: 12.5, color: '#33473c' }}>
                <code>JSONRepository</code> y <code>MemoryAdapter</code> implementan la interfaz <code>application.Repository</code> adaptando el almacenamiento persistente en JSON (o memoria/SQL) sin acoplar las reglas de negocio.
              </p>
            </div>
            <div className="onion-layer">
              <strong>4. Capa de Presentación (API REST + React SPA) — <code>internal/presentation/</code></strong>
              <p style={{ margin: '4px 0 0', fontSize: 12.5, color: '#33473c' }}>
                Expone los endpoints RESTful bajo <code>/api/*</code> con aislamiento multiempresa por cabecera <code>X-Empresa-ID</code> y sirve el build de producción de React incrustado con <code>//go:embed all:web</code>.
              </p>
            </div>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            <div style={{ padding: 14, borderRadius: 12, border: '1px solid #e2e8e3', background: '#f8faf8' }}>
              <strong style={{ color: '#6b38d4' }}>Patrón Creacional: Prototype</strong>
              <p style={{ margin: '6px 0 0', fontSize: 12.5 }}>
                Permite clonar objetos complejos (<code>Product.Clone</code>, <code>Category.Clone</code>, <code>Employee.Clone</code>) asignándoles un nuevo identificador criptográfico en el servidor sin mutar la instancia original.
              </p>
            </div>
            <div style={{ padding: 14, borderRadius: 12, border: '1px solid #e2e8e3', background: '#f8faf8' }}>
              <strong style={{ color: '#1d5bd6' }}>Patrón Estructural: Adapter</strong>
              <p style={{ margin: '6px 0 0', fontSize: 12.5 }}>
                Desacopla los casos de uso de Go del motor de base de datos. <code>JSONRepository</code> traduce las llamadas de dominio a operaciones atómicas protegidas por <code>sync.Mutex</code> en <code>data.json</code>.
              </p>
            </div>
            <div style={{ padding: 14, borderRadius: 12, border: '1px solid #e2e8e3', background: '#f8faf8' }}>
              <strong style={{ color: '#245e3f' }}>Despliegue Continuo en AWS EC2 (Amazon Linux 2023 + Caddy HTTPS)</strong>
              <p style={{ margin: '6px 0 0', fontSize: 12.5 }}>
                El script <code>deploy/compilar-linux.ps1</code> compila el frontend React, ejecuta las pruebas automatizadas de Go y empaqueta un único binario Linux <code>amd64</code> en <code>dist/minimarket</code> listo para subir por SCP a EC2 con <code>systemd</code> y HTTPS automático vía Caddy.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
