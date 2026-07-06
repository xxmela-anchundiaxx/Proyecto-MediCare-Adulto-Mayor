package handlers

import (
	"proyecto-medicare-adulto-mayor/internal/handlers/farmacia"
	"proyecto-medicare-adulto-mayor/internal/handlers/medicacion"
	"proyecto-medicare-adulto-mayor/internal/handlers/monitoreo"
	"proyecto-medicare-adulto-mayor/internal/service"
)

type Server struct {
	Medicacion *medicacion.Server
	Farmacia   *farmacia.ManejadorFarmacia
	Monitoreo  *monitoreo.ManejadorMonitoreo
	Auth       *service.AuthService
}

func NewServer(
	medicacionSrv *medicacion.Server,
	farmaciaSrv *farmacia.ManejadorFarmacia,
	monitoreoSrv *monitoreo.ManejadorMonitoreo,
	authSrv      *service.AuthService,
) *Server {
	return &Server{
		Medicacion: medicacionSrv,
		Farmacia:   farmaciaSrv,
		Monitoreo:  monitoreoSrv,
		Auth:       authSrv,
	}
}