// +build mage

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
		Tags:           getTags(s),
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

func getTags(s Specification) []string {
	tags := []string{s.Imagename + ":" + s.Imagetag}
	if s.Latest {
		tags = append(tags, s.Imagename+":latest")
	}
	if s.Majortag != "" {
		tags = append(tags, s.Imagename+":"+s.Majortag)
	}
	if s.Minortag != "" {
		tags = append(tags, s.Imagename+":"+s.Minortag)
	}
	if s.Lts != "" {
		tags = append(tags, s.Imagename+":"+s.Lts)
	}
	return tags
}

// publish to docker hub
func Publish() error {
	s := ParseEnvVars()
	cli, _ := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))

	authConfig := types.AuthConfig{
		Username: s.Dockeruser,
		Password: s.Dockerpass,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	options := types.ImagePushOptions{
		RegistryAuth: authStr,
		All:          true,
	}
	pushResponse, err := cli.ImagePush(context.Background(), s.Imagename+":"+s.Imagetag, options)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	defer pushResponse.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(pushResponse, os.Stderr, termFd, isTerm, nil)
}

// Clean up sources and the images we created
func Clean() error {
	s := ParseEnvVars()
	cli, _ := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	options := types.ImageRemoveOptions{}
	tags := getTags(s)
	var err error
	for i := 0; i < len(tags); i++ {
		removeResponse, err := cli.ImageRemove(context.Background(), tags[i], options)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		fmt.Println(removeResponse)
	}
	return err
}
