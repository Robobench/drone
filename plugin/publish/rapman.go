package publish

import (
	"fmt"
	"path/filepath"
	"strings"
  "os"
	"ioutil"
	"github.com/drone/drone/plugin/condition"
	"github.com/drone/drone/shared/build/buildfile"
  "github.com/drone/drone/shared/build/dockerfile"
)

type NamedTest struct {
	// Commands for running test
	Command string `yaml:"command"`
	// Name of test
	Name string`yaml:"name"`
	// Description of test
	Description string `yaml:"description"`
	// Flag to set whether test is enabled
	Enabled bool `yaml:"enabled"`
	// Condition under which to run tests
}

type TestSet struct {
	// A set of tests to run
	Tests[]  NamedTest `yaml:"tests"`
}

type RapmanRepo struct
{
	GitRepo string `yaml:"github_repo"`
	Username string `yaml:"username"`
	Password string `yaml:password"`
}

type Rapman struct{
	 Repo RapmanRepo `yaml:"target_repo,omitempty"`
	 Tests[] NamedTest `yaml:"tests"`
	 ImageName string `yaml:"docker_image,omitempty"`
   BaseImageName string `yaml:"base_image"`
	Script[]  string `yaml:build_script"`
	Env[] string `yaml:build_environment"`
  Condition *condition.Condition `yaml:"when,omitempty"`
}


func getDefaultPermissionsString(test *NamedTest) string
{
	var permissionsString = fmt.Sprintf(`{ \n
  "description"                : "%s"\n
  ,"maintainer"                : "robobench"\n
  // Path to executable within the docker image.\n
  ,"executable"                : "%s"\n
  // A list of directories the program should have Read/Write access to.\n
  // Paths are relative to your home. Ex: "Downloads" will access "$HOME/Downloads".\n
  ,"user-dirs"                 : []  // Default: []\n
  // Allowed the program to display x11 windows. This will allow drone to build\n
  // XOrg reliant devices until we iron out the xdummy issues\n
  ,"x11"                       : true        // Default: false\n
  // Allow the program access to your sound playing and recording.\n
  ,"access-working-directory" : false        // Default: false\n
  // Allow the program access to the internet.\n
  ,"allow-network-access"      : true        // Default: false\n
  // Allow privileged access to the devices in /dev. This is never safe,\n
  // but it is by far the simplest way to generically handle \n
  // an arbitrary graphics card\n
  ,"privileged"                : false       // Default: false\n
  // This image must run as root\n
  ,"as-root"                   : false      // Default: false\n
  // This image maps port 0.0.0.0:8000 to port 80\n
  ,"ports"                     : ["8000:80"] // Default: []\n
  ,"use-host-descriptor" : true\n
}\n`, test.Description, test.Command)
}


func getBaseImagefile(rapman *Rapman) dockerfile.Dockerfile
{
	d = dockerfile.New(rapman.BaseImage)

  for _, env := range rapman.Env {
	  	d.WriteEnv(env)
	  }

  for _, cmd := range rapman.Script {
  	d.WriteCmd(cmd)
  }
	return d
}



func createTest(f *buildfile.Buildfile,repoPath string, test *NamedTest, rapman *Rapman)
{
	baseTestPath := filepath.Join(repoPath, test.Name)
	os.MkdirAll(baseTestPath)

	// Write permissions file
	permissions := getDefaultPermissionsString(test)
	permissionsTempPath := filepath.Join("/tmp",test.Name + "permissions.json" )
	permissionsPath := filepath.Join(baseTestPath,"permissions.json" )
	ioutil.WriteFile(permissionsPath, permissions.Bytes(), 700)

	// Write baseimage file
  imageTestPath := filepath.Join(repoPath, "docker-file")
	os.MkdirAll(imageTestPath)
	baseImageTestFilename := filepath.Join(imageTestPath, "SubuserImagefile.base")
	baseDockerfile := getBaseImagefile(rapman)
  ioutil.WriteFile(baseImageTestFilename, baseDockerfile.Bytes(), 700)

  // Write the current imagefile
	dockerfile = dockerfile.New(rapman.ImageName)
  imageTestFilename := filepath.Join(imageTestPath, "SubuserImagefile")
	ioutil.WriteFile(imageTestFilename, dockerfile.Bytes(), 700)
}

func (test *NamedTest) Write(f *buildfile.Buildfile) {
	if test.Enabled {
	  f.WriteCmd(`echo "` + test.Name +`"`)
	  f.WriteCmd(test.Command)
  }
}

func (r *Rapman) Write(f *buildfile.Buildfile) {

	for _, test := range rapman.Tests {
		test.Write(f)
	}

	repopath = `/tmp/rapman_registry_repo`
	if r.Repo != nil {
	   if len(r.Repo.Username) == 0 || len(r.Repo.Password) == 0 || len(r.Repo.GitRepo) == 0{
		   f.WriteCmdSilent(`echo "Rapman: Missing argument(s)"`)
		   return
	   }
		f.WriteCmd(`git clone ` + r.Repo.GitRepo + ` ` + repopath)
		defer f.WriteCmd(`rm -r ` + repopath)
  }

}

func (r *Rapman) GetCondition() *condition.Condition {
	return r.Condition
}
