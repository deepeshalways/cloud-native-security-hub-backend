package usecases

import (
	"database/sql"
	"log"
	"os"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/falcosecurity/cloud-native-security-hub/pkg/vendor"
)

type Factory interface {
	NewRetrieveAllResourcesUseCase() *RetrieveAllResources
	NewRetrieveOneResourceUseCase(resourceID string) *RetrieveOneResource
	NewRetrieveOneResourceByVersionUseCase(resourceID string, version string) *RetrieveOneResourceByVersion
	NewRetrieveFalcoRulesForHelmChartUseCase(resourceID string) *RetrieveFalcoRulesForHelmChart
	NewIncrementDownloadCountUseCase(resourceID string) *IncrementDownloadCount
	NewRetrieveAllVendorsUseCase() *RetrieveAllVendors
	NewRetrieveOneVendorUseCase(vendorID string) *RetrieveOneVendor
	NewRetrieveAllResourcesFromVendorUseCase(vendorID string) *RetrieveAllResourcesFromVendor

	NewResourcesRepository() resource.Repository
	NewVendorRepository() vendor.Repository
}

func NewFactory() Factory {
	factory := &factory{}
	factory.db = factory.newDB()
	factory.resourceRepository = factory.NewResourcesRepository()
	factory.vendorRepository = factory.NewVendorRepository()
	factory.resourceUpdater = factory.NewResourceUpdaterWithDB(factory.db)
	return factory
}

type factory struct {
	db                 *sql.DB
	vendorRepository   vendor.Repository
	resourceRepository resource.Repository
	resourceUpdater    resource.Updater
}

func (f *factory) NewRetrieveAllResourcesUseCase() *RetrieveAllResources {
	return &RetrieveAllResources{
		ResourceRepository: f.resourceRepository,
	}
}

func (f *factory) NewRetrieveOneResourceUseCase(resourceID string) *RetrieveOneResource {
	return &RetrieveOneResource{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
	}
}

func (f *factory) NewRetrieveOneResourceByVersionUseCase(resourceID string, version string) *RetrieveOneResourceByVersion {
	return &RetrieveOneResourceByVersion{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
		Version:            version,
	}
}

func (f *factory) NewRetrieveFalcoRulesForHelmChartUseCase(resourceID string) *RetrieveFalcoRulesForHelmChart {
	return &RetrieveFalcoRulesForHelmChart{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
	}
}

func (f *factory) NewIncrementDownloadCountUseCase(resourceID string) *IncrementDownloadCount {
	return &IncrementDownloadCount{
		Updater:    f.resourceUpdater,
		ResourceID: resourceID,
	}
}

func (f *factory) NewRetrieveAllVendorsUseCase() *RetrieveAllVendors {
	return &RetrieveAllVendors{
		VendorRepository: f.vendorRepository,
	}
}

func (f *factory) NewRetrieveOneVendorUseCase(vendorID string) *RetrieveOneVendor {
	return &RetrieveOneVendor{
		VendorRepository: f.vendorRepository,
		VendorID:         vendorID,
	}
}

func (f *factory) NewRetrieveAllResourcesFromVendorUseCase(vendorID string) *RetrieveAllResourcesFromVendor {
	return &RetrieveAllResourcesFromVendor{
		VendorID:           vendorID,
		VendorRepository:   f.vendorRepository,
		ResourceRepository: f.resourceRepository,
	}
}

func (f *factory) NewResourcesRepository() resource.Repository {
	return resource.NewPostgresRepository(f.db)
}

func (f *factory) NewVendorRepository() vendor.Repository {
	return vendor.NewPostgresRepository(f.db)
}

func (f *factory) newDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return db
}

func (f *factory) NewResourceUpdaterWithDB(db *sql.DB) resource.Updater {
	return resource.NewPostgresUpdater(db)
}
