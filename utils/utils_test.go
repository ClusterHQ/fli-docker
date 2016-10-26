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
    tokenfile: /root/some.token

volumes:
    - name: some-name
      snapshot: example-snapshot
      volumeset: example-vs
`
var fliManifestDataBytes = []byte(fliManifestData)

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

