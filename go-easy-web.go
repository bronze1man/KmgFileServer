package main

import (
  "net/http"
  "net"
  "os"
  "fmt"
  "flag"
)

func main(){
  var listenAddr string
  var servePath string
  var err error
  flag.StringVar(&listenAddr,"http",":18080","http address to listen")
  flag.StringVar(&servePath,"path","","path to servce file,default to current directory.")
  flag.Parse()
  if servePath==""{
    servePath,err=os.Getwd()
    if err!=nil{
      panic(err)
    }
  }
  http.Handle("/",http.FileServer(http.Dir(servePath)))
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