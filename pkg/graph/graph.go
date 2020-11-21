package graph

import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"strconv"
	"time"
)


func parseInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

func convertToSeparateDataSlices(data []GraphPoint) ([]time.Time, []float64) {
	var xvalues []time.Time
	var yvalues []float64

	for _, m := range data {
		xvalues = append(xvalues, m.timestamp)
    yvalues = append(yvalues, float64(m.errorCount))
	}
	return xvalues, yvalues
}

func DrawChart(data map[string][]GraphPoint) {

	// list of stroke colours.
	strokeColours := []drawing.Color {chart.ColorWhite,chart.ColorBlue,chart.ColorCyan,chart.ColorGreen,chart.ColorRed,chart.ColorOrange,chart.ColorYellow,chart.ColorBlack,chart.ColorLightGray,chart.ColorAlternateBlue,chart.ColorAlternateGreen,chart.ColorAlternateGray,chart.ColorAlternateYellow,chart.ColorAlternateLightGray}
	seriesList := []chart.Series{}
	maxY := 0.0

	strokeNumber := 0
	for k,v := range data {
		fmt.Printf("graphing %s\n", k)
		cpuXValues, cpuYValues := convertToSeparateDataSlices(v)

		// get maxY for graphing later.
		for _,y := range cpuYValues {
			if y > maxY {
				maxY = y
			}
		}

		errorSeries := chart.TimeSeries{
			Name: k,
			Style: chart.Style{
				Show:        true,
				StrokeColor: strokeColours[strokeNumber],
				StrokeWidth: chart.Disabled,
				DotWidth:5,

				//FillColor:   chart.ColorBlue.WithAlpha(100),

			},
			XValues: cpuXValues,
			YValues: cpuYValues,

		}
		strokeNumber++
		if strokeNumber >= 14 {
			strokeNumber = 0
		}

		seriesList = append(seriesList, errorSeries)
	}

	fmt.Printf("length of series is %d\n", len(seriesList))

	fmt.Printf("maxY is %f\n", maxY)

	// give it some headroom
	maxY += 20.0

	fmt.Printf("adj maxY is %f\n", maxY)

	graph := chart.Chart{
		Width:  2280,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},
		},
		Canvas:chart.Style {
  		  FillColor:chart.ColorBlack,
		},
		YAxis: chart.YAxis{
			Name:      "Errors",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
			Range: &chart.ContinuousRange{
			  Min: 0,
			  Max: maxY,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("2006-01-02.15-04-05"),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			//GridLines: releases(),
		},
		Series: seriesList,
	}

	graph.Elements = []chart.Renderable{chart.LegendLeft(&graph)}
	//, chart.Style{ FillColor:drawing.Color{R: 50, G: 120, B: 203, A: 255}})}

	// output to file... ?

	f, err := os.Create("graphit.png")
	if err != nil {
		fmt.Printf("unable to create file\n")
	}

	graph.Render(chart.PNG, f)
}

