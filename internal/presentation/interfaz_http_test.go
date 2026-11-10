package presentation_test

import (
	"encoding/json"
	"minimarket/internal/application"
	"minimarket/internal/infrastructure"
	"minimarket/internal/presentation"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestCRUDPrototypeAndPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.json")
	repo, err := infrastructure.NewJSONRepository(path)
	if err != nil {
		t.Fatal(err)
	}

	h := presentation.NewHandler(application.Service{Repo: repo})
	request := func(method, path, body string, status int) map[string]any {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
		if w.Code != status {
			t.Fatalf("%s %s: got %d want %d: %s", method, path, w.Code, status, w.Body.String())
		}
		var v map[string]any
		if status != 204 {
			if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
				t.Fatal(err)
			}
		}
		return v
	}
	product := request("POST", "/api/productos", `{"nombre":"Arroz","precio":4.5,"stock":20}`, 201)
	id := product["id"].(string)
	clone := request("POST", "/api/productos/"+id+"/clonar", "", 201)
	if clone["id"] == id || clone["nombre"] != product["nombre"] {
		t.Fatal("clon inválido")
	}
	request("PUT", "/api/productos/"+id, `{"nombre":"Arroz nuevo","precio":5,"stock":10}`, 200)
	if request("GET", "/api/productos/"+clone["id"].(string), "", 200)["nombre"] != "Arroz" {
		t.Fatal("se modificó el clon")
	}
	supplier := request("POST", "/api/proveedores", `{"nombre":"Proveedor","telefono":"123"}`, 201)
	sid := supplier["id"].(string)
	request("PUT", "/api/proveedores/"+sid, `{"nombre":"Proveedor editado","telefono":"456"}`, 200)
	request("GET", "/api/proveedores/"+sid, "", 200)
	reopened, err := infrastructure.NewJSONRepository(path)
	if err != nil {
		t.Fatal(err)
	}
	p, err := reopened.Product(id)
	if err != nil || p.Name != "Arroz nuevo" {
		t.Fatal("persistencia incorrecta")
	}
	s, err := reopened.Supplier(sid)
	if err != nil || s.Phone != "456" {
		t.Fatal("persistencia proveedor incorrecta")
	}
	for _, route := range []string{"/api/productos", "/api/proveedores"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, route, nil))
		var list []any
		if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &list) != nil || len(list) == 0 {
			t.Fatal("listado incorrecto")
		}
	}
	request("POST", "/api/productos", `{"nombre":"Mal","precio":-1}`, 400)
	request("POST", "/api/proveedores", `{"nombre":""}`, 400)
	request("POST", "/api/productos", `{"nombre":"x"} {}`, 400)
	request("DELETE", "/api/productos/"+id, "", 204)
	request("GET", "/api/productos/"+id, "", 404)
	request("DELETE", "/api/proveedores/"+sid, "", 204)
	request("GET", "/api/proveedores/"+sid, "", 404)
	request("PUT", "/api/productos/no-existe", `{"nombre":"x"}`, 404)
	request("POST", "/api/productos/no-existe/clonar", "", 404)
}

func TestSaaSMultitenantInventorySalesHRAndDashboard(t *testing.T) {
	path := filepath.Join(t.TempDir(), "saas.json")
	repo, err := infrastructure.NewJSONRepository(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.EnsureRichSeed(); err != nil {
		t.Fatal(err)
	}
	h := presentation.NewHandler(application.Service{Repo: repo})

	reqTenant := func(method, route, tenantID, body string, want int) map[string]any {
		t.Helper()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, route, strings.NewReader(body))
		if tenantID != "" {
			r.Header.Set("X-Empresa-ID", tenantID)
		}
		h.ServeHTTP(w, r)
		if w.Code != want {
			t.Fatalf("%s %s: got %d want %d: %s", method, route, w.Code, want, w.Body.String())
		}
		var v map[string]any
		if want != 204 {
			_ = json.Unmarshal(w.Body.Bytes(), &v)
		}
		return v
	}

	// 1. Login SaaS
	auth := reqTenant("POST", "/api/auth/login", "", `{"email":"admin@minimarket.pe","password":"admin123"}`, 200)
	if auth["token"] == nil {
		t.Fatal("se esperaba token de autenticación")
	}

	// 2. Onboarding nueva empresa SaaS
	newComp := reqTenant("POST", "/api/empresas", "", `{"nombre":"Minimarket San Borja","ruc":"20605544332","plan":"Enterprise","direccion":"Av. Aviación 2400"}`, 201)
	compID := newComp["id"].(string)

	// 3. Crear categoría, almacén y producto en la nueva empresa
	cat := reqTenant("POST", "/api/categorias", compID, `{"nombre":"Congelados","descripcion":"Helados y carnes","color":"#0891b2"}`, 201)
	catID := cat["id"].(string)
	reqTenant("POST", "/api/categorias/"+catID+"/clonar", compID, "", 201)

	wh := reqTenant("POST", "/api/almacenes", compID, `{"nombre":"Almacén Frío SB","direccion":"Sótano B","responsable":"Ana","capacidad":500}`, 201)
	whID := wh["id"].(string)

	prod := reqTenant("POST", "/api/productos", compID, `{"nombre":"Helado Artesanal 1L","categoria_id":"`+catID+`","almacen_id":"`+whID+`","precio_compra":12,"precio":20,"stock":15,"stock_minimo":5}`, 201)
	prodID := prod["id"].(string)

	// 4. Movimiento de inventario Kardex (entrada +5 => stock 20)
	mov := reqTenant("POST", "/api/inventario/movimientos", compID, `{"producto_id":"`+prodID+`","tipo":"entrada","cantidad":5,"motivo":"Reposición inicial"}`, 201)
	if int(mov["stock_nuevo"].(float64)) != 20 {
		t.Fatalf("stock esperado 20, obtenido %v", mov["stock_nuevo"])
	}

	// 5. Registro de Venta POS (vende 3 unidades => stock 17, total 60.0)
	sale := reqTenant("POST", "/api/ventas", compID, `{"cliente_nombre":"Juan Pérez","metodo_pago":"yape_plin","items":[{"producto_id":"`+prodID+`","cantidad":3}]}`, 201)
	if sale["total"].(float64) != 60.0 {
		t.Fatalf("total venta esperado 60.0, obtenido %v", sale["total"])
	}

	// 6. Recursos Humanos (Empleado + Prototype Clone + Asistencia)
	emp := reqTenant("POST", "/api/rrhh/empleados", compID, `{"dni":"45678901","nombre":"Sofía Ramos","cargo":"Cajera Senior","turno":"Mañana","salario":1800}`, 201)
	empID := emp["id"].(string)
	empClone := reqTenant("POST", "/api/rrhh/empleados/"+empID+"/clonar", compID, "", 201)
	if empClone["id"] == empID {
		t.Fatal("el clon de empleado debe tener distinto ID")
	}
	reqTenant("POST", "/api/rrhh/asistencias", compID, `{"empleado_id":"`+empID+`","hora_entrada":"08:00","hora_salida":"17:00","estado":"puntual"}`, 201)

	// 7. Dashboard de la nueva empresa
	dash := reqTenant("GET", "/api/dashboard", compID, "", 200)
	if dash["ingresos_totales"].(float64) != 60.0 {
		t.Fatalf("ingresos dashboard esperados 60.0, obtenidos %v", dash["ingresos_totales"])
	}
}

func TestHomePage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.json")
	repo, err := infrastructure.NewJSONRepository(path)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	presentation.NewHandler(application.Service{Repo: repo}).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("GET /: got %d want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("Content-Type: got %q", got)
	}
	for _, content := range []string{"Minimarket", "Productos", "Proveedores", "product-count", "supplier-count"} {
		if !strings.Contains(w.Body.String(), content) {
			t.Errorf("home page does not contain %q", content)
		}
	}
}
