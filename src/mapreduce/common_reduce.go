package mapreduce

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFileName string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	//
	intermediates := make(map[string][]string)

	for i := 0; i < nMap; i++ {
		bytestring, err := ioutil.ReadFile(reduceName(jobName, i, reduceTask))
		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Fields(string(bytestring))

		for _, line := range lines {
			var kv KeyValue
			json.Unmarshal([]byte(line), &kv)

			if intermediates[kv.Key] != nil {
				intermediates[kv.Key] = append(intermediates[kv.Key], kv.Value)
			} else {
				intermediates[kv.Key] = []string{kv.Value}
			}
		}
	}

	var finalResults []KeyValue
	for key, values := range intermediates {
		finalResults = append(finalResults, KeyValue{key, reduceF(key, values)})
	}

	sort.Sort(byKey(finalResults))

	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, finalResult := range finalResults {

		finalResultJSON, err := json.Marshal(finalResult)
		if err != nil {
			log.Fatal(err)
		}

		outFile.WriteString(string(finalResultJSON) + "\n")
	}
}

// Sorting method
type byKey []KeyValue

func (a byKey) Len() int           { return len(a) }
func (a byKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKey) Less(i, j int) bool { return a[i].Key < a[j].Key }
