package models

import (
	"time"

	"gorm.io/gorm"
)

type Paciente struct {
	gorm.Model
	IdPersona          int       `gorm:"unique;type:int(11)"`
	Paterno            string    `gorm:"not null;type:varchar(50)"`
	Materno            string    `gorm:"not null;type:varchar(50)"`
	Nombre             string    `gorm:"not null;type:varchar(100)"`
	Nss                string    `gorm:"not null;type:varchar(10)"`
	AgregadoMedico     string    `gorm:"not null;type:varchar(8)"`
	Curp               string    `gorm:"not null;type:varchar(18)"`
	FechaNacimiento    time.Time `gorm:"not null;column:fecha_nacimiento;type:date"`
	Edad               int       `gorm:"not null;type:int(3)"`
	Sexo               string    `gorm:"not null;type:varchar(1)"`
	DhUMF              int       `gorm:"type:int(3)"`
	Consultorio        int       `gorm:"type:int(2)"`
	Turno              string    `gorm:"type:varchar(1)"`
	Idee               string    `gorm:"type:varchar(18)"`
	AgregadoAfiliacion int       `gorm:"type:int(10)"`
	ConDerechoSm       string    `gorm:"type:varchar(2)"`
	ConDerechoInc      string    `gorm:"type:varchar(2)"`
	ClavePresupuestal  string    `gorm:"type:varchar(13)"`
	DhDeleg            int       `gorm:"type:int(3)"`
	RegistroPatronal   string    `gorm:"type:varchar(10)"`
	Direccion          string    `gorm:"type:varchar(255)"`
	Colonia            string    `gorm:"type:varchar(100)"`
	Telefono           string    `gorm:"type:varchar(10)"`
	TipoPension        string    `gorm:"type:varchar(50)"`
	VigenteHasta       time.Time `gorm:"column:vigencia_hasta;type:date"`
	FechaAcceDer       time.Time `gorm:"column:fecha_acceder;type:date"`
	Celular            string    `gorm:"type:varchar(20);unique" validate:"numeric"`
	Email              string    `gorm:"type:varchar(150);unique" validate:"email"`
}

func (Paciente) TableName() string {
	return "tbl_pacientes"
}
