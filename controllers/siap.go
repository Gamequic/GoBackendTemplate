package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	Matricula      string `json:"matricula"`
	Paterno        string `json:"paterno"`
	Materno        string `json:"materno"`
	Nombre         string `json:"nombre"`
	NombreCompleto string `json:"nombre_completo"`
	Puesto         string `json:"puesto"`
	Descripcion    string `json:"descripcion"`
	Departamento   string `json:"departamento"`
	Descripcion1   string `json:"descripcion1"`
}

type SiapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		WsObtenEmpResponse struct {
			XMLName xml.Name `xml:"wsObtenEmpResponse"`
			Return  string   `xml:"wsObtenEmpResult"`
		} `xml:"wsObtenEmpResponse"`
	} `xml:"Body"`
}

func Siap(c *gin.Context) {
	// Obtener los parámetros de la ruta
	strMatricula := c.Param("strMatricula")
	strDele := c.Param("strDele")

	// Llamar a la función Siap con los parámetros obtenidos
	respuesta, err := SiapFunc(strMatricula, strDele)
	if err != nil {
		// Manejar el error, por ejemplo, devolver un mensaje de error al cliente
		c.String(http.StatusInternalServerError, "Error al procesar la solicitud: %s", err.Error())
		return
	}

	// Parsear la respuesta XML
	var siapResp SiapResponse
	if err := xml.Unmarshal([]byte(respuesta), &siapResp); err != nil {
		c.String(http.StatusInternalServerError, "Error al parsear la respuesta XML: %s", err.Error())
		return
	}

	// Explode de la cadena obtenida
	arrCadena := strings.Split(siapResp.Body.WsObtenEmpResponse.Return, "|")
	if len(arrCadena) < 6 {
		c.String(http.StatusInternalServerError, "Error: Respuesta inválida")
		return
	}

	// Crear el struct Employee y asignar los valores
	emp := Employee{
		Matricula:      arrCadena[0],
		Paterno:        strings.Split(arrCadena[1], "/")[0],
		Materno:        strings.Split(arrCadena[1], "/")[1],
		Nombre:         strings.Split(arrCadena[1], "/")[2],
		NombreCompleto: strings.Join(strings.Split(arrCadena[1], "/"), " "),
		Puesto:         arrCadena[4],
		Descripcion:    arrCadena[5],
		Departamento:   arrCadena[2],
		Descripcion1:   arrCadena[3],
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(emp)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al convertir a JSON: %s", err.Error())
		return
	}

	// Devolver la respuesta en JSON
	c.Data(http.StatusOK, "application/json", jsonData)
}

func SiapFunc(strMatricula, strDele string) (string, error) {
	// Crear el cuerpo del mensaje SOAP
	soapBody := fmt.Sprintf(`<tns:wsObtenEmp xmlns:tns="http://tempuri.org/">
        <tns:strMatricula>%s</tns:strMatricula>
        <tns:strDele>%s</tns:strDele>
    </tns:wsObtenEmp>`, strMatricula, strDele)

	// Crear la solicitud SOAP completa
	soapRequest := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://tempuri.org/">
        <soapenv:Header/>
        <soapenv:Body>%s</soapenv:Body>
    </soapenv:Envelope>`, soapBody)

	// Configurar la solicitud HTTP
	req, err := http.NewRequest("POST", "http://172.26.18.157/biometrico/Biometrico/WebServices/wsBiometrico.asmx", strings.NewReader(soapRequest))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	// Realizar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Leer y devolver la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
