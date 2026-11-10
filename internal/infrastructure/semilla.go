package infrastructure

import "minimarket/internal/domain"

// EnsureRichSeed inicializa datos completos de demostración multiempresa si aún no se han configurado categorías o almacenes.
func (r *JSONRepository) EnsureRichSeed() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.data.Categories) > 0 && len(r.data.Warehouses) > 0 && len(r.data.Employees) > 0 {
		return nil
	}
	next := DemoState()
	return r.persist(next)
}

func DemoState() state {
	return state{
		Companies: []domain.Company{
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
		},
		Categories: []domain.Category{
			{ID: "cat-abarrotes", CompanyID: domain.DefaultCompanyID, Code: "CAT-01", Name: "Abarrotes y Granos", Description: "Arroz, azúcar, menestras, aceites y fideos", Color: "#2f6848"},
			{ID: "cat-bebidas", CompanyID: domain.DefaultCompanyID, Code: "CAT-02", Name: "Bebidas y Aguas", Description: "Gaseosas, aguas minerales, néctares y rehidratantes", Color: "#2563eb"},
			{ID: "cat-lacteos", CompanyID: domain.DefaultCompanyID, Code: "CAT-03", Name: "Lácteos y Embutidos", Description: "Leches, yogures, quesos, mantequillas y jamones", Color: "#d97706"},
			{ID: "cat-snacks", CompanyID: domain.DefaultCompanyID, Code: "CAT-04", Name: "Snacks y Confitería", Description: "Galletas, chocolates, barras energéticas y piqueos", Color: "#7c3aed"},
			{ID: "cat-limpieza", CompanyID: domain.DefaultCompanyID, Code: "CAT-05", Name: "Limpieza y Hogar", Description: "Detergentes, jabones, papel higiénico y desinfectantes", Color: "#0891b2"},
			{ID: "cat-norte-1", CompanyID: "emp-norte", Code: "CAT-N1", Name: "Despensa Familiar", Description: "Abarrotes de alta rotación sede norte", Color: "#2f6848"},
			{ID: "cat-norte-2", CompanyID: "emp-norte", Code: "CAT-N2", Name: "Bebidas Heladas", Description: "Refrescos y bebidas energizantes", Color: "#2563eb"},
		},
		Warehouses: []domain.Warehouse{
			{ID: "alm-central", CompanyID: domain.DefaultCompanyID, Code: "ALM-01", Name: "Almacén Central Miraflores", Address: "Av. Larco 840 - Sótano 1", Manager: "Jorge Quispe", Capacity: 5000, IsMain: true, Active: true},
			{ID: "alm-tienda", CompanyID: domain.DefaultCompanyID, Code: "ALM-02", Name: "Góndolas Sala de Ventas", Address: "Av. Larco 840 - Piso 1", Manager: "Lucía Fernández", Capacity: 1500, IsMain: false, Active: true},
			{ID: "alm-frio", CompanyID: domain.DefaultCompanyID, Code: "ALM-03", Name: "Cámara Frigorífica Lácteos", Address: "Av. Larco 840 - Zona Frío", Manager: "Miguel Salazar", Capacity: 800, IsMain: false, Active: true},
			{ID: "alm-norte-1", CompanyID: "emp-norte", Code: "ALM-N1", Name: "Almacén Principal Los Olivos", Address: "Av. Antúnez de Mayolo 1120", Manager: "Roberto Gómez", Capacity: 3500, IsMain: true, Active: true},
		},
		Suppliers: []domain.Supplier{
			{ID: "prov-alicorp", CompanyID: domain.DefaultCompanyID, RUC: "20100055237", Name: "Alicorp S.A.A.", Phone: "01-315-0800", Email: "pedidos@alicorp.com.pe", Address: "Av. Argentina 4793, Callao", Category: "Abarrotes y Aceites"},
			{ID: "prov-gloria", CompanyID: domain.DefaultCompanyID, RUC: "20100190797", Name: "Leche Gloria S.A.", Phone: "01-470-7170", Email: "distribucion@gloria.com.pe", Address: "Av. República de Panamá 2461", Category: "Lácteos y Derivados"},
			{ID: "prov-lindley", CompanyID: domain.DefaultCompanyID, RUC: "20101024645", Name: "Arca Continental Lindley", Phone: "0800-14020", Email: "ventas@lindley.pe", Address: "Av. Javier Prado Este 6210", Category: "Bebidas"},
			{ID: "prov-andino", CompanyID: domain.DefaultCompanyID, RUC: "20604511229", Name: "Alimentos Andinos Tarwi Perú", Phone: "999888777", Email: "contacto@tarwiperu.pe", Address: "Jr. Los Andes 340, San Luis", Category: "Snacks Saludables"},
			{ID: "prov-norte-1", CompanyID: "emp-norte", RUC: "20554433221", Name: "Distribuidora Mayorista Cono Norte", Phone: "987654321", Email: "ventas@mayoristanorte.pe", Address: "Mercado Unicachi Puesto 42", Category: "Multirubro"},
		},
		Products: []domain.Product{
			{ID: "prod-1", CompanyID: domain.DefaultCompanyID, SKU: "ABR-001", Name: "Arroz Extra Costeño Bolsa 750g", CategoryID: "cat-abarrotes", WarehouseID: "alm-central", SupplierID: "prov-alicorp", CostPrice: 3.40, Price: 4.80, Stock: 64, MinStock: 15, Unit: "bolsa", Active: true},
			{ID: "prod-2", CompanyID: domain.DefaultCompanyID, SKU: "ABR-002", Name: "Aceite Vegetal Primor Premium 900ml", CategoryID: "cat-abarrotes", WarehouseID: "alm-central", SupplierID: "prov-alicorp", CostPrice: 6.50, Price: 9.20, Stock: 38, MinStock: 12, Unit: "botella", Active: true},
			{ID: "prod-3", CompanyID: domain.DefaultCompanyID, SKU: "LAC-001", Name: "Leche Evaporada Gloria Azul 390g", CategoryID: "cat-lacteos", WarehouseID: "alm-frio", SupplierID: "prov-gloria", CostPrice: 3.20, Price: 4.50, Stock: 85, MinStock: 20, Unit: "lata", Active: true},
			{ID: "prod-4", CompanyID: domain.DefaultCompanyID, SKU: "LAC-002", Name: "Yogurt Griego Gloria Frutos Rojos 1Kg", CategoryID: "cat-lacteos", WarehouseID: "alm-frio", SupplierID: "prov-gloria", CostPrice: 6.80, Price: 9.90, Stock: 4, MinStock: 10, Unit: "botella", Active: true},
			{ID: "prod-5", CompanyID: domain.DefaultCompanyID, SKU: "BEB-001", Name: "Inca Kola Sin Azúcar 1.5 Litros", CategoryID: "cat-bebidas", WarehouseID: "alm-tienda", SupplierID: "prov-lindley", CostPrice: 4.50, Price: 6.80, Stock: 42, MinStock: 15, Unit: "botella", Active: true},
			{ID: "prod-6", CompanyID: domain.DefaultCompanyID, SKU: "BEB-002", Name: "Agua Mineral San Mateo Sin Gas 600ml", CategoryID: "cat-bebidas", WarehouseID: "alm-tienda", SupplierID: "prov-lindley", CostPrice: 1.20, Price: 2.50, Stock: 110, MinStock: 24, Unit: "botella", Active: true},
			{ID: "prod-7", CompanyID: domain.DefaultCompanyID, SKU: "SNK-001", Name: "Barra Energética de Tarwi y Cacao 60g", CategoryID: "cat-snacks", WarehouseID: "alm-tienda", SupplierID: "prov-andino", CostPrice: 2.10, Price: 3.80, Stock: 3, MinStock: 12, Unit: "und", Active: true},
			{ID: "prod-8", CompanyID: domain.DefaultCompanyID, SKU: "LIM-001", Name: "Detergente Bolívar Active Duo 800g", CategoryID: "cat-limpieza", WarehouseID: "alm-central", SupplierID: "prov-alicorp", CostPrice: 6.90, Price: 9.80, Stock: 26, MinStock: 8, Unit: "bolsa", Active: true},
			{ID: "prod-n1", CompanyID: "emp-norte", SKU: "NOR-001", Name: "Azúcar Rubia Cartavio 1Kg", CategoryID: "cat-norte-1", WarehouseID: "alm-norte-1", SupplierID: "prov-norte-1", CostPrice: 3.10, Price: 4.40, Stock: 90, MinStock: 20, Unit: "bolsa", Active: true},
			{ID: "prod-n2", CompanyID: "emp-norte", SKU: "NOR-002", Name: "Gaseosa Coca-Cola Original 3L", CategoryID: "cat-norte-2", WarehouseID: "alm-norte-1", SupplierID: "prov-norte-1", CostPrice: 8.20, Price: 11.50, Stock: 4, MinStock: 10, Unit: "botella", Active: true},
		},
		Movements: []domain.InventoryMovement{
			{ID: "mov-1", CompanyID: domain.DefaultCompanyID, ProductID: "prod-1", ProductName: "Arroz Extra Costeño Bolsa 750g", WarehouseID: "alm-central", WarehouseName: "Almacén Central Miraflores", Type: "entrada", Quantity: 50, PreviousStock: 14, NewStock: 64, Reason: "Compra Factura F001-8821 Alicorp", UserID: "usr-admin-1", UserName: "Carlos Mendoza", CreatedAt: "2026-11-09 09:15:00"},
			{ID: "mov-2", CompanyID: domain.DefaultCompanyID, ProductID: "prod-3", ProductName: "Leche Evaporada Gloria Azul 390g", WarehouseID: "alm-frio", WarehouseName: "Cámara Frigorífica Lácteos", Type: "entrada", Quantity: 48, PreviousStock: 37, NewStock: 85, Reason: "Abastecimiento semanal Gloria", UserID: "usr-admin-1", UserName: "Carlos Mendoza", CreatedAt: "2026-11-09 11:30:00"},
			{ID: "mov-3", CompanyID: domain.DefaultCompanyID, ProductID: "prod-4", ProductName: "Yogurt Griego Gloria Frutos Rojos 1Kg", WarehouseID: "alm-frio", WarehouseName: "Cámara Frigorífica Lácteos", Type: "salida", Quantity: 6, PreviousStock: 10, NewStock: 4, Reason: "Venta POS B001-00412", UserID: "usr-cajero-1", UserName: "Lucía Fernández", CreatedAt: "2026-11-10 10:20:00"},
		},
		Sales: []domain.Sale{
			{
				ID:            "vta-101",
				CompanyID:     domain.DefaultCompanyID,
				ReceiptNumber: "B001-00410",
				ReceiptType:   "Boleta",
				CustomerName:  "María Luisa Paredes",
				CustomerDoc:   "44891203",
				WarehouseID:   "alm-tienda",
				UserID:        "usr-cajero-1",
				UserName:      "Lucía Fernández",
				PaymentMethod: "yape_plin",
				Items: []domain.SaleItem{
					{ProductID: "prod-1", ProductName: "Arroz Extra Costeño Bolsa 750g", Quantity: 3, UnitPrice: 4.80, Subtotal: 14.40},
					{ProductID: "prod-2", ProductName: "Aceite Vegetal Primor Premium 900ml", Quantity: 1, UnitPrice: 9.20, Subtotal: 9.20},
					{ProductID: "prod-3", ProductName: "Leche Evaporada Gloria Azul 390g", Quantity: 4, UnitPrice: 4.50, Subtotal: 18.00},
				},
				Subtotal:  35.25,
				Tax:       6.35,
				Total:     41.60,
				Status:    "completada",
				CreatedAt: "2026-11-10 09:42:10",
			},
			{
				ID:            "vta-102",
				CompanyID:     domain.DefaultCompanyID,
				ReceiptNumber: "F001-00108",
				ReceiptType:   "Factura",
				CustomerName:  "Cafetería San Antonio SAC",
				CustomerDoc:   "20518892104",
				WarehouseID:   "alm-tienda",
				UserID:        "usr-cajero-1",
				UserName:      "Lucía Fernández",
				PaymentMethod: "tarjeta",
				Items: []domain.SaleItem{
					{ProductID: "prod-4", ProductName: "Yogurt Griego Gloria Frutos Rojos 1Kg", Quantity: 6, UnitPrice: 9.90, Subtotal: 59.40},
					{ProductID: "prod-5", ProductName: "Inca Kola Sin Azúcar 1.5 Litros", Quantity: 4, UnitPrice: 6.80, Subtotal: 27.20},
					{ProductID: "prod-7", ProductName: "Barra Energética de Tarwi y Cacao 60g", Quantity: 5, UnitPrice: 3.80, Subtotal: 19.00},
				},
				Subtotal:  89.49,
				Tax:       16.11,
				Total:     105.60,
				Status:    "completada",
				CreatedAt: "2026-11-10 11:18:45",
			},
			{
				ID:            "vta-103",
				CompanyID:     domain.DefaultCompanyID,
				ReceiptNumber: "B001-00411",
				ReceiptType:   "Boleta",
				CustomerName:  "Cliente Varios",
				CustomerDoc:   "00000000",
				WarehouseID:   "alm-tienda",
				UserID:        "usr-admin-1",
				UserName:      "Carlos Mendoza",
				PaymentMethod: "efectivo",
				Items: []domain.SaleItem{
					{ProductID: "prod-6", ProductName: "Agua Mineral San Mateo Sin Gas 600ml", Quantity: 6, UnitPrice: 2.50, Subtotal: 15.00},
					{ProductID: "prod-8", ProductName: "Detergente Bolívar Active Duo 800g", Quantity: 2, UnitPrice: 9.80, Subtotal: 19.60},
				},
				Subtotal:  29.32,
				Tax:       5.28,
				Total:     34.60,
				Status:    "completada",
				CreatedAt: "2026-11-10 13:05:22",
			},
			{
				ID:            "vta-n01",
				CompanyID:     "emp-norte",
				ReceiptNumber: "B002-00089",
				ReceiptType:   "Boleta",
				CustomerName:  "Pedro Castillo R.",
				CustomerDoc:   "10457812",
				WarehouseID:   "alm-norte-1",
				UserID:        "usr-norte-1",
				UserName:      "Roberto Gómez",
				PaymentMethod: "efectivo",
				Items: []domain.SaleItem{
					{ProductID: "prod-n1", ProductName: "Azúcar Rubia Cartavio 1Kg", Quantity: 5, UnitPrice: 4.40, Subtotal: 22.00},
					{ProductID: "prod-n2", ProductName: "Gaseosa Coca-Cola Original 3L", Quantity: 2, UnitPrice: 11.50, Subtotal: 23.00},
				},
				Subtotal:  38.14,
				Tax:       6.86,
				Total:     45.00,
				Status:    "completada",
				CreatedAt: "2026-11-10 12:10:00",
			},
		},
		Users: []domain.User{
			{ID: "usr-admin-1", CompanyID: domain.DefaultCompanyID, Name: "Carlos Mendoza (Admin)", Email: "admin@minimarket.pe", Password: "admin123", Role: "admin", Active: true, CreatedAt: "2026-01-15"},
			{ID: "usr-cajero-1", CompanyID: domain.DefaultCompanyID, Name: "Lucía Fernández (Cajera)", Email: "caja@minimarket.pe", Password: "caja123", Role: "cajero", Active: true, CreatedAt: "2026-02-01"},
			{ID: "usr-alm-1", CompanyID: domain.DefaultCompanyID, Name: "Jorge Quispe (Almacén)", Email: "almacen@minimarket.pe", Password: "almacen123", Role: "almacenero", Active: true, CreatedAt: "2026-02-05"},
			{ID: "usr-rrhh-1", CompanyID: domain.DefaultCompanyID, Name: "Elena Vargas (Jefa RRHH)", Email: "rrhh@minimarket.pe", Password: "rrhh123", Role: "rrhh", Active: true, CreatedAt: "2026-02-10"},
			{ID: "usr-norte-1", CompanyID: "emp-norte", Name: "Roberto Gómez (Admin Norte)", Email: "admin@marketnorte.pe", Password: "admin123", Role: "admin", Active: true, CreatedAt: "2026-03-10"},
		},
		Employees: []domain.Employee{
			{ID: "emp-trab-1", CompanyID: domain.DefaultCompanyID, UserID: "usr-admin-1", DNI: "42109834", Name: "Carlos Mendoza Ríos", Position: "Administrador de Tienda", Department: "Gerencia y Operaciones", WarehouseID: "alm-central", Shift: "Completo", Salary: 3400.00, HireDate: "2025-06-01", Status: "activo", Phone: "998112233", Email: "admin@minimarket.pe"},
			{ID: "emp-trab-2", CompanyID: domain.DefaultCompanyID, UserID: "usr-cajero-1", DNI: "71892341", Name: "Lucía Fernández Castro", Position: "Cajera Principal POS", Department: "Caja y Atención al Cliente", WarehouseID: "alm-tienda", Shift: "Mañana", Salary: 1650.00, HireDate: "2025-09-15", Status: "activo", Phone: "991223344", Email: "caja@minimarket.pe"},
			{ID: "emp-trab-3", CompanyID: domain.DefaultCompanyID, UserID: "usr-alm-1", DNI: "46781290", Name: "Jorge Quispe Huamán", Position: "Jefe de Almacén y Kardex", Department: "Logística e Inventario", WarehouseID: "alm-central", Shift: "Mañana", Salary: 2100.00, HireDate: "2025-08-10", Status: "activo", Phone: "987112299", Email: "almacen@minimarket.pe"},
			{ID: "emp-trab-4", CompanyID: domain.DefaultCompanyID, UserID: "usr-rrhh-1", DNI: "43891022", Name: "Elena Vargas Soto", Position: "Coordinadora de RRHH y Planillas", Department: "Recursos Humanos", WarehouseID: "alm-central", Shift: "Completo", Salary: 2600.00, HireDate: "2026-01-10", Status: "activo", Phone: "984556677", Email: "rrhh@minimarket.pe"},
			{ID: "emp-trab-5", CompanyID: domain.DefaultCompanyID, UserID: "", DNI: "74129088", Name: "Miguel Salazar Prado", Position: "Reponedor de Frío y Góndolas", Department: "Operaciones y Tienda", WarehouseID: "alm-frio", Shift: "Tarde", Salary: 1400.00, HireDate: "2026-03-01", Status: "vacaciones", Phone: "977881122", Email: "msalazar@minimarket.pe"},
			{ID: "emp-trab-n1", CompanyID: "emp-norte", UserID: "usr-norte-1", DNI: "41229081", Name: "Roberto Gómez Silva", Position: "Administrador Sede Norte", Department: "Gerencia", WarehouseID: "alm-norte-1", Shift: "Completo", Salary: 2900.00, HireDate: "2026-03-10", Status: "activo", Phone: "966554433", Email: "admin@marketnorte.pe"},
		},
		Attendances: []domain.Attendance{
			{ID: "asi-1", CompanyID: domain.DefaultCompanyID, EmployeeID: "emp-trab-1", EmployeeName: "Carlos Mendoza Ríos", Date: "2026-11-10", CheckIn: "07:55", CheckOut: "17:05", Status: "puntual", Notes: "Apertura de bóveda y arqueo inicial"},
			{ID: "asi-2", CompanyID: domain.DefaultCompanyID, EmployeeID: "emp-trab-2", EmployeeName: "Lucía Fernández Castro", Date: "2026-11-10", CheckIn: "08:00", CheckOut: "16:30", Status: "puntual", Notes: "Turno caja #1 sin diferencias"},
			{ID: "asi-3", CompanyID: domain.DefaultCompanyID, EmployeeID: "emp-trab-3", EmployeeName: "Jorge Quispe Huamán", Date: "2026-11-10", CheckIn: "08:14", CheckOut: "17:00", Status: "tardanza", Notes: "Tráfico por recepción de camión proveedor"},
			{ID: "asi-4", CompanyID: domain.DefaultCompanyID, EmployeeID: "emp-trab-4", EmployeeName: "Elena Vargas Soto", Date: "2026-11-10", CheckIn: "08:02", CheckOut: "17:00", Status: "puntual", Notes: "Cierre de asistencia semanal"},
		},
	}
}
