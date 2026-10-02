package authorization

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestCustomerCouponDistributionJobGeneratedHandlersAlwaysDeny(t *testing.T) {
	handlers := RegisterHandlers(gen.DefaultResolutionHandlers())
	ctx := couponDistributionJobAdminContext()
	value := reflect.ValueOf(handlers)
	typ := value.Type()
	covered := 0
	for index := 0; index < value.NumField(); index++ {
		name := typ.Field(index).Name
		if !strings.Contains(name, "CustomerCouponDistributionJob") &&
			!strings.Contains(name, "CouponDistributionJobs") {
			continue
		}
		covered++
		fn := value.Field(index)
		args := make([]reflect.Value, fn.Type().NumIn())
		args[0] = reflect.ValueOf(ctx)
		for parameter := 1; parameter < len(args); parameter++ {
			args[parameter] = reflect.Zero(fn.Type().In(parameter))
		}
		outputs := fn.Call(args)
		err, _ := outputs[len(outputs)-1].Interface().(error)
		if auth.ErrorCode(err) != auth.CodePermissionDenied {
			t.Errorf("generated job entry %s returned %v", name, err)
		}
	}
	if covered < 9 {
		t.Fatalf("only %d generated job handlers covered", covered)
	}
}

func TestCustomerCouponDistributionJobRelationshipIDsDenyBeforeLoaders(t *testing.T) {
	resolver := &gen.GeneratedResolver{Handlers: RegisterHandlers(gen.DefaultResolutionHandlers())}
	ctx := couponDistributionJobAdminContext()
	memberResolver := &gen.GeneratedCustomerMemberResolver{GeneratedResolver: resolver}
	if ids, err := memberResolver.CouponDistributionJobsIds(ctx, &gen.CustomerMember{ID: "member"}); auth.ErrorCode(err) != auth.CodePermissionDenied || ids != nil {
		t.Fatalf("member job IDs = %v, %v", ids, err)
	}
	templateResolver := &gen.GeneratedCustomerCouponTemplateResolver{GeneratedResolver: resolver}
	if ids, err := templateResolver.DistributionJobsIds(ctx, &gen.CustomerCouponTemplate{ID: "template"}); auth.ErrorCode(err) != auth.CodePermissionDenied || ids != nil {
		t.Fatalf("template job IDs = %v, %v", ids, err)
	}
}

func couponDistributionJobAdminContext() context.Context {
	organizationID := "org-hq"
	permissions := map[string]struct{}{}
	for _, verb := range []string{"read", "create", "update", "delete"} {
		permissions["customerCouponDistributionJob:"+verb] = struct{}{}
	}
	return auth.WithPrincipal(context.Background(), &auth.WorkspacePrincipal{
		AccountID: "hq-admin", OrganizationID: &organizationID,
		WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: permissions,
	})
}
