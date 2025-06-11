package sequence

type Call struct {
	Caller string
	Callee string
	Note   string
}

type edgeKey struct {
	Caller string
	Callee string
	Note   string
}
