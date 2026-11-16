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

var _ resource.Resource = &actionResource{}
var _ resource.ResourceWithConfigure = &actionResource{}
var _ resource.ResourceWithImportState = &actionResource{}

func newActionResource() resource.Resource {
	return &actionResource{}
}

type actionResource struct {
	client opnsense.Client
}

func (r *actionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acmeclient_action"
}

func (r *actionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = actionResourceSchema()
}

func (r *actionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *actionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *actionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, err := convertActionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME action, got error: %s", err))
		return
	}

	id, err := r.client.Acmeclient().AddAction(ctx, action)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)
			if readStruct, readErr := r.client.Acmeclient().GetAction(ctx, id); readErr == nil {
				if readModel, convErr := convertActionStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					data = readModel
				}
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ACME action, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	actionStruct, err := r.client.Acmeclient().GetAction(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME action after creation, got error: %s", err))
		return
	}

	actionModel, err := convertActionStructToSchema(actionStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME action after creation, got error: %s", err))
		return
	}

	actionModel.Id = data.Id

	tflog.Trace(ctx, "created acmeclient_action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &actionModel)...)
}

func (r *actionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *actionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, err := r.client.Acmeclient().GetAction(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "acmeclient action not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
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

func (r *actionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *actionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existing, err := r.client.Acmeclient().GetAction(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME action before update, got error: %s", err))
		return
	}

	planned, err := convertActionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME action, got error: %s", err))
		return
	}

	overlayActionManagedFields(existing, planned)

	err = r.client.Acmeclient().UpdateAction(ctx, data.Id.ValueString(), existing)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update ACME action, got error: %s", err))
		return
	}

	actionStruct, err := r.client.Acmeclient().GetAction(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME action after update, got error: %s", err))
		return
	}

	actionModel, err := convertActionStructToSchema(actionStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME action after update, got error: %s", err))
		return
	}

	actionModel.Id = data.Id

	tflog.Trace(ctx, "updated acmeclient_action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &actionModel)...)
}

func (r *actionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *actionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Acmeclient().DeleteAction(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete ACME action, got error: %s", err))
		return
	}
}

func (r *actionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
