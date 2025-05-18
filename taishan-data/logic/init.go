package logic

func InitJob() {
	go ConsumerSamplingData()
	go ConsumerTransferData()
}
