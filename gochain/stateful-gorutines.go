package main

import (
    "fmt"
    "math/rand"
    "sync/atomic"
    "time"
    "sort"
)

type readOp struct {
    key  int
    resp chan int
}
type writeOp struct {
    key  int
    val  int
    resp chan bool
}


func randDuration(min int, max int) time.Duration {
    rand.Seed(time.Now().UTC().UnixNano())
    randInt := min + rand.Intn(max-min)
    return time.Duration(randInt) * time.Millisecond
}

func main() {

    var readOps uint64
    var writeOps uint64
    var state = make(map[int]int)

    reads := make(chan *readOp)
    writes := make(chan *writeOp)
    var numOfRepeats = 2
    var minWait = 10
    var maxWait = 11
    var totalRunTime = 100

fmt.Println(" ++++ about to go down ++++ ")

    go func() {
        for {
            select {
            case read := <-reads:
                read.resp <- state[read.key]
                keys := make([]int,0,len(state))
                for key := range state {keys = append(keys,key)}
                sort.Ints(keys)
//                for _, key := range keys{fmt.Println("== read:", key, state[key])}
                fmt.Println("==== read:", read.key, state[read.key])
          case write := <-writes:
                state[write.key] = write.val
                write.resp <- true
                keys := make([]int,0,len(state))
                for key := range state {keys = append(keys,key)}
                sort.Ints(keys)
  //              for _, key := range keys{fmt.Println("= write:", key, state[key])}
                fmt.Println("=== write:", write.key, write.val)
            }
        }
    }()

fmt.Println(" ++++ about to do writes  ++++ ")

    for w := 0; w < numOfRepeats; w++ {
        go func() {
            for {
                write := &writeOp{
                    key:  rand.Intn(5),
                    val:  rand.Intn(100),
                    resp: make(chan bool)}
                writes <- write
                <-write.resp
                atomic.AddUint64(&writeOps, 1)
//		fmt.Println(" ++++  a write happend", write.key, write.val)
                time.Sleep(randDuration(minWait,maxWait))
            }
        }()

    }

fmt.Println(" ++++ about to do reads  ++++ ")

    for r := 0; r < numOfRepeats; r++ {
         go func() {
            for {
                read := &readOp{
                    key:  rand.Intn(5),
                    resp: make(chan int)}
                reads <- read
                /*res :=*/ <-read.resp
                atomic.AddUint64(&readOps, 1)
//		fmt.Println(" ==== a read happend", read.key, res)
                time.Sleep(randDuration(minWait,maxWait))
            }
        }()
   }

    fmt.Println(" ++++ about to go to bed  ++++ ")
    time.Sleep(time.Duration(totalRunTime) * time.Millisecond)

    readOpsFinal := atomic.LoadUint64(&readOps)
    fmt.Println("readOps:", readOpsFinal)
    writeOpsFinal := atomic.LoadUint64(&writeOps)
    fmt.Println("writeOps:", writeOpsFinal)
    fmt.Println("== state:", state)
}

