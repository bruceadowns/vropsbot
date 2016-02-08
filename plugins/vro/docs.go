package plugins

// String houses the inner string json
type String struct {
	Value string `json:"value"`
}

// Value houses the inner value json
type Value struct {
	String String `json:"string"`
}

// Parameter houses the inner parameter json
type Parameter struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value Value  `json:"value"`
}

// WorkflowExecutionJSON houses the json to POST a workflow
type WorkflowExecutionJSON struct {
	Parameters []Parameter `json:"parameters"`
}

// WorkflowsJSON houses the json returned from vRO's workflows query
type WorkflowsJSON struct {
	Link []struct {
		Href string
	}
}
