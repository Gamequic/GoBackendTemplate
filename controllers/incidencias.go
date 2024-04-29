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

// IncidenciaResponse es la estructura de datos para parsear la respuesta XML de wsObtenIncidencia
type IncidenciaResponse struct {
	IncidenciaResult string `xml:"Body>wsObtenIncidenciaResponse>wsObtenIncidenciaResult"`
}

// Incidencia representa la información de una incidencia
type Incidencia struct {
	Matricula string `json:"matricula"`
	Periodo   string `json:"periodo"`
}

// ObtenerIncidencia obtiene información sobre una incidencia
func ObtenerIncidencia(c *gin.Context) {
	// Obtener los parámetros de la ruta
	matricula := c.Param("matricula")
	periodo := c.Param("periodo")

	// Llamar a la función ObtenerIncidenciaFunc con los parámetros obtenidos
	respuesta, err := ObtenerIncidenciaFunc(matricula, periodo)
	if err != nil {
		// Manejar el error, por ejemplo, devolver un mensaje de error al cliente
		c.String(http.StatusInternalServerError, "Error al procesar la solicitud: %s", err.Error())
		return
	}

	// Parsear la respuesta SOAP
	var incidenciaResponse IncidenciaResponse
	if err := xml.Unmarshal([]byte(respuesta), &incidenciaResponse); err != nil {
		c.String(http.StatusInternalServerError, "Error al parsear la respuesta XML: %s", err.Error())
		return
	}

	// Crear la estructura Incidencia y asignar los valores
	incidencia := Incidencia{
		Matricula: matricula,
		Periodo:   periodo,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(incidencia)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al convertir a JSON: %s", err.Error())
		return
	}

	// Devolver la respuesta en JSON
	c.Data(http.StatusOK, "application/json", jsonData)
}

// ObtenerIncidenciaFunc obtiene información sobre una incidencia llamando al servicio web wsObtenIncidencia
func ObtenerIncidenciaFunc(matricula, periodo string) (string, error) {
	// Crear el cuerpo del mensaje SOAP
	soapBody := fmt.Sprintf(`<tns:wsObtenIncidencia xmlns:tns="http://tempuri.org/">
        <tns:strMatricula>%s</tns:strMatricula>
        <tns:strPeriodo>%s</tns:strPeriodo>
    </tns:wsObtenIncidencia>`, matricula, periodo)

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
