import React, { useState } from 'react';
import { apiRequest } from '../api';

export default function InventoryView({
  products,
  warehouses,
  suppliers,
  movements,
  currentUser,
  onRefresh,
  notify,
}) {
  const [movForm, setMovForm] = useState({
    producto_id: '',
    almacen_id: '',
    tipo: 'entrada',
    cantidad: '',
    motivo: '',
  });

  const [whForm, setWhForm] = useState({
    id: '',
    codigo: '',
    nombre: '',
    direccion: '',
    responsable: '',
    capacidad: 2000,
    principal: false,
  });

  const [supForm, setSupForm] = useState({
    id: '',
    ruc: '',
    nombre: '',
    telefono: '',
    email: '',
    direccion: '',
    rubro: '',
  });

  async function handleSaveMovement(e) {
    e.preventDefault();
    try {
      await apiRequest('inventario/movimientos', {
        method: 'POST',
        body: {
          producto_id: movForm.producto_id,
          almacen_id: movForm.almacen_id,
          tipo: movForm.tipo,
          cantidad: Number(movForm.cantidad),
          motivo: movForm.motivo || 'Movimiento de almacén',
          usuario_id: currentUser?.id || '',
          usuario_nombre: currentUser?.nombre || 'Jefe de Almacén',
        },
      });
      setMovForm({ producto_id: '', almacen_id: '', tipo: 'entrada', cantidad: '', motivo: '' });
      await onRefresh();
      notify('Movimiento registrado en Kardex y stock actualizado.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleSaveWarehouse(e) {
    e.preventDefault();
    try {
      const payload = {
        codigo: whForm.codigo,
        nombre: whForm.nombre,
        direccion: whForm.direccion,
        responsable: whForm.responsable,
        capacidad: Number(whForm.capacidad || 1000),
        principal: Boolean(whForm.principal),
        activo: true,
      };
      if (whForm.id) {
        await apiRequest(`almacenes/${whForm.id}`, { method: 'PUT', body: payload });
        notify('Almacén actualizado.');
      } else {
        await apiRequest('almacenes', { method: 'POST', body: payload });
        notify('Nuevo almacén registrado.');
      }
      setWhForm({ id: '', codigo: '', nombre: '', direccion: '', responsable: '', capacidad: 2000, principal: false });
      await onRefresh();
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleDeleteWarehouse(id) {
    if (!window.confirm('¿Eliminar este almacén?')) return;
    try {
      await apiRequest(`almacenes/${id}`, { method: 'DELETE' });
      await onRefresh();
      notify('Almacén eliminado.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleSaveSupplier(e) {
    e.preventDefault();
    try {
      const payload = {
        ruc: supForm.ruc,
        nombre: supForm.nombre,
        telefono: supForm.telefono,
        email: supForm.email,
        direccion: supForm.direccion,
        rubro: supForm.rubro,
      };
      if (supForm.id) {
        await apiRequest(`proveedores/${supForm.id}`, { method: 'PUT', body: payload });
        notify('Proveedor actualizado.');
      } else {
        await apiRequest('proveedores', { method: 'POST', body: payload });
        notify('Proveedor registrado.');
      }
      setSupForm({ id: '', ruc: '', nombre: '', telefono: '', email: '', direccion: '', rubro: '' });
      await onRefresh();
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleDeleteSupplier(id) {
    if (!window.confirm('¿Eliminar este proveedor?')) return;
    try {
      await apiRequest(`proveedores/${id}`, { method: 'DELETE' });
      await onRefresh();
      notify('Proveedor eliminado.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  return (
    <div>
      {/* CONTROL DE INVENTARIO KARDEX */}
      <div className="panel">
        <div className="panel-header">
          <div className="panel-title">
            <h2>Control de Inventario y Movimientos (Kardex)</h2>
            <p>Registra entradas por compras, salidas por mermas o ajustes físicos de stock con trazabilidad completa</p>
          </div>
          <span className="badge badge-blue">{movements.length} movimientos</span>
        </div>

        <form className="form-grid" onSubmit={handleSaveMovement}>
          <div className="field" style={{ gridColumn: 'span 2' }}>
            <label>Producto *</label>
            <select
              required
              value={movForm.producto_id}
              onChange={(e) => setMovForm({ ...movForm, producto_id: e.target.value })}
            >
              <option value="">-- Seleccionar producto --</option>
              {products.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.nombre} (Stock actual: {p.stock})
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label>Tipo de Movimiento *</label>
            <select
              value={movForm.tipo}
              onChange={(e) => setMovForm({ ...movForm, tipo: e.target.value })}
            >
              <option value="entrada">➕ Entrada (Compra / Ingreso)</option>
              <option value="salida">➖ Salida (Merma / Retiro)</option>
              <option value="ajuste">⚙ Ajuste Físico Directo</option>
            </select>
          </div>
          <div className="field">
            <label>{movForm.tipo === 'ajuste' ? 'Nuevo Stock Físico *' : 'Cantidad *'}</label>
            <input
              required
              type="number"
              min="0"
              placeholder="Ej. 24"
              value={movForm.cantidad}
              onChange={(e) => setMovForm({ ...movForm, cantidad: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Almacén Destino</label>
            <select
              value={movForm.almacen_id}
              onChange={(e) => setMovForm({ ...movForm, almacen_id: e.target.value })}
            >
              <option value="">Mantener almacén actual</option>
              {warehouses.map((w) => (
                <option key={w.id} value={w.id}>
                  {w.nombre}
                </option>
              ))}
            </select>
          </div>
          <div className="field" style={{ gridColumn: 'span 2' }}>
            <label>Motivo / Documento de Referencia *</label>
            <input
              required
              placeholder="Ej. Factura F001-9210 Alicorp / Inventario mensual"
              value={movForm.motivo}
              onChange={(e) => setMovForm({ ...movForm, motivo: e.target.value })}
            />
          </div>
          <div>
            <button type="submit" className="btn btn-primary">
              Aplicar Movimiento
            </button>
          </div>
        </form>

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Fecha y Hora</th>
                <th>Producto</th>
                <th>Almacén</th>
                <th>Tipo</th>
                <th>Cant.</th>
                <th>Stock Anterior → Nuevo</th>
                <th>Motivo</th>
                <th>Responsable</th>
              </tr>
            </thead>
            <tbody>
              {[...movements].reverse().map((m) => (
                <tr key={m.id}>
                  <td>{m.fecha}</td>
                  <td><strong>{m.producto_nombre}</strong></td>
                  <td>{m.almacen_nombre}</td>
                  <td>
                    <span
                      className={`badge ${
                        m.tipo === 'entrada'
                          ? 'badge-green'
                          : m.tipo === 'salida'
                          ? 'badge-red'
                          : 'badge-amber'
                      }`}
                    >
                      {m.tipo.toUpperCase()}
                    </span>
                  </td>
                  <td><strong>{m.cantidad}</strong></td>
                  <td>
                    {m.stock_anterior} → <strong>{m.stock_nuevo}</strong>
                  </td>
                  <td>{m.motivo}</td>
                  <td>{m.usuario_nombre}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* ALMACENES Y PROVEEDORES */}
      <div className="grid-equal">
        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Almacenes y Sucursales</h2>
              <p>Bodegas, cámaras de frío y salas de exhibición de la empresa</p>
            </div>
          </div>

          <form className="form-grid" onSubmit={handleSaveWarehouse}>
            <div className="field">
              <label>Código</label>
              <input
                placeholder="ALM-04"
                value={whForm.codigo}
                onChange={(e) => setWhForm({ ...whForm, codigo: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Nombre del Almacén *</label>
              <input
                required
                placeholder="Ej. Almacén Secos Piso 2"
                value={whForm.nombre}
                onChange={(e) => setWhForm({ ...whForm, nombre: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Dirección / Ubicación</label>
              <input
                placeholder="Av. Larco 840 - Piso 2"
                value={whForm.direccion}
                onChange={(e) => setWhForm({ ...whForm, direccion: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Responsable</label>
              <input
                placeholder="Nombre encargado"
                value={whForm.responsable}
                onChange={(e) => setWhForm({ ...whForm, responsable: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Capacidad (Unds)</label>
              <input
                type="number"
                min="0"
                value={whForm.capacidad}
                onChange={(e) => setWhForm({ ...whForm, capacidad: e.target.value })}
              />
            </div>
            <div>
              <button type="submit" className="btn btn-primary">
                {whForm.id ? 'Actualizar' : 'Guardar Almacén'}
              </button>
            </div>
          </form>

          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Código</th>
                  <th>Almacén</th>
                  <th>Responsable</th>
                  <th>Capacidad</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {warehouses.map((w) => (
                  <tr key={w.id}>
                    <td><code>{w.codigo}</code></td>
                    <td>
                      <strong>{w.nombre}</strong>
                      <div style={{ fontSize: 11, color: '#68776e' }}>{w.direccion}</div>
                    </td>
                    <td>{w.responsable || '—'}</td>
                    <td><span className="badge badge-green">{w.capacidad} unds</span></td>
                    <td>
                      <div style={{ display: 'flex', gap: 6 }}>
                        <button
                          type="button"
                          className="btn btn-secondary btn-sm"
                          onClick={() => setWhForm(w)}
                        >
                          Editar
                        </button>
                        <button
                          type="button"
                          className="btn btn-danger btn-sm"
                          onClick={() => handleDeleteWarehouse(w.id)}
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

        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Directorio de Proveedores</h2>
              <p>Empresas distribuidoras y fabricantes asociados al abastecimiento</p>
            </div>
          </div>

          <form className="form-grid" onSubmit={handleSaveSupplier}>
            <div className="field">
              <label>RUC</label>
              <input
                placeholder="20100055237"
                value={supForm.ruc}
                onChange={(e) => setSupForm({ ...supForm, ruc: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Razón Social *</label>
              <input
                required
                placeholder="Ej. Distribuidora Central SAC"
                value={supForm.nombre}
                onChange={(e) => setSupForm({ ...supForm, nombre: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Teléfono</label>
              <input
                placeholder="999-888-777"
                value={supForm.telefono}
                onChange={(e) => setSupForm({ ...supForm, telefono: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Email</label>
              <input
                type="email"
                placeholder="ventas@proveedor.pe"
                value={supForm.email}
                onChange={(e) => setSupForm({ ...supForm, email: e.target.value })}
              />
            </div>
            <div className="field">
              <label>Rubro</label>
              <input
                placeholder="Lácteos / Bebidas"
                value={supForm.rubro}
                onChange={(e) => setSupForm({ ...supForm, rubro: e.target.value })}
              />
            </div>
            <div>
              <button type="submit" className="btn btn-primary">
                {supForm.id ? 'Actualizar' : 'Guardar Proveedor'}
              </button>
            </div>
          </form>

          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Proveedor</th>
                  <th>Contacto</th>
                  <th>Rubro</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {suppliers.map((s) => (
                  <tr key={s.id}>
                    <td>
                      <strong>{s.nombre}</strong>
                      {s.ruc && <div style={{ fontSize: 11, color: '#68776e' }}>RUC: {s.ruc}</div>}
                    </td>
                    <td>
                      <div>{s.telefono || '—'}</div>
                      <div style={{ fontSize: 11, color: '#68776e' }}>{s.email}</div>
                    </td>
                    <td><span className="badge badge-blue">{s.rubro || 'General'}</span></td>
                    <td>
                      <div style={{ display: 'flex', gap: 6 }}>
                        <button
                          type="button"
                          className="btn btn-secondary btn-sm"
                          onClick={() => setSupForm(s)}
                        >
                          Editar
                        </button>
                        <button
                          type="button"
                          className="btn btn-danger btn-sm"
                          onClick={() => handleDeleteSupplier(s.id)}
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
    </div>
  );
}
