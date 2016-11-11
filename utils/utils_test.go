/*
 *Copyright ClusterHQ Inc.  See LICENSE file for details.
 *
 */

package utils

import (
	"testing"
	"os"
	"io/ioutil"
	"reflect"
	"regexp"

	"gopkg.in/yaml.v2"

	"github.com/ClusterHQ/fli-docker/logger"
	"github.com/ClusterHQ/fli-docker/types"
)

var fliManifestData = `docker_app: some-compose.yml

flocker_hub:
    endpoint: https://flockerhub.com

volumes:
    - name: some-name
      snapshot: example-snapshot
      volumeset: example-vs
`

var fliBadManifestDataCompose = `docker_app: 

flocker_hub:
    endpoint: https://flockerhub.com

volumes:
    - name: some-name
      snapshot: example-snapshot
      volumeset: example-vs
`

var fliBadManifestDataVSet = `docker_app: some-compose.yml

flocker_hub:
    endpoint: https://flockerhub.com

volumes:
    - name: some-name
      snapshot: example-snapshot
`

var fliBadManifestDataSnapBranch = `docker_app: some-compose.yml

flocker_hub:
    endpoint: https://flockerhub.com

volumes:
    - name: some-name
      volumeset: example-vs
`

var fliDockerVols = `fli-somevolume, some-vs`

var fliDockerVolsBytes = []byte(fliDockerVols)
var fliManifestDataBytes = []byte(fliManifestData)
var fliBadManifestDataBytesCompose = []byte(fliBadManifestDataCompose)
var fliBadManifestDataBytesVSet = []byte(fliBadManifestDataVSet)
var fliBadManifestDataBytesSnapBranch = []byte(fliBadManifestDataSnapBranch)

var composeManifestData = `version: '2'
services:
  web:
    image: mysql/mysql-server
    environment:
      - "MYSQL_ROOT_PASSWORD=my-secret-pw"
    ports:
      - "3306"
    volumes:
      - mysql-data:/var/lib/mysql
`
var composeManifestDataBytes = []byte(composeManifestData)

func init(){
	// Used by function in go
	logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

func TestCheckForPath(t *testing.T) {
	// Most Linux distros will have this binary
	// CheckForPath needs binary
	found, err := CheckForPath("/usr/bin/less")
    if ! found {
    	t.Error(err)
    }
}

func TestCheckForFile(t *testing.T) {
	wrerr := ioutil.WriteFile("/tmp/flitest-file.yml", composeManifestDataBytes, 0644)
    if wrerr != nil {
    	t.Error(wrerr)
	}

	defer os.Remove("/tmp/flitest-file.yml")

    found, err := CheckForFile("/tmp/flitest-file.yml")
    if ! found {
    	t.Error(err)
    }
}

func TestCheckForCmd(t *testing.T) {
	found, err := CheckForCmd("echo 'hello'")
    if ! found {
    	t.Error(err)
    }

}

func TestMakeCopy(t *testing.T) {
	wrerr := ioutil.WriteFile("/tmp/flitest-file.yml", composeManifestDataBytes, 0644)
    if wrerr != nil {
    	t.Error(wrerr)
	}

	defer os.Remove("/tmp/flitest-file.yml")
	defer os.Remove("/tmp/flitest-file.yml-fli.copy")

    MakeCopy("/tmp/flitest-file.yml")
}

func TestCheckForCopy(t *testing.T) {
		wrerr := ioutil.WriteFile("/tmp/flitest-file.yml", composeManifestDataBytes, 0644)
    if wrerr != nil {
    	t.Error(wrerr)
	}

	defer os.Remove("/tmp/flitest-file.yml")
	defer os.Remove("/tmp/flitest-file.yml-fli.copy")

    MakeCopy("/tmp/flitest-file.yml")
    CheckForCopy("/tmp/flitest-file.yml")
}

func TestParseManifest(t *testing.T) {
    m := ParseManifest(fliManifestDataBytes)

    mCompare := types.Manifest{}
	err := yaml.Unmarshal([]byte(fliManifestData), &mCompare)
    if err != nil {
        t.Error("error: %v", err)
    }

    t.Log("Got ", m, "Expected ", &mCompare)
    if ! reflect.DeepEqual(m, &mCompare){
    	t.Error("Got ", m, "Expected ", &mCompare)
    }
}

func TestBadManifestCompose(t *testing.T) {
    var manifest types.Manifest
	err := yaml.Unmarshal(fliBadManifestDataBytesCompose, &manifest)
	if err != nil {
		logger.Error.Fatal(err)
	}
	// Validate manifest.
	valErr := verifyManifest(manifest)
	t.Log(valErr)
	if valErr == nil {
		t.Error("Got nil, expecting, Missing Docker Compose file error")
	}
}

func TestBadManifestVSet(t *testing.T) {
    var manifest types.Manifest
	err := yaml.Unmarshal(fliBadManifestDataBytesVSet, &manifest)
	if err != nil {
		logger.Error.Fatal(err)
	}
	// Validate manifest.
	valErr := verifyManifest(manifest)
	t.Log(valErr)
	if valErr == nil {
		t.Error("Got nil, expecting, Missing volumeset: for volume")
	}
}

func TestBadManifestSnapBranch(t *testing.T) {
    var manifest types.Manifest
	err := yaml.Unmarshal(fliBadManifestDataBytesSnapBranch, &manifest)
	if err != nil {
		logger.Error.Fatal(err)
	}
	// Validate manifest.
	valErr := verifyManifest(manifest)
	t.Log(valErr)
	if valErr == nil {
		t.Error("Got nil, expecting, Need snapshot: or branch: for volume")
	}
}

func TestMapVolumeToCompose(t *testing.T) {
	wrerr := ioutil.WriteFile("/tmp/flitest-mapcompose.yml", composeManifestDataBytes, 0644)
    if wrerr != nil {
    	t.Error(wrerr)
	}

	defer os.Remove("/tmp/flitest-mapcompose.yml")

	MapVolumeToCompose("mysql-data", "/chq/test/vol", "/tmp/flitest-mapcompose.yml")

}

func TestParseCompose(t *testing.T) {
	wrerr := ioutil.WriteFile("/tmp/flitest-compose.yml", composeManifestDataBytes, 0644)
    if wrerr != nil {
    	t.Error(wrerr)
	}

	defer os.Remove("/tmp/flitest-compose.yml")

    ParseCompose("/tmp/flitest-compose.yml")
}

func isValidUUID(uuid string) bool {
    r := regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}")
    return r.MatchString(uuid)
}

func TestGenUUID(t *testing.T) {
	exists, _ := CheckForCmd("uuidgen")
	if ! exists {
    	t.Skip("skipping test; uuidgen not available")
  	}
	uuid, err := GenUUID()
	if err != nil {
		t.Error(err)
	}
	if ! isValidUUID(uuid) {
		t.Error("Not a valid UUIDv4: ", uuid)
	}
}

func TestCleanEnv(t *testing.T) {
	w1rerr := ioutil.WriteFile("flitest-compose.yml", composeManifestDataBytes, 0644)
    if w1rerr != nil {
    	t.Error(w1rerr)
	}
	w2rerr := ioutil.WriteFile("flitest-compose.yml-fli.copy", composeManifestDataBytes, 0644)
    if w2rerr != nil {
    	t.Error(w2rerr)
	}
	w3rerr := ioutil.WriteFile(".flidockervols", fliDockerVolsBytes, 0644)
    if w3rerr != nil {
    	t.Error(w3rerr)
	}

	defer os.Remove("flitest-compose.yml")

	CleanEnv("flitest-compose.yml")

	if _, err := os.Stat("flitest-compose.yml-fli.copy"); err == nil {
		t.Error("Expecting flitest-compose.yml-fli.copy to be deleted")
	}

	if _, err := os.Stat(".flidockervols"); err == nil {
		t.Error("Expecting .flidockervols to be deleted")
	}
}

func TestGetBasePath(t *testing.T) {
	path, err := GetBasePath("/usr/bin/less")
	if err != nil {
		t.Log(err)
    	t.Error("Basepath could not be found")
  	}
	if path != "/usr/bin"{
		t.Error("Expected /usr/bin, got", path)
	}
}

