package models

type UnidadInf struct {
	ID            uint      `gorm:"primarykey"`
	UI            string    `gorm:"index;not null;type:varchar(6)"`
	UnidadOpID    string    `gorm:"not null"`               // Clave externa para UnidadOp, no puede ser nulo
	UnidadOp      UnidadOp  `gorm:"foreignKey:UnidadOpID"`  // Define la relación con UnidadOp
	UbicacionID   string    `gorm:"not null"`               // Clave externa para Ubicacion, no puede ser nulo
	Ubicacion     Ubicacion `gorm:"foreignKey:UbicacionID"` // Define la relación con Ubicacion
	ZonaID        string    `gorm:"not null"`               // Clave externa para Zona, no puede ser nulo
	Zona          Zona      `gorm:"foreignKey:ZonaID"`      // Define la relación con Zona
	NombreCorto   string    `gorm:"not null;type:varchar(70)"`
	Nombre        string    `gorm:"not null;type:varchar(100)"`
	ClavePresup   string    `gorm:"not null;type:varchar(12)"`
	SqlSimf       string    `gorm:"not null;type:varchar(50)"`
	MatTitular    int16     `gorm:"type:int(11);default:0"`
	NombreTitular string    `gorm:"type:varchar(100)"`
	PuestoTitular string    `gorm:"type:varchar(100)"`
	Activo        bool      `gorm:"not null;default:true"`
}

type UnidadOp struct {
	ID         string `gorm:"primaryKey;type:varchar(5)"`
	Delegacion string `gorm:"not null;type:varchar(100)"`
}

type Ubicacion struct {
	ID           string `gorm:"primaryKey;type:varchar(9)"`
	Direccion    string `gorm:"type:varchar(100)"`
	Localidad    string `gorm:"type:varchar(100)"`
	Municipio    string `gorm:"type:varchar(100)"`
	Estado       string `gorm:"type:varchar(100)"`
	LatLng       string `gorm:"type:varchar(100)"`
	MarkerStatus int    `gorm:"type:int"`
	Imagen       string `gorm:"type:varchar(250)"`
}

type Zona struct {
	ID         string `gorm:"primaryKey;type:varchar(6)"`
	NombreZona string `gorm:"primaryKey;type:varchar(100)"`
	Ciudad     int    `gorm:"type:int"`
}

type CentroCostos struct {
	ID            string `gorm:"primaryKey;type:varchar(6)"`
	CentroCostos  string `gorm:"type:varchar(100)"`
	CentroCostos2 string `gorm:"type:varchar(100)"`
}

type Division struct {
	Dv       string `gorm:"primaryKey;type:varchar(3)"`
	Division string `gorm:"type:varchar(100)"`
}

type SubDivision struct {
	CC          string `gorm:"primaryKey;type:varchar(6)"`
	Dv          string `gorm:"primaryKey;type:varchar(3)"`
	Sd          string `gorm:"primaryKey;type:varchar(3)"`
	SubDivision string `gorm:"type:varchar(100)"`
}

func (UnidadInf) TableName() string {
	return "tbl_ui"
}

func (UnidadOp) TableName() string {
	return "tbl_uo"
}

func (Ubicacion) TableName() string {
	return "tbl_ubicacion"
}

func (Zona) TableName() string {
	return "tbl_zonas"
}

func (CentroCostos) TableName() string {
	return "tbl_cc"
}

func (Division) TableName() string {
	return "tbl_div"
}

func (SubDivision) TableName() string {
	return "tbl_sdiv"
}
