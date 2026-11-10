import React from 'react';
import { formatPEN } from '../api';

export default function DashboardView({ dashboard, onNavigate }) {
  if (!dashboard) {
    return <div className="panel">Cargando métricas del Dashboard...</div>;
  }

  const maxCategorySale = Math.max(
    1,
    ...(dashboard.ventas_por_categoria || []).map((c) => c.monto_ventas || 0)
  );

  const paymentLabels = {
    efectivo: 'Efectivo en Caja',
    tarjeta: 'Tarjeta Débito / Crédito',
    yape_plin: 'Billetera Digital (Yape / Plin)',
    transferencia: 'Transferencia Bancaria',
  };

  return (
    <div>
      {/* KPI CARDS */}
      <div className="kpi-grid">
        <div className="kpi-card">
          <div className="kpi-top">
            <span>Ingresos por Ventas</span>
            <span className="kpi-icon">S/</span>
          </div>
          <div className="kpi-value">{formatPEN(dashboard.ingresos_totales)}</div>
          <div className="kpi-sub">
            {dashboard.total_ventas} ventas · Ticket prom: {formatPEN(dashboard.ticket_promedio)}
          </div>
        </div>

        <div className="kpi-card">
          <div className="kpi-top">
            <span>Catálogo y Stock</span>
            <span className="kpi-icon">▦</span>
          </div>
          <div className="kpi-value" id="product-count">
            {dashboard.total_productos} prod.
          </div>
          <div className="kpi-sub">
            <strong id="stock-count">{dashboard.total_unidades_stock}</strong> unds · {dashboard.total_categorias} categorías
          </div>
        </div>

        <div className="kpi-card">
          <div className="kpi-top">
            <span>Valorización Inventario</span>
            <span className="kpi-icon">↗</span>
          </div>
          <div className="kpi-value">{formatPEN(dashboard.valor_inventario)}</div>
          <div className="kpi-sub">
            Costo: {formatPEN(dashboard.costo_inventario)} · {dashboard.total_almacenes} almacenes
          </div>
        </div>

        <div className="kpi-card">
          <div className="kpi-top">
            <span>Alertas de Stock Bajo</span>
            <span className="kpi-icon" style={{ background: '#fef4e6', color: '#c27820' }}>!</span>
          </div>
          <div className="kpi-value" style={{ color: dashboard.alertas_stock_bajo > 0 ? '#c27820' : '#245e3f' }}>
            {dashboard.alertas_stock_bajo}
          </div>
          <div className="kpi-sub">
            Proveedores activos: <strong id="supplier-count">{dashboard.total_proveedores}</strong>
          </div>
        </div>

        <div className="kpi-card">
          <div className="kpi-top">
            <span>Recursos Humanos (RRHH)</span>
            <span className="kpi-icon">☺</span>
          </div>
          <div className="kpi-value">{dashboard.empleados_activos} activos</div>
          <div className="kpi-sub">
            Planilla mensual: {formatPEN(dashboard.planilla_mensual)}
          </div>
        </div>
      </div>

      {/* CHARTS ROW */}
      <div className="grid-2">
        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Rendimiento por Categoría</h2>
              <p>Distribución de ingresos y unidades en stock por línea de producto</p>
            </div>
            <button className="btn btn-secondary btn-sm" onClick={() => onNavigate('catalogo')}>
              Gestionar Categorías
            </button>
          </div>
          {(dashboard.ventas_por_categoria || []).length === 0 ? (
            <p style={{ color: '#68776e' }}>Aún no hay categorías con movimientos.</p>
          ) : (
            (dashboard.ventas_por_categoria || []).map((cat) => {
              const pct = Math.min(100, Math.round(((cat.monto_ventas || 0) / maxCategorySale) * 100));
              return (
                <div className="bar-row" key={cat.categoria_id}>
                  <div className="bar-label">
                    <span>
                      <span
                        style={{
                          display: 'inline-block',
                          width: 10,
                          height: 10,
                          borderRadius: '50%',
                          background: cat.color || '#245e3f',
                          marginRight: 8,
                        }}
                      />
                      {cat.categoria_nombre}{' '}
                      <small style={{ color: '#68776e', fontWeight: 500 }}>
                        ({cat.cantidad_productos} prod · {cat.stock_total} unds)
                      </small>
                    </span>
                    <span>{formatPEN(cat.monto_ventas)}</span>
                  </div>
                  <div className="bar-track">
                    <div
                      className="bar-fill"
                      style={{
                        width: `${Math.max(pct, 6)}%`,
                        background: cat.color || '#245e3f',
                      }}
                    />
                  </div>
                </div>
              );
            })
          )}
        </div>

        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Medios de Pago y Top Productos</h2>
              <p>Preferencia de pago y artículos de mayor rotación</p>
            </div>
            <button className="btn btn-primary btn-sm" onClick={() => onNavigate('ventas')}>
              Ir al POS
            </button>
          </div>

          <div style={{ marginBottom: 16 }}>
            <div style={{ fontSize: 11, fontWeight: 700, textTransform: 'uppercase', color: '#68776e', marginBottom: 8 }}>
              Recaudación por Método de Pago
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {(dashboard.ventas_por_metodo_pago || []).map((pm) => (
                <div
                  key={pm.metodo}
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '8px 12px',
                    background: '#f8faf8',
                    borderRadius: 9,
                    border: '1px solid #e6ede7',
                  }}
                >
                  <div>
                    <strong>{paymentLabels[pm.metodo] || pm.metodo}</strong>
                    <span style={{ marginLeft: 8, fontSize: 11.5, color: '#68776e' }}>
                      ({pm.cantidad} operaciones)
                    </span>
                  </div>
                  <strong style={{ color: '#245e3f' }}>{formatPEN(pm.total)}</strong>
                </div>
              ))}
            </div>
          </div>

          <div>
            <div style={{ fontSize: 11, fontWeight: 700, textTransform: 'uppercase', color: '#68776e', marginBottom: 8 }}>
              Productos Más Vendidos
            </div>
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Producto</th>
                    <th>Unidades</th>
                    <th>Ingresos</th>
                  </tr>
                </thead>
                <tbody>
                  {(dashboard.productos_mas_vendidos || []).map((tp) => (
                    <tr key={tp.producto_id}>
                      <td><strong>{tp.producto_nombre}</strong></td>
                      <td><span className="badge badge-green">{tp.unidades_vendidas} und</span></td>
                      <td><strong>{formatPEN(tp.ingresos)}</strong></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      {/* LOW STOCK ALERTS & RECENT ACTIVITY */}
      <div className="grid-equal">
        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Alertas de Reposición (Stock Crítico)</h2>
              <p>Productos cuyo stock actual está en o por debajo del mínimo permitido</p>
            </div>
            <button className="btn btn-secondary btn-sm" onClick={() => onNavigate('inventario')}>
              Abastecer Kardex
            </button>
          </div>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>SKU</th>
                  <th>Producto</th>
                  <th>Stock Actual</th>
                  <th>Mínimo</th>
                  <th>Precio</th>
                </tr>
              </thead>
              <tbody>
                {(dashboard.productos_stock_bajo || []).length === 0 ? (
                  <tr>
                    <td colSpan={5} style={{ textAlign: 'center', padding: 20, color: '#68776e' }}>
                      Todo el inventario tiene niveles óptimos de stock.
                    </td>
                  </tr>
                ) : (
                  dashboard.productos_stock_bajo.map((p) => (
                    <tr key={p.id}>
                      <td><code>{p.sku || 'S/C'}</code></td>
                      <td><strong>{p.nombre}</strong></td>
                      <td><span className="badge badge-amber">{p.stock} {p.unidad}</span></td>
                      <td>{p.stock_minimo}</td>
                      <td>{formatPEN(p.precio)}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="panel">
          <div className="panel-header">
            <div className="panel-title">
              <h2>Últimas Ventas Registradas</h2>
              <p>Transacciones recientes emitidas en el Punto de Venta</p>
            </div>
          </div>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Comprobante</th>
                  <th>Cliente</th>
                  <th>Método</th>
                  <th>Total</th>
                </tr>
              </thead>
              <tbody>
                {(dashboard.ultimas_ventas || []).map((v) => (
                  <tr key={v.id}>
                    <td>
                      <strong>{v.numero_comprobante}</strong>
                      <div style={{ fontSize: 11, color: '#68776e' }}>{v.fecha}</div>
                    </td>
                    <td>{v.cliente_nombre}</td>
                    <td><span className="badge badge-blue">{v.metodo_pago}</span></td>
                    <td><strong>{formatPEN(v.total)}</strong></td>
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
