package kafka

// BindPending registers all handlers in reg to the pending worker.
func BindPending(p *PendingWorker, reg *Registry) {
	for taskType, fn := range reg.All() {
		p.RegisterHandler(taskType, fn) // 这里依赖 RegisterHandler 接受的是 Handler/HandlerFunc
	}
}
