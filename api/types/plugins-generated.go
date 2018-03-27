// Code generated by swagger-gen. DO NOT EDIT

package types

// A plugin for the Engine API
type Plugin struct {
	// The config of a plugin.
	Config PluginConfig `json:"Config"`
	// True if the plugin is running. False if the plugin is not running, only installed.
	Enabled bool   `json:"Enabled"`
	ID      string `json:"Id,omitempty"`
	Name    string `json:"Name"`
	// plugin remote reference used to push/pull the plugin
	PluginReference string `json:"PluginReference,omitempty"`
	// Settings that can be modified by users.
	Settings PluginSettings `json:"Settings"`
}

// The config of a plugin.
type PluginConfig struct {
	Args        PluginConfigArgs `json:"Args"`
	Description string           `json:"Description"`
	// Docker Version used to create the plugin
	DockerVersion string      `json:"DockerVersion,omitempty"`
	Documentation string      `json:"Documentation"`
	Entrypoint    []string    `json:"Entrypoint"`
	Env           []PluginEnv `json:"Env"`
	// The interface between Docker and the plugin
	Interface       PluginConfigInterface `json:"Interface"`
	IpcHost         bool                  `json:"IpcHost"`
	Linux           PluginConfigLinux     `json:"Linux"`
	Mounts          []PluginMount         `json:"Mounts"`
	Network         PluginConfigNetwork   `json:"Network"`
	PidHost         bool                  `json:"PidHost"`
	PropagatedMount string                `json:"PropagatedMount"`
	User            PluginConfigUser      `json:"User,omitempty"`
	WorkDir         string                `json:"WorkDir"`
	Rootfs          *PluginConfigRootfs   `json:"rootfs,omitempty"`
}
type PluginConfigArgs struct {
	Description string   `json:"Description"`
	Name        string   `json:"Name"`
	Settable    []string `json:"Settable"`
	Value       []string `json:"Value"`
}

// The interface between Docker and the plugin
type PluginConfigInterface struct {
	Socket string                `json:"Socket"`
	Types  []PluginInterfaceType `json:"Types"`
}
type PluginConfigLinux struct {
	AllowAllDevices bool           `json:"AllowAllDevices"`
	Capabilities    []string       `json:"Capabilities"`
	Devices         []PluginDevice `json:"Devices"`
}
type PluginConfigNetwork struct {
	Type string `json:"Type"`
}
type PluginConfigUser struct {
	GID uint32 `json:"GID,omitempty"`
	UID uint32 `json:"UID,omitempty"`
}
type PluginConfigRootfs struct {
	DiffIds []string `json:"diff_ids,omitempty"`
	Type    string   `json:"type,omitempty"`
}

// Settings that can be modified by users.
type PluginSettings struct {
	Args    []string       `json:"Args"`
	Devices []PluginDevice `json:"Devices"`
	Env     []string       `json:"Env"`
	Mounts  []PluginMount  `json:"Mounts"`
}
type PluginDevice struct {
	Description string   `json:"Description"`
	Name        string   `json:"Name"`
	Path        *string  `json:"Path"`
	Settable    []string `json:"Settable"`
}
type PluginEnv struct {
	Description string   `json:"Description"`
	Name        string   `json:"Name"`
	Settable    []string `json:"Settable"`
	Value       *string  `json:"Value"`
}
type PluginInterfaceType struct {
	Capability string `json:"Capability"`
	Prefix     string `json:"Prefix"`
	Version    string `json:"Version"`
}
type PluginMount struct {
	Description string   `json:"Description"`
	Destination string   `json:"Destination"`
	Name        string   `json:"Name"`
	Options     []string `json:"Options"`
	Settable    []string `json:"Settable"`
	Source      *string  `json:"Source"`
	Type        string   `json:"Type"`
}
