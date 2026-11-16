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

var _ datasource.DataSource = &accountDataSource{}
var _ datasource.DataSourceWithConfigure = &accountDataSource{}

func newAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountDataSource struct {
	client opnsense.Client
}

func (d *accountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acmeclient_account"
}

func (d *accountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = accountDataSourceSchema()
}

func (d *accountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *accountResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := d.client.Acmeclient().GetAccount(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			resp.Diagnostics.AddError("Not Found",
				fmt.Sprintf("ACME account with ID %s not found", data.Id.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME account, got error: %s", err))
		return
	}

	accountModel, err := convertAccountStructToSchema(account)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME account, got error: %s", err))
		return
	}

	accountModel.Id = data.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &accountModel)...)
}
