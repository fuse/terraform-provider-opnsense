package routing

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/routing"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type gatewayResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Interface      types.String `tfsdk:"interface"`
	IpProtocol     types.String `tfsdk:"ip_protocol"`
	Gateway        types.String `tfsdk:"gateway"`
	DefaultGateway types.Bool   `tfsdk:"default_gateway"`
	FarGateway     types.Bool   `tfsdk:"far_gateway"`
	MonitorDisable types.Bool   `tfsdk:"monitor_disable"`
	MonitorNoRoute types.Bool   `tfsdk:"monitor_noroute"`
	Monitor        types.String `tfsdk:"monitor"`
	ForceDown      types.Bool   `tfsdk:"force_down"`
	Priority       types.String `tfsdk:"priority"`
	Weight         types.String `tfsdk:"weight"`
	Id             types.String `tfsdk:"id"`
}

func gatewayResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Manages a System gateway (Routing → Gateways) via the OPNsense `routing` API.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this gateway. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Gateway name. Changing this forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description.",
				Optional:            true,
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Interface identifier (e.g. `wan`, `opt1`, `opt2`).",
				Required:            true,
			},
			"ip_protocol": schema.StringAttribute{
				MarkdownDescription: "Address family: `inet` (IPv4) or `inet6` (IPv6). Defaults to `inet`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inet"),
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "Gateway IP address (or `dynamic` when provided by the interface).",
				Required:            true,
			},
			"default_gateway": schema.BoolAttribute{
				MarkdownDescription: "Mark this gateway as a default gateway. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"far_gateway": schema.BoolAttribute{
				MarkdownDescription: "Far gateway (gateway outside the interface subnet). Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"monitor_disable": schema.BoolAttribute{
				MarkdownDescription: "Disable gateway monitoring. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"monitor_noroute": schema.BoolAttribute{
				MarkdownDescription: "Do not create a static route for the monitor IP. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"monitor": schema.StringAttribute{
				MarkdownDescription: "Optional monitor IP address.",
				Optional:            true,
			},
			"force_down": schema.BoolAttribute{
				MarkdownDescription: "Force the gateway to be considered down. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"priority": schema.StringAttribute{
				MarkdownDescription: "Gateway priority (0-255, lower is preferred). Defaults to `255`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("255"),
			},
			"weight": schema.StringAttribute{
				MarkdownDescription: "Gateway weight for multi-WAN. Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the gateway.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertGatewaySchemaToStruct(d *gatewayResourceModel) (*routing.Gateway, error) {
	return &routing.Gateway{
		Disabled:       tools.BoolToString(!d.Enabled.ValueBool()),
		Name:           d.Name.ValueString(),
		Description:    d.Description.ValueString(),
		Interface:      api.SelectedMap(d.Interface.ValueString()),
		IpProtocol:     api.SelectedMap(d.IpProtocol.ValueString()),
		Gateway:        d.Gateway.ValueString(),
		DefaultGateway: tools.BoolToString(d.DefaultGateway.ValueBool()),
		FarGateway:     tools.BoolToString(d.FarGateway.ValueBool()),
		MonitorDisable: tools.BoolToString(d.MonitorDisable.ValueBool()),
		MonitorNoRoute: tools.BoolToString(d.MonitorNoRoute.ValueBool()),
		Monitor:        d.Monitor.ValueString(),
		ForceDown:      tools.BoolToString(d.ForceDown.ValueBool()),
		Priority:       d.Priority.ValueString(),
		Weight:         d.Weight.ValueString(),
	}, nil
}

func convertGatewayStructToSchema(d *routing.Gateway) (*gatewayResourceModel, error) {
	return &gatewayResourceModel{
		Enabled:        types.BoolValue(!tools.StringToBool(d.Disabled)),
		Name:           types.StringValue(d.Name),
		Description:    tools.StringOrNull(d.Description),
		Interface:      types.StringValue(d.Interface.String()),
		IpProtocol:     types.StringValue(d.IpProtocol.String()),
		Gateway:        types.StringValue(d.Gateway),
		DefaultGateway: types.BoolValue(tools.StringToBool(d.DefaultGateway)),
		FarGateway:     types.BoolValue(tools.StringToBool(d.FarGateway)),
		MonitorDisable: types.BoolValue(tools.StringToBool(d.MonitorDisable)),
		MonitorNoRoute: types.BoolValue(tools.StringToBool(d.MonitorNoRoute)),
		Monitor:        tools.StringOrNull(d.Monitor),
		ForceDown:      types.BoolValue(tools.StringToBool(d.ForceDown)),
		Priority:       types.StringValue(d.Priority),
		Weight:         types.StringValue(d.Weight),
	}, nil
}
