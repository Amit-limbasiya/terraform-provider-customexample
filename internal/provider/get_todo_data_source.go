package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GetToDoDataSource{}
	_ datasource.DataSourceWithConfigure = &GetToDoDataSource{}
)

// NewGetToDoDataSource is a helper function to simplify the provider implementation.
func NewGetToDoDataSource() datasource.DataSource {
	return &GetToDoDataSource{}
}

// GetToDoDataSource is the data source implementation.
type GetToDoDataSource struct {
	baseurl string
}

// GetToDoDataSourceModel maps the data source schema data.
type ToDoDataSourceModel struct {
	TodoList []string `tfsdk:"todo_list"`
}

// Metadata returns the data source type name.
func (d *GetToDoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_todo"
}

// Schema defines the schema for the data source.
func (d *GetToDoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"todo_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *GetToDoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ToDoDataSourceModel

	res, err := http.Get(fmt.Sprintf("%s/get", d.baseurl))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Hit the Todo List get endpoint",
			err.Error(),
		)
		return
	}

	// Map response body to model
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the response body",
			err.Error(),
		)
		return
	}

	err = json.Unmarshal(bodyBytes, &state.TodoList)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read/Unmarshal Todo List",
			err.Error(),
		)
		return
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *GetToDoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	baseurl, ok := req.ProviderData.(string)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected string as baseurl, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.baseurl = baseurl
}
