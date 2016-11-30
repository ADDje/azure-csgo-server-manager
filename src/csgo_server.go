package main

type CsgoServer struct {
	ConfigFile string
	Latency    int  `json:"latency"`
	Port       int  `json:"port"`
	Running    bool `json:"running"`
	Settings   CsgoServerSettings
}

func (f *CsgoServer) Run() error {
	//var err error

	return nil
}

func (f *CsgoServer) Stop() error {

	return nil
}
