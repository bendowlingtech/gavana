package graft

var RegisteredModels []interface{}

func RegisterModel(model interface{}) {
	RegisteredModels = append(RegisteredModels, model)
}
