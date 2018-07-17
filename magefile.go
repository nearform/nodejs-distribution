// +build mage

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/magefile/mage/sh"
	// mg contains helpful utility functions, like Deps
)

type Specification struct {
	Os          string `required:"true"`
	Nodeversion string `required:"true"`
	V8version   string `required:"true"`
	Dockerfile  string `required:"true"`
	Imagetag    string `required:"true"`
	Latest      bool   `default:false`
	Lts         string
	Majortag    string `required:"true"`
	Minortag    string `required:"true"`
	Imagename   string `required:"true"`
	Npmversion  string `required:"true"`
	Fromdata    string `required:"true"`
	Prebuilt    string `default:false`
	Dockeruser  string `required:"true"`
	Dockerpass  string `required:"true"`
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

func ParseEnvVars() Specification {
	var s Specification
	err := envconfig.Process("nd", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	format := "Debug:\nNodeversion: %s\nOs: %d\nV8: %s\nDockerfile: %f\nImagetag: %s\nLatest: %v\nLts: %v\nMajor Tag: %s\nMinor Tag: %s\nImage name: %s\nNPM Version: %s\nFrom: %s\nPrebuilt: %v\n"
	_, err = fmt.Printf(
		format,
		s.Nodeversion,
		s.Os,
		s.V8version,
		s.Dockerfile,
		s.Imagetag,
		s.Latest,
		s.Lts,
		s.Majortag,
		s.Minortag,
		s.Imagename,
		s.Npmversion,
		s.Fromdata,
		s.Prebuilt)
	if err != nil {
		log.Fatal(err.Error())
	}
	return s
}

// Get Node.js sources, build
func InstallSources(s Specification) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := dir + "/src/node-v" + s.Nodeversion + "-linux-x64.tar.gz"
	if s.Prebuilt != "" {
		log.Println("checking if " + path + " exists")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			runDownloadScript(s)
		} else {
			log.Println("already exists, no need to download again.")
		}
	} else {
		runDownloadScript(s)
	}
	return err
}

func runDownloadScript(s Specification) error {
	dir, err := os.Getwd()
	var envs = map[string]string{
		"PREBUILT": s.Prebuilt,
		"OS":       s.Os,
	}
	fmt.Println("Installing Node.js sources...")
	_, err = sh.Exec(envs, os.Stdout, os.Stdout, "./contrib/etc/get_node_source.sh", s.Nodeversion, dir+"/src/")
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Squash the image using a shell command
func Squash() {
	fmt.Println("Squashing the image...")
	s := ParseEnvVars()
	tags := getTags(s)
	var envs = map[string]string{}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "docker-squash", tags[0], "-t", tags[0])
	if err != nil {
		log.Fatal(err)
	}
}

// Run a test on the image
func Test() {
	fmt.Println("Squashing the image...")
	sh.Rm("src/.")
	s := ParseEnvVars()
	fmt.Println("Cleanup image " + s.Imagename + ":" + s.Imagetag)
	var envs = map[string]string{}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "docker", "rmi", s.Imagename+":"+s.Imagetag)
	if err != nil {
		log.Fatal(err)
	}
}
