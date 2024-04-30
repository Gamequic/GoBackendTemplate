package models

import (
	"gorm.io/gorm"
)

type Cuenta struct {
	gorm.Model
	PREI              int    `gorm:"unique;not null;type:int(8)"`
	CONAC             int    `gorm:"unique;not null;type:int(8)"`
	PartidaDelCog     int    `gorm:"type:int(6)"`
	DescripcionCuenta string `gorm:"type:varchar(255)"`
	ConceptoGasto     string `gorm:"type:varchar(2500)"`
	SAI               string `gorm:"type:varchar(4)"`
	Capitulo          string `gorm:"type:varchar(50)"`
	Subcapitulo       string `gorm:"type:varchar(50)"`
}

type CuentaUI struct {
	CuentaID uint      `gorm:"index;not null"` // Clave foránea para la cuenta
	Cuenta   Cuenta    `gorm:"foreignKey:CuentaID"`
	UIID     string    `gorm:"index;not null"` // Clave foránea para la unidad de información
	UI       UnidadInf `gorm:"foreignKey:UIID"`
}

func (Cuenta) TableName() string {
	return "tbl_cuentas"
}

func (CuentaUI) TableName() string {
	return "tbl_cuentas_ui"
}
