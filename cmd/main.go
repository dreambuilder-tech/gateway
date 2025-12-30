package main

import (
	"context"
	"gateway/docs"
	"gateway/internal/middleware"
	"gateway/internal/router"
	"log"
	"net/http"
	"strings"
	"time"
	"wallet/common-lib/app"
	"wallet/common-lib/config"
	"wallet/common-lib/mw"
	"wallet/common-lib/natsx"
	"wallet/common-lib/rdb"
	"wallet/common-lib/rpcx/im_rpcx"
	"wallet/common-lib/rpcx/kms_rpcx"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/rpcx/merchant_rpcx"
	"wallet/common-lib/rpcx/trade_rpcx"
	"wallet/common-lib/rpcx/wallet_rpcx"
	"wallet/common-lib/utils/rate_limit"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	svrName = "gateway"
)

func main() {
	app.Run(svrName, entry)
}

func entry(conf *config.BaseConf, svrConf *config.ServiceConfig) func() {
	rdb.Client = rdb.Init(&conf.RedisInstance)

	etcdAddr := conf.Etcd.Endpoints
	trade_rpcx.InitConn(svrName, etcdAddr)
	member_rpcx.InitConn(svrName, etcdAddr)
	wallet_rpcx.InitConn(svrName, etcdAddr)
	kms_rpcx.InitConn(svrName, etcdAddr)
	im_rpcx.InitConn(svrName, etcdAddr)
	merchant_rpcx.InitConn(svrName, etcdAddr)

	natsx.Init(svrName, &conf.Nats)

	rate_limit.Init(rdb.Client, svrName)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		gin.Recovery(),
		mw.Cors(),
		mw.LimitBodySize(15<<20),
		middleware.IPLimiter(),
		otelgin.Middleware(svrName),
		mw.RequestContext(),
		mw.Response(),
		mw.LoggerWithZap(zap.L()),
	)
	
	host := strings.TrimPrefix(svrConf.Service.Http.Address, "0.0.0.0:")
	if host == svrConf.Service.Http.Address {
		host = svrConf.Service.Http.Address
	}
	docs.SwaggerInfo.Host = host
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Init(r)

	s := &http.Server{
		Addr:           svrConf.Service.Http.Address,
		Handler:        r,
		ReadTimeout:    svrConf.HttpReadTimeout(),
		WriteTimeout:   svrConf.HttpWriteTimeout(),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("http listens on %s.", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			log.Printf("http listen error: %v", err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			zap.L().Error("http shutdown error", zap.Error(err))
			_ = s.Close()
		}
		rdb.Close()
		trade_rpcx.CloseConn()
		member_rpcx.CloseConn()
		wallet_rpcx.CloseConn()
		kms_rpcx.CloseConn()
		im_rpcx.CloseConn()
		merchant_rpcx.CloseConn()
		natsx.Close()
	}
}
