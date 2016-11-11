/*
 *Copyright ClusterHQ Inc.  See LICENSE file for details.
 *
 */
 
package types

import (
	"testing"
	"reflect"

	"gopkg.in/yaml.v2"
)

var manifestData = `
docker_app: some-compose.yml

flocker_hub:
    endpoint: https://flockerhub.com

volumes:
    - name: some-name
      snapshot: example-snapshot
      volumeset: example-vs
`

func TestManifest(t *testing.T) {
	m := Manifest{}
	err := yaml.Unmarshal([]byte(manifestData), &m)
    if err != nil {
        t.Error("error: %v", err)
    }
    t.Log("Testing DockerApp, got", m.DockerApp)
 	if m.DockerApp != "some-compose.yml" {
    	t.Error("Expected some-compose.yml, got ", m.DockerApp)
  	}
}


func TestFlockerHub(t *testing.T) {
	m := Manifest{}
	err := yaml.Unmarshal([]byte(manifestData), &m)
    if err != nil {
        t.Error("error: %v", err)
    }
  	t.Log("Testing FlockerHub, Got", m.Hub)
  	if m.Hub != (FlockerHub{"https://flockerhub.com"}) {
  		t.Error("Expected FlockerHub{https://flockerhub.com}, got ", m.Hub)
  	}
}

func TestVolume(t *testing.T) {
	m := Manifest{}
	err := yaml.Unmarshal([]byte(manifestData), &m)
    if err != nil {
        t.Error("error: %v", err)
    }
  	t.Log("Testing Volumes, Got", m.Volumes)
  	vols := []Volume{}
  	vols = append(vols, Volume{"some-name", "example-snapshot", "example-vs", ""})
  	if ! reflect.DeepEqual(m.Volumes, vols){
  		t.Error("Expected [{some-name, example-snapshot, example-vs}], got ", m.Volumes)
  	}
 
}

func TestNewVolume(t *testing.T) {
	nv := NewVolume{"compose-volume-name", "/chq/some/path", "fli-volume-name", "volumeset"}
	t.Log("Testing NewVolume, Got:", nv)
	if nv.Name != "compose-volume-name" {
  		t.Error("Expected compose-volume-name, got ", nv.Name)
  	}
  	if nv.VolumePath != "/chq/some/path" {
  		t.Error("Expected /chq/some/path, got ", nv.VolumePath)
  	}
  	if nv.VolumeName != "fli-volume-name" {
  		t.Error("Expected fli-volume-name, got ", nv.VolumeName)
  	}
  	if nv.VolumeSet != "volumeset" {
  		t.Error("Expected volumeset, got ", nv.VolumeSet)
  	}
}
