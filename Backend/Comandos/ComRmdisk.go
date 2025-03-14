package Comandos

import (
	"Backend/Analizador"
	"os"
	"os/exec"
	"strings"
)

func Rmdisk(commandArray []string) {
	Analizador.Salida_comando += "[MENSAJE] El comando RMDISK aqui inicia\\n"
	val_path := ""
	band_path := false
	band_error := false

	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]
		switch {
		case strings.Contains(data, "path="):
			if band_path {
				Analizador.Salida_comando += "[ERROR] El parametro -path ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_path = true
			val_path = strings.Replace(val_data, "\"", "", 2)
		default:
			Analizador.Salida_comando += "[ERROR] Parametro no valido...\\n"
		}
	}

	if !band_error {
		if band_path {
			_, e := os.Stat(val_path)

			if e != nil {
				if os.IsNotExist(e) {
					Analizador.Salida_comando += "[ERROR] No existe el disco que desea eliminar...\\n"
					band_path = false
				}
			} else {
				cmd := exec.Command("/bin/sh", "-c", "rm \""+val_path+"\"")
				cmd.Dir = "/"
				_, err := cmd.Output()

				if err != nil {
					Analizador.Salida_comando += "[ERROR] Al ejecutar un comando en consola\\n"
				} else {
					Analizador.Salida_comando += "[SUCCES] El Disco fue eliminado!\\n"
				}

				band_path = false
			}
		} else {
			Analizador.Salida_comando += "[ERROR] el parametro -path no fue ingresado...\\n"
		}
	}

	Analizador.Salida_comando += "[MENSAJE] El comando RMDISK aqui finaliza\\n"
}
