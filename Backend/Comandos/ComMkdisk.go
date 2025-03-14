package Comandos

import (
	"Backend/Analizador"
	"Backend/Auxiliares"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Mkdisk(commandArray []string) {
	Analizador.Salida_comando += "[MENSAJE] El comando MKDISK aqui inicia\\n"
	val_size := 0
	val_fit := ""
	val_unit := ""
	val_path := ""

	band_size := false
	band_fit := false
	band_unit := false
	band_path := false
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
				band_error = true
				Analizador.Salida_comando += "[ERROR] En la conversion a entero\\n"
				break
			}

			if val_size < 0 {
				band_error = true
				Analizador.Salida_comando += "[ERROR] El parametro -size es negativo...\\n"
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
		case strings.Contains(data, "unit="):
			if band_unit {
				Analizador.Salida_comando += "[ERROR] El parametro -unit ya fue ingresado...\\n"
				band_error = true
				break
			}
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "k" || val_unit == "m" {
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
		default:
			Analizador.Salida_comando += "[ERROR] Parametro no valido...\\n"
		}
	}

	if !band_error {
		if band_path {
			if band_size {
				total_size := 1024
				master_boot_record := MBR{}
				Auxiliares.Crear_disco(val_path)
				fecha := time.Now()
				str_fecha := fecha.Format("02/01/2006 15:04:05")
				copy(master_boot_record.Mbr_fecha_creacion[:], str_fecha)
				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100
				num_random := rand.Intn(max-min+1) + min
				copy(master_boot_record.Mbr_dsk_signature[:], strconv.Itoa(int(num_random)))
				if band_fit {
					copy(master_boot_record.Dsk_fit[:], val_fit)
				} else {
					copy(master_boot_record.Dsk_fit[:], "f")
				}
				if band_unit {
					if val_unit == "m" {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
						total_size = val_size * 1024
					} else {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024)))
						total_size = val_size
					}
				} else {
					copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
					total_size = val_size * 1024
				}
				for i := 0; i < 4; i++ {
					copy(master_boot_record.Mbr_partition[i].Part_status[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_start[:], "-1")
					copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_name[:], "")
				}
				str_total_size := strconv.Itoa(total_size)
				cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+val_path+"\" bs=1024 count="+str_total_size)
				cmd.Dir = "/"
				_, err := cmd.Output()
				if err != nil {
					Analizador.Salida_comando += "[ERROR] Al ejecuatar comando en consola\\n"
				}
				f, err := os.OpenFile(val_path, os.O_RDWR, 0660)
				if err != nil {
					Analizador.Salida_comando += "[ERROR] Al abrir el archivo\\n"
				} else {
					mbr_byte := struct_a_bytes(master_boot_record)
					f.Seek(0, io.SeekStart)
					f.Write(mbr_byte)
					f.Close()

					Analizador.Salida_comando += "[Exito] El disco fue creado con exito!\\n"
				}
			}
		}
	}

	Analizador.Salida_comando += "[MENSAJE] El comando MKDISK aqui finaliza\\n"
}
