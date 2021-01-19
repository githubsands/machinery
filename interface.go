package machinery

// Machine defines the common order of operations of using a machine
// i.e Notify(), Register(), Run(), WaitForExit(), WaitForGroup are ran from first
// to last
type Machine interface {
	Notify()

	Register()

	Run()

	WaitForExit()

	WaitForGroup()
}
