package Comandos

import (
	"Backend/Analizador"
	"Backend/Auxiliares"
	"strconv"
	"strings"
)

func Fdisk(commandArray []string) {
	Analizador.Salida_comando += "[MENSAJE] El comando FDISK aqui inicia\\n"
	val_size := 0
	val_unit := ""
	val_path := ""
	val_type := ""
	val_fit := ""
	val_name := ""
	band_size := false
	band_unit := false
	band_path := false
	band_type := false
	band_fit := false
	band_name := false
	band_error := false

	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "size="):
			if band_size {
				Analizador.Salida_comando += "[ERROR] El parametro -size ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_size = true

			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			if err != nil {
				Analizador.Salida_comando += "[ERROR] Al convertir a entero\\n"
				band_error = true
				break
			}

			if val_size < 0 {
				band_error = true
				Analizador.Salida_comando += "[ERROR] El parametro -size es negativo...\\n"
			}
		case strings.Contains(data, "unit="):
			if band_unit {
				Analizador.Salida_comando += "[ERROR] El parametro -unit ya fue ingresado...\\n"
				band_error = true
				break
			}
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "b" || val_unit == "k" || val_unit == "m" {
				band_unit = true
			} else {
				Analizador.Salida_comando += "[ERROR] El Valor del parametro -unit no es valido...\\n"
				band_error = true
			}
		case strings.Contains(data, "path="):
			if band_path {
				Analizador.Salida_comando += "[ERROR] El parametro -path ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_path = true
			val_path = strings.Replace(val_data, "\"", "", 2)
		case strings.Contains(data, "type="):
			if band_type {
				Analizador.Salida_comando += "[ERROR] El parametro -type ya fue ingresado...\\n"
				band_error = true
				break
			}
			val_type = strings.Replace(val_data, "\"", "", 2)
			val_type = strings.ToLower(val_type)

			if val_type == "p" || val_type == "e" || val_type == "l" {
				band_type = true
			} else {
				Analizador.Salida_comando += "[ERROR] El Valor del parametro -type no es valido...\\n"
				band_error = true
			}
		case strings.Contains(data, "fit="):
			if band_fit {
				Analizador.Salida_comando += "[ERROR] El parametro -fit ya fue ingresado...\\n"
				band_error = true
				break
			}
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				band_fit = true
				val_fit = "w"
			} else {
				Analizador.Salida_comando += "[ERROR] El Valor del parametro -fit no es valido...\\n"
				band_error = true
				break
			}
		case strings.Contains(data, "name="):
			if band_name {
				Analizador.Salida_comando += "[ERROR] El parametro -name ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_name = true
			val_name = strings.Replace(val_data, "\"", "", 2)
		default:
			Analizador.Salida_comando += "[ERROR] Parametro no valido...\\n"
		}
	}

	if !band_error {
		if band_size {
			if band_path {
				if band_name {
					if band_type {
						if val_type == "p" {
							Auxiliares.Crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
						} else if val_type == "e" {
							Auxiliares.Crear_particion_extendia(val_path, val_name, val_size, val_fit, val_unit)
						} else {
							Auxiliares.Crear_particion_logica(val_path, val_name, val_size, val_fit, val_unit)
						}
					} else {
						Auxiliares.Crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
					}
				} else {
					Analizador.Salida_comando += "[ERROR] El parametro -name no fue ingresado\\n"
				}
			} else {
				Analizador.Salida_comando += "[ERROR] El parametro -path no fue ingresado\\n"
			}
		} else {
			Analizador.Salida_comando += "[ERROR] El parametro -size no fue ingresado\\n"
		}
	}

	Analizador.Salida_comando += "[MENSAJE] El comando FDISK aqui finaliza\\n"
}
