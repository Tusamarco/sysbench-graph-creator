/*
 * Copyright (c) Marco Tusa 2021 - present
 *                     GNU GENERAL PUBLIC LICENSE
 *                        Version 3, 29 June 2007
 *
 *  Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
 *  Everyone is permitted to copy and distribute verbatim copies
 *  of this license document, but changing it is not allowed.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package Global

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MyWaitGroup struct {
	sync.Mutex
	sync.WaitGroup
	count int
}

func (wg *MyWaitGroup) WaitTimeout(timeout time.Duration) bool {
	done := make(chan struct{})

	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		return false

	case <-time.After(timeout):
		return true
	}
}

func (wg *MyWaitGroup) IncreaseCounter() {
	wg.Lock()
	defer wg.Unlock()
	wg.count++
}
func (wg *MyWaitGroup) DecreaseCounter() {
	if wg.count > 0 {
		wg.Lock()
		defer wg.Unlock()
		wg.count--
	}
}

func (wg *MyWaitGroup) ReportCounter() int {
	wg.Lock()
	defer wg.Unlock()
	return wg.count
}

const (
	Separator = string(os.PathSeparator)
)

//==================================

func FromStringToMAp(mystring string, separator string) map[string]string {
	myMap := make(map[string]string)
	if mystring != "" && separator != "" {
		keyValuePairArray := strings.Split(mystring, separator)
		for _, keyValuePair := range keyValuePairArray {
			//keyValuePair = strings.Trim(keyValuePair," ")
			keyValueSplit := strings.Split(keyValuePair, "=")
			if len(keyValueSplit) > 1 {
				var key = strings.TrimSpace(keyValueSplit[0])
				var value = strings.TrimSpace(keyValueSplit[1])
				if len(key) > 0 && key != "" && value != "" {
					myMap[key] = value
				}
			}
		}
	}
	return myMap
}

func ToInt(myString string) int {
	if len(myString) > 0 {
		i, err := strconv.Atoi(myString)
		if err != nil {
			pc, fn, line, _ := runtime.Caller(1)
			log.Error(pc, " ", fn, " ", line, ": ", err)
			return -1
		} else {
			return i
		}
	}
	return 0
}

func ToBool(myString string, boolTrueString string) bool {
	myString = strings.ToLower(myString)
	boolTrueString = strings.ToLower(boolTrueString)
	if myString != "" && myString == boolTrueString {
		return true
	} else {
		return false
	}
}

func Bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func ChompSlice(s []string, index int) []string {
	return s[:index]
}

// Reflect if an interface is either a struct or a pointer to a struct
// and has the defined member field, if error is nil, the given
// FieldName exists and is accessible with reflect.
func ReflectStructField(Iface interface{}, FieldName string) error {
	ValueIface := reflect.ValueOf(Iface)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(Iface))
	}

	// 'dereference' with Elem() and get the field by name
	Field := ValueIface.Elem().FieldByName(FieldName)
	if !Field.IsValid() {
		return fmt.Errorf("Interface `%s` does not have the field `%s`", ValueIface.Type(), FieldName)
	}
	return nil
}

//----------------------
/* =====================
STATS
*/

type StatSyncInfo struct {
	sync.RWMutex
	internal map[string][2]int64
}

func NewRegularIntMap() *StatSyncInfo {
	return &StatSyncInfo{
		internal: make(map[string][2]int64),
	}
}

func (rm *StatSyncInfo) Load(key string) (value [2]int64, ok bool) {
	rm.RLock()
	defer rm.RUnlock()
	result, ok := rm.internal[key]

	return result, ok
}

func (rm *StatSyncInfo) Delete(key string) {
	rm.Lock()
	defer rm.Unlock()
	delete(rm.internal, key)

}

func (rm *StatSyncInfo) get(key string) [2]int64 {
	rm.Lock()
	defer rm.Unlock()
	return rm.internal[key]

}

func (rm *StatSyncInfo) Store(key string, value [2]int64) {
	rm.Lock()
	defer rm.Unlock()
	rm.internal[key] = value

}

func (rm *StatSyncInfo) ExposeMap() map[string][2]int64 {
	return rm.internal
}

//====================================================

// Struct
type OrderedPerfMap struct {
	sync.RWMutex
	store map[string]PerfObject
	keys  []string
}

// Constructor
func NewOrderedMap() *OrderedPerfMap {
	return &OrderedPerfMap{
		store: map[string]PerfObject{},
		keys:  []string{},
	}
}

// Get will return the value associated with the key.
// If the key doesn't exist, the second return value will be false.
func (o *OrderedPerfMap) Get(key string) (PerfObject, bool) {
	o.Lock()
	defer o.Unlock()

	val, exists := o.store[key]
	return val, exists
}

// Set will store a key-value pair. If the key already exists,
// it will overwrite the existing key-value pair.
func (o *OrderedPerfMap) Set(key string, val PerfObject) {
	o.Lock()
	defer o.Unlock()

	if _, exists := o.store[key]; !exists {
		o.keys = append(o.keys, key)
	}
	o.store[key] = val
}

// Delete will remove the key and its associated value.
func (o *OrderedPerfMap) Delete(key string) {
	o.Lock()
	defer o.Unlock()

	delete(o.store, key)

	// Find key in slice
	idx := -1

	for i, val := range o.keys {
		if val == key {
			idx = i
			break
		}
	}
	if idx != -1 {
		o.keys = append(o.keys[:idx], o.keys[idx+1:]...)
	}
}

// Iterator is used to loop through the stored key-value pairs.
// The returned anonymous function returns the index, key and value.
func (o *OrderedPerfMap) Iterator() func() (*int, *string, PerfObject) {
	o.Lock()
	defer o.Unlock()

	var keys = o.keys

	j := 0

	return func() (_ *int, _ *string, _ PerfObject) {
		if j > len(keys)-1 {
			return
		}

		row := keys[j]
		j++

		return &[]int{j - 1}[0], &row, o.store[row]
	}
}

func CheckIfPathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CreatePath(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

//=======

func caseInsenstiveFieldByName(v reflect.Value, name string) reflect.Value {
	name = strings.ToLower(name)
	return v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == name })
}

//====================================

// stats structure
type PerfObject struct {
	Name     string
	Time     [2]int64
	LogLevel log.Level
}

// perfomance settings and structure
var Performance bool
var PerformanceMapOrdered *OrderedPerfMap

var PerformanceMap *StatSyncInfo //map[string][2]int64

func SetPerformanceValue(key string, start bool) {
	valStat := [2]int64{0, 0}
	if start {
		valStat[0] = time.Now().UnixNano()
	} else {
		valStat = PerformanceMap.get(key)
		valStat[1] = time.Now().UnixNano()
	}
	PerformanceMap.Store(key, valStat) //  ExposeMap()[key] = valStat
}

func SetPerformanceObj(key string, start bool, logLevel log.Level) {
	var perfObj PerfObject
	valStat := [2]int64{}

	if val, exists := PerformanceMapOrdered.Get(key); !exists {
		perfObj = val
		perfObj.LogLevel = logLevel
		perfObj.Name = key
		valStat = [2]int64{0, 0}
	} else {
		perfObj = val
		valStat = perfObj.Time
	}

	if start {
		valStat[0] = time.Now().UnixNano()
	} else {
		valStat[1] = time.Now().UnixNano()
	}
	perfObj.Time = valStat

	PerformanceMapOrdered.Set(key, perfObj) //  ExposeMap()[key] = valStat
}

func ReportPerformance() {
	formatter := message.NewPrinter(language.English)

	if log.InfoLevel <= log.GetLevel() {
		log.Info("======== Reporting execution times (nanosec/ms)by phase ============")
	}
	it := PerformanceMapOrdered.Iterator()
	for {
		i, _, perfObj := it()
		if perfObj.Name != "" {
			time := perfObj.Time
			value := formatter.Sprintf("%d", time[1]-time[0])
			if perfObj.LogLevel <= log.GetLevel() {
				originalLogLevel := log.GetLevel()
				log.SetLevel(log.InfoLevel)
				log.Info("Phase: ", perfObj.Name, " = ", value, " ns ", strconv.FormatInt((time[1]-time[0])/1000000, 10), " ms")
				log.SetLevel(originalLogLevel)
			}
		}

		if i == nil {
			break
		}
	}
}
func ReturnDateFromString(stringDate string, stringFormat string) (error, time.Time) {

	myDate, err := time.Parse(stringFormat, stringDate)
	if err != nil {
		return err, myDate
	}
	return nil, myDate
}

// quickly return the number of lines in a file counting end of lines
func LineCount(path string) (int, error) {

	file, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	defer file.Close()

	count := 0
	buf := make([]byte, 8192)
	for {
		c, err := file.Read(buf)
		if err != nil {
			if err == io.EOF && c == 0 {
				break
			} else {
				return 0, err
			}
		}
		for _, b := range buf[:c] {
			if b == '\n' {
				count++
			}
		}
		if err == io.EOF {
			err = nil
		}

	}
	return count, nil
}

func ParsetimeLocal(t1 string, t2 string) (time.Time, error) {
	var myTime time.Time
	var err error
	var toProcess string

	if t1 != "" {
		toProcess = t1
	} else {
		toProcess = t2
	}
	if toProcess == "" {
		return myTime, fmt.Errorf("All values passed are empty")
	}

	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{1,2}_\d{2}_\d{2})`)
	match := re.FindStringSubmatch(toProcess)

	if match[0] != "" {
		//strDate := match[1]

		myTime, err = time.Parse("2006-01-02_15_04", match[0])
		if err != nil {
			log.Warnf("Parsing error ", err)
			//return err
		}

	}
	return myTime, err
}

func ReplaceString(input string, pattern string, replacewith string) string {
	var output string
	re := regexp.MustCompile(pattern)
	output = re.ReplaceAllString(input, replacewith)

	return output
}

func IsNumeric(word string) bool {
	return regexp.MustCompile(`\d`).MatchString(word)
}

func Average(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}

//func StandardDeviation(in []float64) float64 {
//	mean := Average(in)
//	result := 0.0
//
//	for _, element := range in {
//		result += math.Pow(element-mean, 2)
//	}
//	result = result / float64(len(in))
//	result = math.Sqrt(result)
//	return result
//
//}

func DistancePct(in []float64) float64 {
	mean := Average(in)
	standardDev := StandardDeviation(in)
	total := (standardDev / mean) * 100
	if standardDev == 0 {
		return 0
	}
	return total
}

func StandardDeviation(num []float64) float64 {
	var mean, sd float64

	mean = Average(num)
	//fmt.Println("The mean of above array is:", mean)
	for _, j := range num {
		sd += math.Pow(j-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(num)))
	return sd
}

func StringsAppendArrayToArray(receiver []string, giver []string) []string {
	for _, element := range giver {
		receiver = append(receiver, element)
	}
	return receiver
}

func FilterOutliners(in []float64) []float64 {
	/*
		The formula to convert datapoint in points z is:
		z =	(x−μ)/σ
		where
		x original value,
		μ average of the given dataset
		σ is the standard deviation.
	*/
	//return in

	out := []float64{}
	bounderies := []float64{-2.0, 2.0} // boundaries set here are very restrictive standard range is -3, 3
	average := Average(in)
	standardDev := StandardDeviation(in)
	for _, value := range in {
		x := value - average
		if x != 0 {
			if x < 0 {
				x = x * -1
			}
			x = x / standardDev
			if x > bounderies[0] && x < bounderies[1] {
				//log.Debugf("Value of x: %.4f out value %.4F", x, value)
				out = append(out, value)
			} else {
				log.Debugf("Filtering out value %.4F x: %.4f", value, x)
			}
		} else {
			out = append(out, value)
		}

	}

	return out

}
