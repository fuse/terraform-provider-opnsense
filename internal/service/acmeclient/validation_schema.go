package acmeclient

import (
	"github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type validationResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Method         types.String `tfsdk:"method"`
	DnsService     types.String `tfsdk:"dns_service"`
	DnsSleep       types.String `tfsdk:"dns_sleep"`
	DnsCfEmail     types.String `tfsdk:"dns_cf_email"`
	DnsCfKey       types.String `tfsdk:"dns_cf_key"`
	DnsCfToken     types.String `tfsdk:"dns_cf_token"`
	DnsCfAccountId types.String `tfsdk:"dns_cf_account_id"`
	DnsCfZoneId    types.String `tfsdk:"dns_cf_zone_id"`
	Id             types.String `tfsdk:"id"`
}

func validationResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages an ACME Client validation method in the OPNsense os-acme-client plugin. The schema focuses on DNS-01 with Cloudflare (`dns_cf`); other DNS providers remain available via the OPNsense UI and are preserved on update.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this validation method. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Display name for the validation method.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for this validation method.",
				Optional:            true,
			},
			"method": schema.StringAttribute{
				MarkdownDescription: "Challenge type. Defaults to `dns01`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("dns01"),
			},
			"dns_service": schema.StringAttribute{
				MarkdownDescription: "DNS-01 provider identifier. Defaults to `dns_cf` (Cloudflare).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("dns_cf"),
			},
			"dns_sleep": schema.StringAttribute{
				MarkdownDescription: "Seconds to wait after publishing DNS records before validation.",
				Optional:            true,
			},
			"dns_cf_email": schema.StringAttribute{
				MarkdownDescription: "Cloudflare account email (Global API Key auth).",
				Optional:            true,
			},
			"dns_cf_key": schema.StringAttribute{
				MarkdownDescription: "Cloudflare Global API Key (sensitive).",
				Optional:            true,
				Sensitive:           true,
			},
			"dns_cf_token": schema.StringAttribute{
				MarkdownDescription: "Cloudflare API token (sensitive). Preferred over Global API Key.",
				Optional:            true,
				Sensitive:           true,
			},
			"dns_cf_account_id": schema.StringAttribute{
				MarkdownDescription: "Optional Cloudflare account ID.",
				Optional:            true,
			},
			"dns_cf_zone_id": schema.StringAttribute{
				MarkdownDescription: "Optional Cloudflare zone ID to limit token scope.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the validation method.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// convertValidationSchemaToStruct sets only Terraform-managed fields.
// Unmanaged dns_* stubs stay at zero values (Create) or are overlaid onto an
// existing Get result (Update) via overlayValidationManagedFields.
func convertValidationSchemaToStruct(d *validationResourceModel) (*acmeclient.Validation, error) {
	return &acmeclient.Validation{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Name:           d.Name.ValueString(),
		Description:    d.Description.ValueString(),
		Method:         api.SelectedMap(d.Method.ValueString()),
		DnsService:     api.SelectedMap(d.DnsService.ValueString()),
		DnsSleep:       d.DnsSleep.ValueString(),
		DnsCfEmail:     d.DnsCfEmail.ValueString(),
		DnsCfKey:       d.DnsCfKey.ValueString(),
		DnsCfToken:     d.DnsCfToken.ValueString(),
		DnsCfAccountId: d.DnsCfAccountId.ValueString(),
		DnsCfZoneId:    d.DnsCfZoneId.ValueString(),
	}, nil
}

func convertValidationStructToSchema(d *acmeclient.Validation) (*validationResourceModel, error) {
	return &validationResourceModel{
		Enabled:        types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:           types.StringValue(d.Name),
		Description:    tools.StringOrNull(d.Description),
		Method:         types.StringValue(d.Method.String()),
		DnsService:     types.StringValue(d.DnsService.String()),
		DnsSleep:       tools.StringOrNull(d.DnsSleep),
		DnsCfEmail:     tools.StringOrNull(d.DnsCfEmail),
		DnsCfKey:       tools.StringOrNull(d.DnsCfKey),
		DnsCfToken:     tools.StringOrNull(d.DnsCfToken),
		DnsCfAccountId: tools.StringOrNull(d.DnsCfAccountId),
		DnsCfZoneId:    tools.StringOrNull(d.DnsCfZoneId),
	}, nil
}

func overlayValidationManagedFields(dst *acmeclient.Validation, src *acmeclient.Validation) {
	dst.Enabled = src.Enabled
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Method = src.Method
	dst.DnsService = src.DnsService
	dst.DnsSleep = src.DnsSleep
	dst.DnsCfEmail = src.DnsCfEmail
	dst.DnsCfKey = src.DnsCfKey
	dst.DnsCfToken = src.DnsCfToken
	dst.DnsCfAccountId = src.DnsCfAccountId
	dst.DnsCfZoneId = src.DnsCfZoneId
}
