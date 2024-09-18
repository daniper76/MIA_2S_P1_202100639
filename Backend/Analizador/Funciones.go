package Analizador

import (
	"Backend/Comandos"
	"Backend/Estructuras"
	"Backend/Mount"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var Lista_montajes *Mount.Lista = Mount.New_lista()
var Salida_comando string = ""
var GraphDot string = ""

func Analizar() {
	mux := http.NewServeMux()
	mux.HandleFunc("/analizar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var Content Estructuras.Cmd_API
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &Content)
		split_cmd(Content.Cmd)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + Salida_comando + `" }`))
		Salida_comando = ""
	})
	fmt.Println("Servidor en el puerto 5000")
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

// Ejecuta comando linea por linea
func split_cmd(cmd string) {
	arr_com := strings.Split(cmd, "\n")

	for i := 0; i < len(arr_com); i++ {
		if arr_com[i] != "" {
			split_comando(arr_com[i])
			Salida_comando += "\\n"
		}
	}
}

func split_comando(comando string) {
	var commandArray []string
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)
	band_comentario := false

	if strings.Contains(comando, "#") {
		band_comentario = true
		Salida_comando += comando + "\\n"
	} else {
		commandArray = strings.Split(comando, " -")
	}
	if !band_comentario {
		ejecutar_comando(commandArray)
	}
}

func ejecutar_comando(commandArray []string) {

	data := strings.ToLower(commandArray[0])

	if data == "mkdisk" {
		Comandos.Mkdisk(commandArray)
	} else if data == "rmdisk" {
		Comandos.Rmdisk(commandArray)
	} else if data == "fdisk" {
		Comandos.Fdisk(commandArray)
	} else if data == "mount" {
		Comandos.MOunt(commandArray)
	} else if data == "rep" {
		Comandos.Rep(commandArray)
	} else {
		Salida_comando += "[ERROR] El comando no fue reconocido...\\n"
	}
}
