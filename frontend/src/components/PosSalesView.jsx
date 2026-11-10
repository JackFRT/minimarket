import React, { useState } from 'react';
import { apiRequest, formatPEN } from '../api';

export default function PosSalesView({
  products,
  categories,
  warehouses,
  sales,
  currentUser,
  onRefresh,
  notify,
}) {
  const [search, setSearch] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');
  const [cart, setCart] = useState([]);
  const [customerName, setCustomerName] = useState('Cliente Varios');
  const [customerDoc, setCustomerDoc] = useState('');
  const [receiptType, setReceiptType] = useState('Boleta');
  const [paymentMethod, setPaymentMethod] = useState('efectivo');
  const [warehouseId, setWarehouseId] = useState(warehouses[0]?.id || '');
  const [lastReceipt, setLastReceipt] = useState(null);
  const [submitting, setSubmitting] = useState(false);

  const filteredProducts = products.filter((p) => {
    const matchSearch =
      p.nombre.toLowerCase().includes(search.toLowerCase()) ||
      (p.sku || '').toLowerCase().includes(search.toLowerCase());
    const matchCategory = !selectedCategory || p.categoria_id === selectedCategory;
    return matchSearch && matchCategory;
  });

  function addToCart(product) {
    if (product.stock <= 0) {
      notify(`Sin stock disponible para ${product.nombre}`, true);
      return;
    }
    setCart((prev) => {
      const existing = prev.find((item) => item.producto_id === product.id);
      if (existing) {
        if (existing.cantidad + 1 > product.stock) {
          notify(`Stock máximo alcanzado (${product.stock}) para ${product.nombre}`, true);
          return prev;
        }
        return prev.map((item) =>
          item.producto_id === product.id
            ? { ...item, cantidad: item.cantidad + 1 }
            : item
        );
      }
      return [
        ...prev,
        {
          producto_id: product.id,
          producto_nombre: product.nombre,
          precio_unitario: product.precio,
          stock_max: product.stock,
          cantidad: 1,
        },
      ];
    });
  }

  function updateQuantity(productId, delta) {
    setCart((prev) =>
      prev
        .map((item) => {
          if (item.producto_id !== productId) return item;
          const nextQty = item.cantidad + delta;
          if (nextQty > item.stock_max) {
            notify(`Solo hay ${item.stock_max} unidades disponibles`, true);
            return item;
          }
          return { ...item, cantidad: nextQty };
        })
        .filter((item) => item.cantidad > 0)
    );
  }

  const totalAmount = cart.reduce(
    (sum, item) => sum + item.precio_unitario * item.cantidad,
    0
  );
  const subtotalAmount = totalAmount / 1.18;
  const igvAmount = totalAmount - subtotalAmount;

  async function handleCheckout(e) {
    e.preventDefault();
    if (cart.length === 0) {
      notify('Agrega al menos un producto al carrito de venta', true);
      return;
    }
    setSubmitting(true);
    try {
      const created = await apiRequest('ventas', {
        method: 'POST',
        body: {
          tipo_comprobante: receiptType,
          cliente_nombre: customerName || 'Cliente Varios',
          cliente_documento: customerDoc,
          almacen_id: warehouseId || warehouses[0]?.id || '',
          usuario_id: currentUser?.id || '',
          vendedor_nombre: currentUser?.nombre || 'Cajero POS',
          metodo_pago: paymentMethod,
          items: cart.map((c) => ({
            producto_id: c.producto_id,
            cantidad: c.cantidad,
            precio_unitario: c.precio_unitario,
          })),
        },
      });
      setCart([]);
      setCustomerName('Cliente Varios');
      setCustomerDoc('');
      setLastReceipt(created);
      await onRefresh();
      notify(`Venta ${created.numero_comprobante} registrada y stock actualizado.`);
    } catch (err) {
      notify(err.message, true);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div>
      <div className="pos-container">
        {/* CATÁLOGO INTERACTIVO POS */}
        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Terminal Punto de Venta (POS)</h2>
              <p>Haz clic en los productos para añadirlos al comprobante y descontar stock</p>
            </div>
          </div>

          <div style={{ display: 'flex', gap: 10, flexWrap: 'wrap' }}>
            <input
              style={{
                flex: 1,
                minWidth: 200,
                height: 38,
                padding: '0 12px',
                borderRadius: 9,
                border: '1px solid #d4ddd6',
              }}
              placeholder="Buscar producto por nombre o código SKU..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <select
              style={{
                height: 38,
                padding: '0 12px',
                borderRadius: 9,
                border: '1px solid #d4ddd6',
                background: '#fff',
              }}
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
            >
              <option value="">Todas las categorías</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.nombre}
                </option>
              ))}
            </select>
          </div>

          <div className="pos-catalog-grid">
            {filteredProducts.map((p) => (
              <div
                key={p.id}
                className={`pos-item-card ${p.stock <= 0 ? 'disabled' : ''}`}
                onClick={() => addToCart(p)}
              >
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 11, color: '#68776e' }}>
                    <code>{p.sku || 'SKU'}</code>
                    <span className={`badge ${p.stock <= (p.stock_minimo || 5) ? 'badge-amber' : 'badge-green'}`}>
                      Stock: {p.stock}
                    </span>
                  </div>
                  <div style={{ fontWeight: 700, marginTop: 6, color: '#14261d', lineHeight: 1.3 }}>
                    {p.nombre}
                  </div>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 6 }}>
                  <strong style={{ fontSize: 15, color: '#245e3f' }}>{formatPEN(p.precio)}</strong>
                  <span style={{ fontSize: 11, color: '#245e3f', fontWeight: 700 }}>+ Agregar</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* CARRITO DE COMPRA Y EMISIÓN */}
        <form className="pos-cart" onSubmit={handleCheckout}>
          <h3 style={{ margin: '0 0 12px', fontSize: 16 }}>Comprobante de Venta</h3>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10, marginBottom: 10 }}>
            <div className="field">
              <label>Tipo Comprobante</label>
              <select value={receiptType} onChange={(e) => setReceiptType(e.target.value)}>
                <option value="Boleta">Boleta de Venta</option>
                <option value="Factura">Factura Electrónica</option>
                <option value="Ticket">Ticket Interno</option>
              </select>
            </div>
            <div className="field">
              <label>Método de Pago</label>
              <select value={paymentMethod} onChange={(e) => setPaymentMethod(e.target.value)}>
                <option value="efectivo">Efectivo</option>
                <option value="yape_plin">Yape / Plin (QR)</option>
                <option value="tarjeta">Tarjeta POS</option>
                <option value="transferencia">Transferencia</option>
              </select>
            </div>
            <div className="field">
              <label>Cliente / Razón Social</label>
              <input
                value={customerName}
                onChange={(e) => setCustomerName(e.target.value)}
                placeholder="Cliente Varios"
              />
            </div>
            <div className="field">
              <label>DNI / RUC</label>
              <input
                value={customerDoc}
                onChange={(e) => setCustomerDoc(e.target.value)}
                placeholder="Opcional"
              />
            </div>
          </div>

          <div style={{ maxHeight: 240, overflowY: 'auto', borderTop: '1px solid #e2e8e3', paddingTop: 10 }}>
            {cart.length === 0 ? (
              <p style={{ textAlign: 'center', color: '#68776e', padding: '18px 0' }}>
                Selecciona productos del catálogo para armar la venta.
              </p>
            ) : (
              cart.map((item) => (
                <div
                  key={item.producto_id}
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '8px 0',
                    borderBottom: '1px dashed #e2e8e3',
                  }}
                >
                  <div style={{ flex: 1 }}>
                    <div style={{ fontWeight: 650, fontSize: 12.5 }}>{item.producto_nombre}</div>
                    <div style={{ fontSize: 11.5, color: '#68776e' }}>
                      {formatPEN(item.precio_unitario)} x {item.cantidad} ={' '}
                      <strong>{formatPEN(item.precio_unitario * item.cantidad)}</strong>
                    </div>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                    <button
                      type="button"
                      className="btn btn-secondary btn-sm"
                      onClick={() => updateQuantity(item.producto_id, -1)}
                    >
                      −
                    </button>
                    <strong style={{ minWidth: 20, textAlign: 'center' }}>{item.cantidad}</strong>
                    <button
                      type="button"
                      className="btn btn-secondary btn-sm"
                      onClick={() => updateQuantity(item.producto_id, 1)}
                    >
                      +
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>

          <div className="cart-totals">
            <div className="cart-row">
              <span>Op. Gravada (Subtotal)</span>
              <span>{formatPEN(subtotalAmount)}</span>
            </div>
            <div className="cart-row">
              <span>I.G.V. (18%)</span>
              <span>{formatPEN(igvAmount)}</span>
            </div>
            <div className="cart-row grand-total">
              <span>Total a Pagar</span>
              <span>{formatPEN(totalAmount)}</span>
            </div>
          </div>

          <div style={{ display: 'flex', gap: 8 }}>
            <button
              type="submit"
              className="btn btn-primary"
              style={{ flex: 1 }}
              disabled={submitting || cart.length === 0}
            >
              {submitting ? 'Procesando...' : 'Emitir Venta y Descontar Stock'}
            </button>
            <button
              type="button"
              className="btn btn-secondary"
              onClick={() => setCart([])}
            >
              Vaciar
            </button>
          </div>
        </form>
      </div>

      {/* HISTORIAL DE VENTAS */}
      <div className="panel" style={{ marginTop: 22 }}>
        <div className="panel-header">
          <div className="panel-title">
            <h2>Registro Histórico de Ventas</h2>
            <p>Comprobantes emitidos en la empresa activa con desglose de IGV e ítems vendidos</p>
          </div>
          <span className="badge badge-green">{sales.length} comprobantes</span>
        </div>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Comprobante</th>
                <th>Fecha</th>
                <th>Cliente</th>
                <th>Vendedor</th>
                <th>Método Pago</th>
                <th>Detalle Productos</th>
                <th>Subtotal</th>
                <th>IGV (18%)</th>
                <th>Total</th>
              </tr>
            </thead>
            <tbody>
              {sales.length === 0 ? (
                <tr>
                  <td colSpan={9} style={{ textAlign: 'center', padding: 20 }}>
                    No hay ventas registradas aún en esta empresa.
                  </td>
                </tr>
              ) : (
                [...sales].reverse().map((s) => (
                  <tr key={s.id}>
                    <td>
                      <strong>{s.numero_comprobante}</strong>
                      <div style={{ fontSize: 11, color: '#68776e' }}>{s.tipo_comprobante}</div>
                    </td>
                    <td>{s.fecha}</td>
                    <td>
                      {s.cliente_nombre}
                      {s.cliente_documento && (
                        <div style={{ fontSize: 11, color: '#68776e' }}>Doc: {s.cliente_documento}</div>
                      )}
                    </td>
                    <td>{s.vendedor_nombre}</td>
                    <td><span className="badge badge-blue">{s.metodo_pago}</span></td>
                    <td>
                      {(s.items || []).map((it, idx) => (
                        <div key={idx} style={{ fontSize: 12 }}>
                          • {it.cantidad}x {it.producto_nombre} ({formatPEN(it.subtotal)})
                        </div>
                      ))}
                    </td>
                    <td>{formatPEN(s.subtotal)}</td>
                    <td>{formatPEN(s.igv)}</td>
                    <td><strong style={{ color: '#245e3f' }}>{formatPEN(s.total)}</strong></td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* MODAL DE BOLETA EMITIDA */}
      {lastReceipt && (
        <div className="modal-backdrop" onClick={() => setLastReceipt(null)}>
          <div className="modal-card" onClick={(e) => e.stopPropagation()}>
            <div style={{ textAlign: 'center', borderBottom: '1px dashed #ccc', paddingBottom: 12, marginBottom: 12 }}>
              <span className="badge badge-green">Comprobante Electrónico Emitido</span>
              <h3 style={{ margin: '8px 0 2px' }}>
                {lastReceipt.tipo_comprobante} {lastReceipt.numero_comprobante}
              </h3>
              <small style={{ color: '#68776e' }}>Fecha: {lastReceipt.fecha}</small>
            </div>
            <p style={{ margin: '4px 0', fontSize: 13 }}>
              <strong>Cliente:</strong> {lastReceipt.cliente_nombre}
            </p>
            <p style={{ margin: '4px 0 12px', fontSize: 13 }}>
              <strong>Medio de Pago:</strong> {lastReceipt.metodo_pago}
            </p>
            <table style={{ marginBottom: 14 }}>
              <thead>
                <tr>
                  <th>Cant.</th>
                  <th>Descripción</th>
                  <th>Importe</th>
                </tr>
              </thead>
              <tbody>
                {(lastReceipt.items || []).map((it, idx) => (
                  <tr key={idx}>
                    <td>{it.cantidad}</td>
                    <td>{it.producto_nombre}</td>
                    <td>{formatPEN(it.subtotal)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div className="cart-totals">
              <div className="cart-row">
                <span>Subtotal</span>
                <span>{formatPEN(lastReceipt.subtotal)}</span>
              </div>
              <div className="cart-row">
                <span>IGV (18%)</span>
                <span>{formatPEN(lastReceipt.igv)}</span>
              </div>
              <div className="cart-row grand-total">
                <span>Total Pagado</span>
                <span>{formatPEN(lastReceipt.total)}</span>
              </div>
            </div>
            <button
              className="btn btn-primary"
              style={{ width: '100%' }}
              onClick={() => setLastReceipt(null)}
            >
              Cerrar Comprobante
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
