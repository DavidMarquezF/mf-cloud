package modules

type ModuleId int

// Must match mf_device.h
const (
	Ultrasound  ModuleId = 1
	Temperature ModuleId = 2
)

type Platform int

const (
	ESP32 Platform = iota
)

type Module struct {
	Id ModuleId `json:"id"`
}

type FirmwareConfig struct {
	DeviceId string   `json:"device_id"`
	Platform Platform `json:"platform"`
	Modules  []Module
}