package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Contrato struct {
	gorm.Model
	Contrato      string    `gorm:"index;not null;type:varchar(50)"`
	TipoContrato  string    `gorm:"not null;type:varchar(2)"`
	CuentaID      uint      `gorm:"not null"`               // Clave externa para la cuenta, no puede ser nulo
	Cuenta        Cuenta    `gorm:"foreignKey:CuentaID"`    // Define la relación con Cuenta
	ProveedorID   string    `gorm:"not null"`               // Clave externa para Proveedor, no puede ser nulo
	Proveedor     Proveedor `gorm:"foreignKey:ProveedorID"` // Define la relación con Proveedor
	FechaInicio   time.Time `gorm:"not null;column:fecha_inicio;type:date"`
	FechaTermino  time.Time `gorm:"not null;column:fecha_termino;type:date"`
	ImporteMinimo float64   `gorm:"type:decimal(11,2)"`
	ImporteMaximo float64   `gorm:"type:decimal(11,2)"`
	Dictamen      string    `gorm:"type:varchar(50)"`
	IVA           int       `gorm:"not null;type:int(1);default:1"`
	RetIVA        int       `gorm:"not null;type:int(1);default:0"`
	RetIVA2       int       `gorm:"not null;type:int(1);default:0"`
	RetISR        int       `gorm:"not null;type:int(1);default:0"`
	IVAID         string    `gorm:"not null;column:tasa_iva"` // Clave externa para IVA, no puede ser nulo
	TasaIVA       IVA       `gorm:"foreignKey:IVAID"`         // Define la relación con IVA
	RetID         string    `gorm:"not null;column:tasa_ret"` // Clave externa para RetIVA, no puede ser nulo
	TasaRet       IVARet    `gorm:"foreignKey:RetID"`         // Define la relación con RetIVA
	ISRID         string    `gorm:"not null;column:tasa_isr"` // Clave externa para ISR, no puede ser nulo
	TasaISR       ISR       `gorm:"foreignKey:ISRID"`         // Define la relación con ISR
	UserID        string    `gorm:"not null"`                 // Clave externa para Usuario, no puede ser nulo
	User          User      `gorm:"foreignKey:UserID"`        // Define la relación con User
	Activo        bool      `gorm:"not null;default:true"`
}

type ContratoUI struct {
	ContratoID uint            `gorm:"index;not null"` // Clave foránea para el contrato
	Contrato   Contrato        `gorm:"foreignKey:ContratoID"`
	UIID       string          `gorm:"index;not null"` // Clave foránea para la unidad de información
	UI         UnidadInf       `gorm:"foreignKey:UIID"`
	Importe    decimal.Decimal `gorm:"type:decimal(11,2)"`
	Activo     bool            `gorm:"not null;default:true"`
}

func (Contrato) TableName() string {
	return "tbl_contratos"
}

func (ContratoUI) TableName() string {
	return "tbl_contratos_ui"
}
