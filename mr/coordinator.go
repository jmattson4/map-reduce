package mr

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	googleGrpc "google.golang.org/grpc"

	"github.com/jmattson4/map-reduce/grpc"
)

//Coordinator is a datastructure which represents a map reduce coordinator
//server. This is used to coordinate tasks amongst a group of distributed workers
//which are running map then reduce tasks on a files.
type Coordinator struct {
	grpc.UnimplementedTaskServiceServer

	grpcServ *googleGrpc.Server
	//count of reducer jobs.
	nReduce int
	mu      *sync.RWMutex
	//indicator for completed of map tasks
	mapTaskCompleted bool
	//indicator for completed of redice tasks
	reduceTaskCompleted bool

	//ThreadSafeMap which holds the reduce task id and its current status.
	reduceTasks *ThreadSafeMap[int, TaskStatus]
	//ThreadSafeMap which holds the map task id and its current status.
	mapTasks *ThreadSafeMap[string, TaskStatus]

	logger *log.Logger

	rpcMethod string
	port      uint16

	reduceTaskChan chan int
	mapTaskChan    chan string
}

func (c *Coordinator) GetTask(ctx context.Context, taskArgs *grpc.TaskArgs) (*grpc.TaskReply, error) {
	localTaskReply := &TaskReply{}
	err := c.TaskHandler(&TaskArgs{Id: taskArgs.Id, HandlerType: HandlerType(taskArgs.HandlerType)}, localTaskReply)
	if err != nil {
		return nil, err
	}
	return &grpc.TaskReply{
		Id:      localTaskReply.Id,
		NReduce: int32(localTaskReply.NReduce),
		Type:    int32(localTaskReply.Type),
	}, nil
}

// TaskHandler is the single RPC handler used to manage different operations.
// These operations are based ont the incoming TaskArgs HandlerType. The GetTask HandlerType
// will return a task for a worker to process and then set a timer in a seperate routine which will
// re queue the task if not compelted in a timely manner. CompleteMap will complete a map task in c.mapTasks.
// CompleteReduce will complete a reduce task from c.reduceTasks.
func (c *Coordinator) TaskHandler(args *TaskArgs, reply *TaskReply) error {
	c.logger.Printf("Task Handler Starting type: %v", args.HandlerType)
	switch args.HandlerType {
	case GetTask:
		select {
		case filename := <-c.mapTaskChan:
			c.logger.Printf("Issuing Map Task %v", filename)
			// allocate map task
			reply.NReduce = c.nReduce
			reply.Id = filename
			reply.Type = Map

			c.mapTasks.Put(filename, Started)

			go c.timerForWorker(Map, filename)
			return nil

		case reduceNum := <-c.reduceTaskChan:
			c.logger.Printf("Issuing Reduce Task %v", reduceNum)
			reply.Type = Reduce
			reply.NReduce = c.nReduce
			reply.Id = strconv.Itoa(reduceNum)
			reply.Type = Reduce

			c.reduceTasks.Put(reduceNum, Started)
			go c.timerForWorker(Reduce, strconv.Itoa(reduceNum))
			return nil
		}
	case CompleteMap:
		c.logger.Printf("Completing map task %v", args.Id)
		c.mapTasks.Put(args.Id, Finished)
	case CompleteReduce:
		c.logger.Printf("Completing reduce task %v", args.Id)
		index, _ := strconv.Atoi(args.Id)
		c.reduceTasks.Put(index, Finished)
	}

	return nil
}

// timerForWorker : monitor the worker re queuing the task if it is not completed
// within 10 seconds
func (c *Coordinator) timerForWorker(taskType TaskType, identify string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if taskType == Map {
				c.logger.Printf("Re adding %v to map tasks", identify)
				c.mapTasks.Put(identify, NotStarted)
				c.mapTaskChan <- identify
			} else if taskType == Reduce {
				c.logger.Printf("Re adding %v to reduce tasks", identify)
				index, _ := strconv.Atoi(identify)
				c.reduceTasks.Put(index, NotStarted)
				c.reduceTaskChan <- index
			}
			return
		default:
			if taskType == Map {
				if val, ok := c.mapTasks.Get(identify); ok && val == Finished {
					return
				}
			} else if taskType == Reduce {
				index, _ := strconv.Atoi(identify)
				if val, ok := c.reduceTasks.Get(index); ok && val == Finished {
					return
				}
			}
		}
	}
}

//checkAllMapTask will check if all the current map tasks are complete for a
//given coordinator
func checkAllMapTask(c *Coordinator) bool {
	c.mapTasks.RLock()
	defer c.mapTasks.RUnlock()
	for _, v := range c.mapTasks.m {
		if v != Finished {
			return false
		}
	}
	return true
}

//checkAllReduceTask will check if all the current reduce tasks are complete for a
//given coordinator
func checkAllReduceTask(c *Coordinator) bool {
	c.reduceTasks.RLock()
	defer c.reduceTasks.RUnlock()
	for _, v := range c.reduceTasks.m {
		if v != Finished {
			return false
		}
	}
	return true
}

//setDefaults will set the default values of the coordinator alongs side the given files
//and nReduce integer
func (c *Coordinator) setDefaults(files []string, nReduce int) {
	c.nReduce = nReduce
	c.mapTaskCompleted = false
	c.reduceTaskCompleted = false

	c.mu = &sync.RWMutex{}
	c.mapTasks = NewThreadSafeMap[string, TaskStatus]()
	for _, f := range files {
		c.mapTasks.Put(f, NotStarted)
	}
	c.reduceTasks = NewThreadSafeMap[int, TaskStatus]()
	for i := 0; i < nReduce; i++ {
		c.reduceTasks.Put(i, NotStarted)
	}
}

// must be ran in its own go routine. generates tasks for each c.mapTasks
// feeding them into the mapTasks channel. Thin it will continously check until each
// c.mapTasks are finished. Once they are finished it will generate reduce tasks
// for each c.reduceTasks feeding them into the reduceTasks channel
// and then also check them until they are all complete.
func (c *Coordinator) generateTasks() {
	c.logger.Print("Starting Generate Map Tasks")
	c.mapTasks.RLock()
	for fileName, s := range c.mapTasks.m {
		if s == NotStarted {
			c.logger.Printf("Generating Map Task %v", fileName)
			c.mapTaskChan <- fileName
		}
	}
	c.mapTasks.RUnlock()
	c.logger.Print("Finished Generate Map Tasks")
	ok := false
	for !ok {
		ok = checkAllMapTask(c)
	}

	c.mu.Lock()
	c.mapTaskCompleted = true
	c.mu.Unlock()
	c.logger.Print("Starting Generate Reduce Tasks")
	c.reduceTasks.RLock()
	for id, s := range c.reduceTasks.m {
		if s == NotStarted {
			c.reduceTaskChan <- id
		}
	}
	c.reduceTasks.RUnlock()
	c.logger.Print("Finish Generate Reduce Tasks")
	ok = false
	for !ok {
		ok = checkAllReduceTask(c)
	}

	c.mu.Lock()
	c.reduceTaskCompleted = true
	c.logger.Print("All Reduce Tasks Complete")
	c.mu.Unlock()
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	go c.generateTasks()
	if c.rpcMethod == "rpc" {
		c.rpc()
		//l, e := net.Listen("tcp", ":1234")
		sockname := coordinatorSock()
		os.Remove(sockname)
		c.logger.Printf("Now listening on sock %v", sockname)
		l, e := net.Listen("unix", sockname)
		if e != nil {
			log.Fatal("listen error:", e)
		}
		c.logger.Print("Starting RPC server")
		go http.Serve(l, nil)
	} else {
		serviceName := fmt.Sprintf("localhost:%v", c.port)
		l, e := net.Listen("tcp", serviceName)
		c.logger.Printf("Now listening on sock %v", serviceName)
		if e != nil {
			log.Fatal("listen error:", e)
		}
		c.logger.Print("Starting GRPC server")
		go c.grpc(l)
	}
}

func (c *Coordinator) rpc() {
	rpc.Register(c)
	rpc.HandleHTTP()
}

func (c *Coordinator) grpc(lis net.Listener) {
	c.mu.Lock()
	c.grpcServ = googleGrpc.NewServer()
	grpc.RegisterTaskServiceServer(c.grpcServ, c)
	c.mu.Unlock()
	c.grpcServ.Serve(lis)
}

//
// cmd/coordinator/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
	c.mu.Lock()
	ret = c.reduceTaskCompleted
	if ret && c.rpcMethod == "grpc" {
		c.logger.Printf("Shutting down GRPC server")
		c.grpcServ.Stop()
	}
	c.mu.Unlock()

	return ret
}

func (c *Coordinator) SetLogger(log *log.Logger) {
	c.logger = log
}

func (c *Coordinator) SetPort(port uint16) {
	c.port = port
}

//
// create a Coordinator.
// cmd/coordiantor/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(
	files []string,
	nReduce int,
	log *log.Logger,
	rpcMethod string,
	port uint16,
) *Coordinator {
	c := Coordinator{}
	mapTasks := make(chan string, len(files))
	reduceTasks := make(chan int, nReduce)

	c.mapTaskChan = mapTasks
	c.reduceTaskChan = reduceTasks
	c.rpcMethod = rpcMethod

	c.setDefaults(files, nReduce)
	c.SetLogger(log)
	c.SetPort(port)
	c.server()
	return &c
}
