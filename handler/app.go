package handler

import (
	mysql "esmartcare/infra/mysql"

	adminmysql "esmartcare/repository/adminRepository/adminMysql"
	alarmmysql "esmartcare/repository/alarmRepository/alarmMysql"
	"esmartcare/repository/pakarRepository/pakarMysql"
	pemeriksaanmysql "esmartcare/repository/pemeriksaanRepository/pemeriksaanMysql"
	riwayatmysql "esmartcare/repository/riwayatTanyaJawabRepository/riwayatMysql"
	siswamysql "esmartcare/repository/siswaRepository/siswaMysql"
	tanyajawabmysql "esmartcare/repository/tanyaJawabRepository/tanyajawabMysql"

	usermysql "esmartcare/repository/userrepository/userMysql"
	"esmartcare/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func StartApp() {
	db := mysql.GetDBInstance()
	//sgin.SetMode(gin.ReleaseMode)

	port := os.Getenv("PORT")
	route := gin.Default()

	route.MaxMultipartMemory = 8 << 20 //

	userRepo := usermysql.NewUserMySql(db)
	siswaRepo := siswamysql.NewSiswaMySql(db)
	adminRepo := adminmysql.NewAdminMySql(db)
	pakarRepo := pakarMysql.NewpakarMysql(db)
	riwayatTanyaJawabRepo := riwayatmysql.NewRiwayatTanyaJawabRepository(db)
	pemeriksaanRepo := pemeriksaanmysql.NewPemeriksaanRepository(db)

	authService := service.NewAuthService(userRepo, siswaRepo, riwayatTanyaJawabRepo, pemeriksaanRepo)

	userService := service.NewUserService(userRepo, adminRepo, siswaRepo, pakarRepo)
	userHandler := NewUserHandler(userService)

	siswaService := service.NewSiswaService(siswaRepo)
	siswaHandler := NewSiswaHandler(siswaService)

	adminService := service.NewadminService(adminRepo)
	adminHandler := NewAdminHandler(adminService)

	pakarService := service.NewpakarService(pakarRepo)
	pakarHandler := NewPakarHandler(pakarService)

	tanyaJawabRepo := tanyajawabmysql.NewTanyaJawabRepository(db)
	tanyaJawabService := service.NewTanyaJawabService(tanyaJawabRepo)
	tanyaJawabHandler := NewTanyaJawabHandler(tanyaJawabService)

	riwayatTanyaJawabService := service.NewRiwayatTanyaJawabService(riwayatTanyaJawabRepo)
	riwayatTanyaJawabHandler := NewRiwayatTanyaJawabHandler(riwayatTanyaJawabService)

	pemeriksaanService := service.NewPemeriksaanService(pemeriksaanRepo)
	pemeriksaanHandler := NewPemeriksaanHandler(pemeriksaanService)

	alarmRepo := alarmmysql.NewAlarmRepository(db)
	alarmService := service.NewAlarmService(alarmRepo)
	alarmHandler := NewAlarmHandler(alarmService)

	route.Static("/uploads", "./uploads")

	UsersRoute := route.Group("/users")
	{
		UsersRoute.POST("/register", userHandler.Register)
		UsersRoute.POST("/login", userHandler.Login)
		UsersRoute.GET("/", authService.Authentication(), authService.AdminAuthorization(), userHandler.GettAllUsers)
		UsersRoute.GET("/all-data-users", authService.Authentication(), authService.AdminAndPakarAuthorization(), userHandler.GetAllDataUser)
		UsersRoute.GET("/not-validate", authService.Authentication(), authService.AdminAuthorization(), userHandler.GettAllUsersNotValidate)
		UsersRoute.PUT("/update-user/:email", authService.Authentication(), authService.AdminAuthorization(), userHandler.UpdateUser)
		UsersRoute.DELETE("/delete-account/:email", authService.Authentication(), authService.AdminAuthorization(), userHandler.DeleteUser)
	}

	SiswaRoute := route.Group("/siswa")
	{
		SiswaRoute.GET("/", authService.Authentication(), authService.SiswaAuthorization(), siswaHandler.GetSiswa)
		SiswaRoute.POST("/upload-photo", authService.Authentication(), authService.SiswaAuthorization(), siswaHandler.UploadProfileImage)
		SiswaRoute.POST("/update-profile", authService.Authentication(), authService.SiswaAuthorization(), siswaHandler.CreateOrUpdateSiswa)
	}

	AdminRoute := route.Group("/admin")
	{
		AdminRoute.GET("/", authService.Authentication(), authService.AdminAuthorization(), adminHandler.GetAdmin)
		AdminRoute.POST("/upload-photo", authService.Authentication(), authService.AdminAuthorization(), adminHandler.UploadProfileImage)
		AdminRoute.POST("/update-profile", authService.Authentication(), authService.AdminAuthorization(), adminHandler.CreateOrUpdateAdmin)
	}

	PakarRoute := route.Group("/pakar")
	{
		PakarRoute.GET("/", authService.Authentication(), authService.PakarAuthorization(), pakarHandler.GetPakar)
		PakarRoute.POST("/upload-photo", authService.Authentication(), authService.PakarAuthorization(), pakarHandler.UploadProfileImage)
		PakarRoute.POST("/update-profile", authService.Authentication(), authService.PakarAuthorization(), pakarHandler.CreateOrUpdatePakar)
	}

	TanyaJawabRoute := route.Group("/tanya-jawab")
	{
		TanyaJawabRoute.GET("/", authService.Authentication(), authService.AdminAndPakarAuthorization(), tanyaJawabHandler.GetTanyaJawab)
		TanyaJawabRoute.POST("/", authService.Authentication(), authService.AdminAndPakarAuthorization(), tanyaJawabHandler.CreateTanyaJawab)
		TanyaJawabRoute.PUT("/:id", authService.Authentication(), authService.AdminAndPakarAuthorization(), tanyaJawabHandler.UpdateTanyaJawab)
		TanyaJawabRoute.PUT("/validator/:id", authService.Authentication(), authService.PakarAuthorization(), tanyaJawabHandler.UpdateValidator)
		TanyaJawabRoute.DELETE("/:id", authService.Authentication(), authService.AdminAndPakarAuthorization(), tanyaJawabHandler.DeleteTanyaJawab)
	}

	ChatBotRoute := route.Group("/chat-bot")
	{
		ChatBotRoute.POST("/get-all-similar", tanyaJawabHandler.ChatSimmilarityBot)
		ChatBotRoute.POST("/", tanyaJawabHandler.ChatBot)
	}

	RiwayatRoute := route.Group("/riwayat")
	{

		RiwayatRoute.GET("/", authService.Authentication(), riwayatTanyaJawabHandler.GetRiwayatByEmail)
		RiwayatRoute.GET("/all-riwayat", authService.Authentication(), authService.AdminAndPakarAuthorization(), riwayatTanyaJawabHandler.GetAllDataRiwayat)
		RiwayatRoute.POST("/", authService.Authentication(), riwayatTanyaJawabHandler.CreateRiwayat)
		RiwayatRoute.DELETE("/:id", authService.Authentication(), authService.RiwayatAuthorization(), riwayatTanyaJawabHandler.DeleteRiwayatById)
		RiwayatRoute.DELETE("/delete-email/:email", authService.Authentication(), authService.AdminAndPakarAuthorization(), riwayatTanyaJawabHandler.DeleteRiwayatByEmail)
		RiwayatRoute.DELETE("/delete-all", authService.Authentication(), riwayatTanyaJawabHandler.DeleteAllRiwayatByUserLogin)

	}

	PemeriksaanRoute := route.Group("/pemeriksaan")
	{
		//PemeriksaanRoute.GET("/", pemeriksaanHandler.GetAllPemeriksaan)
		PemeriksaanRoute.POST("/", authService.Authentication(), pemeriksaanHandler.CreatePemeriksaan)
		PemeriksaanRoute.POST("/upload-photo-pemeriksaan", authService.Authentication(), pemeriksaanHandler.UploadPhotoPemeriksaan)
		PemeriksaanRoute.GET("/", authService.Authentication(), pemeriksaanHandler.GetPemeriksaanByEmail)
		PemeriksaanRoute.DELETE("/:id", authService.Authentication(), authService.PemeriksaanAuthorization(), pemeriksaanHandler.DeletePemeriksaanById)
	}

	AlarmRoute := route.Group("/alarm")
	{
		AlarmRoute.GET("/", authService.Authentication(), alarmHandler.GetAlarmsByEmail)
		AlarmRoute.POST("/", authService.Authentication(), alarmHandler.CreateAlarm)
		AlarmRoute.PUT("/:id", authService.Authentication(), alarmHandler.UpdateAlarm)
		AlarmRoute.DELETE("/:id", authService.Authentication(), alarmHandler.DeleteAlarmByID)
	}

	route.POST("/update-bot", tanyaJawabHandler.Update_Bot)
	route.GET("/monitoring", siswaHandler.GetAllSiswaWithPemeriksaan)

	log.Fatalln(route.Run(":" + port))
}
