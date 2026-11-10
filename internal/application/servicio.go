package application

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"minimarket/internal/domain"
	"sort"
	"strings"
	"time"
)

// Repository define el puerto de salida (interfaz) que implementan los adaptadores de infraestructura (JSON, Memoria, SQL).
type Repository interface {
	// Empresas (Tenants SaaS)
	Companies() ([]domain.Company, error)
	Company(id string) (domain.Company, error)
	SaveCompany(c domain.Company, update bool) error

	// Categorías
	Categories(companyID string) ([]domain.Category, error)
	Category(id string) (domain.Category, error)
	SaveCategory(c domain.Category, update bool) error
	DeleteCategory(id string) error

	// Almacenes
	Warehouses(companyID string) ([]domain.Warehouse, error)
	Warehouse(id string) (domain.Warehouse, error)
	SaveWarehouse(w domain.Warehouse, update bool) error
	DeleteWarehouse(id string) error

	// Productos (compatibilidad directa y multiempresa)
	Products() ([]domain.Product, error)
	ProductsByCompany(companyID string) ([]domain.Product, error)
	Product(id string) (domain.Product, error)
	SaveProduct(p domain.Product, update bool) error
	DeleteProduct(id string) error

	// Proveedores
	Suppliers() ([]domain.Supplier, error)
	SuppliersByCompany(companyID string) ([]domain.Supplier, error)
	Supplier(id string) (domain.Supplier, error)
	SaveSupplier(s domain.Supplier, update bool) error
	DeleteSupplier(id string) error

	// Movimientos de Inventario (Kardex)
	Movements(companyID string) ([]domain.InventoryMovement, error)
	SaveMovement(m domain.InventoryMovement) error

	// Ventas
	Sales(companyID string) ([]domain.Sale, error)
	Sale(id string) (domain.Sale, error)
	SaveSale(s domain.Sale, update bool) error

	// Usuarios y Autenticación
	Users(companyID string) ([]domain.User, error)
	UserByEmail(email string, companyID string) (domain.User, error)
	User(id string) (domain.User, error)
	SaveUser(u domain.User, update bool) error
	DeleteUser(id string) error

	// Recursos Humanos (Empleados y Asistencia)
	Employees(companyID string) ([]domain.Employee, error)
	Employee(id string) (domain.Employee, error)
	SaveEmployee(e domain.Employee, update bool) error
	DeleteEmployee(id string) error
	Attendances(companyID string) ([]domain.Attendance, error)
	SaveAttendance(a domain.Attendance, update bool) error
}

type Service struct {
	Repo Repository
}

func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}

func nowISO() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func todayISO() string {
	return time.Now().Format("2006-01-02")
}

func normalizeCompany(companyID string) string {
	if strings.TrimSpace(companyID) == "" {
		return domain.DefaultCompanyID
	}
	return strings.TrimSpace(companyID)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// --- EMPRESAS Y ONBOARDING SAAS ---

func (s Service) SaveCompany(c domain.Company, id string) (domain.Company, error) {
	c.ID = id
	if err := c.Validate(); err != nil {
		return c, err
	}
	if c.Plan == "" {
		c.Plan = "Pro"
	}
	if c.Currency == "" {
		c.Currency = "PEN"
	}
	isNew := id == ""
	if isNew {
		c.ID = "emp-" + newID()[:8]
		c.Active = true
		c.CreatedAt = todayISO()
	}
	if err := s.Repo.SaveCompany(c, !isNew); err != nil {
		return c, err
	}
	// Si es una nueva empresa SaaS, inicializamos su almacén principal y categorías base
	if isNew {
		wh, _ := s.SaveWarehouse(domain.Warehouse{
			CompanyID: c.ID,
			Code:      "ALM-01",
			Name:      "Almacén Central - " + c.Name,
			Address:   c.Address,
			Manager:   "Administrador General",
			Capacity:  5000,
			IsMain:    true,
			Active:    true,
		}, "")
		cat, _ := s.SaveCategory(domain.Category{
			CompanyID:   c.ID,
			Code:        "CAT-AB",
			Name:        "Abarrotes y Despensa",
			Description: "Productos básicos de primera necesidad",
			Color:       "#2f6848",
		}, "")
		_, _ = s.SaveCategory(domain.Category{
			CompanyID:   c.ID,
			Code:        "CAT-BEB",
			Name:        "Bebidas y Lácteos",
			Description: "Gaseosas, aguas, jugos y lácteos",
			Color:       "#2563eb",
		}, "")
		_, _ = s.SaveUser(domain.User{
			CompanyID: c.ID,
			Name:      "Admin " + c.Name,
			Email:     "admin@" + strings.ToLower(strings.ReplaceAll(c.Name, " ", "")) + ".pe",
			Password:  "admin123",
			Role:      "admin",
			Active:    true,
		}, "")
		_, _ = s.SaveProduct(domain.Product{
			CompanyID:   c.ID,
			SKU:         "INI-001",
			Name:        "Producto Inaugural " + c.Name,
			CategoryID:  cat.ID,
			WarehouseID: wh.ID,
			CostPrice:   3.00,
			Price:       4.50,
			Stock:       50,
			MinStock:    10,
			Unit:        "und",
			Active:      true,
		}, "")
	}
	return c, nil
}

// --- AUTENTICACIÓN Y USUARIOS ---

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	CompanyID string `json:"empresa_id"`
}

type AuthResponse struct {
	Token   string         `json:"token"`
	User    domain.User    `json:"usuario"`
	Company domain.Company `json:"empresa"`
}

func (s Service) Authenticate(req LoginRequest) (AuthResponse, error) {
	u, err := s.Repo.UserByEmail(strings.TrimSpace(req.Email), req.CompanyID)
	if err != nil {
		return AuthResponse{}, domain.ErrUnauthorized
	}
	if u.Password != "" && u.Password != req.Password {
		return AuthResponse{}, domain.ErrUnauthorized
	}
	comp, err := s.Repo.Company(u.CompanyID)
	if err != nil {
		comp = domain.Company{ID: u.CompanyID, Name: "Minimarket Principal", Plan: "Enterprise", Currency: "PEN", Active: true}
	}
	uSafe := u
	uSafe.Password = ""
	token := fmt.Sprintf("saas-jwt-%s-%s-%s", u.CompanyID, u.ID, newID()[:12])
	return AuthResponse{
		Token:   token,
		User:    uSafe,
		Company: comp,
	}, nil
}

func (s Service) SaveUser(u domain.User, id string) (domain.User, error) {
	u.ID = id
	u.CompanyID = normalizeCompany(u.CompanyID)
	if u.Role == "" {
		u.Role = "cajero"
	}
	if err := u.Validate(); err != nil {
		return u, err
	}
	isNew := id == ""
	if isNew {
		u.ID = newID()
		u.Active = true
		u.CreatedAt = todayISO()
		if u.Password == "" {
			u.Password = "123456"
		}
	} else {
		existing, err := s.Repo.User(id)
		if err != nil {
			return u, err
		}
		if u.Password == "" {
			u.Password = existing.Password
		}
		if u.CreatedAt == "" {
			u.CreatedAt = existing.CreatedAt
		}
	}
	err := s.Repo.SaveUser(u, !isNew)
	u.Password = ""
	return u, err
}

// --- CATEGORÍAS ---

func (s Service) SaveCategory(c domain.Category, id string) (domain.Category, error) {
	c.ID = id
	c.CompanyID = normalizeCompany(c.CompanyID)
	if err := c.Validate(); err != nil {
		return c, err
	}
	if c.Color == "" {
		c.Color = "#2f6848"
	}
	if c.Code == "" {
		c.Code = "CAT-" + strings.ToUpper(newID()[:4])
	}
	if id == "" {
		c.ID = newID()
	}
	return c, s.Repo.SaveCategory(c, id != "")
}

func (s Service) CloneCategory(id string) (domain.Category, error) {
	c, err := s.Repo.Category(id)
	if err != nil {
		return c, err
	}
	clone := c.Clone(newID())
	return clone, s.Repo.SaveCategory(clone, false)
}

// --- ALMACENES ---

func (s Service) SaveWarehouse(w domain.Warehouse, id string) (domain.Warehouse, error) {
	w.ID = id
	w.CompanyID = normalizeCompany(w.CompanyID)
	if err := w.Validate(); err != nil {
		return w, err
	}
	if w.Code == "" {
		w.Code = "ALM-" + strings.ToUpper(newID()[:4])
	}
	if w.Capacity == 0 {
		w.Capacity = 1000
	}
	if id == "" {
		w.ID = newID()
		w.Active = true
	}
	return w, s.Repo.SaveWarehouse(w, id != "")
}

// --- PRODUCTOS Y PROTOTYPE ---

func (s Service) SaveProduct(p domain.Product, id string) (domain.Product, error) {
	p.ID = id
	p.CompanyID = normalizeCompany(p.CompanyID)
	if err := p.Validate(); err != nil {
		return p, err
	}
	if p.Unit == "" {
		p.Unit = "und"
	}
	if p.MinStock == 0 {
		p.MinStock = 5
	}
	if p.SKU == "" {
		p.SKU = "PRD-" + strings.ToUpper(newID()[:6])
	}
	p.Active = true
	if id == "" {
		p.ID = newID()
	}
	return p, s.Repo.SaveProduct(p, id != "")
}

func (s Service) CloneProduct(id string) (domain.Product, error) {
	p, err := s.Repo.Product(id)
	if err != nil {
		return p, err
	}
	clone := p.Clone(newID())
	return clone, s.Repo.SaveProduct(clone, false)
}

// --- PROVEEDORES ---

func (s Service) SaveSupplier(p domain.Supplier, id string) (domain.Supplier, error) {
	p.ID = id
	p.CompanyID = normalizeCompany(p.CompanyID)
	if err := p.Validate(); err != nil {
		return p, err
	}
	if id == "" {
		p.ID = newID()
	}
	return p, s.Repo.SaveSupplier(p, id != "")
}

// --- CONTROL DE INVENTARIO (KARDEX) ---

func (s Service) RegisterInventoryMovement(m domain.InventoryMovement) (domain.InventoryMovement, error) {
	m.CompanyID = normalizeCompany(m.CompanyID)
	if err := m.Validate(); err != nil {
		return m, err
	}
	prod, err := s.Repo.Product(m.ProductID)
	if err != nil {
		return m, err
	}
	prevStock := prod.Stock
	var newStock int
	switch m.Type {
	case "entrada":
		newStock = prevStock + m.Quantity
	case "salida":
		if prevStock < m.Quantity {
			return m, domain.ErrOutOfStock
		}
		newStock = prevStock - m.Quantity
	case "ajuste":
		if m.Quantity < 0 {
			return m, errors.New("el stock ajustado no puede ser negativo")
		}
		newStock = m.Quantity
	}

	prod.Stock = newStock
	if m.WarehouseID != "" {
		prod.WarehouseID = m.WarehouseID
	}
	if err := s.Repo.SaveProduct(prod, true); err != nil {
		return m, err
	}

	whName := "Almacén General"
	if prod.WarehouseID != "" {
		if wh, err := s.Repo.Warehouse(prod.WarehouseID); err == nil {
			whName = wh.Name
		}
	}

	m.ID = newID()
	m.ProductName = prod.Name
	m.WarehouseID = prod.WarehouseID
	m.WarehouseName = whName
	m.PreviousStock = prevStock
	m.NewStock = newStock
	if m.UserName == "" {
		m.UserName = "Administrador"
	}
	m.CreatedAt = nowISO()

	return m, s.Repo.SaveMovement(m)
}

// --- REGISTRO DE VENTAS (POS) ---

func (s Service) RegisterSale(sale domain.Sale) (domain.Sale, error) {
	sale.CompanyID = normalizeCompany(sale.CompanyID)
	if err := sale.Validate(); err != nil {
		return sale, err
	}

	var total float64
	enrichedItems := make([]domain.SaleItem, 0, len(sale.Items))

	// 1. Validar stock y calcular importes
	productsToUpdate := make([]domain.Product, 0, len(sale.Items))
	for _, item := range sale.Items {
		prod, err := s.Repo.Product(item.ProductID)
		if err != nil {
			return sale, fmt.Errorf("producto %s no encontrado", item.ProductID)
		}
		if prod.Stock < item.Quantity {
			return sale, fmt.Errorf("stock insuficiente para %s (disponible: %d, solicitado: %d)", prod.Name, prod.Stock, item.Quantity)
		}
		unitPrice := prod.Price
		if item.UnitPrice > 0 {
			unitPrice = item.UnitPrice
		}
		lineSubtotal := round2(unitPrice * float64(item.Quantity))
		total += lineSubtotal

		enrichedItems = append(enrichedItems, domain.SaleItem{
			ProductID:   prod.ID,
			ProductName: prod.Name,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
			Subtotal:    lineSubtotal,
		})

		prod.Stock -= item.Quantity
		productsToUpdate = append(productsToUpdate, prod)
	}

	// 2. Descontar stock y registrar movimientos Kardex
	for i, prod := range productsToUpdate {
		item := enrichedItems[i]
		prevStock := prod.Stock + item.Quantity
		if err := s.Repo.SaveProduct(prod, true); err != nil {
			return sale, err
		}
		_ = s.Repo.SaveMovement(domain.InventoryMovement{
			ID:            newID(),
			CompanyID:     sale.CompanyID,
			ProductID:     prod.ID,
			ProductName:   prod.Name,
			WarehouseID:   prod.WarehouseID,
			WarehouseName: "Despacho POS",
			Type:          "salida",
			Quantity:      item.Quantity,
			PreviousStock: prevStock,
			NewStock:      prod.Stock,
			Reason:        "Venta POS",
			UserID:        sale.UserID,
			UserName:      sale.UserName,
			CreatedAt:     nowISO(),
		})
	}

	sale.ID = newID()
	if sale.ReceiptType == "" {
		sale.ReceiptType = "Boleta"
	}
	if sale.ReceiptNumber == "" {
		sale.ReceiptNumber = "B001-" + strings.ToUpper(sale.ID[:6])
	}
	if sale.CustomerName == "" {
		sale.CustomerName = "Cliente Varios"
	}
	if sale.PaymentMethod == "" {
		sale.PaymentMethod = "efectivo"
	}
	if sale.UserName == "" {
		sale.UserName = "Cajero Principal"
	}
	sale.Items = enrichedItems
	sale.Total = round2(total)
	sale.Subtotal = round2(sale.Total / 1.18)
	sale.Tax = round2(sale.Total - sale.Subtotal)
	sale.Status = "completada"
	sale.CreatedAt = nowISO()

	return sale, s.Repo.SaveSale(sale, false)
}

// --- RECURSOS HUMANOS (RRHH) ---

func (s Service) SaveEmployee(e domain.Employee, id string) (domain.Employee, error) {
	e.ID = id
	e.CompanyID = normalizeCompany(e.CompanyID)
	if err := e.Validate(); err != nil {
		return e, err
	}
	if e.Shift == "" {
		e.Shift = "Mañana"
	}
	if e.Status == "" {
		e.Status = "activo"
	}
	if e.Department == "" {
		e.Department = "Operaciones y Tienda"
	}
	if e.HireDate == "" {
		e.HireDate = todayISO()
	}
	if id == "" {
		e.ID = newID()
	}
	return e, s.Repo.SaveEmployee(e, id != "")
}

func (s Service) CloneEmployee(id string) (domain.Employee, error) {
	e, err := s.Repo.Employee(id)
	if err != nil {
		return e, err
	}
	clone := e.Clone(newID())
	return clone, s.Repo.SaveEmployee(clone, false)
}

func (s Service) RegisterAttendance(a domain.Attendance) (domain.Attendance, error) {
	a.CompanyID = normalizeCompany(a.CompanyID)
	if err := a.Validate(); err != nil {
		return a, err
	}
	emp, err := s.Repo.Employee(a.EmployeeID)
	if err != nil {
		return a, err
	}
	a.EmployeeName = emp.Name
	if a.Date == "" {
		a.Date = todayISO()
	}
	if a.CheckIn == "" {
		a.CheckIn = time.Now().Format("15:04")
	}
	if a.Status == "" {
		a.Status = "puntual"
	}
	isUpdate := a.ID != ""
	if !isUpdate {
		a.ID = newID()
	}
	return a, s.Repo.SaveAttendance(a, isUpdate)
}

// --- DASHBOARD Y REPORTES ---

func (s Service) GetDashboardSummary(companyID string) (domain.DashboardSummary, error) {
	cid := normalizeCompany(companyID)
	comp, err := s.Repo.Company(cid)
	if err != nil {
		comp = domain.Company{ID: cid, Name: "Minimarket Principal"}
	}

	products, _ := s.Repo.ProductsByCompany(cid)
	categories, _ := s.Repo.Categories(cid)
	warehouses, _ := s.Repo.Warehouses(cid)
	suppliers, _ := s.Repo.SuppliersByCompany(cid)
	sales, _ := s.Repo.Sales(cid)
	movements, _ := s.Repo.Movements(cid)
	employees, _ := s.Repo.Employees(cid)
	attendances, _ := s.Repo.Attendances(cid)

	summary := domain.DashboardSummary{
		CompanyID:        cid,
		CompanyName:      comp.Name,
		TotalProducts:    len(products),
		TotalCategories:  len(categories),
		TotalWarehouses:  len(warehouses),
		TotalSuppliers:   len(suppliers),
		LowStockProducts: []domain.Product{},
		CategoryStats:    []domain.CategoryStat{},
		TopProducts:      []domain.TopProductStat{},
		PaymentStats:     []domain.PaymentMethodStat{},
		RecentSales:      []domain.Sale{},
		RecentMovements:  []domain.InventoryMovement{},
	}

	catMap := make(map[string]*domain.CategoryStat)
	for _, c := range categories {
		catMap[c.ID] = &domain.CategoryStat{
			CategoryID:   c.ID,
			CategoryName: c.Name,
			Color:        c.Color,
		}
	}
	catMap["sin-categoria"] = &domain.CategoryStat{
		CategoryID:   "sin-categoria",
		CategoryName: "General / Sin categoría",
		Color:        "#64748b",
	}

	prodCategory := make(map[string]string)
	for _, p := range products {
		summary.TotalUnitsInStock += p.Stock
		summary.InventoryValue += round2(p.Price * float64(p.Stock))
		summary.InventoryCostValue += round2(p.CostPrice * float64(p.Stock))
		minStock := p.MinStock
		if minStock <= 0 {
			minStock = 5
		}
		if p.Stock <= minStock {
			summary.LowStockCount++
			summary.LowStockProducts = append(summary.LowStockProducts, p)
		}
		catKey := p.CategoryID
		if _, ok := catMap[catKey]; !ok {
			catKey = "sin-categoria"
		}
		prodCategory[p.ID] = catKey
		catMap[catKey].ProductCount++
		catMap[catKey].StockCount += p.Stock
	}
	summary.InventoryValue = round2(summary.InventoryValue)
	summary.InventoryCostValue = round2(summary.InventoryCostValue)

	topProdMap := make(map[string]*domain.TopProductStat)
	payMap := make(map[string]*domain.PaymentMethodStat)

	for _, sale := range sales {
		if sale.Status == "anulada" {
			continue
		}
		summary.TotalSalesCount++
		summary.TotalRevenue += sale.Total

		pm := sale.PaymentMethod
		if pm == "" {
			pm = "efectivo"
		}
		if _, ok := payMap[pm]; !ok {
			payMap[pm] = &domain.PaymentMethodStat{Method: pm}
		}
		payMap[pm].Count++
		payMap[pm].Total = round2(payMap[pm].Total + sale.Total)

		for _, item := range sale.Items {
			if _, ok := topProdMap[item.ProductID]; !ok {
				topProdMap[item.ProductID] = &domain.TopProductStat{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
				}
			}
			topProdMap[item.ProductID].UnitsSold += item.Quantity
			topProdMap[item.ProductID].Revenue = round2(topProdMap[item.ProductID].Revenue + item.Subtotal)

			catKey := prodCategory[item.ProductID]
			if catKey == "" {
				catKey = "sin-categoria"
			}
			if st, ok := catMap[catKey]; ok {
				st.SalesAmount = round2(st.SalesAmount + item.Subtotal)
			}
		}
	}

	summary.TotalRevenue = round2(summary.TotalRevenue)
	if summary.TotalSalesCount > 0 {
		summary.AverageTicket = round2(summary.TotalRevenue / float64(summary.TotalSalesCount))
	}

	for _, e := range employees {
		if e.Status == "activo" || e.Status == "" {
			summary.ActiveEmployees++
			summary.MonthlyPayroll += e.Salary
		}
	}
	summary.MonthlyPayroll = round2(summary.MonthlyPayroll)
	summary.AttendanceToday = len(attendances)

	for _, st := range catMap {
		if st.ProductCount > 0 || st.SalesAmount > 0 || st.CategoryID != "sin-categoria" {
			summary.CategoryStats = append(summary.CategoryStats, *st)
		}
	}
	sort.Slice(summary.CategoryStats, func(i, j int) bool {
		return summary.CategoryStats[i].SalesAmount > summary.CategoryStats[j].SalesAmount
	})

	for _, tp := range topProdMap {
		summary.TopProducts = append(summary.TopProducts, *tp)
	}
	sort.Slice(summary.TopProducts, func(i, j int) bool {
		return summary.TopProducts[i].UnitsSold > summary.TopProducts[j].UnitsSold
	})
	if len(summary.TopProducts) > 6 {
		summary.TopProducts = summary.TopProducts[:6]
	}

	for _, pm := range payMap {
		summary.PaymentStats = append(summary.PaymentStats, *pm)
	}
	sort.Slice(summary.PaymentStats, func(i, j int) bool {
		return summary.PaymentStats[i].Total > summary.PaymentStats[j].Total
	})

	// Últimas ventas y movimientos (orden inverso cronológico)
	for i := len(sales) - 1; i >= 0 && len(summary.RecentSales) < 8; i-- {
		summary.RecentSales = append(summary.RecentSales, sales[i])
	}
	for i := len(movements) - 1; i >= 0 && len(summary.RecentMovements) < 8; i-- {
		summary.RecentMovements = append(summary.RecentMovements, movements[i])
	}

	return summary, nil
}
