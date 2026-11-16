package acmeclient

import (
	"testing"

	opnacme "github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertActionSchemaToStruct(t *testing.T) {
	model := &actionResourceModel{
		Enabled:     types.BoolValue(true),
		Name:        types.StringValue("reload-gui"),
		Description: types.StringValue("restart web gui"),
		Type:        types.StringValue("configd_restart_gui"),
	}

	action, err := convertActionSchemaToStruct(model)
	require.NoError(t, err)
	require.Equal(t, "1", action.Enabled)
	require.Equal(t, "reload-gui", action.Name)
	require.Equal(t, api.SelectedMap("configd_restart_gui"), action.Type)
	require.Empty(t, action.SftpHost)
}

func TestOverlayActionManagedFields(t *testing.T) {
	existing := &opnacme.Action{
		Enabled:  "1",
		Name:     "old",
		Type:     api.SelectedMap("configd_restart_gui"),
		SftpHost: "keep.example.com",
	}
	planned := &opnacme.Action{
		Enabled: "0",
		Name:    "new",
		Type:    api.SelectedMap("configd_restart_gui"),
	}

	overlayActionManagedFields(existing, planned)
	require.Equal(t, "0", existing.Enabled)
	require.Equal(t, "new", existing.Name)
	require.Equal(t, "keep.example.com", existing.SftpHost)
}

func TestConvertActionStructToSchema(t *testing.T) {
	action := &opnacme.Action{
		Enabled:     "1",
		Name:        "reload-gui",
		Description: "restart web gui",
		Type:        api.SelectedMap("configd_restart_gui"),
	}

	model, err := convertActionStructToSchema(action)
	require.NoError(t, err)
	require.True(t, model.Enabled.ValueBool())
	require.Equal(t, "reload-gui", model.Name.ValueString())
	require.Equal(t, "configd_restart_gui", model.Type.ValueString())
}
