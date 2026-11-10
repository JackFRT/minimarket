package infrastructure

import (
	"encoding/json"
	"minimarket/internal/application"
	"minimarket/internal/domain"
	"os"
	"strings"
	"sync"
)

type state struct {
	Companies   []domain.Company           `json:"empresas"`
	Categories  []domain.Category          `json:"categorias"`
	Warehouses  []domain.Warehouse         `json:"almacenes"`
	Products    []domain.Product           `json:"productos"`
	Suppliers   []domain.Supplier          `json:"proveedores"`
	Movements   []domain.InventoryMovement `json:"movimientos"`
	Sales       []domain.Sale              `json:"ventas"`
	Users       []domain.User              `json:"usuarios"`
	Employees   []domain.Employee          `json:"empleados"`
	Attendances []domain.Attendance        `json:"asistencias"`
}

// JSONRepository implementa el patrón Adapter entre la interfaz application.Repository y el almacenamiento en archivos JSON.
type JSONRepository struct {
	mu   sync.Mutex
	path string
	data state
}

var _ application.Repository = (*JSONRepository)(nil)

func emptyState() state {
	return state{
		Companies:   []domain.Company{},
		Categories:  []domain.Category{},
		Warehouses:  []domain.Warehouse{},
		Products:    []domain.Product{},
		Suppliers:   []domain.Supplier{},
		Movements:   []domain.InventoryMovement{},
		Sales:       []domain.Sale{},
		Users:       []domain.User{},
		Employees:   []domain.Employee{},
		Attendances: []domain.Attendance{},
	}
}

func NewJSONRepository(path string) (*JSONRepository, error) {
	r := &JSONRepository{path: path, data: emptyState()}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		r.ensureBaseTenants()
		return r, nil
	}
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &r.data); err != nil {
		return nil, err
	}
	r.normalizeLegacyRecords()
	r.ensureBaseTenants()
	return r, nil
}

// normalizeLegacyRecords asigna la empresa por defecto a registros previos que no tenían empresa_id.
func (r *JSONRepository) normalizeLegacyRecords() {
	if r.data.Products == nil {
		r.data.Products = []domain.Product{}
	}
	if r.data.Suppliers == nil {
		r.data.Suppliers = []domain.Supplier{}
	}
	if r.data.Categories == nil {
		r.data.Categories = []domain.Category{}
	}
	if r.data.Warehouses == nil {
		r.data.Warehouses = []domain.Warehouse{}
	}
	if r.data.Movements == nil {
		r.data.Movements = []domain.InventoryMovement{}
	}
	if r.data.Sales == nil {
		r.data.Sales = []domain.Sale{}
	}
	if r.data.Users == nil {
		r.data.Users = []domain.User{}
	}
	if r.data.Employees == nil {
		r.data.Employees = []domain.Employee{}
	}
	if r.data.Attendances == nil {
		r.data.Attendances = []domain.Attendance{}
	}
	for i := range r.data.Products {
		if r.data.Products[i].CompanyID == "" {
			r.data.Products[i].CompanyID = domain.DefaultCompanyID
		}
		if r.data.Products[i].Unit == "" {
			r.data.Products[i].Unit = "und"
		}
		if r.data.Products[i].MinStock == 0 {
			r.data.Products[i].MinStock = 5
		}
	}
	for i := range r.data.Suppliers {
		if r.data.Suppliers[i].CompanyID == "" {
			r.data.Suppliers[i].CompanyID = domain.DefaultCompanyID
		}
	}
}

func (r *JSONRepository) ensureBaseTenants() {
	if len(r.data.Companies) == 0 {
		r.data.Companies = []domain.Company{
			{
				ID:        domain.DefaultCompanyID,
				Name:      "Minimarket El Sol - Sede Central",
				RUC:       "20601234567",
				Plan:      "Enterprise",
				Address:   "Av. Larco 840, Miraflores, Lima",
				Phone:     "01-445-8920",
				Currency:  "PEN",
				Active:    true,
				CreatedAt: "2026-01-15",
			},
			{
				ID:        "emp-norte",
				Name:      "Market Express Los Olivos",
				RUC:       "20609876543",
				Plan:      "Pro",
				Address:   "Av. Antúnez de Mayolo 1120, Los Olivos",
				Phone:     "01-521-3340",
				Currency:  "PEN",
				Active:    true,
				CreatedAt: "2026-03-10",
			},
		}
	}
	if len(r.data.Users) == 0 {
		r.data.Users = []domain.User{
			{
				ID:        "usr-admin-1",
				CompanyID: domain.DefaultCompanyID,
				Name:      "Carlos Mendoza (Admin)",
				Email:     "admin@minimarket.pe",
				Password:  "admin123",
				Role:      "admin",
				Active:    true,
				CreatedAt: "2026-01-15",
			},
			{
				ID:        "usr-cajero-1",
				CompanyID: domain.DefaultCompanyID,
				Name:      "Lucía Fernández (Cajera)",
				Email:     "caja@minimarket.pe",
				Password:  "caja123",
				Role:      "cajero",
				Active:    true,
				CreatedAt: "2026-02-01",
			},
			{
				ID:        "usr-rrhh-1",
				CompanyID: domain.DefaultCompanyID,
				Name:      "Elena Vargas (RRHH)",
				Email:     "rrhh@minimarket.pe",
				Password:  "rrhh123",
				Role:      "rrhh",
				Active:    true,
				CreatedAt: "2026-02-10",
			},
			{
				ID:        "usr-norte-1",
				CompanyID: "emp-norte",
				Name:      "Roberto Gómez (Gerente Norte)",
				Email:     "admin@marketnorte.pe",
				Password:  "admin123",
				Role:      "admin",
				Active:    true,
				CreatedAt: "2026-03-10",
			},
		}
	}
}

func (r *JSONRepository) persist(next state) error {
	b, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(r.path, b, 0600); err != nil {
		return err
	}
	r.data = next
	return nil
}

func (r *JSONRepository) cloneState() state {
	next := r.data
	next.Companies = append([]domain.Company{}, r.data.Companies...)
	next.Categories = append([]domain.Category{}, r.data.Categories...)
	next.Warehouses = append([]domain.Warehouse{}, r.data.Warehouses...)
	next.Products = append([]domain.Product{}, r.data.Products...)
	next.Suppliers = append([]domain.Supplier{}, r.data.Suppliers...)
	next.Movements = append([]domain.InventoryMovement{}, r.data.Movements...)
	next.Sales = append([]domain.Sale{}, r.data.Sales...)
	next.Users = append([]domain.User{}, r.data.Users...)
	next.Employees = append([]domain.Employee{}, r.data.Employees...)
	next.Attendances = append([]domain.Attendance{}, r.data.Attendances...)
	return next
}

func matchesCompany(recordCompanyID, filterCompanyID string) bool {
	if filterCompanyID == "" || filterCompanyID == "*" {
		return true
	}
	if recordCompanyID == "" {
		recordCompanyID = domain.DefaultCompanyID
	}
	return recordCompanyID == filterCompanyID
}

// --- EMPRESAS ---

func (r *JSONRepository) Companies() ([]domain.Company, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]domain.Company{}, r.data.Companies...), nil
}

func (r *JSONRepository) Company(id string) (domain.Company, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.data.Companies {
		if c.ID == id {
			return c, nil
		}
	}
	return domain.Company{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveCompany(c domain.Company, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Companies = append(next.Companies, c)
		return r.persist(next)
	}
	for i, old := range next.Companies {
		if old.ID == c.ID {
			next.Companies[i] = c
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

// --- CATEGORÍAS ---

func (r *JSONRepository) Categories(companyID string) ([]domain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Category{}
	for _, c := range r.data.Categories {
		if matchesCompany(c.CompanyID, companyID) {
			out = append(out, c)
		}
	}
	return out, nil
}

func (r *JSONRepository) Category(id string) (domain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.data.Categories {
		if c.ID == id {
			return c, nil
		}
	}
	return domain.Category{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveCategory(c domain.Category, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Categories = append(next.Categories, c)
		return r.persist(next)
	}
	for i, old := range next.Categories {
		if old.ID == c.ID {
			next.Categories[i] = c
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

func (r *JSONRepository) DeleteCategory(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	next.Categories = []domain.Category{}
	found := false
	for _, c := range r.data.Categories {
		if c.ID == id {
			found = true
		} else {
			next.Categories = append(next.Categories, c)
		}
	}
	if !found {
		return domain.ErrNotFound
	}
	return r.persist(next)
}

// --- ALMACENES ---

func (r *JSONRepository) Warehouses(companyID string) ([]domain.Warehouse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Warehouse{}
	for _, w := range r.data.Warehouses {
		if matchesCompany(w.CompanyID, companyID) {
			out = append(out, w)
		}
	}
	return out, nil
}

func (r *JSONRepository) Warehouse(id string) (domain.Warehouse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, w := range r.data.Warehouses {
		if w.ID == id {
			return w, nil
		}
	}
	return domain.Warehouse{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveWarehouse(w domain.Warehouse, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Warehouses = append(next.Warehouses, w)
		return r.persist(next)
	}
	for i, old := range next.Warehouses {
		if old.ID == w.ID {
			next.Warehouses[i] = w
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

func (r *JSONRepository) DeleteWarehouse(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	next.Warehouses = []domain.Warehouse{}
	found := false
	for _, w := range r.data.Warehouses {
		if w.ID == id {
			found = true
		} else {
			next.Warehouses = append(next.Warehouses, w)
		}
	}
	if !found {
		return domain.ErrNotFound
	}
	return r.persist(next)
}

// --- PRODUCTOS ---

func (r *JSONRepository) Products() ([]domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]domain.Product{}, r.data.Products...), nil
}

func (r *JSONRepository) ProductsByCompany(companyID string) ([]domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Product{}
	for _, p := range r.data.Products {
		if matchesCompany(p.CompanyID, companyID) {
			out = append(out, p)
		}
	}
	return out, nil
}

func (r *JSONRepository) Product(id string) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.data.Products {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Product{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveProduct(p domain.Product, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Products = append(next.Products, p)
		return r.persist(next)
	}
	for i, old := range next.Products {
		if old.ID == p.ID {
			if p.CompanyID == "" {
				p.CompanyID = old.CompanyID
			}
			next.Products[i] = p
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

func (r *JSONRepository) DeleteProduct(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	next.Products = []domain.Product{}
	found := false
	for _, p := range r.data.Products {
		if p.ID == id {
			found = true
		} else {
			next.Products = append(next.Products, p)
		}
	}
	if !found {
		return domain.ErrNotFound
	}
	return r.persist(next)
}

// --- PROVEEDORES ---

func (r *JSONRepository) Suppliers() ([]domain.Supplier, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]domain.Supplier{}, r.data.Suppliers...), nil
}

func (r *JSONRepository) SuppliersByCompany(companyID string) ([]domain.Supplier, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Supplier{}
	for _, s := range r.data.Suppliers {
		if matchesCompany(s.CompanyID, companyID) {
			out = append(out, s)
		}
	}
	return out, nil
}

func (r *JSONRepository) Supplier(id string) (domain.Supplier, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.data.Suppliers {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Supplier{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveSupplier(p domain.Supplier, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Suppliers = append(next.Suppliers, p)
		return r.persist(next)
	}
	for i, old := range next.Suppliers {
		if old.ID == p.ID {
			if p.CompanyID == "" {
				p.CompanyID = old.CompanyID
			}
			next.Suppliers[i] = p
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

func (r *JSONRepository) DeleteSupplier(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	next.Suppliers = []domain.Supplier{}
	found := false
	for _, p := range r.data.Suppliers {
		if p.ID == id {
			found = true
		} else {
			next.Suppliers = append(next.Suppliers, p)
		}
	}
	if !found {
		return domain.ErrNotFound
	}
	return r.persist(next)
}

// --- MOVIMIENTOS DE INVENTARIO (KARDEX) ---

func (r *JSONRepository) Movements(companyID string) ([]domain.InventoryMovement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.InventoryMovement{}
	for _, m := range r.data.Movements {
		if matchesCompany(m.CompanyID, companyID) {
			out = append(out, m)
		}
	}
	return out, nil
}

func (r *JSONRepository) SaveMovement(m domain.InventoryMovement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	next.Movements = append(next.Movements, m)
	return r.persist(next)
}

// --- VENTAS ---

func (r *JSONRepository) Sales(companyID string) ([]domain.Sale, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Sale{}
	for _, s := range r.data.Sales {
		if matchesCompany(s.CompanyID, companyID) {
			out = append(out, s)
		}
	}
	return out, nil
}

func (r *JSONRepository) Sale(id string) (domain.Sale, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, s := range r.data.Sales {
		if s.ID == id {
			return s, nil
		}
	}
	return domain.Sale{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveSale(s domain.Sale, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Sales = append(next.Sales, s)
		return r.persist(next)
	}
	for i, old := range next.Sales {
		if old.ID == s.ID {
			next.Sales[i] = s
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

// --- USUARIOS ---

func (r *JSONRepository) Users(companyID string) ([]domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.User{}
	for _, u := range r.data.Users {
		if matchesCompany(u.CompanyID, companyID) {
			uCopy := u
			uCopy.Password = ""
			out = append(out, uCopy)
		}
	}
	return out, nil
}

func (r *JSONRepository) UserByEmail(email string, companyID string) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.data.Users {
		if strings.EqualFold(u.Email, email) && (companyID == "" || u.CompanyID == companyID) {
			return u, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (r *JSONRepository) User(id string) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.data.Users {
		if u.ID == id {
			return u, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveUser(u domain.User, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Users = append(next.Users, u)
		return r.persist(next)
	}
	for i, old := range next.Users {
		if old.ID == u.ID {
			next.Users[i] = u
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

func (r *JSONRepository) DeleteUser(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	next.Users = []domain.User{}
	found := false
	for _, u := range r.data.Users {
		if u.ID == id {
			found = true
		} else {
			next.Users = append(next.Users, u)
		}
	}
	if !found {
		return domain.ErrNotFound
	}
	return r.persist(next)
}

// --- RECURSOS HUMANOS (EMPLEADOS Y ASISTENCIA) ---

func (r *JSONRepository) Employees(companyID string) ([]domain.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Employee{}
	for _, e := range r.data.Employees {
		if matchesCompany(e.CompanyID, companyID) {
			out = append(out, e)
		}
	}
	return out, nil
}

func (r *JSONRepository) Employee(id string) (domain.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.data.Employees {
		if e.ID == id {
			return e, nil
		}
	}
	return domain.Employee{}, domain.ErrNotFound
}

func (r *JSONRepository) SaveEmployee(e domain.Employee, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Employees = append(next.Employees, e)
		return r.persist(next)
	}
	for i, old := range next.Employees {
		if old.ID == e.ID {
			next.Employees[i] = e
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}

func (r *JSONRepository) DeleteEmployee(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	next.Employees = []domain.Employee{}
	found := false
	for _, e := range r.data.Employees {
		if e.ID == id {
			found = true
		} else {
			next.Employees = append(next.Employees, e)
		}
	}
	if !found {
		return domain.ErrNotFound
	}
	return r.persist(next)
}

func (r *JSONRepository) Attendances(companyID string) ([]domain.Attendance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Attendance{}
	for _, a := range r.data.Attendances {
		if matchesCompany(a.CompanyID, companyID) {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *JSONRepository) SaveAttendance(a domain.Attendance, update bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.cloneState()
	if !update {
		next.Attendances = append(next.Attendances, a)
		return r.persist(next)
	}
	for i, old := range next.Attendances {
		if old.ID == a.ID {
			next.Attendances[i] = a
			return r.persist(next)
		}
	}
	return domain.ErrNotFound
}
