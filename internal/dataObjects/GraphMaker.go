package dataObjects

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
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
	dataBetter   []int
	isBetter     bool
	totalPoints  float64
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
	chartsStats   []charTest
	labels        []string
	statLabels    []string
	benchTool     string
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

	if Graph.configuration.Render.CsvDestinationPath == "" {
		Graph.configuration.Render.CsvDestinationPath, _ = os.Getwd()
		Graph.configuration.Render.CsvDestinationPath += "/csv/"
	}

	return true

}

func (Graph *GraphGenerator) Init(inConfig global.Configuration, inProducers []Producer, testCollection []TestCollection) {
	if len(testCollection) > 0 {
		Graph.benchTool = testCollection[0].Producer
	}

	Graph.producers = inProducers
	Graph.configuration = inConfig
	Graph.checkConfig()
	Graph.chartsData = []charTest{}
	Graph.chartsStats = []charTest{}
	Graph.labels = strings.Split(inConfig.Render.Labels, ",")
	Graph.statLabels = strings.Split(inConfig.Render.StatsLabels, ",")

}

//
//func (Graph *GraphGenerator) Test() {
//
//	graph := chart.BarChart{
//		Title: "Test Bar Chart",
//
//		//YAxis: chart.YAxis{
//		//	Name: "The YAxis",
//		//	Ticks: []chart.Tick{
//		//		{Value: 0, Label: "0"},
//		//		{Value: 2.0, Label: "2"},
//		//		{Value: 4.0, Label: "4"},
//		//		{Value: 6.0, Label: "6"},
//		//		{Value: 8.0, Label: "8"},
//		//		{Value: 10.0, Label: "10"},
//		//		{Value: 12.0, Label: "12"},
//		//	},
//		//},
//		Background: chart.Style{
//			Padding: chart.Box{
//				Top: 5,
//			},
//		},
//		Height:   512,
//		BarWidth: 6,
//		Bars: []chart.Value{
//			{Value: 10.25, Label: "P"},
//			{Value: 4.88, Label: "P"},
//			{Value: 4.74, Label: "P"},
//			{Value: 3.22, Label: "P"},
//			{Value: 3, Label: "P"},
//			{Value: 2.27, Label: "P"},
//			{Value: 1, Label: "P"},
//		},
//	}
//	f, _ := os.Create("output.png")
//	defer f.Close()
//	graph.Render(chart.PNG, f)
//
//}
//
//func (Graph *GraphGenerator) Test2() {
//	numValues := 102
//	numSeries := 4
//	series := make([]chart.Series, numSeries)
//
//	for i := 0; i < numSeries; i++ {
//		xValues := make([]time.Time, numValues)
//		yValues := make([]float64, numValues)
//
//		for j := 0; j < numValues; j++ {
//			xValues[j] = time.Now().AddDate(0, 0, (numValues-j)*-1)
//			yValues[j] = random(float64(-50), float64(50))
//		}
//
//		series[i] = chart.TimeSeries{
//			Name:    fmt.Sprintf("aaa.bbb.hostname-%v.ccc.ddd.eee.fff.ggg.hhh.iii.jjj.kkk.lll.mmm.nnn.value", i),
//			XValues: xValues,
//			YValues: yValues,
//		}
//	}
//
//	graph := chart.Chart{
//		XAxis: chart.XAxis{
//			Name: "Time",
//		},
//		YAxis: chart.YAxis{
//			Name: "Value",
//		},
//		Series: series,
//	}
//
//	f, _ := os.Create("output.png")
//	defer f.Close()
//	graph.Render(chart.PNG, f)
//}
//
//func random(min, max float64) float64 {
//	return rand.Float64()*(max-min) + min
//}
//
//func (Graph *GraphGenerator) Test3() {
//
//	page := components.NewPage()
//
//	page.AddCharts(barBasic())
//	page.AddCharts(barSetToolbox())
//	page.AddCharts(barShowLabel())
//
//	f, err := os.Create("html/results.html")
//	if err != nil {
//		panic(err)
//	}
//	page.Render(io.MultiWriter(f))
//	fs := http.FileServer(http.Dir("html/"))
//	httpServerCoordimates := Graph.configuration.Render.HttpServerIp + ":" + strconv.Itoa(Graph.configuration.Render.HttpServerPort)
//	log.Println("running server at http://" + httpServerCoordimates)
//	log.Fatal(http.ListenAndServe(httpServerCoordimates, logRequest(fs)))
//
//}
//
//func (Graph *GraphGenerator) Test4() {
//	barSetToolbox := barSetToolbox()
//
//	MakeChartSnapshot(barSetToolbox.RenderContent(), "my-bar-title.png")
//}

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
			newCharTestStat := charTest{
				title:        testType.Name,
				charType:     "bar",
				numProviders: producersLen,
			}
			newCharTestData := charTest{
				title:        testType.Name,
				charType:     "bar",
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

						//filling data object
						newCharTestData.dimension = testResult.Key.Dimension
						newCharTestData.actionType = testResult.Key.ActionType
						newCharTestData.prePost = testResult.Key.SelectPreWrites

						//filling stats object
						newCharTestStat.dimension = testResult.Key.Dimension
						newCharTestStat.actionType = testResult.Key.ActionType
						newCharTestStat.prePost = testResult.Key.SelectPreWrites

						//IF test is select scan add TotalTime label to the list
						mylables := []string{}
						if strings.Contains(testResult.Key.TestName, "select_run_select_scan") {
							mylables = []string{"TotalTime"}
						}
						mylables = global.StringsAppendArrayToArray(mylables, Graph.labels)

						for idx, label := range mylables {
							newThreads := []int{}

							//Filling data
							newCharItem := new(chartItem)
							newCharItem.order = idx + 1
							newCharItem.label = label
							newCharItem.provider = producer.MySQLProducer + producer.MySQLVersion + producer.TestCollectionsName
							newCharItem.labelX = XAXISLABELDEFAULT
							newCharItem.labelY = label
							newThreads, newCharItem.data = Graph.getBarData(testResult, label)
							newCharTestData.chartItems = append(newCharTestData.chartItems, *newCharItem)

							if len(newCharTestData.threads) < len(newThreads) {
								newCharTestData.threads = newThreads

							}

							//filling stats
							newCharStatsItem := new(chartItem)
							newCharStatsItem.order = idx + 1
							newCharStatsItem.label = label
							newCharStatsItem.provider = producer.MySQLProducer + producer.MySQLVersion + producer.TestCollectionsName
							newCharStatsItem.labelX = XAXISLABELDEFAULT
							newCharStatsItem.labelY = label
							newThreads, newCharStatsItem.data = Graph.getBarStats(testResult, label)
							newCharTestStat.chartItems = append(newCharTestStat.chartItems, *newCharStatsItem)

							if len(newCharTestStat.threads) < len(newThreads) {
								newCharTestStat.threads = newThreads

							}

						}

						//log.Debugf(testResult.Key.TestName)

						break
					}
				}

			}

			Graph.chartsData = append(Graph.chartsData, newCharTestData)
			Graph.chartsStats = append(Graph.chartsStats, newCharTestStat)
		}

	}
	//calculate summary results
	Graph.calculateSummary()
	return true
}

func (Graph *GraphGenerator) calculateSummary() bool {
	for _, chartDataTest := range Graph.chartsData {
		var evalLabel string
		if chartDataTest.actionType == 0 {
			evalLabel = Graph.configuration.Render.ReadSummaryLabel
		} else {
			evalLabel = Graph.configuration.Render.WriteSummaryLabel
		}
		if !strings.Contains(strings.ToLower(chartDataTest.title), "warmup") {
			for _, item := range chartDataTest.chartItems {

				if item.label == evalLabel {

				}

			}
		}
	}

	return true
}

func (Graph *GraphGenerator) findLongestTestList() []TestType {
	//lenTestTypes := 0
	outTestType := []TestType{}

	for _, producer := range Graph.producers {
		outTestType = Graph.mergeTestList(outTestType, producer.TestsTypes)

		//if len(producer.TestsTypes) > lenTestTypes {
		//	outTestType = producer.TestsTypes
		//	lenTestTypes = len(producer.TestsTypes)
		//}
	}

	return outTestType
}

func (Graph *GraphGenerator) mergeTestList(in []TestType, toMerge []TestType) []TestType {

	for _, elementToMerge := range toMerge {
		merge := true
		for _, elementToTest := range in {
			if
			//elementToTest.TestCollectionName == elementToMerge.TestCollectionName &&
			elementToTest.Name == elementToMerge.Name &&
				elementToTest.Dimension == elementToMerge.Dimension &&
				elementToTest.SelectPreWrites == elementToMerge.SelectPreWrites &&
				elementToTest.ActionType == elementToMerge.ActionType {
				merge = false
			}
		}
		if merge {
			in = append(in, elementToMerge)
		}

	}

	return in
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
	var pageStats *components.Page

	//we create the html page with the data
	if Graph.configuration.Render.PrintData {
		_ = os.Mkdir(Graph.configuration.Render.DestinationPath, os.ModePerm)
		fileFordata, err := os.Create(Graph.configuration.Render.DestinationPath + "data_" +
			global.ReplaceString(Graph.testName, " ", "") + "_" + Graph.benchTool + ".html")
		if err != nil {
			panic(err)
		}

		pageData = components.NewPage()
		pageData.SetLayout(components.PageFlexLayout)
		pageData.PageTitle = Graph.testName

		Graph.addDataToPage(pageData)

		pageData.Render(io.MultiWriter(fileFordata))

	}

	//we create the html page with stats
	if Graph.configuration.Render.PrintStats {
		_ = os.Mkdir(Graph.configuration.Render.DestinationPath, os.ModePerm)
		fileForStats, err := os.Create(Graph.configuration.Render.DestinationPath + "stats_" +
			global.ReplaceString(Graph.testName, " ", "") + "_" + Graph.benchTool + ".html")
		if err != nil {
			panic(err)
		}

		pageStats = components.NewPage()
		pageStats.SetLayout(components.PageFlexLayout)
		pageStats.PageTitle = Graph.testName + " STATISTICS"

		Graph.addStatsToPage(pageStats)

		pageStats.Render(io.MultiWriter(fileForStats))

	}

	//we create the image files (one for each graph)
	if Graph.configuration.Render.PrintCharts {
		Graph.PrintImages()
	}

	//we create the CSV file with all data
	if Graph.configuration.Render.ConvertChartsToCsv {
		Graph.PrintDataCsv()
	}

	return true
}

func (Graph *GraphGenerator) PrintImages() {

	if _, err := os.Stat(Graph.configuration.Render.DestinationPath + "images/"); os.IsNotExist(err) {
		err = os.Mkdir(Graph.configuration.Render.DestinationPath+"images/", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	for _, chartDataTest := range Graph.chartsData {
		if !strings.Contains(strings.ToLower(chartDataTest.title), "warmup") {
			mylables := []string{}
			if strings.Contains(chartDataTest.title, "select_run_select_scan") {
				mylables = []string{"TotalTime", "latencyPct95(μs)"}
				chartDataTest.title += " (Lower is better) "
			} else {
				mylables = global.StringsAppendArrayToArray(mylables, Graph.labels)
			}

			for _, labelReference := range mylables {
				image := Graph.configuration.Render.DestinationPath + "images/"
				bar := charts.NewBar()

				titleFull := global.ReplaceString(chartDataTest.title, "_", " ") + " " + chartDataTest.dimension
				if chartDataTest.prePost == 0 {
					titleFull += " Pre Writes"
				} else {
					titleFull += " Post Writes"
				}

				titleFull = titleFull + "_" + labelReference
				image = image + global.ReplaceString(titleFull, "[\\s\\/%\\(\\)]", "_") + ".jpg"
				//image = image + global.ReplaceString(titleFull, "__", "_") + ".jpg"

				bar.SetGlobalOptions(
					charts.WithInitializationOpts(opts.Initialization{
						Width:  strconv.Itoa(Graph.configuration.Render.ChartWidth) + "px",
						Height: strconv.Itoa(Graph.configuration.Render.ChartHeight) + "px",
					}),
					charts.WithLegendOpts(opts.Legend{Width: "90%", Height: "300", Bottom: "-1%", Type: "plain"}),
					charts.WithXAxisOpts(opts.XAxis{Name: "Threads", NameGap: 20, NameLocation: "middle", SplitLine: &opts.SplitLine{Show: opts.Bool(true)}}),
					charts.WithAnimation(false),
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
					charts.WithTitleOpts(opts.Title{Title: titleFull, Subtitle: labelReference}),
					charts.WithYAxisOpts(opts.YAxis{Name: labelReference, NameLocation: "middle", NameGap: 60, AxisLabel: &opts.AxisLabel{Rotate: 0.00, Align: "right"}}),
				)
				for _, chartItemInstance := range chartDataTest.chartItems {
					if chartItemInstance.label == labelReference {
						//log.Debugf("Len items data %d  test %s label %s", len(chartItemInstance.data), chartDataTest.title, chartItemInstance.label)
						bar.SetXAxis(chartDataTest.threads).AddSeries(chartItemInstance.provider, chartItemInstance.data)
					}
				}

				path, file := filepath.Split(image)
				suffix := filepath.Ext(file)[1:]
				fileName := file[0 : len(file)-len(suffix)-1]

				config := &SnapshotConfig{
					RenderContent: bar.RenderContent(),
					Path:          path,
					FileName:      fileName,
					Suffix:        suffix,
					Quality:       1,
					KeepHtml:      false,
				}

				errImage := MakeSnapshot(config)
				if errImage != nil {
					log.Errorf("Error printing image %s", image)
				} else {
					log.Debugf("Printing image %s", image)
				}

			}
		}
	}
}

func (Graph *GraphGenerator) ActivateHTTPServer() {

}

func (Graph *GraphGenerator) addDataToPage(page *components.Page) {
	//For each test
	// set global params
	// 	Parse labels
	//	set axis labels based on the label
	//		parse provider
	//			add the data
	for _, chartDataTest := range Graph.chartsData {
		if !strings.Contains(strings.ToLower(chartDataTest.title), "warmup") {

			//IF test is select scan we onl show totaltime and latency
			mylables := []string{}
			if strings.Contains(chartDataTest.title, "select_run_select_scan") {
				mylables = []string{"TotalTime", "latencyPct95(μs)"}
				chartDataTest.title += " (Lower is better) "
			} else {
				mylables = global.StringsAppendArrayToArray(mylables, Graph.labels)
			}

			for _, labelReference := range mylables {

				bar := charts.NewBar()

				titleFull := global.ReplaceString(chartDataTest.title, "_", " ") + " " + chartDataTest.dimension
				if strings.Contains(labelReference, "latencyPct95(μs)") {
					titleFull += " (Lower is better)"
				}

				if chartDataTest.prePost == 0 {
					titleFull += " Pre Writes"
				} else {
					titleFull += " Post Writes"
				}
				//general

				bar.SetGlobalOptions(
					charts.WithInitializationOpts(opts.Initialization{
						Width:  "800px",
						Height: "400px",
					}),
					charts.WithLegendOpts(opts.Legend{Width: "90%", Height: "300", Bottom: "-1%", Type: "scroll"}),
					//charts.WithLegendOpts(opts.Legend{Width: "90%", Height: "300", Bottom: "-1%"}),
					charts.WithXAxisOpts(opts.XAxis{Name: "Threads", NameGap: 20, NameLocation: "middle", SplitLine: &opts.SplitLine{Show: opts.Bool(true)}}),
					//charts.WithColorsOpts(opts.Colors{"blue", "orange"}),
					//charts.WithLegendOpts(opts.Legend{Bottom: "0%"}),
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
					charts.WithTitleOpts(opts.Title{Title: titleFull, Subtitle: labelReference}),
					charts.WithYAxisOpts(opts.YAxis{Name: labelReference, NameLocation: "middle", NameGap: 60, AxisLabel: &opts.AxisLabel{Rotate: 0.00, Align: "right"}}),
				)
				for _, chartItemInstance := range chartDataTest.chartItems {
					if chartItemInstance.label == labelReference {
						//log.Debugf("Len items data %d  test %s label %s", len(chartItemInstance.data), chartDataTest.title, chartItemInstance.label)
						bar.SetXAxis(chartDataTest.threads).AddSeries(chartItemInstance.provider, chartItemInstance.data)
					}
				}

				page.AddCharts(bar)

			}

		}
		//log.Debugf("Len items %d", len(chartDataTest.chartItems))

	}

}

/*
Here we will export each chart as a csv set  where the output should be like
Test

	Label
		threads, Provider1-test, provider2-test,...
*/
type csvDat struct {
	name      string
	dimension string
	data      [][]string
}

func (Graph *GraphGenerator) PrintDataCsv() bool {
	//TODO CSV
	currentTime := time.Now()
	csvOutputPath := Graph.configuration.Render.CsvDestinationPath
	if !global.CheckIfPathExists(csvOutputPath) {
		err := os.Mkdir(csvOutputPath, 0755)
		if err != nil {
			log.Errorf("Creating CSV path: %s", err.Error())
		}

	}
	csvFileName := Graph.configuration.Global.TestName + "_" + currentTime.Format("2006-01-02")
	csvFile, err := os.Create(csvOutputPath + csvFileName + ".csv")
	if err != nil {
		log.Errorf("Creating CSV File: %s", err.Error())
	}
	defer csvFile.Close()

	csvLabels := make(map[string]csvDat)

	for _, chartStatTest := range Graph.chartsData {
		labels := []string{}
		providers := []string{}

		if strings.Contains(chartStatTest.title, "select_run_select_scan") {
		}

		csvFile.WriteString(chartStatTest.title + "," + chartStatTest.dimension + "\n")
		csvFile.Sync()

		// we first prepare the objects and the map
		for _, chart := range chartStatTest.chartItems {
			if !slices.Contains(providers, chart.provider) {
				providers = append(providers, chart.provider)
			}

			if !slices.Contains(labels, chart.label) {
				labels = append(labels, chart.label)
				data := make([][]string, len(chartStatTest.threads)+1)

				for i := 0; i < len(data); i++ {
					data[i] = make([]string, chartStatTest.numProviders+1)
					if i > 0 {
						data[i][0] = strconv.Itoa(chartStatTest.threads[i-1])
					} else {
						data[i][0] = "Threads"
					}

				}
				csvLabels[chart.label] = csvDat{chart.label, chartStatTest.dimension, data}
			}
		}
		// we want to have the order of the providers to remain always the same
		sort.Strings(providers)

		// we now fill the data
		providerPosition := 0
		for _, provider := range providers {
			providerPosition++

			//loop also per label
			for _, kLabel := range labels {

				for _, chart := range chartStatTest.chartItems {
					if chart.provider == provider && chart.label == kLabel {
						label := chart.label
						myCsvLabel := csvLabels[label]
						csvData := myCsvLabel.data

						csvData[0][providerPosition] = strings.ReplaceAll(provider, ",", "")

						for i := 0; i < len(chart.data); i++ {
							csvData[i+1][providerPosition] = fmt.Sprintf("%v", chart.data[i].Value)
						}
						myCsvLabel.data = csvData
						log.Debug(csvData)
					}
				}
			}
		}

		//we now flush all the data of the test to file
		for i := 0; i < len(chartStatTest.threads)+1; i++ {
			lineBuffer := bufio.NewWriter(csvFile)
			// we want labels in order so we ue the label array
			for _, kLabel := range labels {
				csvData := csvLabels[kLabel]
				if i == 0 {
					lineBuffer.WriteString(kLabel + ",")
				} else {
					lineBuffer.WriteString(",")
				}
				for ip := 0; ip < (chartStatTest.numProviders + 1); ip++ {
					lineBuffer.WriteString(fmt.Sprintf("%v,", csvData.data[i][ip]))
				}
				lineBuffer.WriteString(",")

			}
			lineBuffer.WriteString("\n")
			lineBuffer.Flush()
			csvFile.Sync()

		}
		csvFile.WriteString("\n\n")
		log.Infof("Label for CSV test: %s  %s", chartStatTest.title, labels)
	}
	return true
}

func (Graph *GraphGenerator) addStatsToPage(page *components.Page) {
	//For each test
	// set global params
	// 	Parse labels
	//	set axis labels based on the label
	//		parse provider
	//			add the data
	for _, chartStatTest := range Graph.chartsStats {
		if !strings.Contains(strings.ToLower(chartStatTest.title), "warmup") {

			//IF test is select scan we onl show totaltime and latency
			mylables := []string{}
			if strings.Contains(chartStatTest.title, "select_run_select_scan") {
				mylables = []string{"TotalTime", "latencyPct95(μs)"}
				chartStatTest.title += " (Lower is better) "
			} else {
				mylables = global.StringsAppendArrayToArray(mylables, Graph.labels)
			}

			for _, labelReference := range mylables {

				bar := charts.NewBar()

				titleFull := global.ReplaceString(chartStatTest.title, "_", " ") + " " + chartStatTest.dimension
				if chartStatTest.prePost == 0 {
					titleFull += " Pre Writes"
				} else {
					titleFull += " Post Writes"
				}
				//general

				bar.SetGlobalOptions(
					charts.WithInitializationOpts(opts.Initialization{
						Width:  "800px",
						Height: "400px",
					}),
					charts.WithLegendOpts(opts.Legend{Width: "90%", Height: "300", Bottom: "-1%", Type: "scroll"}),
					//charts.WithLegendOpts(opts.Legend{Width: "90%", Height: "300", Bottom: "-1%"}),
					charts.WithXAxisOpts(opts.XAxis{Name: "Threads", NameGap: 20, NameLocation: "middle", SplitLine: &opts.SplitLine{Show: opts.Bool(true)}}),
					charts.WithYAxisOpts(opts.YAxis{Name: "Variation %", NameLocation: "middle", NameGap: 60, AxisLabel: &opts.AxisLabel{Rotate: 0.00, Align: "right"}}),
					//charts.WithColorsOpts(opts.Colors{"blue", "orange"}),
					//charts.WithDataZoomOpts(opts.DataZoom{Type:  "slider",Start: 0,End:   50,}),
					//charts.WithDataZoomOpts(opts.DataZoom{Type: "slider"}),
					//charts.WithTitleOpts(opts.Title{Title: chartDataTest.title}),
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
					charts.WithTitleOpts(opts.Title{Title: titleFull, Subtitle: labelReference}),
				)

				for _, chartItemInstance := range chartStatTest.chartItems {
					if chartItemInstance.label == labelReference {
						bar.SetXAxis(chartStatTest.threads).AddSeries(chartItemInstance.provider, chartItemInstance.data)

					}
				}
				page.AddCharts(bar)
			}
		}

	}
}

func (Graph *GraphGenerator) getBarStats(testResult ResultTest, inLabel string) ([]int, []opts.BarData) {
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
		items = append(items, opts.BarData{Value: value.Lerror, Name: value.Label})
		threads = append(threads, value.ThreadNumber)
	}
	return threads, items
}
