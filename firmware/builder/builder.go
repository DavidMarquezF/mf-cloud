package main

import (
	"strings"

	mfModules "github.com/DavidMarquezF/mf-cloud/firmware/modules"
)

type moduleInfo struct {
	Include          string
	DefaultName      string
	DefaultUrl       string
	CreateResourceCb string
	InitCb           string
	DestroyCb        string
}

const IncludeDir = "../"

var modules = map[mfModules.ModuleId]moduleInfo{
	mfModules.Ultrasound: moduleInfo{
		Include:          "mf_ultrasound.h",
		DefaultName:      "Dist",
		DefaultUrl:       "/dist",
		CreateResourceCb: "mf_ultrasound_create_resource",
		InitCb:           "mf_ultrasound_init",
		DestroyCb:        "mf_ultrasound_destroy",
	},
	mfModules.Temperature: moduleInfo{
		Include:          "mf_temp.h",
		DefaultName:      "temp",
		DefaultUrl:       "/temp",
		CreateResourceCb: "mf_temp_create_resource",
		InitCb:           "mf_temp_init",
		DestroyCb:        "mf_temp_destroy",
	},
	mfModules.Semaphore: moduleInfo{
		Include:          "mf_semaphore.h",
		DefaultName:      "semaphore",
		DefaultUrl:       "/semaph",
		CreateResourceCb: "mf_semaphore_create_resource",
		InitCb:           "mf_semaphore_init",
		DestroyCb:        "mf_semaphore_destroy",
	},
	mfModules.Button: moduleInfo{
		Include:          "mf_button.h",
		DefaultName:      "button",
		DefaultUrl:       "/button",
		CreateResourceCb: "mf_button_create_resource",
		InitCb:           "mf_button_init",
		DestroyCb:        "mf_button_destroy",
	},
}

func writeStringNL(builder *strings.Builder, val string) {
	builder.WriteString(val + "\n")
}

func startHeader(builder *strings.Builder) {
	writeStringNL(builder, "#ifndef _MF_GEN_COMPONENTS_H_")
	writeStringNL(builder, "#define _MF_GEN_COMPONENTS_H_")
}

func endHeader(builder *strings.Builder) {
	builder.WriteString("#endif")
}

func stringify(str string) string {
	return "\"" + str + "\""
}

func addInclude(builder *strings.Builder, id mfModules.ModuleId) {
	writeStringNL(builder, "#include "+stringify(IncludeDir+modules[id].Include))
}

func getPropertyString(name string, value string) string {
	return "." + name + " = " + value
}

func BuildFileString(config mfModules.FirmwareConfig) string {
	var b strings.Builder
	builder := &b
	startHeader(builder)
	writeStringNL(builder, "#include \"../mf_component_handler.h\"")

	// Add modules include
	for _, v := range config.Modules {
		addInclude(builder, v.Id)
	}

	writeStringNL(builder, "mf_component_config_t generated_components[] = {")
	for i, v := range config.Modules {
		module := modules[v.Id]
		writeStringNL(builder, "{")

		builder.WriteString(getPropertyString("url", stringify(module.DefaultUrl)))
		writeStringNL(builder, ",")
		builder.WriteString(getPropertyString("name", stringify(module.DefaultName)))
		writeStringNL(builder, ",")
		builder.WriteString(getPropertyString("create_resource_callback", module.CreateResourceCb))
		writeStringNL(builder, ",")
		builder.WriteString(getPropertyString("init_callback", module.InitCb))
		writeStringNL(builder, ",")
		builder.WriteString(getPropertyString("destroy_callback", module.DestroyCb))

		if i != len(config.Modules)-1 {
			writeStringNL(builder, "},")
		} else {
			writeStringNL(builder, "}")
		}
	}

	writeStringNL(builder, "};")
	endHeader(builder)

	return builder.String()
}
