package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &customExampleProvider{}
)

// customExampleProvider is the provider implementation.
type customExampleProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &customExampleProvider{
			version: version,
		}
	}
}

// customExampleProviderModel maps provider schema data to a Go type.
type customExampleProviderModel struct {
	Username  types.String `tfsdk:"username"`
	Passsword types.String `tfsdk:"password"`
	Baseurl   types.String `tfsdk:"baseurl"`
}

// Metadata returns the provider type name.
func (p *customExampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "customexample"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *customExampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"baseurl": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *customExampleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config customExampleProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Username",
			"The provider cannot create the Custom Example client as there is an unknown configuration value for the Username. ",
		)
	}
	if config.Passsword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Password value",
			"The provider cannot create the Custom Example client as there is an unknown configuration value for the Password. ",
		)
	}
	if config.Baseurl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("baseurl"),
			"Unknown Base Url value",
			"The provider cannot create the Custom Example client as there is an unknown configuration value for the baseurl. ",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	username := os.Getenv("CUSTOM_EXAMPLE_USERNAME")
	password := os.Getenv("CUSTOM_EXAMPLE_PASSWORD")
	baseurl := os.Getenv("CUSTOM_EXAMPLE_BASEURL")

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Passsword.IsNull() {
		password = config.Passsword.ValueString()
	}

	if !config.Baseurl.IsNull() {
		baseurl = config.Baseurl.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing username",
			"The provider cannot create the custom example client as there is a missing or empty value for the Username. ",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Password",
			"The provider cannot create the custom example client as there is a missing or empty value for the Password. ",
		)
	}

	if baseurl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("baseurl"),
			"Missing baseurl",
			"The provider cannot create the custom example client as there is a missing or empty value for the baseUrl in variable. ",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the client based on the username and password
	// It is the dummy example to create the client with the username and password
	// we can change it to appropriate way to login with appropriate credential type
	// Here any client is not made and provided to the data sources and resources, instead we will just provide the base url

	// client, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String(region), // Specify the AWS region
	// 	Credentials: credentials.NewStaticCredentials(access_key, secret_key, ""),
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// resp.DataSourceData = client
	// resp.ResourceData = client
	resp.DataSourceData = baseurl
	resp.ResourceData = baseurl
}

// DataSources defines the data sources implemented in the provider.
func (p *customExampleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGetToDoDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *customExampleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAddTodoResource,
	}
}
