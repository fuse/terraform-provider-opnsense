package routing

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newGatewayResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	// Data source intentionally omitted for the initial Gateway release;
	// import + resource read cover the Constellio cutover use case.
	return nil
}
