package organization

import (
	"context"
	"io"
	"log"
	"os"
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
	file, err := os.Open("/Users/thilinashashimalsenarath/Documents/my_projects/tiny-is/resources/test/db_scripts/organization.sql")
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

func TestCreateOrganization(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "test", Name: "test"}
	err := repo.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
}

func TestGetOrganizationByName(t *testing.T) {
	repo := NewMockOrganizationRepository()
	organization := models.Organization{Id: "test", Name: "test"}
	org, err := repo.GetOrganizationByName(context.Background(), organization.Name)
	if err != nil {
		t.Errorf("failed to get organization by name: %v", err)
	}
	if org.Name != organization.Name {
		t.Errorf("expected organization name: %s, got: %s", organization.Name, org.Name)
	}
}
