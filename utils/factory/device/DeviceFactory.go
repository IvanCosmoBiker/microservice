package device

import (
    interfaceDevice "interface/deviceinterface"
    coolerDevice "device/cooler"
    automatDevice "device/automat"
)
// instance of type device
var cooler coolerDevice.NewCoolerStruct
var automat automatDevice.NewAutomatStruct

func GetDevice(device int) (interfaceDevice.Device) {
    switch device {
        case interfaceDevice.TypeCoffee,
        interfaceDevice.TypeSnack,
        interfaceDevice.TypeHoreca,
        interfaceDevice.TypeSodaWater,
        interfaceDevice.TypeMechanical,
        interfaceDevice.TypeComb:
        return automat.NewDevice()
        fallthrough
        case interfaceDevice.TypeCooler:
            return cooler.NewDevice()
    }
    return nil
}
