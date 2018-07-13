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
	Latest      bool   `default:"false"`
	Majortag    string `required:"true"`
	Minortag    string `required:"true"`
	Imagename   string `required:"true"`
	Npmversion  string `required:"true"`
	Fromdata    string `required:"true"`
	Prebuilt    string `default:"false"`
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
	format := "Debug:\nNodeversion: %s\nOs: %d\nV8: %s\nDockerfile: %f\nImagetag: %s\nLatest: %v\nMajor Tag: %s\nMinor Tag: %s\nImage name: %s\nNPM Version: %s\nFrom: %s\nPrebuilt: %v\n"
	_, err = fmt.Printf(
		format,
		s.Nodeversion,
		s.Os,
		s.V8version,
		s.Dockerfile,
		s.Imagetag,
		s.Latest,
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
	log.Println("checking if " + path + " exists")
	if _, err := os.Stat(path); os.IsNotExist(err) {

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
	} else {
		log.Println("already exists, no need to download again.")
	}
	return err
}

// publish to docker hub
func Publish() {
	s := ParseEnvVars()
	var envs = map[string]string{}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "docker", "login", "--username", "ops@nearform.com", "-p", s.Dockerpass)
	if err != nil {
		log.Fatal(err)
	}
	// 	@echo $(DOCKER_PASS) | docker login --username $(DOCKER_USER) --password-stdin
	// 	docker push $(TARGET)

	// ifdef LATEST
	// 	docker tag $(TARGET) $(IMAGE_NAME):latest
	// 	docker push $(IMAGE_NAME):latest

	// ifdef MAJOR_TAG
	// 	docker tag $(TARGET) $(IMAGE_NAME):$(MAJOR_TAG)
	// 	docker push $(IMAGE_NAME):$(MAJOR_TAG)

	// 	ifdef MINOR_TAG
	// 	docker tag $(TARGET) $(IMAGE_NAME):$(MINOR_TAG)
	// 	docker push $(IMAGE_NAME):$(MINOR_TAG)

	// 	ifdef LTS_TAG
	// 	docker tag $(TARGET) $(IMAGE_NAME):$(LTS_TAG)
	// 	docker push $(IMAGE_NAME):$(LTS_TAG)
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning /src dir...")
	sh.Rm("src/.")
	s := ParseEnvVars()
	fmt.Println("Cleanup image " + s.Imagename + ":" + s.Imagetag)
	var envs = map[string]string{}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "docker", "rmi", s.Imagename+":"+s.Imagetag)
	if err != nil {
		log.Fatal(err)
	}
}
