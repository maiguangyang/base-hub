package authorization

import (
	"context"
	"reflect"

	"base-engine/gen"
)

type membershipRelationshipFilter struct {
	base    gen.EntityFilter
	roleID  *string
	storeID *string
}

func (filter *membershipRelationshipFilter) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if filter.base != nil {
		if err := filter.base.Apply(ctx, wheres, values, joins); err != nil {
			return err
		}
	}
	memberships := gen.TableName("operator_memberships", ctx)
	if filter.roleID != nil {
		roles := gen.TableName("operator_membership_roles", ctx)
		*wheres = append(*wheres, "EXISTS (SELECT 1 FROM "+roles+" membership_role_filter WHERE membership_role_filter.operator_membership_id = "+memberships+".id AND membership_role_filter.operator_role_id = ?)")
		*values = append(*values, *filter.roleID)
	}
	if filter.storeID != nil {
		stores := gen.TableName("operator_membership_stores", ctx)
		*wheres = append(*wheres, "EXISTS (SELECT 1 FROM "+stores+" membership_store_filter WHERE membership_store_filter.operator_membership_id = "+memberships+".id AND membership_store_filter.store_id = ?)")
		*values = append(*values, *filter.storeID)
	}
	return nil
}

func splitMembershipRelationshipFilter(client *gen.OperatorMembershipFilterType) (*gen.OperatorMembershipFilterType, *string, *string) {
	if client == nil {
		return nil, nil, nil
	}
	clean := *client
	var roleID, storeID *string
	if exactRoleIDFilter(client.Roles) {
		roleID = client.Roles.ID
		clean.Roles = nil
	}
	if exactStoreIDFilter(client.Stores) {
		storeID = client.Stores.ID
		clean.Stores = nil
	}
	return &clean, roleID, storeID
}

func exactRoleIDFilter(filter *gen.OperatorRoleFilterType) bool {
	if filter == nil || filter.ID == nil {
		return false
	}
	clean := *filter
	clean.ID = nil
	return reflect.DeepEqual(clean, gen.OperatorRoleFilterType{})
}

func exactStoreIDFilter(filter *gen.StoreFilterType) bool {
	if filter == nil || filter.ID == nil {
		return false
	}
	clean := *filter
	clean.ID = nil
	return reflect.DeepEqual(clean, gen.StoreFilterType{})
}
