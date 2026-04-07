package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	_ "net/http/pprof" // 注册 pprof 路由

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/time/rate"
)

type User struct {
	ID       int    `json:"id" xml:"id"`
	Name     string `json:"name" xml:"name"`
	Password string `json:"password" xml:"password"`
	Age      int    `json:"age" xml:"age"`
}

type Response struct {
	Code    int `json:"code" xml:"code" form:"code"`
	Message string
	Data    any
	Header  any
}

type UserRequest struct {
	PageIndex int `json:"pageIndex" form:"pageIndex,default=1" xml:"pageIndex" binding:"required"`
	PageSize  int `json:"pageSize" form:"pageSize,default=20" xml:"pageSize" binding:"validPageSize"`
}

type LoginRequest struct {
	Username    string `form:"user_name,default=spy" binding:"-"`
	Password    string `form:"pass_word" binding:"-"`
	Permissions []int  `form:"permissions" collection_format:"csv"`
}

type PagedResult struct {
	List  []any `json:"list" xml:"list"`
	Total int   `json:"total"`
}

type cusHeader struct {
	Host        string `header:"Host"`
	UserAgent   string `header:"User-Agent"`
	ContentType string `header:"Content-Type"`
	Token       string `header:"token" binding:"-"`
}

var users []User = []User{
	{
		ID:       1,
		Name:     "Admin",
		Password: "admin123",
		Age:      30,
	},
}

var validPageSize validator.Func = func(fl validator.FieldLevel) bool {
	if pageSize, ok := fl.Field().Interface().(int); ok {
		if pageSize <= 0 || pageSize > 1000 {
			return false
		}
	}
	return true
}

func RateLimiter() gin.HandlerFunc {
	type client struct {
		limiter *rate.Limiter
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, exists := clients[ip]; !exists {
			// Allow 10 requests per second with a burst of 20
			clients[ip] = &client{limiter: rate.NewLimiter(1, 2)}
		}
		cl := clients[ip]
		mu.Unlock()

		if !cl.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func main() {
	// 写入文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// Disable log's color
	gin.DisableConsoleColor()

	// Force log's color
	gin.ForceConsoleColor()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// SkipQueryString indicates that the logger should not log the query string.
	// For example, /path?q=1 will be logged as /path
	loggerConfig := gin.LoggerConfig{SkipQueryString: true}
	// 单独起一个端口暴露 pprof，不要和业务端口混在一起
	go func() {
		// 这个端口绝不能对外暴露，只在内网或通过 SSH 隧道访问
		http.ListenAndServe("localhost:6060", nil)
	}()

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(gin.LoggerWithConfig(loggerConfig))
	router.Use(RateLimiter())
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	// router.Use(csrf.Middleware(csrf.Options{
	// 	Secret: "csrf-token-secret",
	// 	ErrorFunc: func(c *gin.Context) {
	// 		c.String(403, "CSRF token mismatch")
	// 		c.Abort()
	// 	},
	// }))
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"msg": "文件获取失败" + err.Error()})
			return
		}

		// 保存到本地
		dst := "./uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(500, gin.H{"msg": "保存失败"})
			return
		}

		c.JSON(200, gin.H{
			"msg":  "上传成功",
			"name": file.Filename,
			"size": file.Size,
		})
	})
	router.POST("/uploads", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["files"] // 字段名 files

		for _, f := range files {
			dst := "./uploads/" + f.Filename
			c.SaveUploadedFile(f, dst)
		}

		c.JSON(200, gin.H{
			"msg":   "上传成功",
			"count": len(files),
		})
	})
	router.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		count := session.Get("count")
		if cInt, ok := count.(int); ok {
			session.Set("count", cInt+1)

		} else {
			session.Set("count", 1)
		}
		session.Save()
		c.JSON(http.StatusOK, gin.H{"count": 1})
	})

	router.GET("/get", func(c *gin.Context) {
		session := sessions.Default(c)
		count := session.Get("count")
		c.JSON(http.StatusOK, gin.H{"count": count})
	})

	router.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user", "john")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "logged in"})
	})

	router.GET("/profile", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	router.GET("/form", func(c *gin.Context) {
		token := csrf.GetToken(c)
		c.JSON(200, gin.H{"csrf_token": token})
	})

	router.Static("/files", "./tmp")
	router.StaticFS("/static", http.Dir("D:\\"))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validPageSize", validPageSize)
	}
	router.Any("/any")

	var fs http.FileSystem = http.Dir("./tmp")
	router.GET("/fs/file", func(ctx *gin.Context) {
		ctx.FileFromFS("/main.exe", fs)
	})

	router.GET("/file1", func(ctx *gin.Context) {
		ctx.File("./tmp/main.exe")
	})

	router.GET("/download1", func(c *gin.Context) {
		file, _ := os.Open("./tmp/main.exe")
		fileInfo, _ := file.Stat()
		var length = fileInfo.Size()
		extHeader := make(map[string]string)
		extHeader["content-disposition"] = `attachment; filename="main.zip"`
		// io.Copy(c.Writer, file)
		// http.ServeContent()
		c.DataFromReader(http.StatusOK, length, "application/file", file, extHeader)
	})

	router.GET("/download", func(c *gin.Context) {
		c.FileAttachment("./tmp/main.exe", "main.exe")
	})
	router.GET("/users", func(ctx *gin.Context) {
		var request UserRequest
		err := ctx.ShouldBindQuery(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, Response{
			Code:    http.StatusOK,
			Message: "",
			Data:    request,
		})
	})

	router.POST("/login", func(ctx *gin.Context) {
		var request LoginRequest
		var header cusHeader
		ctx.BindHeader(&header)
		err := ctx.ShouldBind(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Code:    http.StatusOK,
			Message: "登录<成功",
			Data:    request,
			Header:  header,
		})
	})

	router.Run()
}
