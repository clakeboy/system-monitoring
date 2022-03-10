package router

import (
	"embed"
	"github.com/clakeboy/golib/components"
	"github.com/gin-gonic/gin"
	"system-monitoring/common"
	"system-monitoring/controllers"
	"system-monitoring/middles"
)

type HttpServer struct {
	server  *gin.Engine
	isDebug bool
	isCross bool
	addr    string
}

func NewHttpServer(addr string, isDebug bool, isCross, isPProf bool) *HttpServer {
	server := &HttpServer{isCross: isCross, isDebug: isDebug, addr: addr}
	server.Init()
	if isPProf {
		server.StartPprof()
	}
	return server
}

func (h *HttpServer) Start() {
	wait := make(chan bool)
	go func() {
		err := h.server.Run(h.addr)
		if err != nil {
			wait <- true
		}
	}()
	<-wait
}

func (h *HttpServer) Init() {
	if h.isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	h.server = gin.New()

	//使用中间件
	if h.isDebug {
		h.server.Use(gin.Logger(), gin.Recovery())
	} else {
		h.server.Use(middles.Logger(), middles.Recovery())
	}

	h.server.Use(middles.Cache())
	h.server.Use(middles.BoltDatabase())
	h.server.Use(middles.Cookie())
	//h.server.Use(middles.Mongo())
	//h.server.Use(middles.Redis())
	//h.server.Use(gzip.Gzip(gzip.DefaultCompression))
	//h.server.Use(middles.Session())
	//跨域调用的OPTIONS
	h.server.OPTIONS("*action", func(c *gin.Context) {
		components.Cross(c, h.isCross, c.Request.Header.Get("Origin"))
	})

	//POST服务接收
	h.server.POST("/serv/:controller/:action", func(c *gin.Context) {
		components.Cross(c, h.isCross, c.Request.Header.Get("Origin"))
		controller := GetController(c.Param("controller"), c)
		components.CallAction(controller, c)
	})
	//GET服务
	h.server.GET("/serv/:controller/:action", func(c *gin.Context) {
		controller := GetController(c.Param("controller"), c)
		components.CallActionGet(controller, c)
	})
	//websocket io
	h.server.GET("/socket.cio/*action", func(c *gin.Context) {
		components.Cross(c, h.isCross, c.Request.Header.Get("Origin"))
		common.SocketIO.Accept(c, controllers.NewSocketController())
	})
	////静态文件访问
	//h.server.GET("/backstage/:filepath", func(c *gin.Context) {
	//	c.FileFromFS("",h.embed)
	//})

	//静态文件API接口
	//h.server.Static("/backstage", "./assets/html")

	//模板页
	//h.server.LoadHTMLGlob("./assets/templates/*")

}

//启动性能探测
func (h *HttpServer) StartPprof() {
	components.InitPprof(h.server)
}

//设置静态文件
func (h *HttpServer) StaticEmbedFS(fs embed.FS) {
	h.server.StaticFS("/backstage", &middles.EmbedFiles{
		Embed: fs,
		Path:  "assets/html",
	})
}
