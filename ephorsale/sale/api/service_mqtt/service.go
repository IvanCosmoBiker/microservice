package service_mqtt

import (
	"context"
	config "ephorservices/config"
	request "ephorservices/ephorsale/sale/api/service_mqtt/request"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	logger "ephorservices/pkg/logger"
	queue_manager "ephorservices/pkg/mqttmanager"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Pay struct {
	Name string
}

func (p *Pay) Consume(ctx context.Context, data []byte) error {
	res := request.ResponsePay{}
	err := res.JsonToStruct(data)
	if err != nil {
		logger.Log.Error(err.Error())
		return errors.New("Fail parse json")
	}
	fmt.Printf("%+v", res)
	transaction_dispetcher.Dispetcher.Send(res.Tid, data)
	return nil
}

func (p *Pay) Shutdown(ctx context.Context) error {
	err := errors.New("end pay")
	logger.Log.Error(err.Error())
	return err
}

type Fiscal struct {
	Name string
}

type Command struct {
	Name string
}
type ServiceMqtt struct {
	Ctx context.Context
}

func New(ctx context.Context) *ServiceMqtt {
	modemApi := &ServiceMqtt{
		Ctx: ctx,
	}
	return modemApi
}

func (m *ServiceMqtt) initConsumers(conf *config.Config, ctx context.Context, errorGroup *errgroup.Group) {
	m.Pay(conf, ctx, errorGroup)
	m.Fiscal(conf, ctx, errorGroup)
	m.Command(conf, ctx, errorGroup)
}

func (m *ServiceMqtt) Pay(conf *config.Config, ctx context.Context, errorGroup *errgroup.Group) error {
	consumer, err := queue_manager.Broker.AddConsumer(ctx, conf.Services.EphorPay.NameQueue, conf.Services.EphorPay.NameQueue)
	if err != nil {
		return err
	}
	pay := &Pay{
		Name: conf.Services.EphorPay.NameQueue,
	}
	err = consumer.Subscribe(ctx, errorGroup, pay)
	return err
}

func (m *ServiceMqtt) Fiscal(conf *config.Config, ctx context.Context, errorGroup *errgroup.Group) error {
	return nil
}

func (m *ServiceMqtt) Command(conf *config.Config, ctx context.Context, errorGroup *errgroup.Group) error {
	return nil
}

func (m *ServiceMqtt) InitApi(conf *config.Config, ctx context.Context, errorGroup *errgroup.Group) {
	m.initConsumers(conf, ctx, errorGroup)
}
