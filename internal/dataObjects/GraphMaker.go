package dataObjects

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/components"
	log "github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart/v2"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	global "sysbench-graph-creator/internal/global"
	"time"
)

type chartItem struct {
	label string
	order int
	axis  int
	data  []float64
	color string
}

type charTest struct {
	title        string
	charType     string
	labelX       string
	labelY       string
	numProviders int
	chartItems   []chartItem
}

const (
	HTTPSERVERIPDEFAULT   = "localhost"
	HTTPSERVERPORTDEFAULT = 8089
	PERCONACOLOR          = "orange"
	MYSQLCOLOR            = "blue"
)

//https://github.com/go-echarts/go-echarts

type GraphGenerator struct {
	configuration global.Configuration
	producers     []Producer
	testName      string
	charts        map[string]charTest
}

func (Graph *GraphGenerator) checkConfig() bool {
	if Graph.configuration.Global.TestName == "" {
		for _, producer := range Graph.producers {
			Graph.testName += producer.TestCollectionsName
		}
	} else {
		Graph.testName = Graph.configuration.Global.TestName
	}

	if Graph.configuration.Render.HttpServerPort == 0 {
		Graph.configuration.Render.HttpServerPort = HTTPSERVERPORTDEFAULT
	}
	if Graph.configuration.Render.HttpServerIp == "" {
		Graph.configuration.Render.HttpServerIp = HTTPSERVERIPDEFAULT
	}

	if Graph.configuration.Render.DestinationPath == "" {
		Graph.configuration.Render.DestinationPath, _ = os.Getwd()
		Graph.configuration.Render.DestinationPath += "/html/"
	}

	return true

}

func (Graph *GraphGenerator) Init(inConfig global.Configuration, inProducers []Producer) {
	Graph.producers = inProducers
	Graph.configuration = inConfig
	Graph.checkConfig()
	Graph.charts = make(map[string]charTest)

}

func (Graph *GraphGenerator) Test() {

	graph := chart.BarChart{
		Title: "Test Bar Chart",

		//YAxis: chart.YAxis{
		//	Name: "The YAxis",
		//	Ticks: []chart.Tick{
		//		{Value: 0, Label: "0"},
		//		{Value: 2.0, Label: "2"},
		//		{Value: 4.0, Label: "4"},
		//		{Value: 6.0, Label: "6"},
		//		{Value: 8.0, Label: "8"},
		//		{Value: 10.0, Label: "10"},
		//		{Value: 12.0, Label: "12"},
		//	},
		//},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 5,
			},
		},
		Height:   512,
		BarWidth: 6,
		Bars: []chart.Value{
			{Value: 10.25, Label: "P"},
			{Value: 4.88, Label: "P"},
			{Value: 4.74, Label: "P"},
			{Value: 3.22, Label: "P"},
			{Value: 3, Label: "P"},
			{Value: 2.27, Label: "P"},
			{Value: 1, Label: "P"},
		},
	}
	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}

func (Graph *GraphGenerator) Test2() {
	numValues := 102
	numSeries := 4
	series := make([]chart.Series, numSeries)

	for i := 0; i < numSeries; i++ {
		xValues := make([]time.Time, numValues)
		yValues := make([]float64, numValues)

		for j := 0; j < numValues; j++ {
			xValues[j] = time.Now().AddDate(0, 0, (numValues-j)*-1)
			yValues[j] = random(float64(-50), float64(50))
		}

		series[i] = chart.TimeSeries{
			Name:    fmt.Sprintf("aaa.bbb.hostname-%v.ccc.ddd.eee.fff.ggg.hhh.iii.jjj.kkk.lll.mmm.nnn.value", i),
			XValues: xValues,
			YValues: yValues,
		}
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Time",
		},
		YAxis: chart.YAxis{
			Name: "Value",
		},
		Series: series,
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func (Graph *GraphGenerator) Test3() {

	page := components.NewPage()

	page.AddCharts(barBasic())
	page.AddCharts(barSetToolbox())
	page.AddCharts(barShowLabel())

	f, err := os.Create("html/results.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	fs := http.FileServer(http.Dir("html/"))
	httpServerCoordimates := Graph.configuration.Render.HttpServerIp + ":" + strconv.Itoa(Graph.configuration.Render.HttpServerPort)
	log.Println("running server at http://" + httpServerCoordimates)
	log.Fatal(http.ListenAndServe(httpServerCoordimates, logRequest(fs)))

}

func (Graph *GraphGenerator) Test4() {
	barSetToolbox := barSetToolbox()
	MakeChartSnapshot(barSetToolbox.RenderContent(), "my-bar-title.png")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (Graph *GraphGenerator) RenderReults() bool {
	if Graph.configuration.Render.PrintStats {
		Graph.printStat()
	}

	if Graph.configuration.Render.PrintData {
		Graph.printData()
	}

	return true
}

func (Graph *GraphGenerator) printStat() {
	//Identify how many providers
	//loop for tests
	//	create chartTest
	//		set labels for axis
	//	for each provider identify the test and collect data
	//		setOrder

}

func (Graph *GraphGenerator) printData() {

}
