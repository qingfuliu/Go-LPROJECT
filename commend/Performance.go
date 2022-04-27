package commend

import "flag"

var (
	CpuProfilePath  string
	HeapProfilePath string
)

func init() {
	//
	flag.StringVar(&CpuProfilePath, "cpuProfile", "./pprof/CpuProfile.pprof", "cpu profile Path")
	flag.StringVar(&HeapProfilePath, "heapProfile", "./pprof/HeapProfile.pprof", "heap profile Path")
}
