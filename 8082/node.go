package main
import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"net"
	"time"
	"strconv"
)

type Node struct {
	ip   string
	port string
	receive bool 
	send bool
}
var Iam Node
var Parent Node
var neighbors []Node
var myLeader int
var leaderPath string

func main() {
	filename := `configuration.conf`
	fmt.Println("Start..."  )
	// Find Iam,Initiator,Neighbors
	readFile(filename)
	fmt.Printf("I am %s:%s \nAll my neighbors are: %v \n" , Iam.ip , Iam.port  , neighbors)	
	go server(Iam)
	if checkNeighborServer(neighbors) {
		wfbm := false
		for {
			r, p, i := allExceptOne()
			if r {
				//set p as Parent
				Parent = p
				if myLeader == 0 {
					myLeader,_ = strconv.Atoi(Iam.port) 
					leaderPath = Iam.port
				}
				fmt.Printf("My Parent: %v \n" , Parent)
				message := "&Iam="+Iam.ip+":"+Iam.port+"&leader="+strconv.Itoa(myLeader)+"&path="+leaderPath+"&broadcast=false"
				sendMessage(message, Parent)
				neighbors[i].send = true
				wfbm = true	
			}
			
			if wfbm {
				for {
					time.Sleep(3000 * time.Millisecond)
					fmt.Printf("." )
				}
				
			}
			
		}
	}
}

func broadcast(){
	fmt.Printf("RECEIVEALL %v \n",neighbors )
	
}
func receiveAll() bool{
	retV := true
	for i:=0; i < len(neighbors);i++{
		if !neighbors[i].receive {
			retV = false	
		}
	}
	return retV
}

func allExceptOne() (bool,Node , int){
	j := 0
	k := 0
	for i:=0; i < len(neighbors);i++{
		if !neighbors[i].receive {
				j = i
		}else{
			k = k+1
		}
	}
	if k+1 == len(neighbors){
		return true,neighbors[j],j
	}else{
		return false,Node{"","",false,false},0
	}
	
}



func readFile(fileName string){
	f, _ := os.Open(fileName)
	defer f.Close()
	r := bufio.NewReaderSize(f, 2*1024)
	line, isPrefix, err := r.ReadLine()
	i := 1
	for err == nil && !isPrefix {
		s := string(line)
		if i == 1 {
				// Find Iam
				t :=strings.Split(s, ":")
				Iam = Node{t[0],t[1],false,false}
		}else{
			k :=strings.Split(s, ":")
			neighbors = append(neighbors, Node{k[0],k[1],false,false})		
		}
		i++
		line, isPrefix, err = r.ReadLine()		
	}

}

func analizMessage(message string) map[string]string{
	ms :=strings.Split(message, "&")
	
	msIam :=strings.Split(ms[1], "=")
	mx :=strings.Split(msIam[1], ":")
	GetMessage := make(map[string]string)
	GetMessage["ip"] =  mx[0]
	GetMessage["port"] =  mx[1]
	
	getLeader :=strings.Split(ms[2], "=")
	GetMessage["leader"] =  getLeader[1]
	
	getPath :=strings.Split(ms[3], "=")
	GetMessage["path"] =  getPath[1]
	
	getBroadcast :=strings.Split(ms[4], "=")
	GetMessage["broadcast"] =  getBroadcast[1]
	
	fmt.Println("GetMessage-> ",GetMessage)
	return GetMessage
}

func server(s Node) {
	fmt.Printf("Launching server... %s:%s \n" , s.ip,s.port)
	ln, _ := net.Listen("tcp", s.ip+":"+s.port)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n') 
		if string(message) != "" {
			fmt.Println("->", string(message))
			doIt(analizMessage(message))
		}
	}	
	
}



func sendMessage(s string, n Node){
	conn, _ := net.Dial("tcp", n.ip+":"+n.port)
	defer conn.Close()
	conn.Write([]byte(s))
	fmt.Printf("Message Sent to %s:%s \n" ,n.ip,n.port )	
}


func doIt( ms map[string]string){
	
	fmt.Println(" ++++  ",  ms )
	_,id := findNodeBtwNeighbors(ms["ip"],ms["port"])
	
	//
	bl, _, _ := allExceptOne()
	if !bl {
		neighbors[id].receive = true
		fmt.Println("Myleader,Path",myLeader,leaderPath)
		l,_ := strconv.Atoi(ms["leader"])
		
		if l > myLeader {
			myLeader = l
			leaderPath = ms["port"] + "," +Iam.port
			fmt.Printf("Myleader %v ,Path %v \n",myLeader,leaderPath)
		}	
		fmt.Printf(" &&&&&&&&&& %v \n",neighbors)
	}else{
		if ( ms["broadcast"] == "false"){
			fmt.Printf(" Check if is from my parent \n")
			if fromMyParent(ms["ip"],ms["port"]){
				fmt.Printf(" HAHA \n")
			}else{
				fmt.Printf(" Ignore the mss \n")
			}
		}else{
			fmt.Printf(" broadcast = true \n")
		}
		
		
	}
	
}

func fromMyParent(i string,p string) bool {
	if Parent.ip == i && Parent.port == p {
		return true
	}else{
		return false
	}
}

func checkNeighborServer(n []Node) bool{
	for i:=0; i < len(n);i++{
		for {
			conn, err := net.Dial("tcp", n[i].ip+":"+n[i].port)
			fmt.Println("Looking for " + n[i].ip+":"+n[i].port)
			time.Sleep(3000 * time.Millisecond)
			if err == nil {
				conn.Close()
				break
			}
		}
	}
	
	return true
}

func findNodeBtwNeighbors(ip string, port string) (Node , int){
	j := 0
	for i:=0; i < len(neighbors);i++{
		if neighbors[i].ip == ip && neighbors[i].port == port {
				j = i
		}
	}
	return neighbors[j],j
}
