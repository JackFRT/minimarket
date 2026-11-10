package infrastructure

import (
	"minimarket/internal/application"
	"path/filepath"
	"os"
)

// NewMemoryAdapter crea un adaptador efímero en memoria temporal que implementa application.Repository.
// Demuestra la intercambiabilidad del patrón Adapter sin modificar la capa de aplicación ni el dominio.
func NewMemoryAdapter() (application.Repository, error) {
	tmp := filepath.Join(os.TempDir(), "minimarket-memory-adapter.json")
	_ = os.Remove(tmp)
	repo, err := NewJSONRepository(tmp)
	if err != nil {
		return nil, err
	}
	_ = repo.EnsureRichSeed()
	return repo, nil
}
