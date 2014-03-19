package main

import (
  "net/http"
  "net"
  "os"
  "fmt"
  "flag"
  "time"
)

func main(){
  var err error
  var listenAddr string
  var servePath string
  flag.StringVar(&listenAddr,"http",":18080","http address to listen")
  flag.StringVar(&servePath,"path","","path to servce file,default to current directory.")
  flag.Parse()
  if servePath==""{
    servePath,err=os.Getwd()
    if err!=nil{
      panic(err)
    }
  }
  http.Handle("/",httpHandler{fileServer:http.FileServer(http.Dir(servePath))} )
  l,err:=listenOnAddr(listenAddr)
  if err!=nil{
    panic(err)
  }
  fmt.Printf("Listen on http://%s\n",l.Addr().String() )
  err=http.Serve(l,nil)
  if err!=nil{
    panic(err)
  }
}

//first try addr,if err happened try random addrss.
func listenOnAddr(addr string)(l net.Listener,err error){
  l,err=net.Listen("tcp",addr)
  if err==nil{
    return
  }
  l,err=net.Listen("tcp",":0")
  return 
}
type httpHandler struct{
  fileServer http.Handler
}
type responseStatisticsWriter struct{
  originWriter http.ResponseWriter
  responseCode int
  responseSize int64
}
func (w *responseStatisticsWriter)Header()http.Header{
  return w.originWriter.Header()
}
func (w *responseStatisticsWriter)Write(input []byte) (size int, err error){
  size,err=w.originWriter.Write(input)
  w.responseSize+=int64(size)
  return
}
func (w *responseStatisticsWriter)WriteHeader(code int){
  w.responseCode = code
  return
}
func (w *responseStatisticsWriter)getResponseCode()int{
  if w.responseCode==0{
    return 200
  }
  return w.responseCode
}
func (handler httpHandler)ServeHTTP(w http.ResponseWriter,r *http.Request){
  wrapperWriter := &responseStatisticsWriter{originWriter:w}
  handler.fileServer.ServeHTTP(wrapperWriter,r)
  
  r.URL.Host=r.Host
  r.URL.Scheme = "http";
  fmt.Printf("m[%s] t[%s] r[%s] c[%d] o[%d] u[%s]\n",
  r.Method,
  time.Now().Format(time.RFC3339Nano),
  r.RemoteAddr,
  wrapperWriter.getResponseCode(),
  wrapperWriter.responseSize,
  r.URL.String())
}