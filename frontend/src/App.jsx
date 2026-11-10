import React, { useState, useEffect, useCallback } from 'react';
import {
  apiRequest,
  getStoredCompanyId,
  setStoredCompanyId,
  getStoredSession,
  setStoredSession,
} from './api';
import DashboardView from './components/DashboardView';
import PosSalesView from './components/PosSalesView';
import CatalogView from './components/CatalogView';
import InventoryView from './components/InventoryView';
import HrView from './components/HrView';
import ArchitectureSaasView from './components/ArchitectureSaasView';
import AuthModal from './components/AuthModal';

export default function App() {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [companyId, setCompanyId] = useState(getStoredCompanyId());
  const [session, setSession] = useState(
    getStoredSession() || {
      usuario: {
        id: 'usr-admin-1',
        nombre: 'Carlos Mendoza (Admin)',
        email: 'admin@minimarket.pe',
        rol: 'admin',
      },
      empresa: {
        id: 'emp-default',
        nombre: 'Minimarket El Sol - Sede Central',
        plan: 'Enterprise',
      },
    }
  );
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [toast, setToast] = useState(null);

  // Datos sincronizados con la API Go por empresa activa
  const [companies, setCompanies] = useState([]);
  const [dashboard, setDashboard] = useState(null);
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [warehouses, setWarehouses] = useState([]);
  const [suppliers, setSuppliers] = useState([]);
  const [movements, setMovements] = useState([]);
  const [sales, setSales] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [attendances, setAttendances] = useState([]);
  const [users, setUsers] = useState([]);

  const notify = useCallback((message, isError = false) => {
    setToast({ message, isError });
    setTimeout(() => {
      setToast((prev) => (prev?.message === message ? null : prev));
    }, 3600);
  }, []);

  const loadAllData = useCallback(
    async (targetCompany = companyId) => {
      try {
        const [
          compList,
          dashData,
          prodList,
          catList,
          whList,
          supList,
          movList,
          saleList,
          empList,
          attList,
          usrList,
        ] = await Promise.all([
          apiRequest('empresas', { companyId: targetCompany }),
          apiRequest('dashboard', { companyId: targetCompany }),
          apiRequest('productos', { companyId: targetCompany }),
          apiRequest('categorias', { companyId: targetCompany }),
          apiRequest('almacenes', { companyId: targetCompany }),
          apiRequest('proveedores', { companyId: targetCompany }),
          apiRequest('inventario/movimientos', { companyId: targetCompany }),
          apiRequest('ventas', { companyId: targetCompany }),
          apiRequest('rrhh/empleados', { companyId: targetCompany }),
          apiRequest('rrhh/asistencias', { companyId: targetCompany }),
          apiRequest('usuarios', { companyId: targetCompany }),
        ]);

        setCompanies(compList || []);
        setDashboard(dashData || null);
        setProducts(prodList || []);
        setCategories(catList || []);
        setWarehouses(whList || []);
        setSuppliers(supList || []);
        setMovements(movList || []);
        setSales(saleList || []);
        setEmployees(empList || []);
        setAttendances(attList || []);
        setUsers(usrList || []);
      } catch (err) {
        notify(err.message, true);
      }
    },
    [companyId, notify]
  );

  useEffect(() => {
    loadAllData(companyId);
  }, [companyId, loadAllData]);

  function handleSelectCompany(newCompanyId) {
    setStoredCompanyId(newCompanyId);
    setCompanyId(newCompanyId);
    const found = companies.find((c) => c.id === newCompanyId);
    if (found) {
      notify(`Cambiado a empresa: ${found.nombre}`);
    }
  }

  function handleLoginSuccess(authRes) {
    setStoredSession(authRes);
    setSession(authRes);
    if (authRes.empresa?.id) {
      setStoredCompanyId(authRes.empresa.id);
      setCompanyId(authRes.empresa.id);
    }
  }

  const activeCompanyObj =
    companies.find((c) => c.id === companyId) || {
      id: companyId,
      nombre: dashboard?.empresa_nombre || 'Minimarket Principal',
      plan: 'Enterprise',
      ruc: '20601234567',
    };

  const tabMeta = {
    dashboard: {
      title: 'Dashboard Ejecutivo y Reportes',
      subtitle: `Indicadores en tiempo real de ${activeCompanyObj.nombre}`,
    },
    ventas: {
      title: 'Punto de Venta (POS) y Facturación',
      subtitle: 'Emisión de boletas/facturas con descuento inmediato de inventario',
    },
    catalogo: {
      title: 'Productos y Categorías (Patrón Prototype)',
      subtitle: 'Gestión de catálogo comercial y clonación rápida de productos',
    },
    inventario: {
      title: 'Control de Inventario, Almacenes y Proveedores',
      subtitle: 'Movimientos Kardex (entradas, salidas, ajustes), sucursales y abastecimiento',
    },
    rrhh: {
      title: 'Recursos Humanos (RRHH) y Control de Asistencia',
      subtitle: 'Planilla de trabajadores, turnos horarios, sueldos y marcaciones',
    },
    saas: {
      title: 'Multiempresa SaaS, Usuarios y Arquitectura Cebolla',
      subtitle: 'Administración de empresas (tenants), roles de usuario y patrones Adapter/Prototype',
    },
  };

  return (
    <div className="saas-layout">
      {/* SIDEBAR DE NAVEGACIÓN */}
      <aside className="sidebar">
        <div className="sidebar-brand">
          <div className="brand-logo">M</div>
          <div>
            <div className="brand-title">Minimarket Cloud</div>
            <div className="brand-subtitle">SaaS Multiempresa · Go + React</div>
          </div>
        </div>

        {/* SELECTOR DE EMPRESA (TENANT SWITCHER) */}
        <div className="tenant-box">
          <div className="tenant-label">
            <span>Empresa Activa (Tenant)</span>
            <span className="plan-pill">{activeCompanyObj.plan || 'Pro'}</span>
          </div>
          <select
            className="tenant-select"
            value={companyId}
            onChange={(e) => handleSelectCompany(e.target.value)}
            aria-label="Seleccionar empresa"
          >
            {companies.map((c) => (
              <option key={c.id} value={c.id}>
                {c.nombre}
              </option>
            ))}
          </select>
          <div className="tenant-meta">
            <span>RUC: {activeCompanyObj.ruc || '20601234567'}</span>
            <span>Moneda: PEN (S/)</span>
          </div>
        </div>

        <div className="nav-section-title">Operaciones Comerciales</div>
        <div className="nav-list">
          <button
            className={`nav-btn ${activeTab === 'dashboard' ? 'active' : ''}`}
            onClick={() => setActiveTab('dashboard')}
          >
            <span className="nav-btn-left">📊 Dashboard y Reportes</span>
          </button>
          <button
            className={`nav-btn ${activeTab === 'ventas' ? 'active' : ''}`}
            onClick={() => setActiveTab('ventas')}
          >
            <span className="nav-btn-left">🛒 Punto de Venta (POS)</span>
            <span className="nav-badge">{sales.length}</span>
          </button>
          <button
            className={`nav-btn ${activeTab === 'catalogo' ? 'active' : ''}`}
            onClick={() => setActiveTab('catalogo')}
          >
            <span className="nav-btn-left">📦 Productos y Categorías</span>
            <span className="nav-badge">{products.length}</span>
          </button>
          <button
            className={`nav-btn ${activeTab === 'inventario' ? 'active' : ''}`}
            onClick={() => setActiveTab('inventario')}
          >
            <span className="nav-btn-left">🏭 Inventario y Almacenes</span>
            {dashboard?.alertas_stock_bajo > 0 && (
              <span className="nav-badge warn">{dashboard.alertas_stock_bajo} alerta</span>
            )}
          </button>
        </div>

        <div className="nav-section-title">Organización y SaaS</div>
        <div className="nav-list">
          <button
            className={`nav-btn ${activeTab === 'rrhh' ? 'active' : ''}`}
            onClick={() => setActiveTab('rrhh')}
          >
            <span className="nav-btn-left">👥 Recursos Humanos</span>
            <span className="nav-badge">{employees.length}</span>
          </button>
          <button
            className={`nav-btn ${activeTab === 'saas' ? 'active' : ''}`}
            onClick={() => setActiveTab('saas')}
          >
            <span className="nav-btn-left">🏢 Usuarios y Arquitectura</span>
            <span className="nav-badge">{companies.length} emp</span>
          </button>
        </div>

        <div className="sidebar-footer">
          <div className="user-card">
            <div>
              <div className="user-info-name">{session?.usuario?.nombre || 'Administrador'}</div>
              <div className="user-info-role">Rol: {session?.usuario?.rol || 'admin'}</div>
            </div>
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              onClick={() => setShowAuthModal(true)}
            >
              Cambiar
            </button>
          </div>
        </div>
      </aside>

      {/* ÁREA PRINCIPAL */}
      <div className="main-area">
        <header className="topbar">
          <div className="topbar-left">
            <h1>{tabMeta[activeTab]?.title}</h1>
            <p>{tabMeta[activeTab]?.subtitle}</p>
          </div>
          <div className="topbar-actions">
            <span className="status-pill">
              <span className="status-dot" /> API Go Conectada · {activeCompanyObj.nombre}
            </span>
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              onClick={() => setActiveTab('saas')}
            >
              + Nueva Empresa
            </button>
            <button
              type="button"
              className="btn btn-primary btn-sm"
              onClick={() => setShowAuthModal(true)}
            >
              🔐 {session?.usuario?.rol?.toUpperCase() || 'ADMIN'}
            </button>
          </div>
        </header>

        <main className="content-wrap">
          {activeTab === 'dashboard' && (
            <DashboardView dashboard={dashboard} onNavigate={setActiveTab} />
          )}

          {activeTab === 'ventas' && (
            <PosSalesView
              products={products}
              categories={categories}
              warehouses={warehouses}
              sales={sales}
              currentUser={session?.usuario}
              onRefresh={() => loadAllData(companyId)}
              notify={notify}
            />
          )}

          {activeTab === 'catalogo' && (
            <CatalogView
              products={products}
              categories={categories}
              warehouses={warehouses}
              suppliers={suppliers}
              onRefresh={() => loadAllData(companyId)}
              notify={notify}
            />
          )}

          {activeTab === 'inventario' && (
            <InventoryView
              products={products}
              warehouses={warehouses}
              suppliers={suppliers}
              movements={movements}
              currentUser={session?.usuario}
              onRefresh={() => loadAllData(companyId)}
              notify={notify}
            />
          )}

          {activeTab === 'rrhh' && (
            <HrView
              employees={employees}
              attendances={attendances}
              warehouses={warehouses}
              onRefresh={() => loadAllData(companyId)}
              notify={notify}
            />
          )}

          {activeTab === 'saas' && (
            <ArchitectureSaasView
              companies={companies}
              activeCompanyId={companyId}
              onSelectCompany={handleSelectCompany}
              users={users}
              onRefresh={() => loadAllData(companyId)}
              notify={notify}
            />
          )}
        </main>
      </div>

      {showAuthModal && (
        <AuthModal
          onClose={() => setShowAuthModal(false)}
          onLoginSuccess={handleLoginSuccess}
          notify={notify}
        />
      )}

      {toast && (
        <div className={`toast ${toast.isError ? 'error' : ''}`} role="status">
          {toast.message}
        </div>
      )}
    </div>
  );
}
