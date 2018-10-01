// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/spf13/viper"
)

const bucketName string = "sourcecode-nearform-export-compliance"
const defaultRepo string = "https://github.com/nodejs/node.git"

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

func config() *viper.Viper {
	configFile := os.Getenv("CONFIG_FILE")
	v := viper.New()
	v.SetDefault("Latest", false)
	v.SetDefault("Lts", "")
	v.SetDefault("Fromdata", "{}")
	v.SetDefault("Prebuilt", false)
	v.SetDefault("Repo", "https://github.com/nodejs/node.git")
	v.SetDefault("Commithash", "")
	if len(configFile) > 0 {
		extension := filepath.Ext(configFile)
		bareName := configFile[0 : len(configFile)-len(extension)]
		fmt.Println("Reading config from file: " + configFile)
		v.SetConfigName(bareName)
		v.AddConfigPath(".")
		err := v.ReadInConfig()
		check(err)
		fmt.Println(v.ConfigFileUsed())
	} else {
		fmt.Println("Reading config from ENV vars")
		v.BindEnv("Os", "OS")
		v.BindEnv("Nodeversion", "VERSION")
		v.BindEnv("Npmversion", "NPM")
		v.BindEnv("V8version", "V8")
		v.BindEnv("Imagetag", "TAG")
		v.BindEnv("Majortag", "MAJOR")
		v.BindEnv("Minortag", "MINOR")
		v.BindEnv("LatestString", "LATEST")
		if len(v.GetString("LatestString")) != 0 {
			latest, err := strconv.ParseBool(v.GetString("LatestString"))
			check(err)
			v.Set("Latest", latest)
		} else {
			v.Set("Latest", false)
		}
		v.BindEnv("Lts", "LTS")
		v.BindEnv("Fromdata", "FROM_DATA")
		v.BindEnv("PrebuiltString", "PREBUILT")
		if len(v.GetString("PrebuiltString")) != 0 {
			prebuilt, err := strconv.ParseBool(v.GetString("PrebuiltString"))
			check(err)
			v.Set("Prebuilt", prebuilt)
		} else {
			v.Set("Prebuilt", false)
		}
		v.BindEnv("Repo", "REPO")
		v.BindEnv("Commithash", "COMMIT_HASH")
	}
	v.BindEnv("Dockeruser", "DOCKER_USER")
	v.BindEnv("Dockerpass", "DOCKER_PASS")
	v.BindEnv("Rhsecret", "RH_SECRET")
	v.BindEnv("Rhendpoint", "RH_ENDPOINT")
	v.BindEnv("Rhproject", "RH_PROJECT")

	v.Set("Dockerfile", "image/"+v.GetString("Os")+"/Dockerfile")
	return v
}

func ShowConfig() {
	v := config()
	fmt.Println("OS: " + v.GetString("Os"))
	fmt.Println("NODE_VERSION: " + v.GetString("Nodeversion"))
	fmt.Println("NPM_VERSION: " + v.GetString("Npmversion"))
	fmt.Println("V8_VERSION: " + v.GetString("V8version"))
	fmt.Println("IMAGE_TAG: " + v.GetString("Imagetag"))
	fmt.Println("MAJOR_TAG: " + v.GetString("Majortag"))
	fmt.Println("MINOR_TAG: " + v.GetString("Minortag"))
	fmt.Println("LATEST: " + strconv.FormatBool(v.GetBool("Latest")))
	fmt.Println("LTS: " + v.GetString("Lts"))
	fmt.Println("FROM_DATA: " + v.GetString("Fromdata"))
	fmt.Println("PREBUILT: " + strconv.FormatBool(v.GetBool("Prebuilt")))
	fmt.Println("REPO: " + v.GetString("Repo"))
	fmt.Println("COMMIT_HASH: " + v.GetString("Commithash"))
}

// get base image name
func imageName(v *viper.Viper) string {
	return "nearform/" + v.GetString("Os") + "-s2i-nodejs"
}

func dockerFile(v *viper.Viper) string {
	return "image/" + v.GetString("Os") + "/Dockerfile"
}

func isLatest(v *viper.Viper) string {
	if v.GetBool("Latest") {
		return "T"
	}
	return ""
}

func preBuiltEnv(v *viper.Viper) string {
	if v.GetBool("Prebuilt") {
		fmt.Println("Prebuilt=true, return 'T' for ENV var")
		return "T"
	}
	fmt.Println("Prebuilt=false, return ' ' for ENV var")
	return " "
}

func archiveName(v *viper.Viper) string {
	dashedImageName := strings.Replace(imageName(config()), "/", "-", -1)
	return "sources-" + dashedImageName + ".tgz"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Get Node.js sources, build
func InstallSources() error {
	v := config()
	dir, err := os.Getwd()
	check(err)
	path := dir + "/src/node-v" + v.GetString("Nodeversion") + "-linux-x64.tar.gz"
	if v.GetBool("Prebuilt") {
		log.Println("checking if " + path + " exists")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			runDownloadScript(v)
		} else {
			log.Println("already exists, no need to download again.")
		}
	} else {
		runDownloadScript(v)
	}
	return err
}

func runDownloadScript(v *viper.Viper) error {
	dir, err := os.Getwd()
	var envs = map[string]string{
		"PREBUILT": preBuiltEnv(v),
		"OS":       v.GetString("Os"),
	}
	fmt.Println("Installing Node.js sources...")
	_, err = sh.Exec(
		envs,
		os.Stdout,
		os.Stdout,
		"./contrib/etc/get_node_source.sh",
		v.GetString("Nodeversion"),
		dir+"/src/",
		v.GetString("Repo"),
		v.GetString("Commithash"),
	)
	check(err)
	return err
}

// Squash the image using a shell command
func Squash() {
	fmt.Println("Squashing the image...")
	tags := getTags()
	fmt.Println(tags)
	var envs = map[string]string{}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "docker-squash", tags[0], "-f", fromImage(), "-t", tags[0], "--output-path", "test")
	check(err)
}

// Run a basic test on the image
func Test() {
	fmt.Println("Testing the image...")
	sh.Rm("src/.")
	v := config()
	fmt.Println("Cleanup image " + imageName(v))
	var envs = map[string]string{
		"BUILDER":      imageName(v) + ":" + v.GetString("Imagetag"),
		"NODE_VERSION": v.GetString("Nodeversion"),
	}
	_, err := sh.Exec(envs, os.Stdout, os.Stdout, "test/run.sh")
	check(err)
}

// create archive with sources
func Archive() error {
	v := config()
	os.MkdirAll("dist", os.ModePerm)
	var envs = map[string]string{
		"ARCHIVE":      archiveName(v),
		"NODE_VERSION": v.GetString("Nodeversion"),
		"OS":           v.GetString("Os"),
		"DOCKERFILE":   v.GetString("Dockerfile"),
		"IMAGE_TAG":    v.GetString("Imagetag"),
		"LATEST":       isLatest(v),
		"MAJOR_TAG":    v.GetString("Majortag"),
		"MINOR_TAG":    v.GetString("Minortag"),
		"IMAGE_NAME":   v.GetString("Imagename"),
		"NPM_VERSION":  v.GetString("Npmversion"),
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
	v := config()
	var envs = map[string]string{}
	_, err := sh.Exec(envs,
		os.Stdout,
		os.Stdout,
		"s3cmd",
		"put",
		archiveName(v),
		"s3://"+bucketName+"/sources/"+archiveName(config()),
	)
	check(err)
	return err
}

// helper functions
// get the FROM image string
func fromImage() string {
	v := config()
	b, err := ioutil.ReadFile(v.GetString("Dockerfile"))
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	re := regexp.MustCompile(`FROM (.*)`)
	matches := re.FindStringSubmatch(str)
	fmt.Println("Base image name: " + matches[1])
	return matches[1]
}
