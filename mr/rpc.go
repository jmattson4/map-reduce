package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"

	"github.com/jmattson4/map-reduce/grpc"
)

// Add your RPC definitions here.

type TaskArgs struct {
	Id          string
	HandlerType HandlerType
}

type TaskReply struct {
	Id      string
	Type    TaskType
	NReduce int
}

func NewTaskReplyFromProtobuf(tr *grpc.TaskReply) *TaskReply {
	return &TaskReply{Id: tr.Id, Type: TaskType(tr.Type), NReduce: int(tr.NReduce)}
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
