package listeners

// listeners defines what can be registered with machinery
// TODO: make this abide to as many different listeners as possible i.e grpc or http servers
type Listener interface {
	Run()
}
