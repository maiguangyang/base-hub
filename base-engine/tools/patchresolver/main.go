/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package main

import (
	"fmt"
	"go/format"
	"os"
	"strings"
)

const resolverPath = "src/resolver_gen.go"
const legacyEncryptionPath = "utils/encrypt.go"
const httpHandlerPath = "gen/http-handler.go"
const generatedExecutionPath = "gen/generated.go"
const generatedQueriesPath = "gen/resolver-queries.go"
const graphqlDocsPath = "docs/api.gql"

type replacement struct {
	old string
	new string
}

func main() {
	if err := patchFile(resolverPath, patchResolver); err != nil {
		panic(err)
	}
	if err := patchFile(legacyEncryptionPath, patchLegacyPassword); err != nil {
		panic(err)
	}
	if err := patchFile(httpHandlerPath, patchHTTPHandler); err != nil {
		panic(err)
	}
	if err := patchFile(generatedExecutionPath, patchInitialAccountProjection); err != nil {
		panic(err)
	}
}

func patchCustomerCouponDistributionJobIDProjections(source string) (string, error) {
	targets := []string{
		"func (r *GeneratedCustomerMemberResolver) CouponDistributionJobsIds(ctx context.Context, obj *CustomerMember) (ids []string, err error) {",
		"func (r *GeneratedCustomerCouponTemplateResolver) DistributionJobsIds(ctx context.Context, obj *CustomerCouponTemplate) (ids []string, err error) {",
	}
	const guard = `
	if _, err := auth.RequirePrincipal(ctx); err != nil {
		return nil, err
	}
	if denied := auth.NewError(auth.CodePermissionDenied); denied != nil {
		return nil, denied
	}
`
	for _, target := range targets {
		patched := target + guard
		if strings.Contains(source, patched) {
			continue
		}
		if strings.Count(source, target) != 1 {
			return "", fmt.Errorf("coupon distribution job ID projection patch contract changed: %q", target)
		}
		source = strings.Replace(source, target, patched, 1)
	}
	return source, nil
}

func patchGraphqlDocsFile() error {
	source, err := os.ReadFile(graphqlDocsPath)
	if err != nil {
		return err
	}
	patched, err := patchSpecificationOrderDocs(string(source))
	if err != nil {
		return err
	}
	if patched == string(source) {
		return nil
	}
	return os.WriteFile(graphqlDocsPath, []byte(patched), 0o644)
}

func patchSpecificationOrderDocs(source string) (string, error) {
	const generated = "mutation hqReorderSpecificationValues($specificationId: ID!, $orderedIds: [ID]!)"
	const corrected = "mutation hqReorderSpecificationValues($specificationId: ID!, $orderedIds: [ID!]!)"
	if strings.Count(source, corrected) == 1 {
		return source, nil
	}
	if strings.Count(source, generated) != 1 {
		return "", fmt.Errorf("specification order documentation contract changed")
	}
	return strings.Replace(source, generated, corrected, 1), nil
}

func patchInitialAccountProjection(source string) (string, error) {
	if strings.Contains(source, "return ec.Resolvers.Organization().InitialAccountID(ctx, obj)") {
		return source, nil
	}
	changes := []replacement{
		{old: "type OrganizationResolver interface {\n\tMemberships(ctx context.Context, obj *Organization)", new: "type OrganizationResolver interface {\n\tInitialAccountID(ctx context.Context, obj *Organization) (*string, error)\n\tMemberships(ctx context.Context, obj *Organization)"},
		{old: "return ec.fieldContext_Organization_initialAccountId(ctx, field)\n\t\t},\n\t\tfunc(ctx context.Context) (any, error) {\n\t\t\treturn obj.InitialAccountID, nil", new: "return ec.fieldContext_Organization_initialAccountId(ctx, field)\n\t\t},\n\t\tfunc(ctx context.Context) (any, error) {\n\t\t\treturn ec.Resolvers.Organization().InitialAccountID(ctx, obj)"},
	}
	for _, change := range changes {
		if strings.Count(source, change.old) != 1 {
			return "", fmt.Errorf("initial account projection patch contract changed: %q", change.old)
		}
		source = strings.Replace(source, change.old, change.new, 1)
	}
	return source, nil
}

func patchHTTPHandler(source string) (string, error) {
	if strings.Contains(source, "gqlHandler.SetErrorPresenter(auth.PresentError)") {
		return source, nil
	}
	changes := []replacement{
		{old: `"github.com/99designs/gqlgen/graphql/handler"`, new: `"github.com/99designs/gqlgen/graphql/handler"
	"base-engine/auth"`},
		{old: `gqlHandler := handler.NewDefaultServer(executableSchema)`, new: `gqlHandler := handler.NewDefaultServer(executableSchema)
	gqlHandler.SetErrorPresenter(auth.PresentError)`},
	}
	for _, change := range changes {
		if strings.Count(source, change.old) != 1 {
			return "", fmt.Errorf("HTTP handler patch contract changed: %q", change.old)
		}
		source = strings.Replace(source, change.old, change.new, 1)
	}
	return source, nil
}

func patchFile(path string, transform func(string) (string, error)) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	patched, err := transform(string(source))
	if err != nil {
		return err
	}
	formatted, err := format.Source([]byte(patched))
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0o644)
}

func patchResolver(source string) (string, error) {
	if strings.Contains(source, "services Dependencies") && strings.Contains(source, "authorization.RegisterHandlers") {
		source = strings.Replace(source, "authorization.RegisterHandlers(defaults)", "authorization.RegisterHandlers(defaults, authorization.HandlerDependencies{Sessions: services.Sessions})", 1)
		return source, nil
	}
	for _, change := range resolverReplacements() {
		if strings.Count(source, change.old) != 1 {
			return "", fmt.Errorf("resolver patch contract changed: %q", change.old)
		}
		source = strings.Replace(source, change.old, change.new, 1)
	}
	return source, nil
}

func resolverReplacements() []replacement {
	return []replacement{
		{old: `import (
	"base-engine/gen"
)`, new: `import (
	"base-engine/gen"
	"base-engine/src/services/authorization"
)`},
		{old: `func NewResolver(db *gen.DB, ec *gen.EventController) *Resolver {
	handlers := gen.DefaultResolutionHandlers()
	return &Resolver{&gen.GeneratedResolver{Handlers: handlers, DB: db, EventController: ec}}
}`, new: `func NewResolver(db *gen.DB, ec *gen.EventController, services Dependencies) *Resolver {
	defaults := gen.DefaultResolutionHandlers()
		handlers := authorization.RegisterHandlers(defaults, authorization.HandlerDependencies{Sessions: services.Sessions})
	if err := authorization.ValidateHandlerCoverage(handlers, defaults); err != nil {
		panic(err)
	}
	return &Resolver{GeneratedResolver: &gen.GeneratedResolver{Handlers: handlers, DB: db, EventController: ec}, Services: services}
}`},
		{old: "type Resolver struct {\n\t*gen.GeneratedResolver\n}", new: "type Resolver struct {\n\t*gen.GeneratedResolver\n\tServices Dependencies\n}"},
		{old: "type MutationResolver struct {\n\t*gen.GeneratedMutationResolver\n}", new: "type MutationResolver struct {\n\t*gen.GeneratedMutationResolver\n\tServices Dependencies\n}"},
		{old: "type QueryResolver struct {\n\t*gen.GeneratedQueryResolver\n}", new: "type QueryResolver struct {\n\t*gen.GeneratedQueryResolver\n\tServices Dependencies\n}"},
		{old: "type SubscriptionResolver struct {\n\t*gen.GeneratedSubscriptionResolver\n}", new: "type SubscriptionResolver struct {\n\t*gen.GeneratedSubscriptionResolver\n\tServices Dependencies\n}"},
		{old: `return &MutationResolver{&gen.GeneratedMutationResolver{GeneratedResolver: r.GeneratedResolver}}`, new: `return &MutationResolver{GeneratedMutationResolver: &gen.GeneratedMutationResolver{GeneratedResolver: r.GeneratedResolver}, Services: r.Services}`},
		{old: `return &QueryResolver{&gen.GeneratedQueryResolver{GeneratedResolver: r.GeneratedResolver}}`, new: `return &QueryResolver{GeneratedQueryResolver: &gen.GeneratedQueryResolver{GeneratedResolver: r.GeneratedResolver}, Services: r.Services}`},
		{old: `return &SubscriptionResolver{&gen.GeneratedSubscriptionResolver{GeneratedResolver: r.GeneratedResolver}}`, new: `return &SubscriptionResolver{GeneratedSubscriptionResolver: &gen.GeneratedSubscriptionResolver{GeneratedResolver: r.GeneratedResolver}, Services: r.Services}`},
	}
}

func patchLegacyPassword(source string) (string, error) {
	if !strings.Contains(source, `"crypto/md5"`) {
		return source, nil
	}
	legacy := `// 目前前端密码，使用EncryptPassword相同加密方法传过来，再次加密
// 密码加密
func EncryptPassword(text string) string {
	password := EncryptMd5(text)
	password = base64.StdEncoding.EncodeToString([]byte(password))
	password = EncryptMd5(password)
	return password
}

// MD5验证
func DecryptPassword(text, md5Text string) bool {
	password := EncryptMd5(text)
	password = base64.StdEncoding.EncodeToString([]byte(password))
	password = EncryptMd5(password)
	return password == md5Text
}

// MD5加密
func EncryptMd5(text string) string {
	md5Data := fmt.Sprintf("%x", md5.Sum([]byte(text)))
	return md5Data
}

`
	if strings.Count(source, legacy) != 1 {
		return "", fmt.Errorf("legacy password patch contract changed")
	}
	source = strings.Replace(source, legacy, "", 1)
	source = strings.Replace(source, "\t\"crypto/md5\"\n", "", 1)
	return strings.Replace(source, "\t\"fmt\"\n", "", 1), nil
}
