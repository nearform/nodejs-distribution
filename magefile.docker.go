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
	"github.com/mholt/archiver"
)

const defaultDockerAPIVersion = "v1.37"

// Build the container using the native docker api
func Build() error {
	v := config()
	InstallSources()
	err := archiver.Tar.Make("/tmp/nodejs-distro.tar", []string{
		"contrib",
		"src",
		"test",
		"s2i",
		"help",
		"image",
		"licenses",
	})
	dockerBuildContext, err := os.Open("/tmp/nodejs-distro.tar")
	defer dockerBuildContext.Close()
	cli, _ := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	pb := preBuiltEnv(v)
	nv := v.GetString("Nodeversion")
	args := map[string]*string{
		"PREBUILT":     &pb,
		"NODE_VERSION": &nv,
	}
	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           getTags(),
		Dockerfile:     dockerFile(v),
		BuildArgs:      args,
	}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	check(err)
	defer buildResponse.Body.Close()
	fmt.Printf("********* %s **********\n", buildResponse.OSType)

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stderr, termFd, isTerm, nil)
}

func getTags() []string {
	v := config()
	tags := []string{imageName(v) + ":" + v.GetString("Imagetag")}
	if v.GetBool("Latest") {
		tags = append(tags, imageName(v)+":latest")
	}
	if v.IsSet("Majortag") {
		tags = append(tags, imageName(v)+":"+v.GetString("Majortag"))
	}
	if v.IsSet("Minortag") {
		tags = append(tags, imageName(v)+":"+v.GetString("Minortag"))
	}
	if v.GetString("Lts") != "" {
		tags = append(tags, imageName(v)+":"+v.GetString("Lts"))
	}
	return tags
}

// publish to docker hub
func Publish() error {
	v := config()
	cli, _ := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))

	authConfig := types.AuthConfig{
		Username: v.GetString("Dockeruser"),
		Password: v.GetString("Dockerpass"),
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
	pushResponse, err := cli.ImagePush(context.Background(), imageName(v), options)
	check(err)
	defer pushResponse.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(pushResponse, os.Stderr, termFd, isTerm, nil)
}

// publish to Red Hat Catalog
func PublishRedHat() error {
	v := config()
	cli, _ := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))

	authConfig := types.AuthConfig{
		Username:      "unused",
		Password:      v.GetString("Rhsecret"),
		ServerAddress: v.GetString("Rhendpoint"),
	}
	// fmt.Println(authConfig)
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	_, err = cli.RegistryLogin(context.Background(), authConfig)
	check(err)

	options := types.ImagePushOptions{
		RegistryAuth: authStr,
		All:          true,
	}
	RedHatImageName := v.GetString("Rhendpoint") + "/" + v.GetString("Rhproject") + "/nearform-s2i-node" + ":" + v.GetString("Imagetag")
	err = cli.ImageTag(context.Background(), imageName(v), RedHatImageName)
	check(err)
	pushResponse, err := cli.ImagePush(context.Background(), RedHatImageName, options)
	check(err)
	defer pushResponse.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(pushResponse, os.Stderr, termFd, isTerm, nil)
}

// Clean up sources and the images we created
func Clean() error {
	cli, _ := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	options := types.ImageRemoveOptions{}
	tags := getTags()
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
