package Comandos

import (
	"Backend/Analizador"
	"Backend/Estructuras"
	"Backend/Mount"
	"io"
	"os"
	"strconv"
	"strings"
)

func MOunt(commandArray []string) {
	Analizador.Salida_comando += "[MENSAJE] El comando MOUNT aqui inicia\\n"
	val_path := ""
	val_name := ""
	band_path := false
	band_name := false
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
		if band_path {
			if band_name {
				index_p := Auxiliares.buscar_particion_p_e(val_path, val_name)
				if index_p != -1 {
					f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

					if err == nil {
						mbr_empty := Estructuras.MBR{}
						mbr2 := Auxiliares.struct_a_bytes(mbr_empty)
						sstruct := len(mbr2)
						lectura := make([]byte, sstruct)
						f.Seek(0, io.SeekStart)
						f.Read(lectura)
						master_boot_record := Auxiliares.bytes_a_struct_mbr(lectura)
						copy(master_boot_record.Mbr_partition[index_p].Part_status[:], "2")
						mbr_byte := Auxiliares.struct_a_bytes(master_boot_record)
						f.Seek(0, io.SeekStart)
						f.Write(mbr_byte)
						f.Close()

						if Mount.Buscar_particion(val_path, val_name, lista_montajes) {
							Analizador.Salida_comando += "[ERROR] La particion ya esta montada...\\n"
						} else {
							num := Mount.Buscar_numero(val_path, lista_montajes)
							letra := Mount.Buscar_letra(val_path, lista_montajes)
							id := "30" + strconv.Itoa(num) + letra

							var n *Mount.Nodo = Mount.New_nodo(id, val_path, val_name, letra, num)
							Mount.Insertar(n, lista_montajes)
							Analizador.Salida_comando += "[SUCCES] Particion montada con exito!\\n"
							Analizador.Salida_comando += Mount.Imprimir_contenido(lista_montajes)
						}
					} else {
						Analizador.Salida_comando += "[ERROR] No se encuentra el disco...\\n"
					}
				} else {
					index_p := Auxiliares.buscar_particion_l(val_path, val_name)
					if index_p != -1 {
						f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

						if err == nil {
							ebr_empty := EBR{}
							ebr2 := Auxiliares.struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)
							lectura := make([]byte, sstruct)
							f.Seek(int64(index_p), io.SeekStart)
							f.Read(lectura)
							extended_boot_record := Auxiliares.bytes_a_struct_ebr(lectura)
							copy(extended_boot_record.Part_status[:], "2")
							mbr_byte := Auxiliares.struct_a_bytes(extended_boot_record)
							f.Seek(int64(index_p), io.SeekStart)
							f.Write(mbr_byte)
							f.Close()

							if Mount.Buscar_particion(val_path, val_name, lista_montajes) {
								Analizador.Salida_comando += "[ERROR] La particion ya esta montada...\\n"
							} else {
								num := Mount.Buscar_numero(val_path, lista_montajes)
								letra := Mount.Buscar_letra(val_path, lista_montajes)
								id := "30" + strconv.Itoa(num) + letra

								var n *Mount.Nodo = Mount.New_nodo(id, val_path, val_name, letra, num)
								Mount.Insertar(n, lista_montajes)
								Analizador.Salida_comando += "[Exito] Particion montada con exito!\\n"
								Analizador.Salida_comando += Mount.Imprimir_contenido(lista_montajes)
							}
						} else {
							Analizador.Salida_comando += "[ERROR] No se encuentra el disco...\\n"
						}

					} else {
						Analizador.Salida_comando += "[ERROR] No se encuentra la particion a montar...\\n"
					}
				}
			} else {
				Analizador.Salida_comando += "[ERROR] Parametro -name no definido...\\n"
			}
		} else {
			Analizador.Salida_comando += "[ERROR] Parametro -path no definido...\\n"
		}
	}

	Analizador.Salida_comando += "[MENSAJE] El comando MOUNT aqui finaliza\\n"
}
