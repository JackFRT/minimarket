import React, { useState } from 'react';
import { apiRequest, formatPEN } from '../api';

const emptyEmp = {
  id: '',
  dni: '',
  nombre: '',
  cargo: '',
  departamento: 'Operaciones y Tienda',
  almacen_id: '',
  turno: 'Mañana',
  salario: '',
  fecha_ingreso: '',
  estado: 'activo',
  telefono: '',
  email: '',
};

export default function HrView({
  employees,
  attendances,
  warehouses,
  onRefresh,
  notify,
}) {
  const [empForm, setEmpForm] = useState(emptyEmp);
  const [attForm, setAttForm] = useState({
    empleado_id: '',
    fecha: new Date().toISOString().slice(0, 10),
    hora_entrada: '08:00',
    hora_salida: '17:00',
    estado: 'puntual',
    observacion: '',
  });

  const whMap = Object.fromEntries(warehouses.map((w) => [w.id, w]));

  const activeCount = employees.filter((e) => e.estado === 'activo').length;
  const totalPayroll = employees
    .filter((e) => e.estado === 'activo')
    .reduce((sum, e) => sum + Number(e.salario || 0), 0);

  async function handleSaveEmployee(e) {
    e.preventDefault();
    try {
      const payload = {
        ...empForm,
        salario: Number(empForm.salario || 0),
      };
      if (empForm.id) {
        await apiRequest(`rrhh/empleados/${empForm.id}`, { method: 'PUT', body: payload });
        notify('Ficha del trabajador actualizada.');
      } else {
        await apiRequest('rrhh/empleados', { method: 'POST', body: payload });
        notify('Trabajador registrado en Recursos Humanos.');
      }
      setEmpForm(emptyEmp);
      await onRefresh();
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleCloneEmployee(id) {
    try {
      await apiRequest(`rrhh/empleados/${id}/clonar`, { method: 'POST' });
      await onRefresh();
      notify('Patrón Prototype: Plantilla de puesto/contrato clonada para nuevo ingreso.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleDeleteEmployee(id) {
    if (!window.confirm('¿Retirar este registro de trabajador?')) return;
    try {
      await apiRequest(`rrhh/empleados/${id}`, { method: 'DELETE' });
      await onRefresh();
      notify('Trabajador eliminado.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  async function handleSaveAttendance(e) {
    e.preventDefault();
    try {
      await apiRequest('rrhh/asistencias', {
        method: 'POST',
        body: attForm,
      });
      setAttForm({
        ...attForm,
        empleado_id: '',
        observacion: '',
      });
      await onRefresh();
      notify('Asistencia registrada correctamente.');
    } catch (err) {
      notify(err.message, true);
    }
  }

  return (
    <div>
      {/* KPIs RRHH */}
      <div className="kpi-grid">
        <div className="kpi-card">
          <div className="kpi-top">
            <span>Colaboradores Activos</span>
            <span className="kpi-icon">☺</span>
          </div>
          <div className="kpi-value">{activeCount}</div>
          <div className="kpi-sub">Total en padrón: {employees.length} trabajadores</div>
        </div>
        <div className="kpi-card">
          <div className="kpi-top">
            <span>Planilla Mensual Estimada</span>
            <span className="kpi-icon">S/</span>
          </div>
          <div className="kpi-value">{formatPEN(totalPayroll)}</div>
          <div className="kpi-sub">Remuneraciones brutas de personal activo</div>
        </div>
        <div className="kpi-card">
          <div className="kpi-top">
            <span>Marcaciones de Asistencia</span>
            <span className="kpi-icon">✓</span>
          </div>
          <div className="kpi-value">{attendances.length}</div>
          <div className="kpi-sub">Registros de ingreso y salida en turno</div>
        </div>
      </div>

      {/* GESTIÓN DE PERSONAL Y PLANILLA */}
      <div className="panel">
        <div className="panel-header">
          <div className="panel-title">
            <h2>Sistema de Recursos Humanos — Gestión de Trabajadores y Turnos</h2>
            <p>
              Control de contratos, cargos, sueldos, asignación de almacén y clonación de perfiles con <code>Employee.Clone()</code>
            </p>
          </div>
        </div>

        <form className="form-grid" onSubmit={handleSaveEmployee}>
          <div className="field">
            <label>DNI / Doc.</label>
            <input
              placeholder="Ej. 45891203"
              value={empForm.dni}
              onChange={(e) => setEmpForm({ ...empForm, dni: e.target.value })}
            />
          </div>
          <div className="field" style={{ gridColumn: 'span 2' }}>
            <label>Apellidos y Nombres *</label>
            <input
              required
              placeholder="Ej. Valeria Paredes Gómez"
              value={empForm.nombre}
              onChange={(e) => setEmpForm({ ...empForm, nombre: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Cargo / Puesto *</label>
            <input
              required
              placeholder="Ej. Cajera POS / Almacenero"
              value={empForm.cargo}
              onChange={(e) => setEmpForm({ ...empForm, cargo: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Departamento</label>
            <select
              value={empForm.departamento}
              onChange={(e) => setEmpForm({ ...empForm, departamento: e.target.value })}
            >
              <option value="Operaciones y Tienda">Operaciones y Tienda</option>
              <option value="Caja y Atención al Cliente">Caja y Atención al Cliente</option>
              <option value="Logística e Inventario">Logística e Inventario</option>
              <option value="Recursos Humanos">Recursos Humanos</option>
              <option value="Gerencia y Operaciones">Gerencia y Operaciones</option>
            </select>
          </div>
          <div className="field">
            <label>Sede / Almacén Asignado</label>
            <select
              value={empForm.almacen_id}
              onChange={(e) => setEmpForm({ ...empForm, almacen_id: e.target.value })}
            >
              <option value="">-- Sede General --</option>
              {warehouses.map((w) => (
                <option key={w.id} value={w.id}>
                  {w.nombre}
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label>Turno Horario</label>
            <select
              value={empForm.turno}
              onChange={(e) => setEmpForm({ ...empForm, turno: e.target.value })}
            >
              <option value="Mañana">Mañana (07:30 - 16:00)</option>
              <option value="Tarde">Tarde (14:00 - 22:30)</option>
              <option value="Completo">Completo (08:00 - 17:30)</option>
              <option value="Noche">Noche / Inventario</option>
            </select>
          </div>
          <div className="field">
            <label>Salario Mensual (S/) *</label>
            <input
              required
              type="number"
              step="0.01"
              min="0"
              placeholder="1650.00"
              value={empForm.salario}
              onChange={(e) => setEmpForm({ ...empForm, salario: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Estado Laboral</label>
            <select
              value={empForm.estado}
              onChange={(e) => setEmpForm({ ...empForm, estado: e.target.value })}
            >
              <option value="activo">Activo en Planilla</option>
              <option value="vacaciones">En Vacaciones</option>
              <option value="licencia">Con Licencia Médica</option>
              <option value="cesado">Cesado</option>
            </select>
          </div>
          <div className="field">
            <label>Teléfono</label>
            <input
              placeholder="999-888-777"
              value={empForm.telefono}
              onChange={(e) => setEmpForm({ ...empForm, telefono: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Correo Corporativo</label>
            <input
              type="email"
              placeholder="trabajador@minimarket.pe"
              value={empForm.email}
              onChange={(e) => setEmpForm({ ...empForm, email: e.target.value })}
            />
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
            <button type="submit" className="btn btn-primary">
              {empForm.id ? 'Actualizar Ficha' : 'Registrar Trabajador'}
            </button>
            {empForm.id && (
              <button
                type="button"
                className="btn btn-secondary"
                onClick={() => setEmpForm(emptyEmp)}
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
                <th>DNI</th>
                <th>Trabajador</th>
                <th>Cargo / Área</th>
                <th>Almacén / Turno</th>
                <th>Salario</th>
                <th>Estado</th>
                <th>Acciones (CRUD + Prototype)</th>
              </tr>
            </thead>
            <tbody>
              {employees.map((emp) => {
                const wh = whMap[emp.almacen_id];
                return (
                  <tr key={emp.id}>
                    <td><code>{emp.dni || 'S/D'}</code></td>
                    <td>
                      <strong>{emp.nombre}</strong>
                      <div style={{ fontSize: 11, color: '#68776e' }}>
                        {emp.email} {emp.telefono ? `· ${emp.telefono}` : ''}
                      </div>
                    </td>
                    <td>
                      <strong>{emp.cargo}</strong>
                      <div style={{ fontSize: 11, color: '#68776e' }}>{emp.departamento}</div>
                    </td>
                    <td>
                      <div>{wh ? wh.nombre : 'Sede Principal'}</div>
                      <span className="badge badge-blue">Turno {emp.turno}</span>
                    </td>
                    <td><strong style={{ color: '#245e3f' }}>{formatPEN(emp.salario)}</strong></td>
                    <td>
                      <span
                        className={`badge ${
                          emp.estado === 'activo'
                            ? 'badge-green'
                            : emp.estado === 'vacaciones'
                            ? 'badge-amber'
                            : 'badge-red'
                        }`}
                      >
                        {emp.estado}
                      </span>
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                        <button
                          type="button"
                          className="btn btn-secondary btn-sm"
                          onClick={() => setEmpForm(emp)}
                        >
                          Editar
                        </button>
                        <button
                          type="button"
                          className="btn btn-prototype btn-sm"
                          title="Clonar puesto, salario y turno para nueva contratación"
                          onClick={() => handleCloneEmployee(emp.id)}
                        >
                          ⧉ Clonar Puesto
                        </button>
                        <button
                          type="button"
                          className="btn btn-danger btn-sm"
                          onClick={() => handleDeleteEmployee(emp.id)}
                        >
                          Retirar
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

      {/* CONTROL DE ASISTENCIA DIARIA */}
      <div className="panel">
        <div className="panel-header">
          <div className="panel-title">
            <h2>Control de Asistencia y Marcación de Turnos</h2>
            <p>Registro de hora de entrada, salida y puntualidad de los trabajadores</p>
          </div>
        </div>

        <form className="form-grid" onSubmit={handleSaveAttendance}>
          <div className="field" style={{ gridColumn: 'span 2' }}>
            <label>Trabajador *</label>
            <select
              required
              value={attForm.empleado_id}
              onChange={(e) => setAttForm({ ...attForm, empleado_id: e.target.value })}
            >
              <option value="">-- Seleccionar colaborador --</option>
              {employees.map((emp) => (
                <option key={emp.id} value={emp.id}>
                  {emp.nombre} — {emp.cargo} ({emp.turno})
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label>Fecha</label>
            <input
              type="date"
              value={attForm.fecha}
              onChange={(e) => setAttForm({ ...attForm, fecha: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Hora Entrada</label>
            <input
              type="time"
              value={attForm.hora_entrada}
              onChange={(e) => setAttForm({ ...attForm, hora_entrada: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Hora Salida</label>
            <input
              type="time"
              value={attForm.hora_salida}
              onChange={(e) => setAttForm({ ...attForm, hora_salida: e.target.value })}
            />
          </div>
          <div className="field">
            <label>Condición</label>
            <select
              value={attForm.estado}
              onChange={(e) => setAttForm({ ...attForm, estado: e.target.value })}
            >
              <option value="puntual">Puntual</option>
              <option value="tardanza">Tardanza</option>
              <option value="justificado">Permiso Justificado</option>
              <option value="falta">Falta Injustificada</option>
            </select>
          </div>
          <div className="field" style={{ gridColumn: 'span 2' }}>
            <label>Observación de Turno</label>
            <input
              placeholder="Ej. Apertura de tienda / Cierre de caja conforme"
              value={attForm.observacion}
              onChange={(e) => setAttForm({ ...attForm, observacion: e.target.value })}
            />
          </div>
          <div>
            <button type="submit" className="btn btn-primary">
              Registrar Asistencia
            </button>
          </div>
        </form>

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Fecha</th>
                <th>Colaborador</th>
                <th>Hora Entrada</th>
                <th>Hora Salida</th>
                <th>Estado</th>
                <th>Observación</th>
              </tr>
            </thead>
            <tbody>
              {[...attendances].reverse().map((a) => (
                <tr key={a.id}>
                  <td>{a.fecha}</td>
                  <td><strong>{a.empleado_nombre}</strong></td>
                  <td><code>{a.hora_entrada}</code></td>
                  <td><code>{a.hora_salida || 'En turno'}</code></td>
                  <td>
                    <span
                      className={`badge ${
                        a.estado === 'puntual'
                          ? 'badge-green'
                          : a.estado === 'tardanza'
                          ? 'badge-amber'
                          : 'badge-red'
                      }`}
                    >
                      {a.estado.toUpperCase()}
                    </span>
                  </td>
                  <td>{a.observacion || '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
