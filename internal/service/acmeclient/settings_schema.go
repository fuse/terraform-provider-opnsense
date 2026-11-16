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

// settingsResourceModel describes the ACME Client settings singleton.
// HAProxy integration and UpdateCron are preserved on update but not managed
// by Terraform in this version.
type settingsResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	AutoRenewal      types.Bool   `tfsdk:"auto_renewal"`
	Environment      types.String `tfsdk:"environment"`
	ChallengePort    types.String `tfsdk:"challenge_port"`
	TLSChallengePort types.String `tfsdk:"tls_challenge_port"`
	RestartTimeout   types.String `tfsdk:"restart_timeout"`
	LogLevel         types.String `tfsdk:"log_level"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages ACME Client plugin settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_acmeclient_settings.settings acmeclient_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.\n\n" +
			"HAProxy integration and the automatic renewal cron job reference are preserved on update but are not managed by this resource.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `acmeclient_settings`. Use this value when importing: `terraform import opnsense_acmeclient_settings.settings acmeclient_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the ACME Client plugin. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auto_renewal": schema.BoolAttribute{
				MarkdownDescription: "Automatically renew certificates when close to expiry. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Let's Encrypt environment: `prod` or `stg`. Defaults to `prod`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("prod"),
			},
			"challenge_port": schema.StringAttribute{
				MarkdownDescription: "Local HTTP-01 challenge port. Defaults to `43580`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("43580"),
			},
			"tls_challenge_port": schema.StringAttribute{
				MarkdownDescription: "Local TLS-ALPN-01 challenge port. Defaults to `43581`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("43581"),
			},
			"restart_timeout": schema.StringAttribute{
				MarkdownDescription: "Seconds to wait for restart actions after issuance or renewal. Defaults to `600`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("600"),
			},
			"log_level": schema.StringAttribute{
				MarkdownDescription: "Acme.sh log level (`normal`, `extended`, `debug`, `debug2`, `debug3`). Defaults to `normal`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("normal"),
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads ACME Client plugin settings from the upstream system.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `acmeclient_settings`.",
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the ACME Client plugin is enabled.",
				Computed:            true,
			},
			"auto_renewal": dschema.BoolAttribute{
				MarkdownDescription: "Whether automatic certificate renewal is enabled.",
				Computed:            true,
			},
			"environment": dschema.StringAttribute{
				MarkdownDescription: "Let's Encrypt environment (`prod` or `stg`).",
				Computed:            true,
			},
			"challenge_port": dschema.StringAttribute{
				MarkdownDescription: "Local HTTP-01 challenge port.",
				Computed:            true,
			},
			"tls_challenge_port": dschema.StringAttribute{
				MarkdownDescription: "Local TLS-ALPN-01 challenge port.",
				Computed:            true,
			},
			"restart_timeout": dschema.StringAttribute{
				MarkdownDescription: "Seconds to wait for restart actions after issuance or renewal.",
				Computed:            true,
			},
			"log_level": dschema.StringAttribute{
				MarkdownDescription: "Acme.sh log level.",
				Computed:            true,
			},
		},
	}
}

func convertSettingsStructToSchema(d *acmeclient.Settings) (*settingsResourceModel, error) {
	return &settingsResourceModel{
		Enabled:          types.BoolValue(tools.StringToBool(d.Enabled)),
		AutoRenewal:      types.BoolValue(tools.StringToBool(d.AutoRenewal)),
		Environment:      types.StringValue(string(d.Environment)),
		ChallengePort:    types.StringValue(d.ChallengePort),
		TLSChallengePort: types.StringValue(d.TLSChallengePort),
		RestartTimeout:   types.StringValue(d.RestartTimeout),
		LogLevel:         types.StringValue(string(d.LogLevel)),
	}, nil
}

// applySettingsManagedFields overlays Terraform-managed fields onto an existing
// Settings struct so HAProxy and UpdateCron values are preserved on set.
func applySettingsManagedFields(dst *acmeclient.Settings, d *settingsResourceModel) {
	dst.Enabled = tools.BoolToString(d.Enabled.ValueBool())
	dst.AutoRenewal = tools.BoolToString(d.AutoRenewal.ValueBool())
	dst.Environment = api.SelectedMap(d.Environment.ValueString())
	dst.ChallengePort = d.ChallengePort.ValueString()
	dst.TLSChallengePort = d.TLSChallengePort.ValueString()
	dst.RestartTimeout = d.RestartTimeout.ValueString()
	dst.LogLevel = api.SelectedMap(d.LogLevel.ValueString())
}
