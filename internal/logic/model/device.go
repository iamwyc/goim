package model

// DeviceRegisterInDTO device register dto
type DeviceRegisterInDTO struct {
	Key      string
	Sn       string `from:"sn" binding:"required"`
	Platform int32  `from:"sn" binding:"required"`
	Serias   int32  `from:"sn" binding:"required"`
}
