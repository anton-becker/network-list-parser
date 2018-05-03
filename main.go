package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var aaSettings AutoAggregationSettings
	var source = flag.String("src-file", "dump.csv", "csv file in zapret-info format")
	var destination = flag.String("dst-file", "-", "file to store result, value \"-\" means write result to stdout")
	var maskNotation = flag.Bool("mask-notation", false, "use mask notation (\"255.255.255.0\" instead of prefix notation (\"24\")")
	var prefix = flag.String("prefix", "", "prefix each network with string")
	var delimiter = flag.String("delimiter", "/", "insert string between address and mask/prefix")
	var suffix = flag.String("suffix", "", "suffix each network with string")
	var netDelimiter = flag.String("network-delimiter", "\n", "insert string between networks")
	var logStatistic = flag.Bool("log-statistic", true, "print statistics")
	flag.IntVar(&aaSettings.LogMaxPrefix, "log-aggregation-max-prefix", AutoAggregationDefaultSettings.LogMaxPrefix, "log each aggregation into network with prefix up to value; must be positive <=30 or -1 to disable aggregation logging at all")
	flag.UintVar(&aaSettings.IntensiveMinPrefix, "intensive-aggregation-min-prefix", AutoAggregationDefaultSettings.IntensiveMinPrefix, "minimal prefix for intensive auto aggregation; must be <=31 where 31 means disable intensive auto aggregation")
	flag.IntVar(&aaSettings.IntensiveMinNets, "intensive-aggregation-min-nets", AutoAggregationDefaultSettings.IntensiveMinNets, "minimal networks count for intensive auto aggregation; must be >=2")
	flag.Uint64Var(&aaSettings.IntensiveMaxFake, "intensive-aggregation-max-fake-ips", AutoAggregationDefaultSettings.IntensiveMaxFake, "maximum number of fake IPs for intensive auto aggregation")
	flag.Uint64Var(&aaSettings.MaxFake, "aggregation-max-fake-ips", AutoAggregationDefaultSettings.MaxFake, "maximum number of fake IPs for auto aggregation; 0 means disable auto aggregation (but keep intensive auto aggregation)")
	flag.Float64Var(&aaSettings.LoFakePercent, "aggregation-lo-fake-percent", AutoAggregationDefaultSettings.LoFakePercent, "acceptable percent (\"1\"=100%) of fake IPs when aggregating 2 networks; must be >0")
	flag.Float64Var(&aaSettings.HiFakePercent, "aggregation-hi-fake-percent", AutoAggregationDefaultSettings.HiFakePercent, "acceptable percent (\"1\"=100%) of fake IPs when aggregating \"aggregation-hi-fake-percent-nets\" or more networks; must be > aggregation-lo-fake-percent")
	flag.IntVar(&aaSettings.HiFakePercentNets, "aggregation-hi-fake-percent-nets", AutoAggregationDefaultSettings.HiFakePercentNets, "see \"aggregation-hi-fake-percent\" description; must be >2")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Program version: %s\n\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), helpMessage)
		flag.PrintDefaults()
	}
	flag.Parse()
	aaSettings.Validate()

	var src *os.File
	if *source == "-" {
		src = os.Stdin
	} else {
		var err error
		src, err = os.Open(*source)
		if err != nil {
			panic(err.Error())
		}
		defer src.Close()
	}

	//
	//
	//
	nets := ParseIPs(src)
	if *logStatistic {
		log.Printf("File contains: %v networks/IPs with %v IPs coverage.", len(nets), nets.Count())
	}
	nets = NormalizeIPs(nets)
	netsL := len(nets)
	netsC := nets.Count()
	if *logStatistic {
		log.Printf("After normalization: %v networks/IPs with %v IPs coverage.", netsL, netsC)
	}
	if aaSettings.Enabled() {
		AutoAggregate(nets, aaSettings)
		nets = PackNETs(nets)
		if *logStatistic {
			log.Printf("After auto aggregate: %v networks/IPs with %v IPs coverage. Removed %.2v%% records (%v). Add %.2v%% IPs coverage (%v).\n", len(nets), nets.Count(), float32(netsL-len(nets))*100/float32(netsL), netsL-len(nets), float32(nets.Count()-netsC)*100/float32(netsC), nets.Count()-netsC)
		}
	}

	//
	//
	//
	var dst *os.File
	if *destination == "-" {
		dst = os.Stdout
	} else {
		var err error
		dst, err = os.OpenFile(*destination, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err.Error())
		}
		defer dst.Close()
	}

	for i := range nets {
		if i > 0 {
			dst.WriteString(*netDelimiter)
		}
		dst.WriteString(*prefix)
		dst.WriteString(nets[i].IP.String())
		dst.WriteString(*delimiter)
		if *maskNotation {
			dst.WriteString(nets[i].Prefix.Mask().String())
		} else {
			dst.WriteString(nets[i].Prefix.String())
		}
		dst.WriteString(*suffix)
	}
}
