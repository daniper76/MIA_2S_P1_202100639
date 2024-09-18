package Comandos

import (
	"Backend/Analizador"
	"Backend/Graficar"
	"Backend/Mount"
	"strings"
)

func Rep(commandArray []string) {
	Analizador.Salida_comando += "[MENSAJE] El comando REP aqui inicia...\\n"
	val_name := ""
	val_path := ""
	val_id := ""
	band_name := false
	band_path := false
	band_id := false
	band_ruta := false
	band_error := false
	Analizador.GraphDot = ""

	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]
		switch {
		case strings.Contains(data, "name="):
			if band_name {
				Analizador.Salida_comando += "[ERROR] El parametro -name ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_name = true
			val_name = strings.Replace(val_data, "\"", "", 2)
		case strings.Contains(data, "path="):
			if band_path {
				Analizador.Salida_comando += "[ERROR] El parametro -path ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_path = true
			val_path = strings.Replace(val_data, "\"", "", 2)
		case strings.Contains(data, "id="):
			if band_id {
				Analizador.Salida_comando += "[ERROR] El parametro -id ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_id = true
			val_id = val_data
		case strings.Contains(data, "ruta="):
			if band_ruta {
				Analizador.Salida_comando += "[ERROR] El parametro -ruta ya fue ingresado...\\n"
				band_error = true
				break
			}
			band_ruta = true
		default:
			Analizador.Salida_comando += "[ERROR] Parametro no valido...\\n"
		}
	}

	if !band_error {
		if band_path {
			if band_name {
				if band_id {
					var aux *Mount.Nodo = Mount.Obtener_nodo(val_id, lista_montajes)

					if aux != nil {
						// Reportes validos
						if val_name == "disk" {
							Graficar.Graficar_disk(aux.Direccion, val_path)
						} else {
							Analizador.Salida_comando += "[ERROR] Reporte no valido...\\n"
						}
					} else {
						Analizador.Salida_comando += "[ERROR] No encuentra la particion...\\n"
					}
				} else {
					Analizador.Salida_comando += "[ERROR] El parametro -id no fue ingresado...\\n"
				}
			} else {
				Analizador.Salida_comando += "[ERROR] El parametro -name no fue ingresado...\n"
			}
		} else {
			Analizador.Salida_comando += "[ERROR] El parametro -path no fue ingresado...\\n"
		}
	}
	Analizador.Salida_comando += "\\n[MENSAJE] El comando REP aqui finaliza...\\n"
}
