package parser

import (
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func hell_rings(s, version, app, name string) string {
	s = strings.ReplaceAll(s, "version_default", version)
	s = strings.ReplaceAll(s, "app_default", app)
	s = strings.ReplaceAll(s, "deploy_name", name)
	return s
}
func make_deploy(core_number, cluster_name string) {

	if _, err := os.Stat("./cluster-" + cluster_name); os.IsNotExist(err) {
		service_make(cluster_name)
	}

	f, err := os.Create("cluster-" + cluster_name + "/" + core_number + "-core.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	yaml_file, err := os.ReadFile("deploy_sample.yaml")
	if err != nil {
		panic(err)
	}
	var i interface{}
	yaml.Unmarshal(yaml_file, &i)

	s := string(yaml_file)
	s = hell_rings(s, core_number+"-core", "cluster-"+cluster_name, core_number+"-core")

	_, err = io.WriteString(f, s)
	if err != nil {
		panic(err)
	}

}
func service_make(name string) {
	dir := "cluster-" + name
	os.Mkdir(dir, os.FileMode(0755))
	f, err := os.Create(dir + "/service.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	yaml_file, err := os.ReadFile("service_sample.yaml")
	if err != nil {
		panic(err)
	}
	var i interface{}
	yaml.Unmarshal(yaml_file, &i)
	s := string(yaml_file)
	s = strings.ReplaceAll(s, "default", name)
	_, err = io.WriteString(f, s)
	if err != nil {
		panic(err)
	}

}
func Make_yaml() {
	f, _ := os.ReadFile("markup.txt")
	arr := strings.Split(strings.TrimSpace(string(f)), "\n")
	for _, elem := range arr {
		values := strings.Split(elem, ":")
		make_deploy(values[0], values[1])
	}
}
