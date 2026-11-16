package acmeclient

import (
	"github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type actionResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Id          types.String `tfsdk:"id"`
}

func actionResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages an ACME Client restart action in the OPNsense os-acme-client plugin. The schema focuses on `configd_restart_gui`; other action types remain available via the OPNsense UI and are preserved on update.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this restart action. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Display name for the restart action.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for this restart action.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Action type. Defaults to `configd_restart_gui`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("configd_restart_gui"),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the restart action.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func actionDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Retrieves an ACME Client restart action from the OPNsense os-acme-client plugin by UUID.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the restart action.",
			},
			"enabled": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this restart action is enabled.",
			},
			"name": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Display name for the restart action.",
			},
			"description": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional description for this restart action.",
			},
			"type": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Action type.",
			},
		},
	}
}

// convertActionSchemaToStruct sets only Terraform-managed fields.
func convertActionSchemaToStruct(d *actionResourceModel) (*acmeclient.Action, error) {
	return &acmeclient.Action{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Name:        d.Name.ValueString(),
		Description: d.Description.ValueString(),
		Type:        api.SelectedMap(d.Type.ValueString()),
	}, nil
}

func convertActionStructToSchema(d *acmeclient.Action) (*actionResourceModel, error) {
	return &actionResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:        types.StringValue(d.Name),
		Description: tools.StringOrNull(d.Description),
		Type:        types.StringValue(d.Type.String()),
	}, nil
}

func overlayActionManagedFields(dst *acmeclient.Action, src *acmeclient.Action) {
	dst.Enabled = src.Enabled
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Type = src.Type
}
