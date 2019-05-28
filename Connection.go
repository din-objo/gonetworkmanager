package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	ConnectionInterface = SettingsInterface + ".Connection"

	/* Methods */
	ConnectionUpdate        = ConnectionInterface + ".Update"
	ConnectionUpdateUnsaved = ConnectionInterface + ".UpdateUnsaved"
	ConnectionDelete        = ConnectionInterface + ".Delete"
	ConnectionGetSettings   = ConnectionInterface + ".GetSettings"
	ConnectionGetSecrets    = ConnectionInterface + ".GetSecrets"
	ConnectionClearSecrets  = ConnectionInterface + ".ClearSecrets"
	ConnectionSave          = ConnectionInterface + ".Save"
	ConnectionUpdate2       = ConnectionInterface + ".Update2"

	/* Properties */
	ConnectionPropertyUnsaved  = ConnectionInterface + ".Unsaved"  // readable   b
	ConnectionPropertyFlags    = ConnectionInterface + ".Flags"    // readable   u
	ConnectionPropertyFilename = ConnectionInterface + ".Filename" // readable   s
)

//type ConnectionSettings map[string]map[string]interface{}
type ConnectionSettings map[string]map[string]interface{}

type Connection interface {
	GetPath() dbus.ObjectPath

	// GetSettings gets the settings maps describing this network configuration.
	// This will never include any secrets required for connection to the
	// network, as those are often protected. Secrets must be requested
	// separately using the GetSecrets() callWithReturnAndPanic.
	GetSettings() ConnectionSettings

	MarshalJSON() ([]byte, error)
}

func NewConnection(objectPath dbus.ObjectPath) (Connection, error) {
	var c connection
	return &c, c.init(NetworkManagerInterface, objectPath)
}

type connection struct {
	dbusBase
}

func (c *connection) GetPath() dbus.ObjectPath {
	return c.obj.Path()
}

func (c *connection) GetSettings() ConnectionSettings {
	var settings map[string]map[string]dbus.Variant
	c.callWithReturnAndPanic(&settings, ConnectionGetSettings)

	rv := make(ConnectionSettings)

	for k1, v1 := range settings {
		rv[k1] = make(map[string]interface{})

		for k2, v2 := range v1 {
			rv[k1][k2] = v2.Value()
		}
	}

	return rv
}

func (c *connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.GetSettings())
}
