package main

type PlanetSchematic struct {
	CycleTime     int `yaml:"cycleTime"`
	SchematicID   int `yaml:"schematicID"`
	SchematicName string `yaml:"schematicName"`
}

type PlanetSchematicTypeMap struct {
	IsInput     bool `yaml:"isInput"`
	Quantity    int `yaml:"quantity"`
	SchematicID int `yaml:"schematicID"`
	TypeID      int `yaml:"typeID"`
}
