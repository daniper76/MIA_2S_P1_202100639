package Auxiliares

import (
	"Backend/Analizador"
	"Backend/Estructuras"
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func Crear_disco(ruta string) {
	aux, err := filepath.Abs(ruta)
	if err != nil {
		Analizador.Salida_comando += "[ERROR] Al abrir el archivo\\n"
	}
	cmd1 := exec.Command("/bin/sh", "-c", "echo 253097 | sudo -S mkdir -p '"+filepath.Dir(aux)+"'")
	cmd1.Dir = "/"
	_, err = cmd1.Output()

	if err != nil {
		Analizador.Salida_comando += "[ERROR] Al ejecutar el comando\\n"
	}
	cmd2 := exec.Command("/bin/sh", "-c", "echo 253097 | sudo -S chmod -R 777 '"+filepath.Dir(aux)+"'")
	cmd2.Dir = "/"
	_, err = cmd2.Output()
	if err != nil {
		Analizador.Salida_comando += "[ERROR] Error al ejecutar el comando\\n"
	}
	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			Analizador.Salida_comando += "[FAILURE] No se pudo crear el disco...\\n"
		}
	}
}

func Crear_particion_primaria(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := Estructuras.MBR{}
	var empty [100]byte

	if fit != "" {
		aux_fit = fit
	} else {
		aux_fit = "w"
	}

	if unit != "" {
		aux_unit = unit

		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			size_bytes = size * 1024
		} else {
			size_bytes = size * 1024 * 1024
		}
	} else {
		size_bytes = size * 1024
	}

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err != nil {
		Analizador.Salida_comando += "[ERROR] No existe un disco duro con ese nombre...\\n"
	} else {
		band_particion := false
		num_particion := 0

		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_start := ""

			for i := 0; i < 4; i++ {
				s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")

				if s_part_start == "-1" {
					band_particion = true
					num_particion = i
					break
				}
			}

			if band_particion {
				espacio_usado := 0
				s_part_size := ""
				i_part_size := 0
				s_part_status := ""

				for i := 0; i < 4; i++ {
					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					if s_part_status != "1" {
						espacio_usado += i_part_size
					}
				}

				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

				espacio_disponible := i_tamaño_disco - espacio_usado

				Analizador.Salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(espacio_disponible) + " Bytes\\n"
				Analizador.Salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes\\n"

				if espacio_disponible >= size_bytes {
					if !Existe_particion(direccion, nombre) {
						s_dsk_fit := string(master_boot_record.Dsk_fit[:])
						s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

						if s_dsk_fit == "f" {
							copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

							if num_particion == 0 {
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
							} else {
								s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
								s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
								i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

								s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
								s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
								i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant))
							}

							copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

							mbr_byte := Struct_a_bytes(master_boot_record)

							f.Seek(0, io.SeekStart)
							f.Write(mbr_byte)

							s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							f.Seek(int64(i_part_start), io.SeekStart)

							for i := 0; i < size_bytes; i++ {
								f.Write([]byte{1})
							}

							Analizador.Salida_comando += "[SUCCES] La Particion primaria fue creada con exito!\\n"
						} else if s_dsk_fit == "b" {
							best_index := num_particion

							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_best := ""
							i_part_start_best := 0
							s_part_start_best_ant := ""
							i_part_start_best_ant := 0
							s_part_size_best := ""
							i_part_size_best := 0
							s_part_size_best_ant := ""
							i_part_size_best_ant := 0

							for i := 0; i < 4; i++ {
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
										s_part_size_best = strings.Trim(s_part_size_best, "\x00")
										i_part_size_best, _ = strconv.Atoi(s_part_size_best)

										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_best > i_part_size_act {
											best_index = i
											break
										}
									}
								}
							}

							copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

							if best_index == 0 {
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
							} else {
								s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
								s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
								i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

								s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
								s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
								i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
							}

							copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

							mbr_byte := Struct_a_bytes(master_boot_record)

							f.Seek(0, io.SeekStart)
							f.Write(mbr_byte)

							s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])
							s_part_start_best = strings.Trim(s_part_start_best, "\x00")
							i_part_start_best, _ = strconv.Atoi(s_part_start_best)

							f.Seek(int64(i_part_start_best), io.SeekStart)

							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}

							Analizador.Salida_comando += "[Exito] La Particion primaria fue creada con exito!\\n"
						}
					} else {
						Analizador.Salida_comando += "[ERROR] Ya existe una partición con ese nombre\\n"
					}
				} else {
					Analizador.Salida_comando += "[ERROR] No hay espacio suficiente para la partición\\n"
				}
			} else {
				Analizador.Salida_comando += "[ERROR] No hay espacio para crear una partición\\n"
			}
		} else {
			Analizador.Salida_comando += "[ERROR] El disco está vacío\\n"
		}

		f.Close()
	}
}

func Crear_particion_extendia(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := Estructuras.MBR{}
	var empty [100]byte

	if fit != "" {
		aux_fit = fit
	} else {
		aux_fit = "w"
	}

	if unit != "" {
		aux_unit = unit

		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			size_bytes = size * 1024
		} else {
			size_bytes = size * 1024 * 1024
		}
	} else {
		size_bytes = size * 1024
	}

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err != nil {
		Analizador.Salida_comando += "[ERROR] Al abrir el archivo\\n"
	} else {
		band_particion := false
		band_extendida := false
		num_particion := 0

		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""

			for i := 0; i < 4; i++ {
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					band_extendida = true
					break
				}
			}

			if !band_extendida {
				s_part_start := ""
				s_part_status := ""
				s_part_size := ""
				i_part_size := 0

				for i := 0; i < 4; i++ {
					s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")

					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					if s_part_start == "-1" || (s_part_status == "1" && i_part_size >= size_bytes) {
						band_particion = true
						num_particion = i
						break
					}
				}

				if band_particion {
					espacio_usado := 0

					for i := 0; i < 4; i++ {
						s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
						s_part_status = strings.Trim(s_part_status, "\x00")

						if s_part_status != "1" {
							s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
							s_part_size = strings.Trim(s_part_size, "\x00")
							i_part_size, _ = strconv.Atoi(s_part_size)

							espacio_usado += i_part_size
						}
					}

					s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
					s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
					i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

					espacio_disponible := i_tamaño_disco - espacio_usado

					Analizador.Salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(espacio_disponible) + " Bytes\\n"
					Analizador.Salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes\\n"

					if espacio_disponible >= size_bytes {
						if !Existe_particion(direccion, nombre) {
							s_dsk_fit := string(master_boot_record.Dsk_fit[:])
							s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

							if s_dsk_fit == "f" {
								copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

								if num_particion == 0 {
									mbr_empty_byte := struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
								} else {
									s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
									s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
									i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

									s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
									s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
									i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

									copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant))
								}

								copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

								mbr_byte := struct_a_bytes(master_boot_record)

								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								f.Seek(int64(i_part_start), io.SeekStart)

								extended_boot_record := EBR{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}

								Analizador.Salida_comando += "[SUCCES] La Particion extendida fue creada con exito!\\n"
							} else if s_dsk_fit == "b" {
								best_index := num_particion

								s_part_start_act := ""
								s_part_status_act := ""
								s_part_size_act := ""
								i_part_size_act := 0
								s_part_start_best := ""
								i_part_start_best := 0
								s_part_start_best_ant := ""
								i_part_start_best_ant := 0
								s_part_size_best := ""
								i_part_size_best := 0
								s_part_size_best_ant := ""
								i_part_size_best_ant := 0

								for i := 0; i < 4; i++ {
									s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
									s_part_start_act = strings.Trim(s_part_start_act, "\x00")

									s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
									s_part_status_act = strings.Trim(s_part_status_act, "\x00")

									s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
									s_part_size_act = strings.Trim(s_part_size_act, "\x00")
									i_part_size_act, _ = strconv.Atoi(s_part_size_act)

									if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
										if i != num_particion {
											s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
											s_part_size_best = strings.Trim(s_part_size_best, "\x00")
											i_part_size_best, _ = strconv.Atoi(s_part_size_best)

											s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
											s_part_size_act = strings.Trim(s_part_size_act, "\x00")
											i_part_size_act, _ = strconv.Atoi(s_part_size_act)

											if i_part_size_best > i_part_size_act {
												best_index = i
												break
											}
										}
									}
								}

								copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

								if best_index == 0 {
									mbr_empty_byte := Struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
								} else {
									s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
									s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
									i_part_start_best_ant, _ := strconv.Atoi(s_part_start_best_ant)

									s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
									s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
									i_part_size_best_ant, _ := strconv.Atoi(s_part_size_best_ant)

									copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
								}

								copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

								mbr_byte := Struct_a_bytes(master_boot_record)

								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								s_part_start = string(master_boot_record.Mbr_partition[best_index].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								f.Seek(int64(i_part_start), io.SeekStart)

								extended_boot_record := EBR{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := Struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}

								Analizador.Salida_comando += "[SUCCES] La Particion extendida fue creada con exito!\\n"
							}
						} else {
							Analizador.Salida_comando += "[ERROR] El nombre de la particion ya existe\\n"
						}
					} else {
						Analizador.Salida_comando += "[ERROR] No hay suficiente espacio disponible\\n"
					}
				} else {
					Analizador.Salida_comando += "[ERROR] No se encontró un espacio adecuado para la partición\\n"
				}
			} else {
				Analizador.Salida_comando += "[ERROR] El disco ya tiene una partición extendida\\n"
			}
		} else {
			Analizador.Salida_comando += "[ERROR] El disco está vacío\\n"
		}

		f.Close()
	}
}

func Crear_particion_logica(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := Estructuras.MBR{}
	ebr_empty := Estructuras.EBR{}
	var empty [100]byte

	if fit != "" {
		aux_fit = fit
	} else {
		aux_fit = "w"
	}

	if unit != "" {
		aux_unit = unit

		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			size_bytes = size * 1024
		} else {
			size_bytes = size * 1024 * 1024
		}
	} else {
		size_bytes = size * 1024
	}

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err != nil {
		Analizador.Salida_comando += "[ERROR] No existe el disco duro con ese nombre...\\n"
	} else {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""
			num_extendida := -1

			for i := 0; i < 4; i++ {
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					num_extendida = i
					break
				}
			}

			if !Existe_particion(direccion, nombre) {
				if num_extendida != -1 {
					s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")
					i_part_start, _ := strconv.Atoi(s_part_start)

					cont := i_part_start

					f.Seek(int64(cont), io.SeekStart)

					ebr2 := Struct_a_bytes(ebr_empty)
					sstruct := len(ebr2)

					lectura := make([]byte, sstruct)
					f.Read(lectura)

					extended_boot_record := Bytes_a_struct_ebr(lectura)

					s_part_size_ext := string(extended_boot_record.Part_size[:])
					s_part_size_ext = strings.Trim(s_part_size_ext, "\x00")

					if s_part_size_ext == "0" {
						s_part_size := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
						s_part_size = strings.Trim(s_part_size, "\x00")
						i_part_size, _ := strconv.Atoi(s_part_size)

						Analizador.Salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(i_part_size) + " Bytes\\n"
						Analizador.Salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes\\n"

						if i_part_size < size_bytes {
							Analizador.Salida_comando += "[ERROR] La particion logica a crear excede el espacio disponible de la particion extendida...\\n"
						} else {
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_fit[:], aux_fit)

							pos_actual, _ := f.Seek(0, io.SeekCurrent)
							ebr_empty_byte := Struct_a_bytes(ebr_empty)

							copy(extended_boot_record.Part_start[:], strconv.Itoa(int(pos_actual)-len(ebr_empty_byte)))
							copy(extended_boot_record.Part_size[:], strconv.Itoa(size_bytes))
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], nombre)

							s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Seek(int64(i_part_start), io.SeekStart)
							f.Write(ebr_byte)

							Analizador.Salida_comando += "[SUCCES] La Particion logica fue creada con exito!\\n"
						}
					} else {
						s_part_size := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
						s_part_size = strings.Trim(s_part_size, "\x00")
						i_part_size, _ := strconv.Atoi(s_part_size)

						s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
						s_part_start = strings.Trim(s_part_start, "\x00")
						i_part_start, _ := strconv.Atoi(s_part_start)

						Analizador.Salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(i_part_size+i_part_start) + " Bytes\\n"
						Analizador.Salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes\\n"

						s_part_next := string(extended_boot_record.Part_next[:])
						s_part_next = strings.Trim(s_part_next, "\x00")
						i_part_next, _ := strconv.Atoi(s_part_next)

						pos_actual, _ := f.Seek(0, io.SeekCurrent)

						for (i_part_next != -1) && (int(pos_actual) < (i_part_size + i_part_start)) {
							f.Seek(int64(i_part_next), io.SeekStart)

							ebr2 := Struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)

							lectura := make([]byte, sstruct)
							f.Read(lectura)

							pos_actual, _ = f.Seek(0, io.SeekCurrent)

							extended_boot_record = Bytes_a_struct_ebr(lectura)

							if extended_boot_record.Part_next == empty {
								break
							}

							s_part_next = string(extended_boot_record.Part_next[:])
							s_part_next = strings.Trim(s_part_next, "\x00")
							i_part_next, _ = strconv.Atoi(s_part_next)
						}

						s_part_start_ext := string(extended_boot_record.Part_start[:])
						s_part_start_ext = strings.Trim(s_part_start_ext, "\x00")
						i_part_start_ext, _ := strconv.Atoi(s_part_start_ext)

						s_part_size_ext := string(extended_boot_record.Part_size[:])
						s_part_size_ext = strings.Trim(s_part_size_ext, "\x00")
						i_part_size_ext, _ := strconv.Atoi(s_part_size_ext)

						s_part_size_mbr := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
						s_part_size_mbr = strings.Trim(s_part_size_mbr, "\x00")
						i_part_size_mbr, _ := strconv.Atoi(s_part_size_mbr)

						s_part_start_mbr := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
						s_part_start_mbr = strings.Trim(s_part_start_mbr, "\x00")
						i_part_start_mbr, _ := strconv.Atoi(s_part_start_mbr)

						espacio_necesario := i_part_start_ext + i_part_size_ext + size_bytes

						if espacio_necesario <= (i_part_size_mbr + i_part_start_mbr) {
							copy(extended_boot_record.Part_next[:], strconv.Itoa(i_part_start_ext+i_part_size_ext))

							pos_actual, _ = f.Seek(0, io.SeekCurrent)
							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Seek(int64(int(pos_actual)-len(ebr_byte)), io.SeekStart)
							f.Write(ebr_byte)

							f.Seek(int64(i_part_start_ext+i_part_size_ext), io.SeekStart)
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_fit[:], aux_fit)
							pos_actual, _ = f.Seek(0, io.SeekCurrent)
							copy(extended_boot_record.Part_start[:], strconv.Itoa(int(pos_actual)))
							copy(extended_boot_record.Part_size[:], strconv.Itoa(size_bytes))
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], nombre)
							ebr_byte = struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)

							Analizador.Salida_comando += "[EXITO] La Particion logica fue creada con exito!\\n"
						} else {
							Analizador.Salida_comando += "[ERROR] La particion logica a crear excede el espacio disponible de la particion extendida...\\n"
						}
					}
				} else {
					Analizador.Salida_comando += "[ERROR] No se puede crear una particion logica si no hay una extendida...\\n"
				}
			} else {
				Analizador.Salida_comando += "[ERROR] Ya existe una particion con ese nombre...\\n"
			}
		} else {
			Analizador.Salida_comando += "[ERROR] el disco se encuentra vacio...\\n"
		}
		f.Close()
	}
}

func Existe_particion(direccion string, nombre string) bool {
	extendida := -1
	mbr_empty := Estructuras.MBR{}
	ebr_empty := Estructuras.EBR{}
	var empty [100]byte

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_name := ""
			s_part_type := ""

			for i := 0; i < 4; i++ {
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")

				if s_part_name == nombre {
					f.Close()
					return true
				}

				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					extendida = i
				}
			}

			if extendida != -1 {
				s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")
				i_part_start, _ := strconv.Atoi(s_part_start)

				s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_size[:])
				s_part_size = strings.Trim(s_part_size, "\x00")
				i_part_size, _ := strconv.Atoi(s_part_size)

				ebr2 := Struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				lectura := make([]byte, sstruct)
				n_leidos, _ := f.Read(lectura)

				f.Seek(int64(i_part_start), io.SeekStart)

				pos_actual, _ := f.Seek(0, io.SeekCurrent)

				for n_leidos != 0 && (pos_actual < int64(i_part_size+i_part_start)) {
					lectura := make([]byte, sstruct)
					n_leidos, _ = f.Read(lectura)

					pos_actual, _ = f.Seek(0, io.SeekCurrent)

					extended_boot_record := Bytes_a_struct_ebr(lectura)

					if extended_boot_record.Part_size == empty {
						break
					} else {
						s_part_name = string(extended_boot_record.Part_name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")

						if s_part_name == nombre {
							f.Close()
							return true
						}

						s_part_next := string(extended_boot_record.Part_next[:])
						s_part_next = strings.Trim(s_part_next, "\x00")

						if s_part_next != "-1" {
							f.Close()
							return false
						}
					}
				}
			}
		} else {
			Analizador.Salida_comando += "[ERROR] el disco se encuentra vacio...\\n"
		}
	}

	f.Close()
	return false
}

func Buscar_particion_p_e(direccion string, nombre string) int {

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)
	if err == nil {
		mbr_empty := Estructuras.MBR{}
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)
		master_boot_record := Bytes_a_struct_mbr(lectura)

		s_part_status := ""
		s_part_name := ""
		for i := 0; i < 4; i++ {
			s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
			s_part_status = strings.Trim(s_part_status, "\x00")

			if s_part_status != "1" {
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")
				if s_part_name == nombre {
					return i
				}
			}

		}
	}

	return -1
}

func Buscar_particion_l(direccion string, nombre string) int {
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		extendida := -1
		mbr_empty := Estructuras.MBR{}
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)
		master_boot_record := Bytes_a_struct_mbr(lectura)

		s_part_type := ""
		for i := 0; i < 4; i++ {
			s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
			s_part_type = strings.Trim(s_part_type, "\x00")

			if s_part_type != "e" {
				extendida = i
				break
			}
		}
		if extendida != -1 {
			ebr_empty := Estructuras.EBR{}
			ebr2 := Struct_a_bytes(ebr_empty)
			sstruct := len(ebr2)
			lectura := make([]byte, sstruct)

			s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
			s_part_start = strings.Trim(s_part_start, "\x00")
			i_part_start, _ := strconv.Atoi(s_part_start)

			f.Seek(int64(i_part_start), io.SeekStart)

			n_leidos, _ := f.Read(lectura)
			pos_actual, _ := f.Seek(0, io.SeekCurrent)

			s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
			s_part_size = strings.Trim(s_part_size, "\x00")
			i_part_size, _ := strconv.Atoi(s_part_size)

			for (n_leidos != 0) && (pos_actual < int64(i_part_start+i_part_size)) {
				n_leidos, _ = f.Read(lectura)
				pos_actual, _ = f.Seek(0, io.SeekCurrent)
				extended_boot_record := Bytes_a_struct_ebr(lectura)

				s_part_name_ext := string(extended_boot_record.Part_name[:])
				s_part_name_ext = strings.Trim(s_part_name_ext, "\x00")

				ebr_empty_byte := Struct_a_bytes(ebr_empty)

				if s_part_name_ext == nombre {
					return int(pos_actual) - len(ebr_empty_byte)
				}
			}
		}
		f.Close()
	}

	return -1
}

func Struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil && err != io.EOF {
		Analizador.Salida_comando += "[ERROR] Al codificar de struct a bytes \n"
	}
	return buf.Bytes()
}

func Bytes_a_struct_mbr(s []byte) Estructuras.MBR {
	p := Estructuras.MBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		Analizador.Salida_comando += "[ERROR] Al decodificar a MBR\\n"
	}

	return p
}

func Bytes_a_struct_ebr(s []byte) Estructuras.EBR {
	p := Estructuras.EBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		Analizador.Salida_comando += "[ERROR] AL decodificar a EBR\n"
	}

	return p
}
