package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Define las estructuras de datos para representar la información de los titulares y beneficiarios
type Titular struct {
	//CodigoError        string `xml:"codigoError"`
	//MensajeError       string `xml:"mensajeError"`
	IdPersona          string `xml:"IdPersona"`
	Paciente           string `json:"Paciente,omitempty"`
	Nss                string `xml:"Nss"`
	AgregadoMedico     string `xml:"AgregadoMedico"`
	Curp               string `xml:"Curp"`
	FechaNacimiento    string `xml:"FechaNacimiento"`
	Edad               int    `json:"Edad,omitempty"`
	Sexo               string `xml:"Sexo"`
	DhUMF              string `xml:"DhUMF"`
	Consultorio        string `xml:"Consultorio"`
	Turno              string `xml:"Turno"`
	Paterno            string `xml:"Paterno"`
	Materno            string `xml:"Materno"`
	Nombre             string `xml:"Nombre"`
	Idee               string `xml:"Idee"`
	AgregadoAfiliacion string `xml:"AgregadoAfiliacion"`
	ConDerechoSm       string `xml:"ConDerechoSm"`
	ConDerechoInc      string `xml:"ConDerechoInc"`
	ClavePresupuestal  string `xml:"ClavePresupuestal"`
	DhDeleg            string `xml:"DhDeleg"`
	RegistroPatronal   string `xml:"RegistroPatronal"`
	Direccion          string `xml:"Direccion"`
	Colonia            string `xml:"Colonia"`
	Telefono           string `xml:"Telefono"`
	TipoPension        string `xml:"TipoPension"`
	VigenteHasta       string `xml:"VigenteHasta"`

	// Agrega más campos según la estructura XML
}

type Beneficiario struct {
	IdPersona          string `xml:"IdPersona"`
	Paciente           string `json:"Paciente,omitempty"`
	Nss                string `xml:"Nss"`
	AgregadoMedico     string `xml:"AgregadoMedico"`
	Curp               string `xml:"Curp"`
	FechaNacimiento    string `xml:"FechaNacimiento"`
	Edad               int    `json:"Edad,omitempty"`
	Sexo               string `xml:"Sexo"`
	DhUMF              string `xml:"DhUMF"`
	Consultorio        string `xml:"Consultorio"`
	Turno              string `xml:"Turno"`
	Paterno            string `xml:"Paterno"`
	Materno            string `xml:"Materno"`
	Nombre             string `xml:"Nombre"`
	Idee               string `xml:"Idee"`
	AgregadoAfiliacion string `xml:"AgregadoAfiliacion"`
	ConDerechoSm       string `xml:"ConDerechoSm"`
	ConDerechoInc      string `xml:"ConDerechoInc"`
	ClavePresupuestal  string `xml:"ClavePresupuestal"`
	DhDeleg            string `xml:"DhDeleg"`
	RegistroPatronal   string `xml:"RegistroPatronal"`
	Direccion          string `xml:"Direccion"`
	Colonia            string `xml:"Colonia"`
	Telefono           string `xml:"Telefono"`
	TipoPension        string `xml:"TipoPension"`
	VigenteHasta       string `xml:"VigenteHasta"`
	// Agrega más campos según la estructura XML
}

// Define las estructuras de datos para representar la información de los titulares y beneficiarios
type TitularesResponse struct {
	Titulares []Titular `xml:"Body>getInfoResponse>return"`
}

type BeneficiariosResponse struct {
	Beneficiarios []Beneficiario `xml:"Body>getInfoResponse>return>Beneficiarios"`
}

// Función para calcular la edad
func calcularEdad(fechaNacimiento string) int {
	// Parsear la fecha de nacimiento en formato "1980/02/14"
	nacimiento, err := time.Parse("2006/01/02", fechaNacimiento)
	if err != nil {
		log.Fatalf("Error al parsear la fecha de nacimiento: %v", err)
	}
	// Calcular la diferencia de años entre la fecha actual y la de nacimiento
	edad := time.Now().Year() - nacimiento.Year()
	// Ajustar la edad si el cumpleaños todavía no ha pasado este año
	if time.Now().Before(nacimiento.AddDate(edad, 0, 0)) {
		edad--
	}
	return edad
}

// Acceder es una función pública que puede ser accedida desde otros paquetes
func Acceder(c *gin.Context) {
	nssValue := c.Param("nss")
	soapBody := fmt.Sprintf(`<ns:getInfo xmlns:ns="%s"><nss>%s</nss></ns:getInfo>`, "http://vigenciaderechos.imss.gob.mx/", nssValue)
	soapRequest := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="%s">
    <soapenv:Header/>
    <soapenv:Body>%s</soapenv:Body>
    </soapenv:Envelope>`, "http://vigenciaderechos.imss.gob.mx/", soapBody)
	req, err := http.NewRequest("POST", "http://vigenciaderechos.imss.gob.mx/WSConsVigGpoFamComXNssService/WSConsVigGpoFamComXNss", strings.NewReader(soapRequest))
	if err != nil {
		log.Fatalf("Error creando la solicitud HTTP: %v", err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error haciendo la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Procesar la respuesta SOAP
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error leyendo la respuesta SOAP: %v", err)
	}

	// Parsear la respuesta XML
	var titularesResponse TitularesResponse
	if err := xml.Unmarshal(respBody, &titularesResponse); err != nil {
		log.Fatalf("Error parseando la respuesta XML de titulares: %v", err)
	}

	var beneficiariosResponse BeneficiariosResponse
	if err := xml.Unmarshal(respBody, &beneficiariosResponse); err != nil {
		log.Fatalf("Error parseando la respuesta XML de beneficiarios: %v", err)
	}

	// Separar pacientes en vigentes y no vigentes
	var pacientesVigentes, pacientesNoVigentes []interface{}

	// Convertir titulares a tipo interface{} y agregarlos al array de pacientes correspondiente
	for _, titular := range titularesResponse.Titulares {
		titular.Paciente = titular.Nombre + " " + titular.Paterno + " " + titular.Materno
		titular.Edad = calcularEdad(titular.FechaNacimiento)
		if titular.ConDerechoSm == "SI" {
			pacientesVigentes = append(pacientesVigentes, titular)
		} else {
			pacientesNoVigentes = append(pacientesNoVigentes, titular)
		}
	}

	// Convertir beneficiarios a tipo interface{} y agregarlos al array de pacientes correspondiente
	for _, beneficiario := range beneficiariosResponse.Beneficiarios {
		beneficiario.Paciente = beneficiario.Nombre + " " + beneficiario.Paterno + " " + beneficiario.Materno
		beneficiario.Edad = calcularEdad(beneficiario.FechaNacimiento)
		if beneficiario.ConDerechoSm == "SI" {
			pacientesVigentes = append(pacientesVigentes, beneficiario)
		} else {
			pacientesNoVigentes = append(pacientesNoVigentes, beneficiario)
		}
	}

	// Convertir estructuras a JSON
	pacientesVigentesJSON, err := json.Marshal(pacientesVigentes)
	if err != nil {
		log.Fatalf("Error convirtiendo pacientes vigentes a JSON: %v", err)
	}

	pacientesNoVigentesJSON, err := json.Marshal(pacientesNoVigentes)
	if err != nil {
		log.Fatalf("Error convirtiendo pacientes no vigentes a JSON: %v", err)
	}

	// Construir respuesta JSON final
	responseJSON := fmt.Sprintf(`{"vigentes": %s, "no_vigentes": %s}`, string(pacientesVigentesJSON), string(pacientesNoVigentesJSON))

	// Opcionalmente, enviar la respuesta al cliente
	c.Data(http.StatusOK, "application/json", []byte(responseJSON))
}
