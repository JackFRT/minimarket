import React, { useState } from 'react';
import { apiRequest } from '../api';

export default function AuthModal({ onClose, onLoginSuccess, notify }) {
  const [email, setEmail] = useState('admin@minimarket.pe');
  const [password, setPassword] = useState('admin123');
  const [loading, setLoading] = useState(false);

  async function submitLogin(targetEmail, targetPassword) {
    setLoading(true);
    try {
      const res = await apiRequest('auth/login', {
        method: 'POST',
        body: { email: targetEmail, password: targetPassword },
      });
      onLoginSuccess(res);
      notify(`Bienvenido(a), ${res.usuario.nombre} (${res.empresa.nombre})`);
      onClose();
    } catch (err) {
      notify(err.message, true);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 14 }}>
          <div>
            <span className="badge badge-green">SaaS Multiempresa Auth</span>
            <h3 style={{ margin: '6px 0 0', fontSize: 18 }}>Iniciar Sesión o Cambiar Rol</h3>
          </div>
          <button type="button" className="btn btn-secondary btn-sm" onClick={onClose}>
            ✕
          </button>
        </div>

        <form
          onSubmit={(e) => {
            e.preventDefault();
            submitLogin(email, password);
          }}
          style={{ display: 'flex', flexDirection: 'column', gap: 12, marginBottom: 18 }}
        >
          <div className="field">
            <label>Correo Corporativo</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div className="field">
            <label>Contraseña</label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Verificando credenciales...' : 'Ingresar al Sistema'}
          </button>
        </form>

        <div style={{ borderTop: '1px solid #e2e8e3', paddingTop: 14 }}>
          <div style={{ fontSize: 11, fontWeight: 700, textTransform: 'uppercase', color: '#68776e', marginBottom: 8 }}>
            Accesos Rápidos de Demostración (1 Clic)
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              onClick={() => submitLogin('admin@minimarket.pe', 'admin123')}
            >
              👑 Admin (Sede Central)
            </button>
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              onClick={() => submitLogin('caja@minimarket.pe', 'caja123')}
            >
              🛒 Cajera POS (Sede Central)
            </button>
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              onClick={() => submitLogin('rrhh@minimarket.pe', 'rrhh123')}
            >
              👥 Jefa RRHH (Sede Central)
            </button>
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              onClick={() => submitLogin('admin@marketnorte.pe', 'admin123')}
            >
              🏢 Admin (Market Los Olivos)
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
