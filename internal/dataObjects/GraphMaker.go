package dataObjects

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	log "github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart/v2"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	global "sysbench-graph-creator/internal/global"
	"time"
)

type chartItem struct {
	label    string
	provider string
	order    int
	axis     int
	data     []opts.BarData
	color    string
	labelX   string
	labelY   string
}

type charTest struct {
	title    string
	charType string
	//labelX       string
	//labelY       string
	numProviders int
	chartItems   []chartItem
	prePost      int
	dimension    string
	actionType   int
	threads      []int
}

const (
	HTTPSERVERIPDEFAULT   = "localhost"
	HTTPSERVERPORTDEFAULT = 8089
	PERCONACOLOR          = "orange"
	MYSQLCOLOR            = "blue"
	XAXISLABELDEFAULT     = "Threads"
)

//https://github.com/go-echarts/go-echarts

type GraphGenerator struct {
	configuration global.Configuration
	producers     []Producer
	testName      string
	chartsData    []charTest
	chartsStats   map[string]charTest
	labels        []string
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
	Graph.chartsData = []charTest{}
	Graph.chartsStats = make(map[string]charTest)
	Graph.labels = strings.Split(inConfig.Render.Labels, ",")

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
	producersLen := len(Graph.producers)
	//emptyResult := ResultTest{}

	if producersLen > 0 {
		testTypes := Graph.findLongestTestList()

		for _, testType := range testTypes {
			//newCharTestStat := charTest{
			//	title:        testType.Name,
			//	charType:     "bar",
			//	labelX:       "Threads",
			//	labelY:       "Distance",
			//	numProviders: producersLen,
			//	chartItems:   nil,
			//}
			newCharTestData := charTest{
				title:    testType.Name,
				charType: "bar",
				//labelX:       "Threads",
				//labelY:       "",
				numProviders: producersLen,
			}
			newCharTestData.chartItems = []chartItem{}
			for _, producer := range Graph.producers {
				var testResult ResultTest

				testKey := TestKey{testType.ActionType,
					producer.TestCollectionsName,
					producer.MySQLProducer,
					producer.MySQLVersion,
					testType.SelectPreWrites,
					testType.Name,
					testType.Dimension}

				for _, tmpResult := range producer.TestsResults {
					if tmpResult.Key == testKey {
						testResult = tmpResult
						break
					}
				}

				newCharTestData.dimension = testResult.Key.Dimension
				newCharTestData.actionType = testResult.Key.ActionType
				newCharTestData.prePost = testResult.Key.SelectPreWrites

				for idx, label := range Graph.labels {
					newThreads := []int{}
					newCharItem := new(chartItem)
					newCharItem.order = idx + 1
					newCharItem.label = label
					newCharItem.provider = producer.MySQLProducer + producer.MySQLVersion
					newCharItem.labelX = XAXISLABELDEFAULT
					newCharItem.labelY = label
					newThreads, newCharItem.data = Graph.getBarData(testResult, label)
					newCharTestData.chartItems = append(newCharTestData.chartItems, *newCharItem)

					if len(newCharTestData.threads) < len(newThreads) {
						newCharTestData.threads = newThreads

					}
				}

				log.Debugf(testResult.Key.TestName)

			}

			Graph.chartsData = append(Graph.chartsData, newCharTestData)

		}

	}

	return true
}

func (Graph *GraphGenerator) printStat(testKey TestKey) {

}

func (Graph *GraphGenerator) printData(testKey TestKey) {

}

func (Graph *GraphGenerator) findLongestTestList() []TestType {
	lenTestTypes := 0
	outTestType := []TestType{}

	for _, producer := range Graph.producers {
		if len(producer.TestsTypes) > lenTestTypes {
			outTestType = producer.TestsTypes
			lenTestTypes = len(producer.TestsTypes)
		}
	}

	return outTestType
}

func (Graph *GraphGenerator) checkForThreadInThreads(in []int, value int) bool {
	if len(in) > 0 {
		for _, th := range in {
			if th == value {
				return true
			}
		}
	}

	return false
}

func (Graph *GraphGenerator) getBarData(testResult ResultTest, inLabel string) ([]int, []opts.BarData) {
	values := []ResultValue{}
	threads := []int{}
	for key, labelValues := range testResult.Labels {
		key = strings.TrimSpace(key)
		if key == strings.TrimSpace(inLabel) {
			values = labelValues
			break
		}
	}
	items := make([]opts.BarData, 0)
	for _, value := range values {
		items = append(items, opts.BarData{Value: value.Value, Name: value.Label})
		threads = append(threads, value.ThreadNumber)
	}
	return threads, items
}

func (Graph *GraphGenerator) BuildPage() bool {
	// Identify if what we need to print (stats/data both)
	var pageData *components.Page
	//var pageStats *components.Page

	if Graph.configuration.Render.PrintData {
		_ = os.Mkdir(Graph.configuration.Render.DestinationPath, os.ModePerm)
		fileFordata, err := os.Create(Graph.configuration.Render.DestinationPath + "data_" + global.ReplaceString(Graph.testName, " ", "") + ".html")
		if err != nil {
			panic(err)
		}

		pageData = components.NewPage()
		pageData.SetLayout(components.PageFlexLayout)
		pageData.PageTitle = Graph.testName

		Graph.addDataToPage(pageData)

		pageData.Render(io.MultiWriter(fileFordata))

	}

	return true
}

func (Graph *GraphGenerator) ActivateHTTPServer() {

}

func (Graph *GraphGenerator) addDataToPage(data *components.Page) {
	//For each test
	// set global params
	// 	Parse labels
	//	set axis labels based on the label
	//		parse provider
	//			add the data
	for _, chartDataTest := range Graph.chartsData {

		bar := charts.NewBar()

		//general
		bar.SetGlobalOptions(
			charts.WithLegendOpts(opts.Legend{Bottom: "0%"}),
			//charts.WithDataZoomOpts(opts.DataZoom{Type:  "slider",Start: 0,End:   50,}),
			//charts.WithDataZoomOpts(opts.DataZoom{Type: "slider"}),
			charts.WithTitleOpts(opts.Title{Title: chartDataTest.title, Subtitle: "Date to add"}),
			charts.WithToolboxOpts(opts.Toolbox{
				Right: "20%",
				Feature: &opts.ToolBoxFeature{
					SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
						Type:  "jpg",
						Title: "Save File",
					},
					DataView: &opts.ToolBoxFeatureDataView{
						Title: "DataView",
						Lang:  []string{"data view", "turn off", "refresh"},
					},
				}},
			),
		)

		for _, labelReference := range Graph.labels {
			for _, chartItemInstance := range chartDataTest.chartItems {
				if chartItemInstance.label == labelReference {

					break
				}
			}
		}

	}

	// to assign by label
	//bar.SetXAxis(weeks).
	//	AddSeries("Category A", generateBarItems()).
	//	AddSeries("Category B", generateBarItems())
	//bar.SetGlobalOptions(
	//	charts.WithYAxisOpts(opts.YAxis{Name: "The YAxis", NameLocation: "middle", NameGap: 50, AxisLabel: &opts.AxisLabel{Rotate: 0.00, Align: "right"}}),
	//	charts.WithXAxisOpts(opts.XAxis{Name: "Threads", NameLocation: "middle", SplitLine: &opts.SplitLine{Show: opts.Bool(true)}}),
	//)

}
