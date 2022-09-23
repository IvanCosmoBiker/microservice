package factory

import (
	//ofdFerma "ephorservices/ephorsale/fiscal/fermaOfd"
	"ephorservices/ephorsale/fiscal/interfaceFiscal"
	//ofdNanokass "ephorservices/ephorsale/fiscal/nanokass"
	//ofdOrange "ephorservices/ephorsale/fiscal/orange"
)

// instance of type fiscal
//var ofd ofdFerma.NewFermaOfdStruct
//var Orange ofdOrange.NewOrangeStruct
//var Nanokass ofdNanokass.NewNanokassStruct

func GetFiscal(typeFiscal uint8) interfaceFiscal.Fiscal {
	switch typeFiscal {
	case interfaceFiscal.Fr_EphorOrangeData,
		interfaceFiscal.Fr_ServerOrangeData,
		interfaceFiscal.Fr_EphorServerOrangeData,
		interfaceFiscal.Fr_OrangeData:
		return nil
		//return Orange.NewFiscal()
	case interfaceFiscal.Fr_NanoKassa,
		interfaceFiscal.Fr_ServerNanoKassa:
		return nil
		//return Nanokass.NewFiscal()
	case interfaceFiscal.Fr_OFD:
		return nil
		//return ofd.NewFiscal()
	}
	return nil
}
