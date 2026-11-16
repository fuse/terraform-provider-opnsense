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

var _ resource.Resource = &validationResource{}
var _ resource.ResourceWithConfigure = &validationResource{}
var _ resource.ResourceWithImportState = &validationResource{}

func newValidationResource() resource.Resource {
	return &validationResource{}
}

type validationResource struct {
	client opnsense.Client
}

func (r *validationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acmeclient_validation"
}

func (r *validationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = validationResourceSchema()
}

func (r *validationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *validationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *validationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	validation, err := convertValidationSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME validation, got error: %s", err))
		return
	}

	id, err := r.client.Acmeclient().AddValidation(ctx, validation)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)
			if readStruct, readErr := r.client.Acmeclient().GetValidation(ctx, id); readErr == nil {
				if readModel, convErr := convertValidationStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					preserveValidationSensitive(readModel, data)
					data = readModel
				}
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ACME validation, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	validationStruct, err := r.client.Acmeclient().GetValidation(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME validation after creation, got error: %s", err))
		return
	}

	validationModel, err := convertValidationStructToSchema(validationStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME validation after creation, got error: %s", err))
		return
	}

	validationModel.Id = data.Id
	preserveValidationSensitive(validationModel, data)

	tflog.Trace(ctx, "created acmeclient_validation resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &validationModel)...)
}

func (r *validationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *validationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	validation, err := r.client.Acmeclient().GetValidation(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "acmeclient validation not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME validation, got error: %s", err))
		return
	}

	validationModel, err := convertValidationStructToSchema(validation)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME validation, got error: %s", err))
		return
	}

	validationModel.Id = data.Id
	preserveValidationSensitive(validationModel, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &validationModel)...)
}

func (r *validationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *validationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existing, err := r.client.Acmeclient().GetValidation(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME validation before update, got error: %s", err))
		return
	}

	planned, err := convertValidationSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME validation, got error: %s", err))
		return
	}

	overlayValidationManagedFields(existing, planned)

	err = r.client.Acmeclient().UpdateValidation(ctx, data.Id.ValueString(), existing)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update ACME validation, got error: %s", err))
		return
	}

	validationStruct, err := r.client.Acmeclient().GetValidation(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ACME validation after update, got error: %s", err))
		return
	}

	validationModel, err := convertValidationStructToSchema(validationStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ACME validation after update, got error: %s", err))
		return
	}

	validationModel.Id = data.Id
	preserveValidationSensitive(validationModel, data)

	tflog.Trace(ctx, "updated acmeclient_validation resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &validationModel)...)
}

func (r *validationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *validationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Acmeclient().DeleteValidation(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete ACME validation, got error: %s", err))
		return
	}
}

func (r *validationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func preserveValidationSensitive(model *validationResourceModel, state *validationResourceModel) {
	if model.DnsCfKey.IsNull() || model.DnsCfKey.ValueString() == "" {
		if !state.DnsCfKey.IsNull() && !state.DnsCfKey.IsUnknown() && state.DnsCfKey.ValueString() != "" {
			model.DnsCfKey = state.DnsCfKey
		}
	}
	if model.DnsCfToken.IsNull() || model.DnsCfToken.ValueString() == "" {
		if !state.DnsCfToken.IsNull() && !state.DnsCfToken.IsUnknown() && state.DnsCfToken.ValueString() != "" {
			model.DnsCfToken = state.DnsCfToken
		}
	}
}
