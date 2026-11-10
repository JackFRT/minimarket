package domain

import (
	"errors"
	"strings"
)

var (
	ErrNotFound     = errors.New("registro no encontrado")
	ErrUnauthorized = errors.New("credenciales inválidas o no autorizado")
	ErrOutOfStock   = errors.New("stock insuficiente para realizar la operación")
)

const DefaultCompanyID = "emp-default"

// Company representa una empresa (tenant) dentro de la arquitectura SaaS multiempresa.
type Company struct {
	ID        string `json:"id"`
	Name      string `json:"nombre"`
	RUC       string `json:"ruc"`
	Plan      string `json:"plan"` // Starter, Pro, Enterprise
	Address   string `json:"direccion"`
	Phone     string `json:"telefono"`
	Currency  string `json:"moneda"`
	Active    bool   `json:"activa"`
	CreatedAt string `json:"creado_en"`
}

func (c Company) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("el nombre de la empresa es obligatorio")
	}
	return nil
}

// Category representa una categoría de productos por empresa.
type Category struct {
	ID          string `json:"id"`
	CompanyID   string `json:"empresa_id"`
	Code        string `json:"codigo"`
	Name        string `json:"nombre"`
	Description string `json:"descripcion"`
	Color       string `json:"color"`
}

// Clone implementa el patrón Prototype para categorías.
func (c Category) Clone(id string) Category {
	c.ID = id
	c.Name = c.Name + " (Copia)"
	return c
}

func (c Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("el nombre de la categoría es obligatorio")
	}
	return nil
}

// Warehouse representa un almacén o sucursal física de la empresa.
type Warehouse struct {
	ID        string `json:"id"`
	CompanyID string `json:"empresa_id"`
	Code      string `json:"codigo"`
	Name      string `json:"nombre"`
	Address   string `json:"direccion"`
	Manager   string `json:"responsable"`
	Capacity  int    `json:"capacidad"`
	IsMain    bool   `json:"principal"`
	Active    bool   `json:"activo"`
}

func (w Warehouse) Validate() error {
	if strings.TrimSpace(w.Name) == "" {
		return errors.New("el nombre del almacén es obligatorio")
	}
	if w.Capacity < 0 {
		return errors.New("la capacidad del almacén no puede ser negativa")
	}
	return nil
}

// Product representa un artículo del catálogo e inventario de la empresa.
type Product struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"empresa_id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"nombre"`
	CategoryID  string  `json:"categoria_id"`
	WarehouseID string  `json:"almacen_id"`
	SupplierID  string  `json:"proveedor_id"`
	CostPrice   float64 `json:"precio_compra"`
	Price       float64 `json:"precio"`
	Stock       int     `json:"stock"`
	MinStock    int     `json:"stock_minimo"`
	Unit        string  `json:"unidad"`
	Active      bool    `json:"activo"`
}

// Clone implementa el patrón Prototype. Cada copia recibe su propia identidad manteniendo sus atributos.
func (p Product) Clone(id string) Product {
	p.ID = id
	if p.SKU != "" {
		p.SKU = p.SKU + "-CLON"
	}
	return p
}

func (p Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" || p.Price < 0 || p.Stock < 0 || p.CostPrice < 0 || p.MinStock < 0 {
		return errors.New("nombre obligatorio; precio y stock deben ser no negativos")
	}
	return nil
}

// Supplier representa un proveedor de mercadería.
type Supplier struct {
	ID        string `json:"id"`
	CompanyID string `json:"empresa_id"`
	RUC       string `json:"ruc"`
	Name      string `json:"nombre"`
	Phone     string `json:"telefono"`
	Email     string `json:"email"`
	Address   string `json:"direccion"`
	Category  string `json:"rubro"`
}

func (s Supplier) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("nombre obligatorio")
	}
	return nil
}

// InventoryMovement registra entradas, salidas y ajustes de stock (Kardex).
type InventoryMovement struct {
	ID            string `json:"id"`
	CompanyID     string `json:"empresa_id"`
	ProductID     string `json:"producto_id"`
	ProductName   string `json:"producto_nombre"`
	WarehouseID   string `json:"almacen_id"`
	WarehouseName string `json:"almacen_nombre"`
	Type          string `json:"tipo"` // entrada, salida, ajuste
	Quantity      int    `json:"cantidad"`
	PreviousStock int    `json:"stock_anterior"`
	NewStock      int    `json:"stock_nuevo"`
	Reason        string `json:"motivo"`
	UserID        string `json:"usuario_id"`
	UserName      string `json:"usuario_nombre"`
	CreatedAt     string `json:"fecha"`
}

func (m InventoryMovement) Validate() error {
	if m.ProductID == "" {
		return errors.New("debe seleccionar un producto")
	}
	if m.Quantity <= 0 && m.Type != "ajuste" {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	if m.Type != "entrada" && m.Type != "salida" && m.Type != "ajuste" {
		return errors.New("el tipo de movimiento debe ser entrada, salida o ajuste")
	}
	return nil
}

// SaleItem representa una línea de detalle dentro de una venta.
type SaleItem struct {
	ProductID   string  `json:"producto_id"`
	ProductName string  `json:"producto_nombre"`
	Quantity    int     `json:"cantidad"`
	UnitPrice   float64 `json:"precio_unitario"`
	Subtotal    float64 `json:"subtotal"`
}

// Sale representa una transacción de venta en el minimarket.
type Sale struct {
	ID            string     `json:"id"`
	CompanyID     string     `json:"empresa_id"`
	ReceiptNumber string     `json:"numero_comprobante"`
	ReceiptType   string     `json:"tipo_comprobante"` // Boleta, Factura, Ticket
	CustomerName  string     `json:"cliente_nombre"`
	CustomerDoc   string     `json:"cliente_documento"`
	WarehouseID   string     `json:"almacen_id"`
	UserID        string     `json:"usuario_id"`
	UserName      string     `json:"vendedor_nombre"`
	PaymentMethod string     `json:"metodo_pago"` // efectivo, tarjeta, yape_plin, transferencia
	Items         []SaleItem `json:"items"`
	Subtotal      float64    `json:"subtotal"`
	Tax           float64    `json:"igv"`
	Total         float64    `json:"total"`
	Status        string     `json:"estado"` // completada, anulada
	CreatedAt     string     `json:"fecha"`
}

func (s Sale) Validate() error {
	if len(s.Items) == 0 {
		return errors.New("la venta debe incluir al menos un producto")
	}
	for _, item := range s.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return errors.New("cada producto vendido debe tener cantidad mayor a cero")
		}
	}
	return nil
}

// User representa un usuario con acceso al sistema dentro de una empresa.
type User struct {
	ID        string `json:"id"`
	CompanyID string `json:"empresa_id"`
	Name      string `json:"nombre"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Role      string `json:"rol"` // admin, gerente, cajero, almacenero, rrhh
	Active    bool   `json:"activo"`
	CreatedAt string `json:"creado_en"`
}

func (u User) Validate() error {
	if strings.TrimSpace(u.Name) == "" || strings.TrimSpace(u.Email) == "" {
		return errors.New("nombre y correo electrónico son obligatorios")
	}
	return nil
}

// Employee representa un trabajador en el módulo de Recursos Humanos (RRHH).
type Employee struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"empresa_id"`
	UserID      string  `json:"usuario_id"`
	DNI         string  `json:"dni"`
	Name        string  `json:"nombre"`
	Position    string  `json:"cargo"`
	Department  string  `json:"departamento"`
	WarehouseID string  `json:"almacen_id"`
	Shift       string  `json:"turno"` // Mañana, Tarde, Noche, Completo
	Salary      float64 `json:"salario"`
	HireDate    string  `json:"fecha_ingreso"`
	Status      string  `json:"estado"` // activo, vacaciones, licencia, cesado
	Phone       string  `json:"telefono"`
	Email       string  `json:"email"`
}

// Clone implementa el patrón Prototype para duplicar plantillas de puesto/contrato de trabajadores.
func (e Employee) Clone(id string) Employee {
	e.ID = id
	e.Name = e.Name + " (Nuevo Ingreso)"
	e.DNI = ""
	return e
}

func (e Employee) Validate() error {
	if strings.TrimSpace(e.Name) == "" || strings.TrimSpace(e.Position) == "" {
		return errors.New("nombre y cargo del trabajador son obligatorios")
	}
	if e.Salary < 0 {
		return errors.New("el salario no puede ser negativo")
	}
	return nil
}

// Attendance registra la asistencia diaria de los trabajadores (RRHH).
type Attendance struct {
	ID           string `json:"id"`
	CompanyID    string `json:"empresa_id"`
	EmployeeID   string `json:"empleado_id"`
	EmployeeName string `json:"empleado_nombre"`
	Date         string `json:"fecha"`
	CheckIn      string `json:"hora_entrada"`
	CheckOut     string `json:"hora_salida"`
	Status       string `json:"estado"` // puntual, tardanza, falta, justificado
	Notes        string `json:"observacion"`
}

func (a Attendance) Validate() error {
	if a.EmployeeID == "" {
		return errors.New("debe indicar el trabajador")
	}
	return nil
}

// CategoryStat resume métricas por categoría para el Dashboard.
type CategoryStat struct {
	CategoryID   string  `json:"categoria_id"`
	CategoryName string  `json:"categoria_nombre"`
	Color        string  `json:"color"`
	ProductCount int     `json:"cantidad_productos"`
	StockCount   int     `json:"stock_total"`
	SalesAmount  float64 `json:"monto_ventas"`
}

// TopProductStat resume los productos más vendidos.
type TopProductStat struct {
	ProductID   string  `json:"producto_id"`
	ProductName string  `json:"producto_nombre"`
	UnitsSold   int     `json:"unidades_vendidas"`
	Revenue     float64 `json:"ingresos"`
}

// PaymentMethodStat resume las ventas por método de pago.
type PaymentMethodStat struct {
	Method string  `json:"metodo"`
	Count  int     `json:"cantidad"`
	Total  float64 `json:"total"`
}

// DashboardSummary agrupa todos los indicadores clave (KPIs) de una empresa.
type DashboardSummary struct {
	CompanyID          string              `json:"empresa_id"`
	CompanyName        string              `json:"empresa_nombre"`
	TotalProducts      int                 `json:"total_productos"`
	TotalCategories    int                 `json:"total_categorias"`
	TotalWarehouses    int                 `json:"total_almacenes"`
	TotalSuppliers     int                 `json:"total_proveedores"`
	TotalUnitsInStock  int                 `json:"total_unidades_stock"`
	InventoryValue     float64             `json:"valor_inventario"`
	InventoryCostValue float64             `json:"costo_inventario"`
	LowStockCount      int                 `json:"alertas_stock_bajo"`
	LowStockProducts   []Product           `json:"productos_stock_bajo"`
	TotalSalesCount    int                 `json:"total_ventas"`
	TotalRevenue       float64             `json:"ingresos_totales"`
	AverageTicket      float64             `json:"ticket_promedio"`
	ActiveEmployees    int                 `json:"empleados_activos"`
	MonthlyPayroll     float64             `json:"planilla_mensual"`
	AttendanceToday    int                 `json:"asistencias_registradas"`
	CategoryStats      []CategoryStat      `json:"ventas_por_categoria"`
	TopProducts        []TopProductStat    `json:"productos_mas_vendidos"`
	PaymentStats       []PaymentMethodStat `json:"ventas_por_metodo_pago"`
	RecentSales        []Sale              `json:"ultimas_ventas"`
	RecentMovements    []InventoryMovement `json:"ultimos_movimientos"`
}
