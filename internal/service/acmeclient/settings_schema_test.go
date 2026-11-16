package acmeclient

import (
	"testing"

	opnacme "github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertSettingsStructToSchema(t *testing.T) {
	settings := &opnacme.Settings{
		Enabled:            "1",
		AutoRenewal:        "1",
		UpdateCron:         api.SelectedMap("cron-uuid"),
		Environment:        api.SelectedMap("prod"),
		ChallengePort:      "43580",
		TLSChallengePort:   "43581",
		RestartTimeout:     "600",
		HaproxyIntegration: "0",
		LogLevel:           api.SelectedMap("normal"),
		ShowIntro:          "1",
	}

	model, err := convertSettingsStructToSchema(settings)
	require.NoError(t, err)
	require.True(t, model.Enabled.ValueBool())
	require.True(t, model.AutoRenewal.ValueBool())
	require.Equal(t, "prod", model.Environment.ValueString())
	require.Equal(t, "43580", model.ChallengePort.ValueString())
	require.Equal(t, "43581", model.TLSChallengePort.ValueString())
	require.Equal(t, "600", model.RestartTimeout.ValueString())
	require.Equal(t, "normal", model.LogLevel.ValueString())
}

func TestApplySettingsManagedFieldsPreservesUnmanaged(t *testing.T) {
	dst := &opnacme.Settings{
		Enabled:            "0",
		AutoRenewal:        "0",
		UpdateCron:         api.SelectedMap("keep-cron"),
		Environment:        api.SelectedMap("stg"),
		ChallengePort:      "1",
		TLSChallengePort:   "2",
		RestartTimeout:     "10",
		HaproxyIntegration: "1",
		HaproxyAclRef:      api.SelectedMap("acl"),
		HaproxyActionRef:   api.SelectedMap("action"),
		HaproxyServerRef:   api.SelectedMap("server"),
		HaproxyBackendRef:  api.SelectedMap("backend"),
		LogLevel:           api.SelectedMap("debug"),
		ShowIntro:          "1",
	}

	plan := &settingsResourceModel{
		Enabled:          types.BoolValue(true),
		AutoRenewal:      types.BoolValue(true),
		Environment:      types.StringValue("prod"),
		ChallengePort:    types.StringValue("43580"),
		TLSChallengePort: types.StringValue("43581"),
		RestartTimeout:   types.StringValue("600"),
		LogLevel:         types.StringValue("extended"),
	}

	applySettingsManagedFields(dst, plan)

	require.Equal(t, "1", dst.Enabled)
	require.Equal(t, "1", dst.AutoRenewal)
	require.Equal(t, api.SelectedMap("prod"), dst.Environment)
	require.Equal(t, "43580", dst.ChallengePort)
	require.Equal(t, "43581", dst.TLSChallengePort)
	require.Equal(t, "600", dst.RestartTimeout)
	require.Equal(t, api.SelectedMap("extended"), dst.LogLevel)

	require.Equal(t, api.SelectedMap("keep-cron"), dst.UpdateCron)
	require.Equal(t, "1", dst.HaproxyIntegration)
	require.Equal(t, api.SelectedMap("acl"), dst.HaproxyAclRef)
	require.Equal(t, api.SelectedMap("action"), dst.HaproxyActionRef)
	require.Equal(t, api.SelectedMap("server"), dst.HaproxyServerRef)
	require.Equal(t, api.SelectedMap("backend"), dst.HaproxyBackendRef)
	require.Equal(t, "1", dst.ShowIntro)
}
