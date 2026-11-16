package acmeclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &actionDataSource{}
var _ datasource.DataSourceWithConfigure = &actionDataSource{}

func newActionDataSource() datasource.DataSource {
	return &actionDataSource{}
}

type actionDataSource struct {
	client opnsense.Client
}

func (d *actionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acmeclient_action"
}

func (d *actionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = actionDataSourceSchema()
}

func (d *actionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = opnsense.NewClient(apiClient)
}

func (d *actionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *actionResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, err := d.client.Acmeclient().GetAction(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			resp.Diagnostics.AddError("Not Found",
				fmt.Sprintf("ACME action with ID %s not found", data.Id.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME action, got error: %s", err))
		return
	}

	actionModel, err := convertActionStructToSchema(action)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME action, got error: %s", err))
		return
	}

	actionModel.Id = data.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &actionModel)...)
}
