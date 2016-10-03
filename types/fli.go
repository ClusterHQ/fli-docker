package types

type Manifest struct {
	DockerApp string    `yaml:"docker_app"`
	Hub FlockerHub      `yaml:"flocker_hub"`
	Volumes []Volume    `yaml:"volumes"`
}

type FlockerHub struct { 
	Endpoint string   `yaml:"endpoint"`
	AuthToken string  `yaml:"tokenfile"`
}

type Volume struct {
	Name string      `yaml:"name"`
	Snapshot string  `yaml:"snapshot"`
	VolumeSet string `yaml:"volumeset"`
}

// Represents {compose_volume_name : "/chq/<vol_path>"}
// for volume names in compose to their new path
// after fli creates them.
type NewVolume struct {
	Name string
	VolumePath string
}