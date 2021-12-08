package metric

//type APIServerMetricClient struct {
//	Name string
//	SupportedMetrics map[string]struct{}
//	Command string
//	ServerKey string
//	ClientCert string
//	Endpoint string
//	WorkDir string
//}
//
//func (m APIServerMetricClient) Plot() ([]string,error) {
//	return []string{m.ServerKey,m.Name}, nil
//}
//
//func (m APIServerMetricClient) GetCommandOutput(serverKey,clientCert, endpoint string) (string) {
//	if serverKey == "" {
//		m.ServerKey = EtcdDefaultKeyFile
//	} else {
//		m.ServerKey = serverKey
//	}
//	if clientCert == "" {
//		m.ClientCert = EtcdDefaultClientCert
//	} else {
//		m.ClientCert = clientCert
//	}
//	if endpoint == "" {
//		m.Endpoint = EtcdDefaultEndpoint
//	} else {
//		m.Endpoint = endpoint
//	}
//	curlCmd := "sudo curl -sk"
//	m.Command = fmt.Sprintf("%s --cert %s --key %s %s", curlCmd, m.ClientCert, m.ServerKey, m.Endpoint)
//	return m.Command
//}
//
//func (m APIServerMetricClient) IsKnownMetric(metric string) bool {
//	_, ok := m.SupportedMetrics[metric]
//	return ok
//}
