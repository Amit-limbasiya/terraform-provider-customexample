package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		Description: "Fetches the todo list.",
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
	var items orderResourceModel
	diags := req.Plan.Get(ctx, &items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// // Generate API request body from plan --- already holding the items to add
	// var items []string
	// for _, item := range plan.Items {
	// 	items = append(items, hashicups.OrderItem{
	// 		Coffee: hashicups.Coffee{
	// 			ID: int(item.Coffee.ID.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }

	// Create new order
	rb, err := json.Marshal(items.TodoList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Marshal the items",
			err.Error(),
		)
		return
	}

	addRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/add", r.baseurl), bytes.NewBuffer([]byte(rb)))

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create new request of add endpoint",
			err.Error(),
		)
		return
	}
	addRequest.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(addRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to hit /add endpoint",
			err.Error(),
		)
		return
	}
	// Map response body to model
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&items.TodoList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the response body",
			err.Error(),
		)
		return
	}
	// Set state to fully populated data
	diags = resp.State.Set(ctx, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *addTodoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ToDoDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := http.Get(fmt.Sprintf("%s/get", r.baseurl))
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
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *addTodoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// // Retrieve values from plan
	// var plan orderResourceModel
	// diags := req.Plan.Get(ctx, &plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Generate API request body from plan
	// var hashicupsItems []hashicups.OrderItem
	// for _, item := range plan.Items {
	// 	hashicupsItems = append(hashicupsItems, hashicups.OrderItem{
	// 		Coffee: hashicups.Coffee{
	// 			ID: int(item.Coffee.ID.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }

	// // Update existing order
	// _, err := r.client.UpdateOrder(plan.ID.ValueString(), hashicupsItems)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Updating HashiCups Order",
	// 		"Could not update order, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// // Fetch updated items from GetOrder as UpdateOrder items are not
	// // populated.
	// order, err := r.client.GetOrder(plan.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading HashiCups Order",
	// 		"Could not read HashiCups order ID "+plan.ID.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }

	// // Update resource state with updated items and timestamp
	// plan.Items = []orderItemModel{}
	// for _, item := range order.Items {
	// 	plan.Items = append(plan.Items, orderItemModel{
	// 		Coffee: orderItemCoffeeModel{
	// 			ID:          types.Int64Value(int64(item.Coffee.ID)),
	// 			Name:        types.StringValue(item.Coffee.Name),
	// 			Teaser:      types.StringValue(item.Coffee.Teaser),
	// 			Description: types.StringValue(item.Coffee.Description),
	// 			Price:       types.Float64Value(item.Coffee.Price),
	// 			Image:       types.StringValue(item.Coffee.Image),
	// 		},
	// 		Quantity: types.Int64Value(int64(item.Quantity)),
	// 	})
	// }
	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// diags = resp.State.Set(ctx, plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

func (r *addTodoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// // Retrieve values from state
	// var state orderResourceModel
	// diags := req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Delete existing order
	// err := r.client.DeleteOrder(state.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Deleting HashiCups Order",
	// 		"Could not delete order, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
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
