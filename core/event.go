package core

func (k *KumaEngine) AddEvent(event interface{}) func() {
	return k.Session.AddHandler(event)
}

func (k *KumaEngine) AddEventOnce(event interface{}) func() {
	return k.Session.AddHandlerOnce(event)
}
