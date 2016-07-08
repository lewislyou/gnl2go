package main
import (
	"github.com/kataras/iris"
    "github.com/tedcy/gnl2go"                                                                
	"fmt"
)

func main() {
	ipvs := new(gnl2go.IpvsClient)                                                 
	err := ipvs.Init()                                                             
	if err != nil {                                                                
		fmt.Printf("Cant initialize client, erro is %#v\n", err)                   
		return                                                                     
	}                                                                              
	defer ipvs.Exit()
    iris.Get("/service", func(c *iris.Context) {
		p, err := ipvs.GetPools()                                                      
		if err != nil {                                                                
			fmt.Printf("Error while running GetPools method %#v\n", err)               
			return                                                                     
		}
        c.JSON(iris.StatusOK, p)
    })
	iris.Get("/service/:address", func(c *iris.Context) {
		address := c.Param("address")
		pools, err := ipvs.GetPools()                                                      
		if err != nil {                                                                
			fmt.Printf("Error while running GetPools method %#v\n", err)               
			return                                                                     
		}
		var needPool *gnl2go.Pool
		for _, p := range pools {
			if p.Service.ToString() == address {
				needPool = &p	
			}
        }
        c.JSON(iris.StatusOK, needPool)
    })
    iris.Listen("0.0.0.0:8088")
}
