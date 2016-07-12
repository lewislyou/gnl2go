package main                                                                    
                                                                                
import (                                                                        
    log "github.com/golang/glog"                                                
	"github.com/tedcy/gnl2go"
	"net/http"
    "os"                                                                        
	"fmt"
	"encoding/json"
    "os/signal"                                                                 
    "runtime"                                                                   
    "syscall"                                                                   
	"time"
)

type Server struct {
	ipvs 			*gnl2go.IpvsClient
}

func (this *Server) httpDo(wr http.ResponseWriter, r *http.Request) {
	p, err := this.ipvs.GetPools()                                                      
	if err != nil {                                                                
		log.Infof("Error while running GetPools method %#v\n", err)               
		return                                                                     
	}
	b,err := json.Marshal(p)
	if err != nil {
		log.Infof("json marshal %#v\n", err)               
		return                                                                     
    }
	_,err = wr.Write(b[:])
	if err != nil {
		log.Infof("write %#v\n", err)               
		return
    }
	return
}

func main() {                                                                      
    defer log.Flush()                                                              
    log.Infof("gnl2go proxy [version: %s] start", "1.0")                            
    runtime.GOMAXPROCS(runtime.NumCPU())                                           
    // init http                                                                   
	go func() {                                                                 
		var s *Server
		s = &Server{
			ipvs:			new(gnl2go.IpvsClient),
        }
		err := s.ipvs.Init()                                                             
		if err != nil {                                                                
			fmt.Printf("Cant initialize client, erro is %#v\n", err)                   
			return                                                                     
		}                                                                              
		mux := http.NewServeMux()                                                
        mux.HandleFunc("/service", s.httpDo)                                           
        server := &http.Server {                                                
            Addr:               "127.0.0.1:8088",                                        
            Handler:            mux,                                            
            ReadTimeout:        500000 * time.Second,                         
            WriteTimeout:       500000 * time.Second,                        
        }                                                                       
        if err = server.ListenAndServe();err != nil {                           
            log.Infof("server port serve failed: %s",err)                       
            return                                                              
        }                                                                       
    }()

    ch := make(chan os.Signal, 1)                                                  
    signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
    for {                                                                          
        s := <-ch                                                                  
        log.Infof("get a signal %s", s.String())                                   
        switch s {                                                                 
        case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT: 
            return                                                                 
        case syscall.SIGHUP:                                                       
            // TODO reload                                                         
        default:                                                                   
            return                                                                 
        }                                                                          
    }                                                                              
}
