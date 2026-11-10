const STORAGE_TENANT_KEY = 'minimarket_active_company';
const STORAGE_AUTH_KEY = 'minimarket_auth_session';

export function getStoredCompanyId() {
  return localStorage.getItem(STORAGE_TENANT_KEY) || 'emp-default';
}

export function setStoredCompanyId(companyId) {
  localStorage.setItem(STORAGE_TENANT_KEY, companyId);
}

export function getStoredSession() {
  try {
    const raw = localStorage.getItem(STORAGE_AUTH_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function setStoredSession(session) {
  if (!session) {
    localStorage.removeItem(STORAGE_AUTH_KEY);
    return;
  }
  localStorage.setItem(STORAGE_AUTH_KEY, JSON.stringify(session));
  if (session.empresa?.id) {
    setStoredCompanyId(session.empresa.id);
  }
}

export async function apiRequest(path, { method = 'GET', body, companyId } = {}) {
  const activeCompany = companyId || getStoredCompanyId();
  const session = getStoredSession();

  const headers = {
    'Content-Type': 'application/json',
    'X-Empresa-ID': activeCompany,
  };
  if (session?.token) {
    headers['Authorization'] = `Bearer ${session.token}`;
  }

  const response = await fetch(`/api/${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (response.status === 204) {
    return null;
  }

  let data = null;
  try {
    data = await response.json();
  } catch {
    throw new Error('Respuesta inválida del servidor');
  }

  if (!response.ok) {
    throw new Error(data?.error || `Error HTTP ${response.status}`);
  }

  return data;
}

export const formatPEN = (value) =>
  new Intl.NumberFormat('es-PE', {
    style: 'currency',
    currency: 'PEN',
    minimumFractionDigits: 2,
  }).format(Number(value || 0));
