package v2

var Responses = make(chan Response, 1000)
var FinishWithCollectOfStatistics = make(chan struct{})

type Response struct {
	StatusCode int
	Endpoint string
	TimeDuration int64
}

type Metrics struct {
	TotalNumberOfSentRequests     int64
	TotalTimeOfWaitingForResponse int64
	MaxTimeForResponse            int64
	EndpointWithTheSlowestResponse string
	MinTimeForResponse            int64
}

func (m *Metrics) Collect() {
	go func() {
		for  resp := range Responses {
			m.TotalNumberOfSentRequests++
			m.TotalTimeOfWaitingForResponse += resp.TimeDuration

			if m.MaxTimeForResponse < resp.TimeDuration {
				m.EndpointWithTheSlowestResponse = resp.Endpoint
				m.MaxTimeForResponse = resp.TimeDuration
			}

			if (m.MinTimeForResponse > resp.TimeDuration) || m.MinTimeForResponse == 0 {
				m.MinTimeForResponse = resp.TimeDuration
			}
		}

		FinishWithCollectOfStatistics <- struct{}{}
	}()
}