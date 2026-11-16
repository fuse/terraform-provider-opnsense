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

type accountResourceModel struct {
	Enabled          types.Bool   `tfsdk:"enabled"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Email            types.String `tfsdk:"email"`
	Ca               types.String `tfsdk:"ca"`
	CustomCa         types.String `tfsdk:"custom_ca"`
	EabKid           types.String `tfsdk:"eab_kid"`
	EabHmac          types.String `tfsdk:"eab_hmac"`
	AccountKey       types.String `tfsdk:"account_key"`
	StatusCode       types.String `tfsdk:"status_code"`
	StatusLastUpdate types.String `tfsdk:"status_last_update"`
	Id               types.String `tfsdk:"id"`
}

func accountResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages an ACME Client account in the OPNsense os-acme-client plugin.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this ACME account. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Display name for the ACME account.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for this ACME account.",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Contact email address registered with the ACME CA.",
				Required:            true,
			},
			"ca": schema.StringAttribute{
				MarkdownDescription: "Certificate authority identifier (e.g. `letsencrypt`, `letsencrypt_test`, `zerossl`). Defaults to `letsencrypt`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("letsencrypt"),
			},
			"custom_ca": schema.StringAttribute{
				MarkdownDescription: "Custom ACME CA directory URL when `ca` is set to a custom value.",
				Optional:            true,
			},
			"eab_kid": schema.StringAttribute{
				MarkdownDescription: "External Account Binding (EAB) key identifier, if required by the CA.",
				Optional:            true,
			},
			"eab_hmac": schema.StringAttribute{
				MarkdownDescription: "External Account Binding (EAB) HMAC key, if required by the CA.",
				Optional:            true,
				Sensitive:           true,
			},
			"account_key": schema.StringAttribute{
				MarkdownDescription: "ACME account private key. Managed by OPNsense when omitted; may be imported or read back from the API.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status_code": schema.StringAttribute{
				MarkdownDescription: "Last known ACME account status code from OPNsense (read-only).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status_last_update": schema.StringAttribute{
				MarkdownDescription: "Timestamp of the last ACME account status update (read-only).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the ACME account.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertAccountSchemaToStruct(d *accountResourceModel) (*acmeclient.Account, error) {
	return &acmeclient.Account{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Name:        d.Name.ValueString(),
		Description: d.Description.ValueString(),
		Email:       d.Email.ValueString(),
		Ca:          api.SelectedMap(d.Ca.ValueString()),
		CustomCa:    d.CustomCa.ValueString(),
		EabKid:      d.EabKid.ValueString(),
		EabHmac:     d.EabHmac.ValueString(),
		Key:         d.AccountKey.ValueString(),
	}, nil
}

func convertAccountStructToSchema(d *acmeclient.Account) (*accountResourceModel, error) {
	return &accountResourceModel{
		Enabled:          types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:             types.StringValue(d.Name),
		Description:      tools.StringOrNull(d.Description),
		Email:            types.StringValue(d.Email),
		Ca:               types.StringValue(d.Ca.String()),
		CustomCa:         tools.StringOrNull(d.CustomCa),
		EabKid:           tools.StringOrNull(d.EabKid),
		EabHmac:          tools.StringOrNull(d.EabHmac),
		AccountKey:       tools.StringOrNull(d.Key),
		StatusCode:       tools.StringOrNull(d.StatusCode),
		StatusLastUpdate: tools.StringOrNull(d.StatusLastUpdate),
	}, nil
}
