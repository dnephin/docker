package container

// ----------------------------------------------------------------------------
// DO NOT EDIT THIS FILE
// This file was generated by `swagger generate operation`
//
// See hack/swagger-gen.sh
// ----------------------------------------------------------------------------

// ContainerChangeResponseItem container change response item
// swagger:model ContainerChangeResponseItem
type ContainerChangeResponseItem struct {

	// Kind of change
	// Required: true
	Kind uint8 `json:"Kind"`

	// Path to file that has changed
	// Required: true
	Path string `json:"Path"`
}
