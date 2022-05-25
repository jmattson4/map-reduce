package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	packageLogger *log.Logger = nil
)

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//completeMapTask makes a rpc call to Coordinator.TaskHandler with the HandlerType of CompleteMap along with
//the taskId which needs to be completed.
func completeMapTask(taskId string) bool {
	tr := &TaskReply{}
	return call("Coordinator.TaskHandler", &TaskArgs{Id: taskId, HandlerType: CompleteMap}, tr)
}

//completeReduceTask makes a rpc call to Coordinator.TaskHandler with the HandlerType of CompleteReduce along with
//the taskId which needs to be completed.
func completeReduceTask(taskId string) bool {
	tr := &TaskReply{}
	return call("Coordinator.TaskHandler", &TaskArgs{Id: taskId, HandlerType: CompleteReduce}, tr)
}

//getTask makes an RPC call to Coordinator.TaskHandler with HandlerType GetTask.
//It then returns a TaskReply and whether the request was ok.
func getTask() (tr *TaskReply, ok bool) {
	tr = &TaskReply{}
	ok = call("Coordinator.TaskHandler", &TaskArgs{HandlerType: GetTask}, tr)
	return
}

//
// cm/worker/mrworker.go calls this function.
// passing in a mapF and a reduceF. It will then continously loop attempting
// to RPC the coordinator server getting tasks for it to complete. These tasks
// are either to Map or Reduce and will call an underlying run function injecting the
// mapf or reducef parameter function into it. These run functions then handle executing
// the tasks.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string, log *log.Logger) {
	// Your worker implementation here.
	packageLogger = log
	ok := true
	for ok {
		tr, ok := getTask()
		if !ok {
			packageLogger.Printf("failed getting task from coordinator %v", os.Getpid())
			break
		}

		switch tr.Type {
		case Map:
			{
				runMapf(tr, mapf)
			}
		case Reduce:
			{
				runReducef(tr, reducef)
			}
		}

	}
}

//runMapf will read the file given in the TaskReply, it will then use the mapf function
// on its contents returning a set of key values which it will then write to intermediary
// files which will later be reduced.
func runMapf(mapTaskReply *TaskReply, mapf func(string, string) []KeyValue) {
	if mapTaskReply.Id == "" {
		return
	}

	packageLogger.Printf("starting map task %v read file operation", mapTaskReply.Id)
	filename := mapTaskReply.Id
	fileContents, err := readFile(filename)
	if err != nil {
		packageLogger.Print(fmt.Errorf("error reading file %v for task %v. Error: %w", filename, mapTaskReply.Id, err))
		return
	}

	packageLogger.Printf("starting map task %v map & sort operation", mapTaskReply.Id)
	kva := mapf(filename, fileContents)

	packageLogger.Printf("starting map task %v intermediary file write", mapTaskReply.Id)
	tempFiles, err := writeIntermediaryFiles(
		mapTaskReply.NReduce,
		mapTaskReply.Id,
		kva,
	)

	if err != nil {
		packageLogger.Print(fmt.Errorf("error writing intermediary files for task %v. Error: %w", mapTaskReply.Id, err))
		return
	}

	ok := completeMapTask(mapTaskReply.Id)
	if !ok {
		packageLogger.Printf("failed map task %v call to CompleteMapTask rpc call to coordinator", mapTaskReply.Id)
		return
	}

	if err = renameTempFiles(tempFiles, mapTaskReply.Id); err != nil {
		packageLogger.Printf("failed map task %v rewrite files error %v", mapTaskReply.Id, err)
		return
	}
}

//runReducef will read the intermediary files based off the id given in the TaskReply.
// it will then attempt to decode these files in parallel by the use of a channel and the function
// decodeIntermediaryFile. After it will finally sort then reduce the given key value pairs before writing to an
// output file.
func runReducef(reduceTaskReply *TaskReply, reducef func(string, []string) string) {
	packageLogger.Printf("starting reduce task %v intermediary file read", reduceTaskReply.Id)
	index, _ := strconv.Atoi(reduceTaskReply.Id)
	files, err := readIntermediaryFiles(index)
	if err != nil {
		packageLogger.Print(fmt.Errorf("error occured while attempting to read reduce files %w", err))
		return
	}

	packageLogger.Printf("starting reduce task %v decode intermediary file", index)
	kvChan := make(chan []KeyValue, reduceTaskReply.NReduce+1)
	wg := &sync.WaitGroup{}
	packageLogger.Printf("reduce task %v returned files length %v", index, len(files))
	wg.Add(len(files))
	for _, f := range files {
		go decodeIntermediaryFile(f, wg, kvChan)
	}

	wg.Wait()
	close(kvChan)

	var intermediate []KeyValue
	for kv := range kvChan {
		intermediate = append(intermediate, kv...)
	}

	sort.Sort(ByKey(intermediate))
	packageLogger.Printf("reduce task %v intermediate arr size %v", index, len(intermediate))
	packageLogger.Printf("starting reduce task %v open output file", index)
	oname := fmt.Sprintf("mr-out-%v", index)
	ofile, err := os.Create(oname)

	if err != nil {
		packageLogger.Print(fmt.Errorf("error occured while attempting to read reduce files %w", err))
		return
	}

	packageLogger.Printf("starting reduce task %v reduce function call", index)
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}
	ofile.Close()
	ok := completeReduceTask(reduceTaskReply.Id)
	if !ok {
		packageLogger.Printf("failed to message coordinator with reduce complete for task %v", index)
		return
	}
}

func decodeIntermediaryFile(f *os.File, wg *sync.WaitGroup, kvChan chan<- []KeyValue) {
	defer wg.Done()
	dec := json.NewDecoder(f)
	var kvArr []KeyValue
	for {
		var kv KeyValue
		if err := dec.Decode(&kv); err != nil {
			break
		}
		kvArr = append(kvArr, kv)
	}
	kvChan <- kvArr
}

func readIntermediaryFiles(reduceId int) ([]*os.File, error) {
	var reduceFiles []*os.File
	files, err := filepath.Glob(fmt.Sprintf("mr-*-%v", reduceId))

	if err != nil {
		return nil, err
	}

	for _, path := range files {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		reduceFiles = append(reduceFiles, file)
	}

	return reduceFiles, nil
}

func createTempFiles(nReduce int) (map[int]*os.File, error) {
	tempFiles := make(map[int]*os.File)
	for i := 0; i < nReduce; i++ {
		file, err := os.CreateTemp("", fmt.Sprintf("temp-%v", i+1))
		if err != nil {
			return nil, fmt.Errorf("cannot create test file")
		}
		tempFiles[i] = file
	}
	return tempFiles, nil
}

func createJsonEncoders(tempFiles map[int]*os.File) map[int]*json.Encoder {
	jsonEncoders := make(map[int]*json.Encoder)
	for i, tf := range tempFiles {
		enc := json.NewEncoder(tf)
		jsonEncoders[i] = enc
	}
	return jsonEncoders
}

func renameTempFiles(tempFiles map[int]*os.File, taskId string) error {
	strippedId := strings.Split(taskId, "/")
	for i, tf := range tempFiles {
		err := os.Rename(tf.Name(), fmt.Sprintf("mr-%v-%v", strippedId[len(strippedId)-1], i))
		if err != nil {
			return err
		}
	}
	return nil
}

func closeTempFiles(tempFiles map[int]*os.File) error {
	for _, tf := range tempFiles {
		if err := tf.Close(); err != nil {
			return err
		}
	}
	return nil
}

func writeIntermediaryFiles(nReduce int, taskId string, keyVArr []KeyValue) (map[int]*os.File, error) {
	tempFiles, err := createTempFiles(nReduce)
	if err != nil {
		return nil, err
	}
	defer closeTempFiles(tempFiles)

	encoders := createJsonEncoders(tempFiles)

	for _, kv := range keyVArr {
		reduceFileNum := ihash(kv.Key) % nReduce
		if err := encoders[reduceFileNum].Encode(kv); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return tempFiles, nil
}

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("cannot open %v", filename)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("cannot read %v", filename)
	}
	return string(content), nil
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		packageLogger.Print("dialing:", err)
		return false
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	packageLogger.Printf("Error making RPC call %v", err)
	return false
}
