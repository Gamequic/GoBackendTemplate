package models

import (
	"time"
)

type Especialidad struct {
	ID           uint   `gorm:"primaryKey;type:int(5);autoIncrement:false"`
	Especialidad string `gorm:"type:varchar(100)"`
	Activa       bool   `gorm:"not null;default:true"`
}

type Turno struct {
	ID    uint   `gorm:"primaryKey;type:int(1);autoIncrement:false"`
	Turno string `gorm:"type:varchar(10)"`
}

type Estatus struct {
	ID        uint   `gorm:"primaryKey;type:int(2)"`
	IDEstatus string `gorm:"primaryKey;type:varchar(2)"`
	Estatus   string `gorm:"type:varchar(30)"`
}

type Siap struct {
	Matricula      uint      `gorm:"primaryKey;type:int(10) UNSIGNED ZEROFILL;autoIncrement:false"`
	Paterno        string    `gorm:"type:varchar(50)"`
	Materno        string    `gorm:"type:varchar(50)"`
	Nombre         string    `gorm:"type:varchar(100)"`
	NombreCompleto string    `gorm:"type:varchar(200)"`
	Puesto         int       `gorm:"type:int(11)"`
	Descripcion    string    `gorm:"type:varchar(100)"`
	Departamento   string    `gorm:"type:varchar(50)"`
	Descripcion1   string    `gorm:"type:varchar(100)"`
	FechaSiap      time.Time `gorm:"column:fecha_siap;type:date"`
}

type IVA struct {
	TasaIVA float64 `gorm:"primaryKey;type:decimal(4,4)"`
}

type IVARet struct {
	TasaRet float64 `gorm:"primaryKey;type:decimal(4,4)"`
}

type ISR struct {
	TasaISR float64 `gorm:"primaryKey;type:decimal(4,4)"`
}

func (Especialidad) TableName() string {
	return "tbl_especialidades"
}

func (Turno) TableName() string {
	return "tbl_turnos"
}

func (Estatus) TableName() string {
	return "tbl_estatus"
}

func (Siap) TableName() string {
	return "tbl_siap"
}

func (IVA) TableName() string {
	return "tbl_iva"
}

func (IVARet) TableName() string {
	return "tbl_iva_ret"
}

func (ISR) TableName() string {
	return "tbl_isr"
}
