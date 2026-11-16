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
				MarkdownDescription: "Seconds to wait after publishing DNS records before validation. Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
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

func validationDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Retrieves an ACME Client validation method from the OPNsense os-acme-client plugin by UUID.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the validation method.",
			},
			"enabled": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this validation method is enabled.",
			},
			"name": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Display name for the validation method.",
			},
			"description": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional description for this validation method.",
			},
			"method": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Challenge type.",
			},
			"dns_service": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DNS-01 provider identifier.",
			},
			"dns_sleep": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Seconds to wait after publishing DNS records before validation.",
			},
			"dns_cf_email": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Cloudflare account email (Global API Key auth).",
			},
			"dns_cf_key": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Cloudflare Global API Key.",
			},
			"dns_cf_token": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Cloudflare API token.",
			},
			"dns_cf_account_id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional Cloudflare account ID.",
			},
			"dns_cf_zone_id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional Cloudflare zone ID to limit token scope.",
			},
		},
	}
}

// convertValidationSchemaToStruct sets only Terraform-managed fields.
// Unmanaged dns_* stubs stay at zero values (Create) or are overlaid onto an
// existing Get result (Update) via overlayValidationManagedFields.
func convertValidationSchemaToStruct(d *validationResourceModel) (*acmeclient.Validation, error) {
	dnsSleep := d.DnsSleep.ValueString()
	if dnsSleep == "" {
		dnsSleep = "0"
	}

	// OPNsense requires http_service / tlsalpn_service even for dns01 methods.
	return &acmeclient.Validation{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Name:           d.Name.ValueString(),
		Description:    d.Description.ValueString(),
		Method:         api.SelectedMap(d.Method.ValueString()),
		HttpService:    api.SelectedMap("opnsense"),
		TlsalpnService: api.SelectedMap("acme"),
		DnsService:     api.SelectedMap(d.DnsService.ValueString()),
		DnsSleep:       dnsSleep,
		DnsCfEmail:     d.DnsCfEmail.ValueString(),
		DnsCfKey:       d.DnsCfKey.ValueString(),
		DnsCfToken:     d.DnsCfToken.ValueString(),
		DnsCfAccountId: d.DnsCfAccountId.ValueString(),
		DnsCfZoneId:    d.DnsCfZoneId.ValueString(),
	}, nil
}

func convertValidationStructToSchema(d *acmeclient.Validation) (*validationResourceModel, error) {
	dnsSleep := d.DnsSleep
	if dnsSleep == "" {
		dnsSleep = "0"
	}

	return &validationResourceModel{
		Enabled:        types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:           types.StringValue(d.Name),
		Description:    tools.StringOrNull(d.Description),
		Method:         types.StringValue(d.Method.String()),
		DnsService:     types.StringValue(d.DnsService.String()),
		DnsSleep:       types.StringValue(dnsSleep),
		DnsCfEmail:     tools.StringOrNull(d.DnsCfEmail),
		DnsCfKey:       tools.StringOrNull(d.DnsCfKey),
		DnsCfToken:     tools.StringOrNull(d.DnsCfToken),
		DnsCfAccountId: tools.StringOrNull(d.DnsCfAccountId),
		DnsCfZoneId:    tools.StringOrNull(d.DnsCfZoneId),
	}, nil
}

// overlayValidationManagedFields copies Terraform-managed fields onto an
// existing Validation. HttpService and TlsalpnService are left untouched so
// values set in the OPNsense UI are preserved on update (Create still sets
// dns01-compatible defaults via convertValidationSchemaToStruct).
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
