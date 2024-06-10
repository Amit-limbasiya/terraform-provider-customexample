terraform {
  required_providers {
    customexample = {
      source = "Amit-limbasiya/customexample"
      version = "1.0.4"
    }
  }
}

provider "customexample"{
	username  =  "amit"
	password  =  "abc"
	baseurl   =  "http://localhost:8080"
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