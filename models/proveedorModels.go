package models

import (
	"time"

	"gorm.io/gorm"
)

type ProveedorSucursal struct {
	gorm.Model
	ProveedorID string    `gorm:"not null"`               // Clave externa para Proveedor, no puede ser nulo
	Proveedor   Proveedor `gorm:"foreignKey:ProveedorID"` // Define la relación con Proveedor
	Calle       string    `gorm:"not null;type:varchar(255)"`
	NumExt      string    `gorm:"type:varchar(10)"`
	NumInt      string    `gorm:"type:varchar(10)"`
	Colonia     string    `gorm:"not null;type:varchar(100)"`
	Localidad   string    `gorm:"type:varchar(100)"`
	Municipio   string    `gorm:"not null;type:varchar(100)"`
	Estado      string    `gorm:"not null;type:varchar(100)"`
	CodPostal   string    `gorm:"not null;type:varchar(5)"`
	Pais        string    `gorm:"not null;type:varchar(100)"`
	Telefono    string    `gorm:"type:varchar(10)"`
	Email       string    `gorm:"type:varchar(70)"`
}

type Proveedor struct {
	ID          uint `gorm:"primaryKey;type:int(10) UNSIGNED ZEROFILL;autoIncrement:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Proveedor   string         `gorm:"type:varchar(255)"`
	RFC         string         `gorm:"type:varchar(13)"`
	RegPatronal string         `gorm:"type:varchar(11)"`
	Calle       string         `gorm:"not null;type:varchar(255)"`
	NumExt      string         `gorm:"type:varchar(10)"`
	NumInt      string         `gorm:"type:varchar(10)"`
	Colonia     string         `gorm:"not null;type:varchar(100)"`
	Localidad   string         `gorm:"type:varchar(100)"`
	Municipio   string         `gorm:"not null;type:varchar(100)"`
	Estado      string         `gorm:"not null;type:varchar(100)"`
	CodPostal   string         `gorm:"not null;type:varchar(5)"`
	Pais        string         `gorm:"not null;type:varchar(100)"`
	Telefono    string         `gorm:"type:varchar(10)"`
	Email       string         `gorm:"type:varchar(70)"`
}

func (ProveedorSucursal) TableName() string {
	return "tbl_proveedores_sucursales"
}

func (Proveedor) TableName() string {
	return "tbl_proveedores"
}
