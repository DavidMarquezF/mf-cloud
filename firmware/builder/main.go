package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	mfModules "github.com/DavidMarquezF/mf-cloud/firmware/modules"
)

func errToJsonRes(err error) map[string]string {
	return map[string]string{"err": err.Error()}
}

func writeError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	b, _ := json.Marshal(errToJsonRes(err))
	http.Error(w, string(b), http.StatusBadRequest)
}

func createFirmware(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		writeError(w, fmt.Errorf("can't read body: %w", err))
		return
	}

	var conf mfModules.FirmwareConfig
	if err := json.Unmarshal(body, &conf); err != nil {
		writeError(w, fmt.Errorf("invalid json body: %w", err))
		return
	}

	genH := BuildFileString(conf)

	path := os.Getenv("GEN_H_PATH")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700) // Create your file
	}
	f, err := os.Create(filepath.Join(path, "gen.h"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error creating file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = f.WriteString(genH)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error creating file", http.StatusInternalServerError)
		return
	}
	idfPath := os.Getenv("IDF_PATH")
	mfEmbbedPath := os.Getenv("MF_EMBEDDED_SRC_PATH")
	cmd := exec.Command(filepath.Join(idfPath, "tools", "idf.py"), "build", "-DMF_NUMBER_COMPONENTS="+fmt.Sprintf("%v", len(conf.Modules)))
	cmd.Dir = mfEmbbedPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot build successfully", http.StatusInternalServerError)
		return
	}

	fileName := "mf_embedded.bin"
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filepath.Join(mfEmbbedPath, "build", fileName))

	//fmt.Fprintf(w, , html.EscapeString(r.URL.Path))

}

func main() {
	log.Print(os.Getenv("IDF_PATH"))
	http.HandleFunc("/create-firmware", createFirmware)

	http.ListenAndServe(":8091", nil)

	config := mfModules.FirmwareConfig{
		DeviceName: "Test",
		DeviceId:   "sdasd",
		Platform:   mfModules.ESP32,
		Modules: []mfModules.Module{
			mfModules.Module{
				Id: mfModules.Ultrasound,
			},
			mfModules.Module{
				Id: mfModules.Temperature,
			},
		},
	}
	log.Print(BuildFileString(config))
}
