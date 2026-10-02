package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli"

	"base-engine/config"
	"base-engine/gen"
	"base-engine/src"
	"base-engine/src/dbup"
	"base-engine/src/middleware"
	"base-engine/src/services/ai"
	aitools "base-engine/src/services/ai/tools"
	"base-engine/src/services/paymentconfig"
	sessionservice "base-engine/src/services/session"
	"gorm.io/gorm"
)

// const defaultPort = "8080"

func main() {
	app := cli.NewApp()
	app.Name = "dolphin"
	app.Usage = "This tool is for generating "
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		startCmd,
		migrateCmd,
		bootstrapHQCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

var startCmd = cli.Command{
	Name:  "start",
	Usage: "start api server",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "cors",
			Usage: "Enable cors",
		},
		cli.StringFlag{
			Name:   "p,port",
			Usage:  "Port to listen to",
			Value:  "80",
			EnvVar: "PORT",
		},
	},
	Action: func(ctx *cli.Context) error {
		cors := ctx.Bool("cors")
		port := ctx.String("port")
		if err := startServer(cors, port); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	},
}

var migrateCmd = cli.Command{
	Name:  "migrate",
	Usage: "migrate schema database",
	Action: func(ctx *cli.Context) error {
		fmt.Println("starting migration")
		if err := automigrate(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		fmt.Println("migration complete")
		return nil
	},
}

func automigrate() error {
	db, err := openEngineDBFromEnv()
	if err != nil {
		return err
	}
	defer db.Close()
	return runMigrations(db)
}

// runMigrations 按安全依赖顺序迁移实体、约束、角色与权限。
func runMigrations(db *gen.DB) error {
	if err := autoMigrateEngineDB(db); err != nil {
		return err
	}
	steps := []func(*gorm.DB) error{
		reportUnresolvedInitialAccounts, dbup.MigrateSecurityTables, dbup.EnsureGovernanceIndexes,
		dbup.InitRoles, dbup.InitPermissions,
		dbup.NormalizePermissionScopes, dbup.InitRolePermissions,
	}
	for _, step := range steps {
		if err := step(db.Query()); err != nil {
			return err
		}
	}
	return ensureSQLiteSharedIndexes(db.Query())
}

func reportUnresolvedInitialAccounts(db *gorm.DB) error {
	unresolved, err := dbup.ReportUnresolvedFranchiseInitialAccounts(db)
	if len(unresolved) > 0 {
		log.Printf("FRANCHISE_INITIAL_ACCOUNT_VERIFICATION_REQUIRED: %v", unresolved)
	}
	return err
}

// startServer 启动 GraphQL 服务并在中断信号后释放资源。
func startServer(_ bool, port string) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	db, err := openEngineDBFromEnv()
	if err != nil {
		return err
	}
	defer db.Close()
	securityConfig, err := config.LoadSecurityConfig()
	if err != nil {
		return err
	}
	handler, err := newEngineHandler(db, securityConfig)
	if err != nil {
		return err
	}


	h := newHTTPServer(port, handler)
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("connect to http://localhost:%s/graphql for GraphQL playground", port)
		serverErrors <- h.ListenAndServe()
	}()

	select {
	case <-stop:
	case serveErr := <-serverErrors:
		if serveErr != nil && serveErr != http.ErrServerClosed {
			return serveErr
		}
		return nil
	}

	log.Println("\nShutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}
	log.Println("Server gracefully stopped")

	return nil

}

func newEngineHandler(db *gen.DB, security config.SecurityConfig) (http.Handler, error) {
	aiSecurity, err := config.LoadAIModelSecurityConfig()
	if err != nil {
		log.Printf("AI_MODEL_SECURITY_CONFIG_INVALID: model configuration disabled")
		aiSecurity = config.AIModelSecurityConfig{}
	}
	router, dependencies, aiService, err := newEngineRouter(db, security, aiSecurity)
	if err != nil {
		return nil, err
	}
	handler := middleware.New(middleware.Dependencies{Config: security, PrincipalResolver: dependencies.Principal})(router)
	if err := aiService.SetProtectedHandler(handler); err != nil {
		return nil, err
	}
	items, err := aitools.Build(aiService.FixedToolRuntime())
	if err != nil {
		return nil, err
	}
	if err := aiService.SetTools(items); err != nil {
		return nil, err
	}
	return handler, nil
}

func newEngineRouter(db *gen.DB, security config.SecurityConfig, aiSecurity config.AIModelSecurityConfig) (*mux.Router, src.Dependencies, *ai.Service, error) {
	events, err := gen.NewEventController()
	if err != nil {
		return nil, src.Dependencies{}, nil, err
	}
	dependencies := src.NewDependencies(db.Query(), security, sessionservice.NewPublisher())
	if err := dependencies.Stores.SweepDocumentFiles(context.Background()); err != nil {
		log.Printf("STORE_DOCUMENT_SWEEP_FAILED: %v", err)
	}
	modelStore := ai.NewModelConfigStore(db.Query(), aiSecurity, nil)
	aiService, err := newAIService(modelStore, dependencies)
	if err != nil {
		return nil, src.Dependencies{}, nil, err
	}
	router := gen.GetHTTPServeMux(src.New(db, &events, dependencies), db)
	src.RegisterSystemInitializationRoutes(router, src.SystemInitializationDependencies{DB: db.Query(), Config: security})
	src.RegisterFranchiseInitialAccountRoute(router, db.Query(), security)
	src.RegisterAIModelConfigRoutes(router, modelStore, security)
	paymentKeys, keyErr := config.LoadPaymentConfigSecurity()
	if keyErr != nil {
		log.Printf("PAYMENT_CONFIG_SECURITY_UNAVAILABLE")
	}
	src.RegisterPaymentConfigRoutes(router, paymentconfig.NewStore(db.Query(), paymentKeys), security)
	src.RegisterStoreDocumentRoutes(router, dependencies.Stores, security)
	aiService.RegisterRoutes(router, security)
	return router, dependencies, aiService, nil
}

func newHTTPServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
