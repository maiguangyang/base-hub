package productcatalog

import "base-engine/gen"

func activePackagesMatchTemplate(packages []gen.ProductPackage, version int64, templates []gen.ProductPackageTemplate) bool {
	byTemplate, count, valid := indexActiveTemplatePackages(packages, version)
	if !valid || count != len(templates) || len(byTemplate) != len(templates) {
		return false
	}
	for _, template := range templates {
		pack, found := byTemplate[template.ID]
		if !found || !matchesTemplatePackage(pack, template, byTemplate) {
			return false
		}
	}
	return true
}

func indexActiveTemplatePackages(packages []gen.ProductPackage, version int64) (map[string]gen.ProductPackage, int, bool) {
	byTemplate := make(map[string]gen.ProductPackage, len(packages))
	activeCount := 0
	for _, pack := range packages {
		if !pack.Enabled || pack.PackageSetVersion != version {
			continue
		}
		activeCount++
		if pack.TemplateID == nil {
			return nil, 0, false
		}
		byTemplate[*pack.TemplateID] = pack
	}
	return byTemplate, activeCount, true
}

func matchesTemplatePackage(pack gen.ProductPackage, template gen.ProductPackageTemplate, byTemplate map[string]gen.ProductPackage) bool {
	if pack.Name != template.Name || !sameQuantity(pack.ContainsQuantity, template.ContainsQuantity) {
		return false
	}
	if template.ContainsPackageID == nil {
		return pack.ContainsPackageID == nil
	}
	child, found := byTemplate[*template.ContainsPackageID]
	return found && pack.ContainsPackageID != nil && *pack.ContainsPackageID == child.ID
}

func sameQuantity(a, b *int64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}
