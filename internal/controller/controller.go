package controller
import(
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/client"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/model"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/services"
	"github.com/gorilla/websocket"
)

type Controller interface {
	WebsocketHandler(c *gin.Context)
}

type controller struct {
	DB *leveldb.DB
	Serv services.SendTransactionService
}

func NewController(
	DB *leveldb.DB,
	Serv services.SendTransactionService,
) Controller {
	return &controller{
		DB : DB,
		Serv : Serv,
	}
}

func (controller *controller) WebsocketHandler (c *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := client.Client{
		Conn	: conn,
		MsgChan : make(chan model.Message),
		DB : controller.DB,
		Serv: controller.Serv,
	}
	go client.ReadMessage()
	go client.WriteMessage()
}