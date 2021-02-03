package debug

import (
	"net/http"
	"os"
	"runtime/pprof"
)

var (
	cpuProfile = "/etc/profile/cpuProfile"
	memProfile = "/etc/profile/memProfile"
	listenAddr = ""
)

type Debug struct {
	cpu *os.File
	mem *os.File
}

func NewDebug() *Debug {
	// Write cpu profile if requested.
	f, err := os.Create(cpuProfile)
	if err != nil {
		panic("Unable to create cpu profile")
	}

	g, err := os.Create(memProfile)
	if err != nil {
		panic("Unable to create mem profile")
	}

	return &Debug{cpu: f, mem: g}
}

func (d *Debug) Run() {
	// Write cpu profile if requested.
	pprof.StartCPUProfile(d.cpu)
	defer d.cpu.Close()
	defer pprof.StopCPUProfile()

	go func() {
		listenAddr := ""
		profileRedirect := http.RedirectHandler("/debug/pprof",
			http.StatusSeeOther)
		http.Handle("/", profileRedirect)
		err := http.ListenAndServe(listenAddr, nil)
		if err != nil {
			panic("Unable to create profile")
		}
	}()

	defer d.cpu.Close()
	defer pprof.StopCPUProfile()

	defer d.mem.Close()
	defer pprof.WriteHeapProfile(d.mem)
}
