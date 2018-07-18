// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/magefile/mage/sh"
	// mg contains helpful utility functions, like Deps
)

type Specification struct {
	Os                 string `required:"true"`
	Nodeversion        string `required:"true"`
	V8version          string `required:"true"`
	Dockerfile         string `required:"true"`
	Imagetag           string `required:"true"`
	Latest             bool   `default:false`
	Lts                string
	Majortag           string `required:"true"`
	Minortag           string `required:"true"`
	Imagename          string `required:"true"`
	Npmversion         string `required:"true"`
	Fromdata           string `required:"true"`
	Prebuilt           bool   `default:false`
	Dockeruser         string `required:"true"`
	Dockerpass         string `required:"true"`
	Rhsecret           string `required:"true"`
	Rhendpoint         string `required:"true"`
	Rhproject          string `required:"true"`
	AwsAccessKey       string `required:"true" split_words:"true"`
	AwsSecretAccessKey string `required:"true" split_words:"true"`
}

const bucketName string = "sourcecode-nearform-export-compliance"

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
		s.Prebuilt,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	return s
}

// get base image name
func imageName(s Specification) string {
	return s.Imagename + ":" + s.Imagetag
}

func isLatest(s Specification) string {
	if s.Latest {
		return "T"
	}
	return ""
}

func preBuiltEnv(s Specification) string {
	if s.Prebuilt {
		return "T"
	}
	return " "
}

func archiveName(s Specification) string {
	dashedImageName := strings.Replace(s.Imagename, "/", "-", -1)
	return "sources-" + dashedImageName + ".tgz"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// show configuratio
func ShowConfig() {
	fmt.Println(ParseEnvVars())
}

// Get Node.js sources, build
func InstallSources() error {
	s := ParseEnvVars()
	dir, err := os.Getwd()
	check(err)
	path := dir + "/src/node-v" + s.Nodeversion + "-linux-x64.tar.gz"
	if s.Prebuilt {
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
		"PREBUILT": preBuiltEnv(s),
		"OS":       s.Os,
	}
	fmt.Println("Installing Node.js sources...")
	_, err = sh.Exec(envs, os.Stdout, os.Stdout, "./contrib/etc/get_node_source.sh", s.Nodeversion, dir+"/src/")
	check(err)
	return err
}

// Squash the image using a shell command
func Squash() {
	fmt.Println("Squashing the image...")
	s := ParseEnvVars()
	tags := getTags(s)
	var envs = map[string]string{}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "docker-squash", tags[0], "-f", fromImage(s), "-t", tags[0])
	check(err)
}

// Run a test on the image
func Test() {
	fmt.Println("Testing the image...")
	sh.Rm("src/.")
	s := ParseEnvVars()
	fmt.Println("Cleanup image " + imageName(s))
	var envs = map[string]string{
		"BUILDER":      imageName(s),
		"NODE_VERSION": s.Nodeversion,
	}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "test/run.sh")
	check(err)
}

// create archive with sources
func Archive() error {
	s := ParseEnvVars()
	os.MkdirAll("dist", os.ModePerm)
	var envs = map[string]string{
		"ARCHIVE":      archiveName(s),
		"NODE_VERSION": s.Nodeversion,
		"OS":           s.Os,
		"DOCKERFILE":   s.Dockerfile,
		"IMAGE_TAG":    s.Imagetag,
		"LATEST":       isLatest(s),
		"MAJOR_TAG":    s.Majortag,
		"MINOR_TAG":    s.Minortag,
		"IMAGE_NAME":   s.Imagename,
		"NPM_VERSION":  s.Npmversion,
	}
	_, err := sh.Exec(envs,
		os.Stdout,
		os.Stdout,
		"contrib/etc/archive.sh",
	)
	check(err)
	return err
}

// upload archive to S3
// Upload input parameters
func Upload() error {
	// s3cmd put $(ARCHIVE) "$(S3BUCKET)/$(ARCHIVE)"
	s := ParseEnvVars()
	var envs = map[string]string{
		"AWS_ACCESS_KEY_ID":     s.AwsAccessKey,
		"AWS_SECRET_ACCESS_KEY": s.AwsSecretAccessKey,
	}
	_, err := sh.Exec(envs,
		os.Stdout,
		os.Stdout,
		"s3cmd",
		"put",
		archiveName(s),
		"s3://"+bucketName+"/sources/"+archiveName(s),
	)
	check(err)
	return err
}

// helper functions
// get the FROM image string
func fromImage(s Specification) string {
	b, err := ioutil.ReadFile(s.Dockerfile) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	re := regexp.MustCompile(`FROM (.*)`)
	matches := re.FindStringSubmatch(str)
	fmt.Println("Base image name: " + matches[1])
	return matches[1]
}
