package src

import "context"

func (r *QueryResolver) HqGenerateProductPackageBarcode(ctx context.Context) (string, error) {
 principal,err:=requireResolverPrincipal(ctx)
 if err!=nil {return "",err}
 return r.Services.ProductCatalog.GeneratePackageBarcode(ctx,principal)
}
