package presentation

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"minimarket/internal/application"
	"minimarket/internal/domain"
	"net/http"
	"strings"
)

func respond(w http.ResponseWriter, value any, err error, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		status = 400
		if errors.Is(err, domain.ErrNotFound) {
			status = 404
		} else if errors.Is(err, domain.ErrUnauthorized) {
			status = 401
		}
		value = map[string]string{"error": err.Error()}
	}
	w.WriteHeader(status)
	if status != 204 {
		_ = json.NewEncoder(w).Encode(value)
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return errors.New("JSON inválido: " + err.Error())
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return errors.New("envía un único objeto JSON")
	}
	return nil
}

func companyFromRequest(r *http.Request) string {
	if cid := strings.TrimSpace(r.Header.Get("X-Empresa-ID")); cid != "" {
		return cid
	}
	if cid := strings.TrimSpace(r.URL.Query().Get("empresa_id")); cid != "" {
		return cid
	}
	return ""
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Empresa-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewHandler(s application.Service) http.Handler {
	m := http.NewServeMux()

	// Archivos estáticos del build de React embebidos
	subWeb, err := fs.Sub(webFS, "web")
	var fileServer http.Handler
	if err == nil {
		fileServer = http.FileServer(http.FS(subWeb))
	}

	m.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if b, err := webFS.ReadFile("web/index.html"); err == nil {
			_, _ = w.Write(b)
			return
		}
		_, _ = io.WriteString(w, home)
	})

	if fileServer != nil {
		m.Handle("GET /assets/", fileServer)
	}

	// --- SAAS MULTIEMPRESA ---
	m.HandleFunc("GET /api/empresas", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Companies()
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/empresas", func(w http.ResponseWriter, r *http.Request) {
		var c domain.Company
		if e := decode(w, r, &c); e != nil {
			respond(w, nil, e, 400)
			return
		}
		v, e := s.SaveCompany(c, "")
		respond(w, v, e, 201)
	})
	m.HandleFunc("PUT /api/empresas/{id}", func(w http.ResponseWriter, r *http.Request) {
		var c domain.Company
		if e := decode(w, r, &c); e != nil {
			respond(w, nil, e, 400)
			return
		}
		v, e := s.SaveCompany(c, r.PathValue("id"))
		respond(w, v, e, 200)
	})

	// --- AUTENTICACIÓN Y USUARIOS ---
	m.HandleFunc("POST /api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req application.LoginRequest
		if e := decode(w, r, &req); e != nil {
			respond(w, nil, e, 400)
			return
		}
		v, e := s.Authenticate(req)
		respond(w, v, e, 200)
	})
	m.HandleFunc("GET /api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Users(companyFromRequest(r))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		var u domain.User
		if e := decode(w, r, &u); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if u.CompanyID == "" {
			u.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveUser(u, "")
		respond(w, v, e, 201)
	})
	m.HandleFunc("PUT /api/usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		var u domain.User
		if e := decode(w, r, &u); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if u.CompanyID == "" {
			u.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveUser(u, r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("DELETE /api/usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, nil, s.Repo.DeleteUser(r.PathValue("id")), 204)
	})

	// --- DASHBOARD Y REPORTES ---
	m.HandleFunc("GET /api/dashboard", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.GetDashboardSummary(companyFromRequest(r))
		respond(w, v, e, 200)
	})

	// --- CATEGORÍAS ---
	m.HandleFunc("GET /api/categorias", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Categories(companyFromRequest(r))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/categorias", func(w http.ResponseWriter, r *http.Request) {
		var c domain.Category
		if e := decode(w, r, &c); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if c.CompanyID == "" {
			c.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveCategory(c, "")
		respond(w, v, e, 201)
	})
	m.HandleFunc("PUT /api/categorias/{id}", func(w http.ResponseWriter, r *http.Request) {
		var c domain.Category
		if e := decode(w, r, &c); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if c.CompanyID == "" {
			c.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveCategory(c, r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("DELETE /api/categorias/{id}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, nil, s.Repo.DeleteCategory(r.PathValue("id")), 204)
	})
	m.HandleFunc("POST /api/categorias/{id}/clonar", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.CloneCategory(r.PathValue("id"))
		respond(w, v, e, 201)
	})

	// --- ALMACENES ---
	m.HandleFunc("GET /api/almacenes", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Warehouses(companyFromRequest(r))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/almacenes", func(w http.ResponseWriter, r *http.Request) {
		var wh domain.Warehouse
		if e := decode(w, r, &wh); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if wh.CompanyID == "" {
			wh.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveWarehouse(wh, "")
		respond(w, v, e, 201)
	})
	m.HandleFunc("PUT /api/almacenes/{id}", func(w http.ResponseWriter, r *http.Request) {
		var wh domain.Warehouse
		if e := decode(w, r, &wh); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if wh.CompanyID == "" {
			wh.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveWarehouse(wh, r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("DELETE /api/almacenes/{id}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, nil, s.Repo.DeleteWarehouse(r.PathValue("id")), 204)
	})

	// --- PRODUCTOS Y PROTOTYPE ---
	m.HandleFunc("GET /api/productos", func(w http.ResponseWriter, r *http.Request) {
		cid := companyFromRequest(r)
		if cid == "" {
			v, e := s.Repo.Products()
			respond(w, v, e, 200)
			return
		}
		v, e := s.Repo.ProductsByCompany(cid)
		respond(w, v, e, 200)
	})
	m.HandleFunc("GET /api/productos/{id}", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Product(r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/productos", func(w http.ResponseWriter, r *http.Request) {
		var p domain.Product
		if e := decode(w, r, &p); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if p.CompanyID == "" {
			p.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveProduct(p, "")
		respond(w, v, e, 201)
	})
	m.HandleFunc("PUT /api/productos/{id}", func(w http.ResponseWriter, r *http.Request) {
		var p domain.Product
		if e := decode(w, r, &p); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if p.CompanyID == "" {
			p.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveProduct(p, r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("DELETE /api/productos/{id}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, nil, s.Repo.DeleteProduct(r.PathValue("id")), 204)
	})
	m.HandleFunc("POST /api/productos/{id}/clonar", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.CloneProduct(r.PathValue("id"))
		respond(w, v, e, 201)
	})

	// --- PROVEEDORES ---
	m.HandleFunc("GET /api/proveedores", func(w http.ResponseWriter, r *http.Request) {
		cid := companyFromRequest(r)
		if cid == "" {
			v, e := s.Repo.Suppliers()
			respond(w, v, e, 200)
			return
		}
		v, e := s.Repo.SuppliersByCompany(cid)
		respond(w, v, e, 200)
	})
	m.HandleFunc("GET /api/proveedores/{id}", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Supplier(r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/proveedores", func(w http.ResponseWriter, r *http.Request) {
		var p domain.Supplier
		if e := decode(w, r, &p); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if p.CompanyID == "" {
			p.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveSupplier(p, "")
		respond(w, v, e, 201)
	})
	m.HandleFunc("PUT /api/proveedores/{id}", func(w http.ResponseWriter, r *http.Request) {
		var p domain.Supplier
		if e := decode(w, r, &p); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if p.CompanyID == "" {
			p.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveSupplier(p, r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("DELETE /api/proveedores/{id}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, nil, s.Repo.DeleteSupplier(r.PathValue("id")), 204)
	})

	// --- CONTROL DE INVENTARIO (KARDEX) ---
	m.HandleFunc("GET /api/inventario/movimientos", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Movements(companyFromRequest(r))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/inventario/movimientos", func(w http.ResponseWriter, r *http.Request) {
		var mov domain.InventoryMovement
		if e := decode(w, r, &mov); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if mov.CompanyID == "" {
			mov.CompanyID = companyFromRequest(r)
		}
		v, e := s.RegisterInventoryMovement(mov)
		respond(w, v, e, 201)
	})

	// --- VENTAS (POS) ---
	m.HandleFunc("GET /api/ventas", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Sales(companyFromRequest(r))
		respond(w, v, e, 200)
	})
	m.HandleFunc("GET /api/ventas/{id}", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Sale(r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/ventas", func(w http.ResponseWriter, r *http.Request) {
		var sale domain.Sale
		if e := decode(w, r, &sale); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if sale.CompanyID == "" {
			sale.CompanyID = companyFromRequest(r)
		}
		v, e := s.RegisterSale(sale)
		respond(w, v, e, 201)
	})

	// --- RECURSOS HUMANOS (EMPLEADOS Y ASISTENCIA) ---
	m.HandleFunc("GET /api/rrhh/empleados", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Employees(companyFromRequest(r))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/rrhh/empleados", func(w http.ResponseWriter, r *http.Request) {
		var emp domain.Employee
		if e := decode(w, r, &emp); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if emp.CompanyID == "" {
			emp.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveEmployee(emp, "")
		respond(w, v, e, 201)
	})
	m.HandleFunc("PUT /api/rrhh/empleados/{id}", func(w http.ResponseWriter, r *http.Request) {
		var emp domain.Employee
		if e := decode(w, r, &emp); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if emp.CompanyID == "" {
			emp.CompanyID = companyFromRequest(r)
		}
		v, e := s.SaveEmployee(emp, r.PathValue("id"))
		respond(w, v, e, 200)
	})
	m.HandleFunc("DELETE /api/rrhh/empleados/{id}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, nil, s.Repo.DeleteEmployee(r.PathValue("id")), 204)
	})
	m.HandleFunc("POST /api/rrhh/empleados/{id}/clonar", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.CloneEmployee(r.PathValue("id"))
		respond(w, v, e, 201)
	})

	m.HandleFunc("GET /api/rrhh/asistencias", func(w http.ResponseWriter, r *http.Request) {
		v, e := s.Repo.Attendances(companyFromRequest(r))
		respond(w, v, e, 200)
	})
	m.HandleFunc("POST /api/rrhh/asistencias", func(w http.ResponseWriter, r *http.Request) {
		var att domain.Attendance
		if e := decode(w, r, &att); e != nil {
			respond(w, nil, e, 400)
			return
		}
		if att.CompanyID == "" {
			att.CompanyID = companyFromRequest(r)
		}
		v, e := s.RegisterAttendance(att)
		respond(w, v, e, 201)
	})

	return withCORS(m)
}
