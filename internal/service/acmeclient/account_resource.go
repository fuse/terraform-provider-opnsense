package acmeclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &accountResource{}
var _ resource.ResourceWithConfigure = &accountResource{}
var _ resource.ResourceWithImportState = &accountResource{}

func newAccountResource() resource.Resource {
	return &accountResource{}
}

type accountResource struct {
	client opnsense.Client
}

func (r *accountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acmeclient_account"
}

func (r *accountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = accountResourceSchema()
}

func (r *accountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = opnsense.NewClient(apiClient)
}

func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *accountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := convertAccountSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME account, got error: %s", err))
		return
	}

	id, err := r.client.Acmeclient().AddAccount(ctx, account)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)
			if readStruct, readErr := r.client.Acmeclient().GetAccount(ctx, id); readErr == nil {
				if readModel, convErr := convertAccountStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					preserveAccountSensitive(readModel, data)
					data = readModel
				}
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ACME account, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	accountStruct, err := r.client.Acmeclient().GetAccount(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME account after creation, got error: %s", err))
		return
	}

	accountModel, err := convertAccountStructToSchema(accountStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME account after creation, got error: %s", err))
		return
	}

	accountModel.Id = data.Id
	preserveAccountSensitive(accountModel, data)

	tflog.Trace(ctx, "created acmeclient_account resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &accountModel)...)
}

func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *accountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := r.client.Acmeclient().GetAccount(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "acmeclient account not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
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
	preserveAccountSensitive(accountModel, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &accountModel)...)
}

func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *accountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := convertAccountSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME account, got error: %s", err))
		return
	}

	err = r.client.Acmeclient().UpdateAccount(ctx, data.Id.ValueString(), account)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update ACME account, got error: %s", err))
		return
	}

	accountStruct, err := r.client.Acmeclient().GetAccount(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME account after update, got error: %s", err))
		return
	}

	accountModel, err := convertAccountStructToSchema(accountStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME account after update, got error: %s", err))
		return
	}

	accountModel.Id = data.Id
	preserveAccountSensitive(accountModel, data)

	tflog.Trace(ctx, "updated acmeclient_account resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &accountModel)...)
}

func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *accountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Acmeclient().DeleteAccount(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete ACME account, got error: %s", err))
		return
	}
}

func (r *accountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// preserveAccountSensitive keeps prior sensitive values when the API omits them,
// while preferring a non-empty key returned by the API.
func preserveAccountSensitive(model *accountResourceModel, state *accountResourceModel) {
	if model.AccountKey.IsNull() || model.AccountKey.ValueString() == "" {
		if !state.AccountKey.IsNull() && !state.AccountKey.IsUnknown() && state.AccountKey.ValueString() != "" {
			model.AccountKey = state.AccountKey
		}
	}
	if model.EabHmac.IsNull() || model.EabHmac.ValueString() == "" {
		if !state.EabHmac.IsNull() && !state.EabHmac.IsUnknown() && state.EabHmac.ValueString() != "" {
			model.EabHmac = state.EabHmac
		}
	}
}
