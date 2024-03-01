package dataObjects

import (
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"math/rand"
	"os"
	global "sysbench-graph-creator/internal/global"
	"time"
)

//https://github.com/go-echarts/go-echarts

type GraphGenerator struct {
	configuration global.Configuration
	producers     []Producer
}

func (Graph *GraphGenerator) Init(inConfig global.Configuration, inProducers []Producer) {
	Graph.producers = inProducers
	Graph.configuration = inConfig
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
