package acmeclient

import (
	"github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type certificateResourceModel struct {
	Enabled          types.Bool   `tfsdk:"enabled"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	AltNames         types.Set    `tfsdk:"alt_names"`
	Account          types.String `tfsdk:"account"`
	ValidationMethod types.String `tfsdk:"validation_method"`
	KeyLength        types.String `tfsdk:"key_length"`
	Ocsp             types.Bool   `tfsdk:"ocsp"`
	AutoRenewal      types.Bool   `tfsdk:"auto_renewal"`
	RenewInterval    types.String `tfsdk:"renew_interval"`
	Aliasmode        types.String `tfsdk:"aliasmode"`
	Domainalias      types.String `tfsdk:"domainalias"`
	Challengealias   types.String `tfsdk:"challengealias"`
	RestartActions   types.Set    `tfsdk:"restart_actions"`
	CertRefId        types.String `tfsdk:"cert_ref_id"`
	LastUpdate       types.String `tfsdk:"last_update"`
	StatusCode       types.String `tfsdk:"status_code"`
	StatusLastUpdate types.String `tfsdk:"status_last_update"`
	Id               types.String `tfsdk:"id"`
}

func certificateResourceSchema() schema.Schema {
	emptySet := setdefault.StaticValue(tools.EmptySetValue(types.StringType))

	return schema.Schema{
		MarkdownDescription: "Manages an ACME Client certificate in the OPNsense os-acme-client plugin. Referenced accounts, validations and restart actions must be enabled (OPNsense relation filters).",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this certificate. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Common Name (CN) for the certificate.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for this certificate.",
				Optional:            true,
			},
			"alt_names": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Subject Alternative Names (DNS names).",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
			},
			"account": schema.StringAttribute{
				MarkdownDescription: "UUID of the ACME account used to request the certificate.",
				Required:            true,
			},
			"validation_method": schema.StringAttribute{
				MarkdownDescription: "UUID of the ACME validation method used for challenges.",
				Required:            true,
			},
			"key_length": schema.StringAttribute{
				MarkdownDescription: "Private key length or type (`key_2048`, `key_3072`, `key_4096`, `key_ec256`, `key_ec384`). Defaults to `key_2048`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("key_2048"),
			},
			"ocsp": schema.BoolAttribute{
				MarkdownDescription: "Enable OCSP Must-Staple. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auto_renewal": schema.BoolAttribute{
				MarkdownDescription: "Automatically renew the certificate. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"renew_interval": schema.StringAttribute{
				MarkdownDescription: "Days before expiry to attempt renewal. Defaults to `60`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("60"),
			},
			"aliasmode": schema.StringAttribute{
				MarkdownDescription: "DNS alias mode (`none`, `automatic`, `domain`, `challenge`). Defaults to `none`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("none"),
			},
			"domainalias": schema.StringAttribute{
				MarkdownDescription: "Domain alias used when `aliasmode` is `domain`.",
				Optional:            true,
			},
			"challengealias": schema.StringAttribute{
				MarkdownDescription: "Challenge alias used when `aliasmode` is `challenge`.",
				Optional:            true,
			},
			"restart_actions": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "UUIDs of ACME restart actions to run after issuance or renewal.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
			},
			"cert_ref_id": schema.StringAttribute{
				MarkdownDescription: "Trust store certificate reference ID after successful issuance (read-only).",
				Computed:            true,
			},
			"last_update": schema.StringAttribute{
				MarkdownDescription: "Timestamp of the last certificate update (read-only).",
				Computed:            true,
			},
			"status_code": schema.StringAttribute{
				MarkdownDescription: "Last known certificate status code (read-only).",
				Computed:            true,
			},
			"status_last_update": schema.StringAttribute{
				MarkdownDescription: "Timestamp of the last certificate status update (read-only).",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the ACME certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func certificateDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Retrieves an ACME Client certificate from the OPNsense os-acme-client plugin by UUID.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the ACME certificate.",
			},
			"enabled": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this certificate is enabled.",
			},
			"name": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Common Name (CN) for the certificate.",
			},
			"description": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional description for this certificate.",
			},
			"alt_names": dschema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Subject Alternative Names (DNS names).",
			},
			"account": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the ACME account used to request the certificate.",
			},
			"validation_method": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the ACME validation method used for challenges.",
			},
			"key_length": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Private key length or type.",
			},
			"ocsp": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether OCSP Must-Staple is enabled.",
			},
			"auto_renewal": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether automatic renewal is enabled.",
			},
			"renew_interval": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Days before expiry to attempt renewal.",
			},
			"aliasmode": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DNS alias mode.",
			},
			"domainalias": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Domain alias used when `aliasmode` is `domain`.",
			},
			"challengealias": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Challenge alias used when `aliasmode` is `challenge`.",
			},
			"restart_actions": dschema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "UUIDs of ACME restart actions to run after issuance or renewal.",
			},
			"cert_ref_id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Trust store certificate reference ID after successful issuance.",
			},
			"last_update": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of the last certificate update.",
			},
			"status_code": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last known certificate status code.",
			},
			"status_last_update": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of the last certificate status update.",
			},
		},
	}
}

func convertCertificateSchemaToStruct(d *certificateResourceModel) (*acmeclient.Certificate, error) {
	return &acmeclient.Certificate{
		Enabled:          tools.BoolToString(d.Enabled.ValueBool()),
		Name:             d.Name.ValueString(),
		Description:      d.Description.ValueString(),
		AltNames:         api.SelectedMapList(tools.SetToStringSlice(d.AltNames)),
		Account:          api.SelectedMap(d.Account.ValueString()),
		ValidationMethod: api.SelectedMap(d.ValidationMethod.ValueString()),
		KeyLength:        api.SelectedMap(d.KeyLength.ValueString()),
		Ocsp:             tools.BoolToString(d.Ocsp.ValueBool()),
		RestartActions:   api.SelectedMapList(tools.SetToStringSlice(d.RestartActions)),
		AutoRenewal:      tools.BoolToString(d.AutoRenewal.ValueBool()),
		RenewInterval:    d.RenewInterval.ValueString(),
		Aliasmode:        api.SelectedMap(d.Aliasmode.ValueString()),
		Domainalias:      d.Domainalias.ValueString(),
		Challengealias:   d.Challengealias.ValueString(),
	}, nil
}

func convertCertificateStructToSchema(d *acmeclient.Certificate) (*certificateResourceModel, error) {
	return &certificateResourceModel{
		Enabled:          types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:             types.StringValue(d.Name),
		Description:      tools.StringOrNull(d.Description),
		AltNames:         tools.StringSliceToSet([]string(d.AltNames)),
		Account:          types.StringValue(d.Account.String()),
		ValidationMethod: types.StringValue(d.ValidationMethod.String()),
		KeyLength:        types.StringValue(d.KeyLength.String()),
		Ocsp:             types.BoolValue(tools.StringToBool(d.Ocsp)),
		AutoRenewal:      types.BoolValue(tools.StringToBool(d.AutoRenewal)),
		RenewInterval:    types.StringValue(d.RenewInterval),
		Aliasmode:        types.StringValue(d.Aliasmode.String()),
		Domainalias:      tools.StringOrNull(d.Domainalias),
		Challengealias:   tools.StringOrNull(d.Challengealias),
		RestartActions:   tools.StringSliceToSet([]string(d.RestartActions)),
		CertRefId:        tools.StringOrNull(d.CertRefId),
		LastUpdate:       tools.StringOrNull(d.LastUpdate),
		StatusCode:       tools.StringOrNull(d.StatusCode),
		StatusLastUpdate: tools.StringOrNull(d.StatusLastUpdate),
	}, nil
}
