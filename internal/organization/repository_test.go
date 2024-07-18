package organization

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
)

var (
	testDB     *sqlx.DB
	dbOnce     sync.Once
	schema     []byte
	schemaOnce sync.Once
)

func loadSchema() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	path := filepath.Join(cwd, "..", "..", "resources", "test", "db_scripts", "organization.sql")
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open schema file: %v", err)
	}
	defer file.Close()
	schema, err = io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read schema file: %v", err)
	}
}

func setupTestDB() {
	schemaOnce.Do(loadSchema)
	var err error
	testDB, err = sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	_, err = testDB.Exec("PRAGMA ignore_check_constraints=OFF;")
	if err != nil {
		log.Fatalf("failed to set PRAGMA: %v", err)
	}
	_, err = testDB.Exec(string(schema))
	if err != nil {
		log.Fatalf("failed to execute schema: %v", err)
	}
}

func getTestDB() *sqlx.DB {
	dbOnce.Do(setupTestDB)
	return testDB
}

func NewMockOrganizationRepository() OrganizationRepository {
	return &organizationRepository{db: getTestDB()}
}

func TestMain(m *testing.M) {
	getTestDB()
	code := m.Run()
	if testDB != nil {
		testDB.Close()
	}
	os.Exit(code)
}

func TestRepoCreateOrganization(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "org-1", Name: "org-1"}
	err := repo.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	err = repo.DeleteOrganization(context.Background(), organization.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}

}

func TestRepoCreateOrganizationWithNullName(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "org-2"}
	err := repo.CreateOrganization(context.Background(), organization)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRepoCreateOrganizationWithEmptyName(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "org-2", Name: ""}
	err := repo.CreateOrganization(context.Background(), organization)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRepoCreateOrganizationWithEmptyId(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Name: "org-2"}
	err := repo.CreateOrganization(context.Background(), organization)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRepoCreateOrganizationWithDuplicateName(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "org-3", Name: "org-3"}
	err := repo.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	organization2 := models.Organization{Id: "org-4", Name: "org-3"}
	err = repo.CreateOrganization(context.Background(), organization2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	err = repo.DeleteOrganization(context.Background(), organization.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
	err = repo.DeleteOrganization(context.Background(), organization2.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
}

func TestRepoGetOrganizationByName(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "org-4", Name: "org-4"}
	err := repo.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	org, err := repo.GetOrganizationByName(context.Background(), organization.Name)
	if err != nil {
		t.Errorf("failed to get organization by name: %v", err)
	}
	if org.Name != organization.Name {
		t.Errorf("expected organization name: %s, got: %s", organization.Name, org.Name)
	}
	err = repo.DeleteOrganization(context.Background(), organization.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
}

func TestRepoGetOrganizationByNameNotFound(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "org-5", Name: "org-5"}
	_, err := repo.GetOrganizationByName(context.Background(), organization.Name)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
