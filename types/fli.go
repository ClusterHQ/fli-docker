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

// The idea is that we could use the manifest
// to create a volume from a snapshot, branch, or volumeset.
// Having only a VolumeSet in the manifest can indicate
// creating from a VolumeSet or branch or snapshot respectively.
type Volume struct {
	Name string      `yaml:"name"`
	Snapshot string  `yaml:"snapshot"`
	VolumeSet string `yaml:"volumeset"`
	Branch string    `yaml:"branch"`
}

// Represents {compose_volume_name : "/chq/<vol_path>"}
// for volume names in compose to their new path
// after fli creates them.
type NewVolume struct {
	Name string
	VolumePath string
	VolumeName string
}