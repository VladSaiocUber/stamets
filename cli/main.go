package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/vladsaiocuber/stamets"
	"golang.org/x/exp/constraints"
)

func main() {
	var dir string
	var pta, cg bool
	flag.StringVar(&dir, "dir", os.Getenv("PWD"), "Target directory.")
	flag.BoolVar(&pta, "pta", false, "Aggregate PTA results.")
	flag.BoolVar(&cg, "cg", false, "Aggregate call graph results.")
	flag.Parse()

	if pta {
		ptas := stamets.AggregatePTAResults(dir)

		PrintSeries(
			"PTA Duration",
			stamets.MakeSeries(func(m stamets.PTAMetrics) time.Duration {
				return m.Duration
			}, ptas...))
		PrintSeries(
			"PTA P50 points-to-set size",
			stamets.MakeSeries(func(m stamets.PTAMetrics) int {
				return m.PointsToSetSizeP50
			}, ptas...))
		PrintSeries(
			"PTA P90 points-to-set size",
			stamets.MakeSeries(func(m stamets.PTAMetrics) int {
				return m.PointsToSetSizeP90
			}, ptas...))
		PrintSeries(
			"PTA P99 points-to-set size",
			stamets.MakeSeries(func(m stamets.PTAMetrics) int {
				return m.PointsToSetSizeP99
			}, ptas...))
		PrintSeries(
			"PTA Max points-to-set size",
			stamets.MakeSeries(func(m stamets.PTAMetrics) int {
				return m.PointsToSetSizeMax
			}, ptas...))
	}

	if cg {
		cgs := stamets.AggregateCallGraphResults(dir)

		PrintSeries(
			"Call graph number of functions",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.Functions
			}, cgs...))
		PrintSeries(
			"P50 in-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.InDegreeP50
			}, cgs...))
		PrintSeries(
			"P90 in-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.InDegreeP90
			}, cgs...))
		PrintSeries(
			"P99 in-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.InDegreeP99
			}, cgs...))
		PrintSeries(
			"Max in-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.InDegreeMax
			}, cgs...))
		PrintSeries(
			"P50 out-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.OutDegreeP50
			}, cgs...))
		PrintSeries(
			"P90 out-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.OutDegreeP90
			}, cgs...))
		PrintSeries(
			"P99 out-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.OutDegreeP99
			}, cgs...))
		PrintSeries(
			"Max out-degree",
			stamets.MakeSeries(func(m stamets.CallGraphMetrics) int {
				return m.OutDegreeMax
			}, cgs...))
	}
}

func PrintSeries[T constraints.Ordered](name string, s stamets.Series[T]) {
	fmt.Println(name+" aggregate metrics over", len(s), "instances:")
	fmt.Println("- P50:", s.P50())
	fmt.Println("- P90:", s.P90())
	fmt.Println("- P99:", s.P99())
	fmt.Println("- Max:", s.Max())
	fmt.Println("- Mode:", s.Mode())
}
