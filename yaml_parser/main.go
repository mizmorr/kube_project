package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Person struct {
	Name    string   `yaml:"name"`
	Age     int      `yaml:"age"`
	Hobbies []string `yaml:"hobbies"`
}

func hell_rings(s, version, app, name string) string {
	s = strings.ReplaceAll(s, "version_default", version)
	s = strings.ReplaceAll(s, "app_default", app)
	s = strings.ReplaceAll(s, "deploy_name", name)
	return s
}
func main() {
	// Person := Person{
	// 	Name:    "Kevin Kelche",
	// 	Age:     20,
	// 	Hobbies: []string{"Programming", "Reading", "Writing"},
	// }

	// yamlFile, err := yaml.Marshal(&Person)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(yamlFile))

	f, err := os.Create("six_core.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	yaml_file, err := os.ReadFile("sample.yaml")
	if err != nil {
		panic(err)
	}
	var i interface{}
	yaml.Unmarshal(yaml_file, &i)

	s := string(yaml_file)
	s = hell_rings(s, "six-core", "cluster-second", "six-core-deploy")
	fmt.Println(s)
	// yamlFile, err := yaml.Marshal(&i)
	// if err != nil {
	// 	panic(err)
	// }
	_, err = io.WriteString(f, s)
	if err != nil {
		panic(err)
	}
}
