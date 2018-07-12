// +build mage

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"
	"github.com/magefile/mage/sh"
	// mg contains helpful utility functions, like Deps
)

const defaultDockerAPIVersion = "v1.37"

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
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	s := ParseEnvVars()
	// InstallSources(s)
	fmt.Println("Building Docker Image...")
	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		panic(err)
	}

	// dockerBuildContext, err := os.Open(".")
	// if err != nil {
	// 	panic(err)
	// }

	os.Setenv("PREBUILT", s.Prebuilt)

	opt := types.ImageBuildOptions{
		BuildArgs: map[string]*string{
			"NODE_VERSION": &s.Nodeversion,
			"NPM_VERSION":  &s.Npmversion,
			"V8_VERSION":   &s.V8version,
			"PREBUILT":     &s.Prebuilt,
			"FROM_DATA":    &s.Fromdata,
		},
		Tags: []string{s.Imagename + ":" + s.Imagetag},
		// Context:    dockerBuildContext,
		// Dockerfile: "image/" + s.Os + "/Dockerfile",
	}

	_, err = cli.ImageBuild(ctx, nil, opt)
	if err != nil {
		panic(err)
	}

	// app := "docker"
	// args := []string{
	// 	"build", "-f", dir + "/image/" + s.Os + "/Dockerfile",
	// 	"--build-arg", "NODE_VERSION=" + s.Nodeversion,
	// 	"--build-arg", "NPM_VERSION=" + s.Npmversion,
	// 	"--build-arg", "V8_VERSION=" + s.V8version,
	// 	"--build-arg", "PREBUILT=" + s.Prebuilt,
	// 	"--build-arg", "FROM_DATA=" + s.Fromdata,
	// 	"-t", s.Imagename + ":" + s.Imagetag, "."}
	// _, err = sh.Exec(envs, os.Stdout, os.Stdout, app, args...)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return err
}

func ParseEnvVars() Specification {
	var s Specification
	err := envconfig.Process("nd", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	format := "Debug:\nNodeversion: %s\nOs: %d\nDockerfile: %f\nImagetag: %s\nLatest: %v\nMajor Tag: %s\nMinor Tag: %s\nImage name: %s\nNPM Version: %s\nFrom: %s\nPrebuilt: %v\n"
	_, err = fmt.Printf(
		format,
		s.Os,
		s.Nodeversion,
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

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning /src dir...")
	sh.Rm("src/.")
	// fmt.Println("Cleanup image")
	// docker rmi `docker images $(TARGET) -q`
}
