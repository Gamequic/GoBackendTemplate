package initializers

import (
	"github.com/calleros/sich/models"
)

func SyncDatabase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.UnidadInf{},
		&models.UnidadOp{},
		&models.Ubicacion{},
		&models.Zona{},
		&models.CentroCostos{},
		&models.Division{},
		&models.SubDivision{},
		&models.Proveedor{},
		&models.ProveedorSucursal{},
		&models.Paciente{},
		&models.Cuenta{},
		&models.CuentaUI{},
		&models.Contrato{},
		&models.ContratoUI{},
		&models.Especialidad{},
		&models.Turno{},
		&models.Estatus{},
		&models.Siap{},
		&models.IVA{},
		&models.IVARet{},
		&models.ISR{},
	)
}
