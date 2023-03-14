package factory

import (
	ofdFerma "ephorfiscal/fiscal/fermaOfd"
	"ephorfiscal/fiscal/interface/fr"
	ofdNanokass "ephorfiscal/fiscal/nanokass"
	ofdOrange "ephorfiscal/fiscal/orange"
	config "ephorservices/config"
)

// instance of type fiscal
var ofd ofdFerma.NewFermaOfdStruct
var Orange ofdOrange.NewOrangeStruct
var Nanokass ofdNanokass.NewNanokassaStruct

func GetFiscal(typeFiscal uint8, cfg *config.Config) fr.Fiscal {
	switch typeFiscal {
	case fr.Fr_EphorOrangeData,
		fr.Fr_ServerOrangeData,
		fr.Fr_EphorServerOrangeData,
		fr.Fr_OrangeData:
		return Orange.New(cfg)
	case fr.Fr_NanoKassa,
		fr.Fr_ServerNanoKassa:
		return Nanokass.New(cfg)
	case fr.Fr_OFD:
		return ofd.New(cfg)
	}
	return nil
}
