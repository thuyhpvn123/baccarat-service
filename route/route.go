package route
import(
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/controller"
	"net/http"
)
func InitialRoutes(engine *gin.Engine, controller controller.Controller) {
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	engine.MaxMultipartMemory = 8 << 20 // 8 MiB
	v1 := engine.Group("/api/v1")
	{
		v1.StaticFS("", http.Dir("../frontend"))
	}
	engine.GET("/ws", func(c *gin.Context) {
		controller.WebsocketHandler(c )})
}
	//http://localhost:8999/api/v1/template/
