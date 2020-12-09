package main

import (
	"fmt"
	"net/http"
	"container/list"
	"io/ioutil"
	"strconv"
)

const rutaServidor = ":9000"
const metodoPost = "POST"
const metodoGet = "GET"

type Alumno struct {
	Nombre       string
	Calificacion float64
}

type Mensaje struct {
	Alumno       string
	Materia      string
	Calificacion float64
}

type Materia struct {
	Nombre  string
	Alumnos list.List
}

type AdminMaterias struct {
	Materias list.List
}
var admin AdminMaterias

func (admin *AdminMaterias) existeMateria(nombreMateria string) bool {
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.Nombre == nombreMateria {
			return true
		}
	}
	return false
}

func (admin *AdminMaterias) obtenerMateria(nombreMateria string) *Materia {
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.Nombre == nombreMateria {
			return m
		}
	}
	return nil
}

func (admin *AdminMaterias) existeAlumnoMateria(nombreAlumno string, nombreMateria string) bool {
	m := admin.obtenerMateria(nombreMateria)
	return m.existeAlumno(nombreAlumno)
}

func (admin *AdminMaterias) existeAlumno(nombreAlumno string) bool {
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.existeAlumno(nombreAlumno) {
			return true
		}
	}
	return false
}

func (admin *AdminMaterias) AgregarCalificacionAlumno(mensaje Mensaje) string {
	var respuesta string
	nombreMateria := mensaje.Materia
	nombreAlumno := mensaje.Alumno
	calificacion := mensaje.Calificacion

	if len(nombreMateria) > 0 && len(nombreAlumno) > 0 {
		if admin.existeMateria(nombreMateria) {
			if admin.existeAlumnoMateria(nombreAlumno, nombreMateria) {
				return "El alumno ya tiene registrada una calificación en esta materia: " + nombreMateria
			}
			admin.agregarAlumno(nombreAlumno, nombreMateria, calificacion)
			respuesta = "Ha sido agregado a la materia, un nuevo alumno" + " (" + nombreAlumno + ")"
		} else {
			admin.agregarMateria(nombreMateria)
			admin.agregarAlumno(nombreAlumno, nombreMateria, calificacion)
			respuesta = "Una nueva materia ha sido creada" + "<br>" + "Ha sido agregado a la materia, un nuevo alumno"
		}
		return respuesta + "<br> El alumno: " + nombreAlumno + ", en la materia: " + nombreMateria + ", calificación= " + fmt.Sprintf("%f", calificacion)
	}
	return "Para que sea posible registrar la calificación del alumno, ingrese todos los datos"
}

func (admin *AdminMaterias) ObtenerPromedioAlumno(nombreAlumno string) string {
	if len(nombreAlumno) > 0 {
		if admin.existeAlumno(nombreAlumno) {
			return admin.obtenerCalificacionesAlumno(nombreAlumno) + "<br>Promedio = " + fmt.Sprintf("%f", admin.obtenerPromedioAlumno(nombreAlumno))
		}
		return "No existe el alumno: " + nombreAlumno
	}
	return "Ingrese nombre del alumno"
}

func (admin *AdminMaterias) ObtenerPromedioGeneral() string {
	if admin.Materias.Len() > 0 {
		return admin.obtenerCalificacionesMaterias() + "<br>Promedio = " + fmt.Sprintf("%f", admin.obtenerPromedioGeneral())
	}
	return "No existen materias registradas"
}

func (admin *AdminMaterias) ObtenerPromedioMateria(nombreMateria string) string {
	if len(nombreMateria) > 0 {
		if admin.existeMateria(nombreMateria) {
			return admin.obtenerCalificacionesAlumnosMateria(nombreMateria) + "<br>Promedio = " + fmt.Sprintf("%f", admin.obtenerPromedioMateria(nombreMateria))
		}
		return "No existe la materia: " + nombreMateria
	}
	return "Ingrese nombre de la materia"
}

func (admin *AdminMaterias) agregarAlumno(nombreAlumno string, nombreMateria string, calificacion float64) {
	m := admin.obtenerMateria(nombreMateria)
	m.agregarAlumno(nombreAlumno, calificacion)
}

func (admin *AdminMaterias) agregarMateria(nombreMateria string) {
	m := new(Materia)
	m.Nombre = nombreMateria
	admin.Materias.PushBack(m)
}

func (admin *AdminMaterias) obtenerPromedioAlumno(nombreAlumno string) float64 {
	cantMaterias := 0.0
	sumaCalificaciones := 0.0
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.existeAlumno(nombreAlumno) {
			cantMaterias++
			sumaCalificaciones += m.obtenerAlumno(nombreAlumno).Calificacion
		}
	}
	return sumaCalificaciones / cantMaterias
}

func (admin *AdminMaterias) obtenerCalificacionesAlumno(nombreAlumno string) string {
	calificaciones := ""
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.existeAlumno(nombreAlumno) {
			calificaciones += "<tr> <th>" + m.Nombre + "</th> <th>" + fmt.Sprintf("%f", m.obtenerAlumno(nombreAlumno).Calificacion) + "</th> </tr>"
		}
	}
	return "Alumno: " + nombreAlumno + "<br> <br> <table> <tr> <th>Materia</th> <th>Calificación</th> </tr> " + calificaciones + "</table>"
}

func (admin *AdminMaterias) obtenerCalificacionesMaterias() string {
	calificaciones := ""
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		calificaciones += "<tr> <th>" + m.Nombre + "</th> <th>" + fmt.Sprintf("%f", m.obtenerPromedio()) + "</th> </tr>"
	}
	return "<table> <tr> <th>Materia</th> <th>Promedio</th> </tr> " + calificaciones + "</table>"
}


func (m *Materia) obtenerPromedio() float64 {
	contAlumnos := 0.0
	sumaCalificaciones := 0.0
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		contAlumnos++
		sumaCalificaciones += a.Calificacion
	}
	return sumaCalificaciones / contAlumnos
}

func (m *Materia) obtenerAlumno(nombreAlumno string) *Alumno {
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		if a.Nombre == nombreAlumno {
			return a
		}
	}
	return nil
}

func (admin *AdminMaterias) obtenerCalificacionesAlumnosMateria(nombreMateria string) string {
	m := admin.obtenerMateria(nombreMateria)
	return m.obtenerCalificacionesAlumnos()
}

func (admin *AdminMaterias) obtenerPromedioMateria(nombreMateria string) float64 {
	m := admin.obtenerMateria(nombreMateria)
	return m.obtenerPromedio()
}

func (admin *AdminMaterias) obtenerPromedioGeneral() float64 {
	cantMaterias := 0.0
	sumaPromediosMaterias := 0.0
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		cantMaterias++
		sumaPromediosMaterias += m.obtenerPromedio()
	}
	return sumaPromediosMaterias / cantMaterias
}

func (m *Materia) existeAlumno(nombreAlumno string) bool {
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		if a.Nombre == nombreAlumno {
			return true
		}
	}
	return false
}

func (m *Materia) agregarAlumno(nombreAlumno string, calificacion float64) {
	a := new(Alumno)
	a.Calificacion = calificacion
	a.Nombre = nombreAlumno
	m.Alumnos.PushBack(a)
}

func (m *Materia) obtenerCalificacionesAlumnos() string {
	calificaciones := ""
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		calificaciones += "<tr> <th>" + a.Nombre + "</th> <th>" + fmt.Sprintf("%f", m.obtenerAlumno(a.Nombre).Calificacion) + "</th> </tr>"
	}
	return "Materia: " + m.Nombre + "<br> <br> <table> <tr> <th>Alumno</th> <th>Calificación</th> </tr> " + calificaciones + "</table>"
}

func cargarHTML(a string) string {
	html, _ := ioutil.ReadFile(a)
	return string(html)
}

func index(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHTML("index.html"),
	)
}

func promedioMateria(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case metodoPost:
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("respuesta.html"),
			admin.ObtenerPromedioMateria(req.FormValue("nombreMateria")),
		)
	case metodoGet:
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("promedioMateria.html"),
		)
	}
}

func promedioGeneral(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHTML("respuesta.html"),
		admin.ObtenerPromedioGeneral(),
	)
}

func promedioAlumno(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case metodoPost:
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("respuesta.html"),
			admin.ObtenerPromedioAlumno(req.FormValue("nombreAlumno")),
		)
	case metodoGet:
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("promedioAlumno.html"),
		)
	}
}

func alumno(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case metodoPost:
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		calificacion, _ := strconv.ParseFloat(req.FormValue("calificacion"), 64)
		mensaje := Mensaje{Alumno: req.FormValue("nombreAlumno"), Materia: req.FormValue("materia"), Calificacion: calificacion}
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("respuesta.html"),
			admin.AgregarCalificacionAlumno(mensaje),
		)
	case metodoGet:
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("alumno.html"),
		)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/alumno", alumno)
	http.HandleFunc("/promedioAlumno", promedioAlumno)
	http.HandleFunc("/promedioGeneral", promedioGeneral)
	http.HandleFunc("/promedioMateria", promedioMateria)
	fmt.Println("Ejecutando servidor...")
	http.ListenAndServe(rutaServidor, nil)
}
