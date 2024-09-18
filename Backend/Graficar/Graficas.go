package Graficar

import (
	"Backend/Analizador"
	"Backend/Auxiliares"
	"Backend/Estructuras"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Graficar_disk(direccion string, destino string) {
	mbr_empty := Estructuras.MBR{}
	var empty [100]byte
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)
	mbr2 := Auxiliares.Struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)
	lectura := make([]byte, sstruct)
	f.Seek(0, io.SeekStart)
	f.Read(lectura)
	master_boot_record := Auxiliares.Bytes_a_struct_mbr(lectura)

	if master_boot_record.Mbr_tamano != empty {
		if err == nil {
			Analizador.GraphDot += "digraph G{\n\n"
			Analizador.GraphDot += "  tbl [\n    shape=box\n    label=<\n"
			Analizador.GraphDot += "     <table border='0' cellborder='2' width='600' height='150' color='dodgerblue1'>\n"
			Analizador.GraphDot += "     <tr>\n"
			Analizador.GraphDot += "     <td height='150' width='110'> MBR </td>\n"
			s_mbr_tamano := string(master_boot_record.Mbr_tamano[:])
			s_mbr_tamano = strings.Trim(s_mbr_tamano, "\x00")
			i_mbr_tamano, _ := strconv.Atoi(s_mbr_tamano)
			total := i_mbr_tamano
			var espacioUsado float64
			espacioUsado = 0

			for i := 0; i < 4; i++ {
				s_part_s := string(master_boot_record.Mbr_partition[i].Part_size[:])
				s_part_s = strings.Trim(s_part_s, "\x00")
				i_part_s, _ := strconv.Atoi(s_part_s)
				parcial := i_part_s
				s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")

				if s_part_start != "-1" {
					var porcentaje_real float64
					porcentaje_real = (float64(parcial) * 100) / float64(total)
					var porcentaje_aux float64
					porcentaje_aux = (porcentaje_real * 500) / 100
					espacioUsado += porcentaje_real
					s_part_status := string(master_boot_record.Mbr_partition[i].Part_status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					if s_part_status != "1" {
						s_part_type := string(master_boot_record.Mbr_partition[i].Part_type[:])
						s_part_type = strings.Trim(s_part_type, "\x00")

						if s_part_type == "p" {
							Analizador.GraphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Primaria <br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"

							if i != 3 {
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s
								s_part_start = string(master_boot_record.Mbr_partition[i+1].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)

								p2 := i_part_start

								if s_part_start != "-1" {
									if (p2 - p1) != 0 {
										fragmentacion := p2 - p1
										porcentaje_real = float64(fragmentacion) * 100 / float64(total)
										porcentaje_aux = (porcentaje_real * 500) / 100

										Analizador.GraphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
									}
								}
							} else {
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								mbr_empty_byte := Auxiliares.Struct_a_bytes(mbr_empty)
								mbr_size := total + len(mbr_empty_byte)
								if (mbr_size - p1) != 0 {
									libre := (float64(mbr_size) - float64(p1)) + float64(len(mbr_empty_byte))
									porcentaje_real = (float64(libre) * 100) / float64(total)
									porcentaje_aux = (porcentaje_real * 500) / 100
									Analizador.GraphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
								}
							}
						} else {
							Analizador.GraphDot += "     <td  height='200' width='" + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + "'>\n     <table border='0'  height='200' WIDTH='" + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + "' cellborder='1'>\n"
							Analizador.GraphDot += "     <tr>  <td height='60' colspan='15'>Extendida</td>  </tr>\n     <tr>\n"
							s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)
							f.Seek(int64(i_part_start), io.SeekStart)
							ebr_empty := Estructuras.EBR{}
							ebr2 := Auxiliares.Struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)
							lectura := make([]byte, sstruct)
							f.Read(lectura)
							extended_boot_record := Auxiliares.Bytes_a_struct_ebr(lectura)
							s_part_size := string(extended_boot_record.Part_size[:])
							s_part_size = strings.Trim(s_part_size, "\x00")
							i_part_size, _ := strconv.Atoi(s_part_size)
							if i_part_size != 0 {
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)
								f.Seek(int64(i_part_start), io.SeekStart)
								band := true
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)
								s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)
								for band {
									ebr2 := Auxiliares.Struct_a_bytes(ebr_empty)
									sstruct := len(ebr2)
									lectura := make([]byte, sstruct)
									f.Seek(0, io.SeekStart)
									n, _ := f.Read(lectura)
									pos_actual, _ := f.Seek(0, io.SeekCurrent)
									if n != 0 && pos_actual < int64(i_part_start)+int64(i_part_s) {
										band = false
										break
									}
									s_part_s = string(extended_boot_record.Part_size[:])
									s_part_s = strings.Trim(s_part_s, "\x00")
									i_part_s, _ = strconv.Atoi(s_part_s)
									parcial = i_part_start
									porcentaje_real = float64(parcial) * 100 / float64(total)
									if porcentaje_real != 0 {
										s_part_status = string(extended_boot_record.Part_status[:])
										s_part_status = strings.Trim(s_part_status, "\x00")
										if s_part_status != "1" {
											Analizador.GraphDot += "     <td height='140'>EBR</td>\n"
											Analizador.GraphDot += "     <td height='140'>Logica<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
										} else {
											Analizador.GraphDot += "      <td height='150'>Libre 1 <br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
										}
										s_part_next := string(extended_boot_record.Part_next[:])
										s_part_next = strings.Trim(s_part_next, "\x00")
										i_part_next, _ := strconv.Atoi(s_part_next)

										if i_part_next == -1 {
											s_part_start := string(extended_boot_record.Part_start[:])
											s_part_start = strings.Trim(s_part_start, "\x00")
											i_part_start, _ := strconv.Atoi(s_part_start)
											s_part_size := string(extended_boot_record.Part_size[:])
											s_part_size = strings.Trim(s_part_size, "\x00")
											i_part_size, _ := strconv.Atoi(s_part_size)
											s_part_start_mbr := string(master_boot_record.Mbr_partition[i].Part_start[:])
											s_part_start_mbr = strings.Trim(s_part_start_mbr, "\x00")
											i_part_start_mbr, _ := strconv.Atoi(s_part_start_mbr)
											s_part_s_mbr := string(master_boot_record.Mbr_partition[i].Part_size[:])
											s_part_s_mbr = strings.Trim(s_part_s_mbr, "\x00")
											i_part_s_mbr, _ := strconv.Atoi(s_part_s_mbr)

											parcial = (i_part_start_mbr + i_part_s_mbr) - (i_part_size + i_part_start)
											porcentaje_real = (float64(parcial) * 100) / float64(total)

											if porcentaje_real != 0 {
												Analizador.GraphDot += "     <td height='150'>Libre 2<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
											}
											break

										} else {
											s_part_next := string(extended_boot_record.Part_next[:])
											s_part_next = strings.Trim(s_part_next, "\x00")
											i_part_next, _ := strconv.Atoi(s_part_next)

											f.Seek(int64(i_part_next), io.SeekStart)
										}
									}

								}
							} else {
								Analizador.GraphDot += "     <td height='140'> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
							}
							Analizador.GraphDot += "     </tr>\n     </table>\n     </td>\n"
							if i != 3 {
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)
								p1 := i_part_start + i_part_s
								s_part_start = string(master_boot_record.Mbr_partition[i+1].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)
								p2 := i_part_start
								if s_part_start != "-1" {
									if (p2 - p1) != 0 {
										fragmentacion := p2 - p1
										porcentaje_real = float64(fragmentacion) * 100 / float64(total)
										porcentaje_aux = (porcentaje_real * 500) / 100

										Analizador.GraphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
									}
								}
							} else {
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)
								p1 := i_part_start + i_part_s
								mbr_empty_byte := Auxiliares.Struct_a_bytes(mbr_empty)
								mbr_size := total + len(mbr_empty_byte)
								if (mbr_size - p1) != 0 {
									libre := (float64(mbr_size) - float64(p1)) + float64(len(mbr_empty_byte))
									porcentaje_real = (float64(libre) * 100) / float64(total)
									porcentaje_aux = porcentaje_real * 500 / 100
									Analizador.GraphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
								}
							}
						}
					} else {
						Analizador.GraphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
					}
				}
			}

			Analizador.GraphDot += "     </tr> \n     </table>        \n>];\n\n}"
			err := ioutil.WriteFile("reporte.dot", []byte(Analizador.GraphDot), 0644)
			Analizador.Salida_comando += "[EXITO] Reporte generado con exito!\\n"
			if err != nil {
				Analizador.Salida_comando += "[ERROR] Error al escribir en el archivo\\n"
				return
			}
			cmd := exec.Command("dot", "-Tpng", "reporte.dot", "-o", destino)
			if err := cmd.Run(); err != nil {
				Analizador.Salida_comando += "[ERROR] Error al generar la imagen\\n"
				return
			}
		} else {
			Analizador.Salida_comando += "[ERROR] El disco no fue encontrado...\\n"
		}
	} else {
		Analizador.Salida_comando += "[ERROR] Disco vacio...\\n"
	}
}
