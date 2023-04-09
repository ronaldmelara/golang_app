package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	//"fmt"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD()(conexion *sql.DB){
	Driver:="mysql"
	Usuario:="root"
	Contrasenia:="123456789"
	NombreBD:="sistema"

	conexion,err:= sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+NombreBD)

	if(err != nil){
		panic(err.Error())
	}

	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main(){
	http.HandleFunc("/", Home)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)

	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo.....")
	http.ListenAndServe(":8081", nil)
}

func Home(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hola Develoteca")
	
	arrEmpleados:=Consultar()
	fmt.Println(arrEmpleados)

	plantillas.ExecuteTemplate(w, "inicio", arrEmpleados)
}

func Crear(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hola Develoteca")
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request){

	if r.Method=="POST"{
		nombre:= r.FormValue("nombre")
		correo:= r.FormValue("correo")

		conexionEstablecida:= conexionBD()
		insertarRegistro,err:=conexionEstablecida.Prepare("INSERT INTO sistema.empleados(nombre, email)VALUES(?,?)")
		if(err != nil){
			panic(err.Error())
		}else{
			insertarRegistro.Exec(nombre, correo)
		}

		http.Redirect(w,r,"/",301)

	}
}

type Empleado struct{
	Id int
	Nombre string
	Correo string
}

func Consultar()[]Empleado{
	conexionEstablecida:= conexionBD()
	registros,err:=conexionEstablecida.Query("SELECT * FROM empleados")
	if(err != nil){
		panic(err.Error())
	}

	empleado:= Empleado{}
	arregloEmpleado:= []Empleado{}

	for registros.Next(){
		var id int
		var nombre, email string

		err=registros.Scan(&id,&nombre,&email)

		if(err!=nil){
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = email

		arregloEmpleado = append(arregloEmpleado,empleado)
	}

	return arregloEmpleado
}


func Borrar(w http.ResponseWriter, r *http.Request){

	
		idEmpleado:= r.URL.Query().Get("id")
		fmt.Println(idEmpleado)

		conexionEstablecida:= conexionBD()
		eliminarRegistro,err:=conexionEstablecida.Prepare("DELETE FROM sistema.empleados WHERE id=?")
		if(err != nil){
			panic(err.Error())
		}else{
			eliminarRegistro.Exec(idEmpleado)
		}

		http.Redirect(w,r,"/",301)

	
}

func Editar(w http.ResponseWriter, r *http.Request){

	
	idEmpleado:= r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	empleado:=ConsultarPorId(idEmpleado)
	fmt.Println(empleado)

	plantillas.ExecuteTemplate(w, "editar", empleado)


}


func ConsultarPorId(id string)Empleado{
	conexionEstablecida:= conexionBD()
	registros,err:=conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", id)
	
	if(err != nil){
		panic(err.Error())
	}

	empleado:= Empleado{}
	

	for registros.Next(){
		var id int
		var nombre, email string

		err=registros.Scan(&id,&nombre,&email)

		if(err!=nil){
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = email

		
	}

	return empleado
}


func Actualizar(w http.ResponseWriter, r *http.Request){

	if r.Method=="POST"{
		nombre:= r.FormValue("nombre")
		correo:= r.FormValue("correo")
		id:= r.FormValue("id")

		conexionEstablecida:= conexionBD()
		insertarRegistro,err:=conexionEstablecida.Prepare("UPDATE sistema.empleados SET nombre=?,  email=? WHERE id=?")
		if(err != nil){
			panic(err.Error())
		}else{
			insertarRegistro.Exec(nombre, correo, id)
		}

		http.Redirect(w,r,"/",301)

	}
}