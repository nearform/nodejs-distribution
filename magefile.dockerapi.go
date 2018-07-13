// +build mage

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/jhoonb/archivex"
)

// mg contains helpful utility functions, like Deps

const defaultDockerAPIVersion = "v1.37"

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Build the container using the native docker api
func Build() error {
	s := ParseEnvVars()
	InstallSources(s)
	tar := new(archivex.TarFile)
	tar.Create("/tmp/nodejs-distro.tar")
	tar.AddAll("contrib", true)
	tar.AddAll("src", true)
	tar.AddAll("test", true)
	tar.AddAll("s2i", true)
	tar.AddAll("help", true)
	tar.AddAll("image", true)
	tar.AddAll("licenses", true)
	tar.Close()
	dockerBuildContext, err := os.Open("/tmp/nodejs-distro.tar")
	defer dockerBuildContext.Close()
	cli, _ := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	args := map[string]*string{
		"PREBUILT":     &s.Prebuilt,
		"NODE_VERSION": &s.Nodeversion,
	}
	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{s.Imagename + ":" + s.Imagetag},
		Dockerfile:     "image/" + s.Os + "/Dockerfile",
		BuildArgs:      args,
	}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	defer buildResponse.Body.Close()
	fmt.Printf("********* %s **********\n", buildResponse.OSType)

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stderr, termFd, isTerm, nil)
}
