terraform{
	required_providers{
		customexample={
			source = "terraform.registry.io/edu/custom-example"
			version = "0.1.0"
		}
	}
}

provider "customexample"{
	username  =  "amit"
	password  =  "abc"
	baseurl   =  "http://localhost:5019"
}

resource customexample_add_todo_items "addingtodos"{
	todo_list=["A","B","C"]
}

data customexample_todo "todo"{
	depends_on = [ customexample_add_todo_items.addingtodos ]
}

output "todolist"{
	value = data.customexample_todo.todo
}