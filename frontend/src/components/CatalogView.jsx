import React, { useState } from 'react';
import { apiRequest, formatPEN } from '../api';

const emptyProduct = {
  id: '',
  sku: '',
  nombre: '',
  categoria_id: '',
  almacen_id: '',
  proveedor_id: '',
  precio_compra: '',
  precio: '',
  stock: '',
  stock_minimo: 5,
  unidad: 'und',
};

const emptyCategory = {
  id: '',
  codigo: '',
  nombre: '',
  descripcion: '',
  color: '#245e3f',
};

export default function CatalogView({
  products,
  categories,
  warehouses,
  suppliers,
  onRefresh,
  notify,
}) {
  const [prodForm, setProdForm] = useState(emptyProduct);
  const [catForm, setCatForm] = useState(emptyCategory);
  const [filterCat, setFilterCat] = useState('');

  const catMap = Object.fromEntries(categories.map((c) => [c.id, c]));
  const whMap = Object.fromEntries(warehouses.map((w) => [w.id, w]));
  const supMap = Object.fromEntries(suppliers.map((s) => [s.id, s]));

  async function handleSaveProduct(e) {
    e.preventDefault();
    try {
      const payload = {
        sku: prodForm.sku,
        nombre: prodForm.nombre,
        categoria_id: prodForm.categoria_id,
        almacen_id: prodForm.almacen_id,
        proveedor_id: prodForm.proveedor_id,
        precio_compra: Number(prodForm.precio_compra || 0),
        precio: Number(prodForm.precio || 0),
        stock: Number(prodForm.stock || 0),
        stock_minimo: Number(prodForm.stock_minimo || 5),
        unidad: prodForm.unidad || 'und',
      };
      if (prodForm.id) {
        await apiRequest(`productos/${prodForm.id}`, { method: 'PUT', body: payload });
        notify('Producto actualizado correctamente.');
      } else {
        await apiRequest('productos', { method: 'POST', body: payload });
        notify('Producto registrado en el catálogo.');
      }
      setProdForm(emptyProduct);
      await onRefresh();
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleCloneProduct(id) {
    try {
      const cloned = await apiRequest(`productos/${id}/clonar`, { method: 'POST' });
      await onRefresh();
      notify(`Patrón Prototype aplicado: Producto clonado (${cloned.sku || cloned.id.slice(0, 6)}).`);
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleDeleteProduct(id) {
    if (!window.confirm('¿Deseas eliminar este producto del catálogo?')) return;
    try {
      await apiRequest(`productos/${id}`, { method: 'DELETE' });
      await onRefresh();
      notify('Producto eliminado.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleSaveCategory(e) {
    e.preventDefault();
    try {
      const payload = {
        codigo: catForm.codigo,
        nombre: catForm.nombre,
        descripcion: catForm.descripcion,
        color: catForm.color || '#245e3f',
      };
      if (catForm.id) {
        await apiRequest(`categorias/${catForm.id}`, { method: 'PUT', body: payload });
        notify('Categoría actualizada.');
      } else {
        await apiRequest('categorias', { method: 'POST', body: payload });
        notify('Categoría creada.');
      }
      setCatForm(emptyCategory);
      await onRefresh();
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleCloneCategory(id) {
    try {
      await apiRequest(`categorias/${id}/clonar`, { method: 'POST' });
      await onRefresh();
      notify('Categoría clonada mediante patrón Prototype.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleDeleteCategory(id) {
    if (!window.confirm('¿Eliminar esta categoría?')) return;
    try {
      await apiRequest(`categorias/${id}`, { method: 'DELETE' });
      await onRefresh();
      notify('Categoría eliminada.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  const displayedProducts = filterCat
    ? products.filter((p) => p.categoria_id === filterCat)
    : products;

  return (
    <div>
      {/* GESTIÓN DE PRODUCTOS */}
      <div className="panel">
        <div className="panel-header">
          <div className="panel-title">
            <h2>Catálogo Maestro de Productos (Patrón Prototype)</h2>
            <p>
              Administra precios, costos, categorías, almacén asignado y clona productos con un clic usando <code>Product.Clone()</code>
            </p>
          </div>
          <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
            <select
              style={{ height: 34, borderRadius: 8, border: '1px solid #d4ddd6', padding: '0 10px' }}
              value={filterCat}
              onChange={(e) => setFilterCat(e.target.value)}
            >
              <option value="">Todas las categorías ({products.length})</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.nombre}
                </option>
              ))}
            </select>
          </div>
        </div>

        <form className="form-grid" onSubmit={handleSaveProduct}>
          <div className="field">
            <label>Código SKU</label>
            <input
              placeholder="Ej. ABR-010"
              value={prodForm.sku}
              onChange={(e) => setProdForm({ ...prodForm, sku: e.target.value })}
            />
          </div>
          <div className="field" style={{ gridColumn: 'span 2' }}>
            <label>Nombre del Producto *</label>
            <input
              required
              placeholder="Ej. Arroz Extra Costeño 750g"
              value={prodForm.nombre}
              onChange={(e) => setProdForm({ ...prodForm, nombre: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Categoría</label>
            <select
              value={prodForm.categoria_id}
              onChange={(e) => setProdForm({ ...prodForm, categoria_id: e.target.value })}
            >
              <option value="">-- Seleccionar --</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.nombre}
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label>Almacén</label>
            <select
              value={prodForm.almacen_id}
              onChange={(e) => setProdForm({ ...prodForm, almacen_id: e.target.value })}
            >
              <option value="">-- Seleccionar --</option>
              {warehouses.map((w) => (
                <option key={w.id} value={w.id}>
                  {w.nombre}
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label>Proveedor</label>
            <select
              value={prodForm.proveedor_id}
              onChange={(e) => setProdForm({ ...prodForm, proveedor_id: e.target.value })}
            >
              <option value="">-- Seleccionar --</option>
              {suppliers.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.nombre}
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label>Costo Compra (S/)</label>
            <input
              type="number"
              step="0.01"
              min="0"
              placeholder="0.00"
              value={prodForm.precio_compra}
              onChange={(e) => setProdForm({ ...prodForm, precio_compra: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Precio Venta (S/) *</label>
            <input
              required
              type="number"
              step="0.01"
              min="0"
              placeholder="0.00"
              value={prodForm.precio}
              onChange={(e) => setProdForm({ ...prodForm, precio: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Stock Inicial *</label>
            <input
              required
              type="number"
              min="0"
              placeholder="0"
              value={prodForm.stock}
              onChange={(e) => setProdForm({ ...prodForm, stock: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Stock Mínimo</label>
            <input
              type="number"
              min="0"
              value={prodForm.stock_minimo}
              onChange={(e) => setProdForm({ ...prodForm, stock_minimo: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Unidad</label>
            <select
              value={prodForm.unidad}
              onChange={(e) => setProdForm({ ...prodForm, unidad: e.target.value })}
            >
              <option value="und">Unidad (und)</option>
              <option value="bolsa">Bolsa</option>
              <option value="botella">Botella</option>
              <option value="lata">Lata</option>
              <option value="caja">Caja</option>
              <option value="kg">Kilogramo (kg)</option>
            </select>
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
            <button type="submit" className="btn btn-primary">
              {prodForm.id ? 'Actualizar Producto' : 'Guardar Producto'}
            </button>
            {prodForm.id && (
              <button
                type="button"
                className="btn btn-secondary"
                onClick={() => setProdForm(emptyProduct)}
              >
                Cancelar
              </button>
            )}
          </div>
        </form>

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>SKU</th>
                <th>Producto</th>
                <th>Categoría</th>
                <th>Almacén / Proveedor</th>
                <th>Costo</th>
                <th>Precio Venta</th>
                <th>Stock</th>
                <th>Acciones (CRUD + Prototype)</th>
              </tr>
            </thead>
            <tbody>
              {displayedProducts.map((p) => {
                const cat = catMap[p.categoria_id];
                const wh = whMap[p.almacen_id];
                const sup = supMap[p.proveedor_id];
                const isLow = p.stock <= (p.stock_minimo || 5);
                return (
                  <tr key={p.id}>
                    <td><code>{p.sku || p.id.slice(0, 6)}</code></td>
                    <td><strong>{p.nombre}</strong></td>
                    <td>
                      {cat ? (
                        <span
                          className="badge"
                          style={{
                            background: `${cat.color || '#245e3f'}18`,
                            color: cat.color || '#245e3f',
                          }}
                        >
                          {cat.nombre}
                        </span>
                      ) : (
                        <span style={{ color: '#68776e' }}>Sin categoría</span>
                      )}
                    </td>
                    <td>
                      <div>{wh ? wh.nombre : 'Almacén General'}</div>
                      {sup && <div style={{ fontSize: 11, color: '#68776e' }}>Prov: {sup.nombre}</div>}
                    </td>
                    <td>{formatPEN(p.precio_compra)}</td>
                    <td><strong style={{ color: '#245e3f' }}>{formatPEN(p.precio)}</strong></td>
                    <td>
                      <span className={`badge ${isLow ? 'badge-amber' : 'badge-green'}`}>
                        {p.stock} {p.unidad || 'und'}
                      </span>
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                        <button
                          type="button"
                          className="btn btn-secondary btn-sm"
                          onClick={() =>
                            setProdForm({
                              id: p.id,
                              sku: p.sku || '',
                              nombre: p.nombre,
                              categoria_id: p.categoria_id || '',
                              almacen_id: p.almacen_id || '',
                              proveedor_id: p.proveedor_id || '',
                              precio_compra: p.precio_compra || 0,
                              precio: p.precio,
                              stock: p.stock,
                              stock_minimo: p.stock_minimo || 5,
                              unidad: p.unidad || 'und',
                            })
                          }
                        >
                          Editar
                        </button>
                        <button
                          type="button"
                          className="btn btn-prototype btn-sm"
                          title="Clonar producto con nueva identidad (Patrón Prototype)"
                          onClick={() => handleCloneProduct(p.id)}
                        >
                          ⧉ Clonar (Prototype)
                        </button>
                        <button
                          type="button"
                          className="btn btn-danger btn-sm"
                          onClick={() => handleDeleteProduct(p.id)}
                        >
                          Eliminar
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* GESTIÓN DE CATEGORÍAS */}
      <div className="panel">
        <div className="panel-header">
          <div className="panel-title">
            <h2>Categorías de Productos</h2>
            <p>Clasificación de líneas de productos por empresa</p>
          </div>
          <span className="badge badge-green">{categories.length} categorías</span>
        </div>

        <form className="form-grid" onSubmit={handleSaveCategory}>
          <div className="field">
            <label>Código</label>
            <input
              placeholder="Ej. CAT-06"
              value={catForm.codigo}
              onChange={(e) => setCatForm({ ...catForm, codigo: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Nombre Categoría *</label>
            <input
              required
              placeholder="Ej. Panadería y Pastelería"
              value={catForm.nombre}
              onChange={(e) => setCatForm({ ...catForm, nombre: e.target.value })}
            />
          </div>
          <div className="field" style={{ gridColumn: 'span 2' }}>
            <label>Descripción</label>
            <input
              placeholder="Panes artesanales, kekes y postres"
              value={catForm.descripcion}
              onChange={(e) => setCatForm({ ...catForm, descripcion: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Color Distintivo</label>
            <input
              type="color"
              value={catForm.color}
              onChange={(e) => setCatForm({ ...catForm, color: e.target.value })}
            />
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
            <button type="submit" className="btn btn-primary">
              {catForm.id ? 'Actualizar Categoría' : 'Crear Categoría'}
            </button>
            {catForm.id && (
              <button
                type="button"
                className="btn btn-secondary"
                onClick={() => setCatForm(emptyCategory)}
              >
                Cancelar
              </button>
            )}
          </div>
        </form>

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Código</th>
                <th>Categoría</th>
                <th>Descripción</th>
                <th>Productos Asociados</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {categories.map((c) => {
                const count = products.filter((p) => p.categoria_id === c.id).length;
                return (
                  <tr key={c.id}>
                    <td><code>{c.codigo || 'CAT'}</code></td>
                    <td>
                      <span
                        style={{
                          display: 'inline-block',
                          width: 10,
                          height: 10,
                          borderRadius: '50%',
                          background: c.color || '#245e3f',
                          marginRight: 8,
                        }}
                      />
                      <strong>{c.nombre}</strong>
                    </td>
                    <td>{c.descripcion || '—'}</td>
                    <td><span className="badge badge-green">{count} productos</span></td>
                    <td>
                      <div style={{ display: 'flex', gap: 6 }}>
                        <button
                          type="button"
                          className="btn btn-secondary btn-sm"
                          onClick={() => setCatForm(c)}
                        >
                          Editar
                        </button>
                        <button
                          type="button"
                          className="btn btn-prototype btn-sm"
                          onClick={() => handleCloneCategory(c.id)}
                        >
                          ⧉ Clonar
                        </button>
                        <button
                          type="button"
                          className="btn btn-danger btn-sm"
                          onClick={() => handleDeleteCategory(c.id)}
                        >
                          Eliminar
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
