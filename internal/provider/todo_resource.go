package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &addTodoResource{}
	_ resource.ResourceWithConfigure   = &addTodoResource{}
	_ resource.ResourceWithImportState = &addTodoResource{}
)

// NewAddTodoResource is a helper function to simplify the provider implementation.
func NewAddTodoResource() resource.Resource {
	return &addTodoResource{}
}

// addTodoResource is the resource implementation.
type addTodoResource struct {
	baseurl string
}

// orderResourceModel maps the resource schema data.
type orderResourceModel struct {
	TodoList []string `tfsdk:"todo_list"`
}

// Metadata returns the resource type name.
func (r *addTodoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_add_todo_items"
}

// Schema defines the schema for the resource.
func (r *addTodoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Adds items to the todo list.",
		Attributes: map[string]schema.Attribute{
			"todo_list": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

// Create a new resource.
func (r *addTodoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan orderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	rb, err := json.Marshal(plan.TodoList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to marshal the items",
			err.Error(),
		)
		return
	}

	addRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/create", r.baseurl), bytes.NewBuffer(rb))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create new request for add endpoint",
			err.Error(),
		)
		return
	}
	addRequest.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(addRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to hit /create endpoint",
			err.Error(),
		)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Received non-OK response from /create endpoint",
			fmt.Sprintf("Status code: %d", res.StatusCode),
		)
		return
	}

	// new orderResourceModel state is used to store the response data
	state := orderResourceModel{}
	err = json.NewDecoder(res.Body).Decode(&state.TodoList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the response body",
			err.Error(),
		)
		return
	}
	// Set state to the values sent in the request
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *addTodoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state orderResourceModel
	// diags := req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	res, err := http.Get(fmt.Sprintf("%s/get", r.baseurl))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to hit the Todo List get endpoint",
			err.Error(),
		)
		return
	}
	defer res.Body.Close()

	var responseItems []string
	err = json.NewDecoder(res.Body).Decode(&responseItems)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the response body",
			err.Error(),
		)
		return
	}

	// Set state
	state.TodoList = responseItems
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource information.
func (r *addTodoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan orderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	rb, err := json.Marshal(plan.TodoList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to marshal the items",
			err.Error(),
		)
		return
	}

	updateRequest, err := http.NewRequest("PUT", fmt.Sprintf("%s/update", r.baseurl), bytes.NewBuffer(rb))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create new request for add endpoint to update the todo list",
			err.Error(),
		)
		return
	}
	updateRequest.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to hit /update endpoint to update todo list",
			err.Error(),
		)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Received non-OK response from /update endpoint",
			fmt.Sprintf("Status code: %d", res.StatusCode),
		)
		return
	}

	state := orderResourceModel{}
	err = json.NewDecoder(res.Body).Decode(&state.TodoList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the response body",
			fmt.Sprint(state.TodoList, "\n", err.Error()),
		)
		return
	}
	// Set state to the values sent in the request
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource information.
func (r *addTodoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// making delete request
	deleteRequest, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", r.baseurl), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create new request for delete endpoint to delete the todo list",
			err.Error(),
		)
		return
	}
	deleteRequest.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(deleteRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to hit /delete endpoint to update todo list",
			err.Error(),
		)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Received non-OK response from /delete endpoint",
			fmt.Sprintf("Status code: %d", res.StatusCode),
		)
		return
	}

	state := orderResourceModel{}
	err = json.NewDecoder(res.Body).Decode(&state.TodoList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the response body",
			fmt.Sprint(state.TodoList, "\n", err.Error()),
		)
		return
	}
	// Set state to the values sent in the request
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *addTodoResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.baseurl = baseurl
}

func (r *addTodoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// package provider

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"

// 	"github.com/hashicorp/terraform-plugin-framework/path"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// )

// // Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.Resource                = &addTodoResource{}
// 	_ resource.ResourceWithConfigure   = &addTodoResource{}
// 	_ resource.ResourceWithImportState = &addTodoResource{}
// )

// // NewAddTodoResource is a helper function to simplify the provider implementation.
// func NewAddTodoResource() resource.Resource {
// 	return &addTodoResource{}
// }

// // addTodoResource is the resource implementation.
// type addTodoResource struct {
// 	baseurl string
// }

// // orderResourceModel maps the resource schema data.
// type orderResourceModel struct {
// 	TodoList []string `tfsdk:"todo_list"`
// }

// // Metadata returns the resource type name.
// func (r *addTodoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_add_todo_items"
// }

// // Schema defines the schema for the resource.
// func (r *addTodoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		Description: "Fetches the todo list.",
// 		Attributes: map[string]schema.Attribute{
// 			"todo_list": schema.ListAttribute{
// 				ElementType: types.StringType,
// 				Required:    true,
// 			},
// 		},
// 	}
// }

// // Create a new resource.
// func (r *addTodoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	// Retrieve values from plan
// 	var items orderResourceModel
// 	diags := req.Plan.Get(ctx, &items)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Create new order
// 	rb, err := json.Marshal(items.TodoList)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Marshal the items",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	addRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/create", r.baseurl), bytes.NewBuffer([]byte(rb)))

// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to create new request of create endpoint",
// 			err.Error(),
// 		)
// 		return
// 	}
// 	addRequest.Header.Add("Content-Type", "application/json")
// 	client := &http.Client{}
// 	res, err := client.Do(addRequest)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to hit /create endpoint",
// 			err.Error(),
// 		)
// 		return
// 	}
// 	// Map response body to model
// 	defer res.Body.Close()

// 	err = json.NewDecoder(res.Body).Decode(&items.TodoList)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to read the response body",
// 			err.Error(),
// 		)
// 		return
// 	}
// 	// Set state to fully populated data
// 	diags = resp.State.Set(ctx, items)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// // Read resource information.
// func (r *addTodoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var state ToDoDataSourceModel
// 	diags := req.State.Get(ctx, &state)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	res, err := http.Get(fmt.Sprintf("%s/get", r.baseurl))
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Hit the Todo List get endpoint",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	// Map response body to model
// 	defer res.Body.Close()
// 	bodyBytes, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to read the response body",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	err = json.Unmarshal(bodyBytes, &state.TodoList)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Unable to Read/Unmarshal Todo List",
// 			err.Error(),
// 		)
// 		return
// 	}

// 	// Set state
// 	diags = resp.State.Set(ctx, &state)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// func (r *addTodoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	// // Retrieve values from plan
// 	// var plan orderResourceModel
// 	// diags := req.Plan.Get(ctx, &plan)
// 	// resp.Diagnostics.Append(diags...)
// 	// if resp.Diagnostics.HasError() {
// 	// 	return
// 	// }

// 	// // Generate API request body from plan
// 	// var hashicupsItems []hashicups.OrderItem
// 	// for _, item := range plan.Items {
// 	// 	hashicupsItems = append(hashicupsItems, hashicups.OrderItem{
// 	// 		Coffee: hashicups.Coffee{
// 	// 			ID: int(item.Coffee.ID.ValueInt64()),
// 	// 		},
// 	// 		Quantity: int(item.Quantity.ValueInt64()),
// 	// 	})
// 	// }

// 	// // Update existing order
// 	// _, err := r.client.UpdateOrder(plan.ID.ValueString(), hashicupsItems)
// 	// if err != nil {
// 	// 	resp.Diagnostics.AddError(
// 	// 		"Error Updating HashiCups Order",
// 	// 		"Could not update order, unexpected error: "+err.Error(),
// 	// 	)
// 	// 	return
// 	// }

// 	// // Fetch updated items from GetOrder as UpdateOrder items are not
// 	// // populated.
// 	// order, err := r.client.GetOrder(plan.ID.ValueString())
// 	// if err != nil {
// 	// 	resp.Diagnostics.AddError(
// 	// 		"Error Reading HashiCups Order",
// 	// 		"Could not read HashiCups order ID "+plan.ID.ValueString()+": "+err.Error(),
// 	// 	)
// 	// 	return
// 	// }

// 	// diags = resp.State.Set(ctx, plan)
// 	// resp.Diagnostics.Append(diags...)
// 	// if resp.Diagnostics.HasError() {
// 	// 	return
// 	// }
// }

// func (r *addTodoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	// // Retrieve values from state
// 	// var state orderResourceModel
// 	// diags := req.State.Get(ctx, &state)
// 	// resp.Diagnostics.Append(diags...)
// 	// if resp.Diagnostics.HasError() {
// 	// 	return
// 	// }

// 	// // Delete existing order
// 	// err := r.client.DeleteOrder(state.ID.ValueString())
// 	// if err != nil {
// 	// 	resp.Diagnostics.AddError(
// 	// 		"Error Deleting HashiCups Order",
// 	// 		"Could not delete order, unexpected error: "+err.Error(),
// 	// 	)
// 	// 	return
// 	// }
// }

// // Configure adds the provider configured client to the resource.
// func (r *addTodoResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	baseurl, ok := req.ProviderData.(string)

// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Data Source Configure Type",
// 			fmt.Sprintf("Expected string as baseurl, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)

// 		return
// 	}

// 	r.baseurl = baseurl
// }

// func (r *addTodoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
