package main

import (
	"encoding/json"
	"log"
	"os"

	"base-engine/config"
	"base-engine/gen"
	"base-engine/src"
	"base-engine/src/services/ai"
)

func main() {
	schema := gen.NewExecutableSchema(gen.Config{}).Schema()
	router := gen.GetHTTPServeMux(gen.Config{}, &gen.DB{})
	src.RegisterSystemInitializationRoutes(router, src.SystemInitializationDependencies{})
	src.RegisterFranchiseInitialAccountRoute(router, nil, config.SecurityConfig{})
	src.RegisterAIModelConfigRoutes(router, nil, config.SecurityConfig{})
	src.RegisterPaymentConfigRoutes(router, nil, config.SecurityConfig{})
	src.RegisterStoreDocumentRoutes(router, nil, config.SecurityConfig{})
	var service *ai.Service
	service.RegisterRoutes(router, config.SecurityConfig{})
	operations, err := ai.DiscoverContracts(schema, router)
	if err != nil {
		log.Fatal(err)
	}
	document := struct {
		Version    int                 `json:"version"`
		Operations []ai.ContractRecord `json:"operations"`
	}{Version: 1, Operations: operations}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(document); err != nil {
		log.Fatal(err)
	}
}
