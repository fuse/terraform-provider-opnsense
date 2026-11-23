package interfaces

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// vipResourceModel describes the resource data model.
type vipResourceModel struct {
	Mode        types.String `tfsdk:"mode"`
	Interface   types.String `tfsdk:"interface"`
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`
	Password    types.String `tfsdk:"password"`
	Vhid        types.String `tfsdk:"vhid"`
	Advbase     types.String `tfsdk:"advbase"`
	Advskew     types.String `tfsdk:"advskew"`
	Nobind      types.Bool   `tfsdk:"nobind"`
	Noexpand    types.Bool   `tfsdk:"noexpand"`
	Nosync      types.Bool   `tfsdk:"nosync"`
	Peer        types.String `tfsdk:"peer"`
	Peer6       types.String `tfsdk:"peer6"`
	Id          types.String `tfsdk:"id"`
}

func vipResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Virtual IPs allow an OPNsense firewall to assign multiple IP addresses to the same network interface, including CARP.",

		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode of the VIP. One of `ipalias`, `carp`, or `proxyarp`. `proxyarp` cannot be bound to by anything running on the firewall, such as IPsec, OpenVPN, etc. In most cases an `ipalias` should be used. CARP requires `password` and `vhid`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ipalias"),
				Validators: []validator.String{
					stringvalidator.OneOf("ipalias", "carp", "proxyarp"),
				},
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Choose which interface this VIP applies to.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("wan"),
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "Provide an address and subnet to use. (e.g 192.168.0.1/24)",
				Required:            true,
				Validators: []validator.String{
					validators.CIDR(),
				},
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "For some interface types a gateway is required to configure an IP Alias (ppp/pppoe/tun), leave this field empty for all other interface types.",
				Optional:            true,
				Validators: []validator.String{
					validators.IpOrCIDR(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "CARP password (shared across peers with the same VHID). Required when `mode` is `carp`.",
				Optional:            true,
				Sensitive:           true,
			},
			"vhid": schema.StringAttribute{
				MarkdownDescription: "CARP Virtual Host ID (1-255). Required when `mode` is `carp`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"advbase": schema.StringAttribute{
				MarkdownDescription: "CARP advertisement base (seconds). Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"advskew": schema.StringAttribute{
				MarkdownDescription: "CARP advertisement skew (0-254). Lower values prefer MASTER. Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
			},
			"nobind": schema.BoolAttribute{
				MarkdownDescription: "Prevent binding services to this VIP address. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"noexpand": schema.BoolAttribute{
				MarkdownDescription: "Do not expand this VIP into firewall rules automatically. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"nosync": schema.BoolAttribute{
				MarkdownDescription: "Exclude this VIP from XMLRPC sync. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"peer": schema.StringAttribute{
				MarkdownDescription: "Optional CARP peer IPv4 address (leave empty for multicast).",
				Optional:            true,
			},
			"peer6": schema.StringAttribute{
				MarkdownDescription: "Optional CARP peer IPv6 address (leave empty for multicast).",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the VIP.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func vipDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Virtual IPs allow an OPNsense firewall to assign multiple IP addresses to the same network interface.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode of the VIP. One of `ipalias`, `carp`, or `proxyarp`. `proxyarp` cannot be bound to by anything running on the firewall, such as IPsec, OpenVPN, etc. In most cases an `ipalias` should be used.",
				Computed:            true,
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Choose which interface this VIP applies to.",
				Computed:            true,
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "Provide an address and subnet to use. (e.g 192.168.0.1/24)",
				Computed:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "For some interface types a gateway is required to configure an IP Alias (ppp/pppoe/tun), leave this field empty for all other interface types.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "CARP password.",
				Computed:            true,
				Sensitive:           true,
			},
			"vhid": schema.StringAttribute{
				MarkdownDescription: "CARP Virtual Host ID.",
				Computed:            true,
			},
			"advbase": schema.StringAttribute{
				MarkdownDescription: "CARP advertisement base.",
				Computed:            true,
			},
			"advskew": schema.StringAttribute{
				MarkdownDescription: "CARP advertisement skew.",
				Computed:            true,
			},
			"nobind": schema.BoolAttribute{
				MarkdownDescription: "Prevent binding services to this VIP address.",
				Computed:            true,
			},
			"noexpand": schema.BoolAttribute{
				MarkdownDescription: "Do not expand this VIP into firewall rules automatically.",
				Computed:            true,
			},
			"nosync": schema.BoolAttribute{
				MarkdownDescription: "Exclude this VIP from XMLRPC sync.",
				Computed:            true,
			},
			"peer": schema.StringAttribute{
				MarkdownDescription: "Optional CARP peer IPv4 address.",
				Computed:            true,
			},
			"peer6": schema.StringAttribute{
				MarkdownDescription: "Optional CARP peer IPv6 address.",
				Computed:            true,
			},
		},
	}
}

func convertVipSchemaToStruct(d *vipResourceModel) (*interfaces.Vip, error) {
	return &interfaces.Vip{
		Description: d.Description.ValueString(),
		Mode:        api.SelectedMap(d.Mode.ValueString()),
		Interface:   api.SelectedMap(d.Interface.ValueString()),
		Network:     d.Network.ValueString(),
		Gateway:     d.Gateway.ValueString(),
		Password:    d.Password.ValueString(),
		Vhid:        d.Vhid.ValueString(),
		Advbase:     d.Advbase.ValueString(),
		Advskew:     d.Advskew.ValueString(),
		Nobind:      tools.BoolToString(d.Nobind.ValueBool()),
		Noexpand:    tools.BoolToString(d.Noexpand.ValueBool()),
		Nosync:      tools.BoolToString(d.Nosync.ValueBool()),
		Peer:        d.Peer.ValueString(),
		Peer6:       d.Peer6.ValueString(),
	}, nil
}

func convertVipStructToSchema(d *interfaces.Vip) (*vipResourceModel, error) {
	advbase := d.Advbase
	if advbase == "" {
		advbase = "1"
	}
	advskew := d.Advskew
	if advskew == "" {
		advskew = "0"
	}

	return &vipResourceModel{
		Mode:        types.StringValue(d.Mode.String()),
		Interface:   types.StringValue(d.Interface.String()),
		Network:     types.StringValue(d.Network),
		Gateway:     tools.StringOrNull(d.Gateway),
		Description: tools.StringOrNull(d.Description),
		Password:    tools.StringOrNull(d.Password),
		Vhid:        types.StringValue(d.Vhid),
		Advbase:     types.StringValue(advbase),
		Advskew:     types.StringValue(advskew),
		Nobind:      types.BoolValue(tools.StringToBool(d.Nobind)),
		Noexpand:    types.BoolValue(tools.StringToBool(d.Noexpand)),
		Nosync:      types.BoolValue(tools.StringToBool(d.Nosync)),
		Peer:        tools.StringOrNull(d.Peer),
		Peer6:       tools.StringOrNull(d.Peer6),
	}, nil
}
