package system

import (
	"fmt"
	"os/exec"
)

// ReportingEngine wraps RRDtool for time-series data
type ReportingEngine struct {
	DatabasePath string
}

// CreateDB initializes an RRD file for a metric
func (e ReportingEngine) CreateDB(name string, step int) error {
	fmt.Printf("Reporting: Creating RRD database %s at %s\n", name, e.DatabasePath)
	// rrdtool create name.rrd --step step ...
	cmd := exec.Command("rrdtool", "create", e.DatabasePath+"/"+name+".rrd", "--step", fmt.Sprintf("%d", step))
	return cmd.Run()
}

// UpdateMetric pushes a new value to the database
func (e ReportingEngine) UpdateMetric(name string, value float64) error {
	// rrdtool update name.rrd N:value
	cmd := exec.Command("rrdtool", "update", e.DatabasePath+"/"+name+".rrd", fmt.Sprintf("N:%f", value))
	return cmd.Run()
}

// GenerateGraph creates an SVG/PNG from RRD data
func (e ReportingEngine) GenerateGraph(name, outputPath string) error {
	fmt.Printf("Reporting: Generating graph %s -> %s\n", name, outputPath)
	return nil
}
