package simulated

func (fvmSimulator *FvmSimulator) GasCharge(_ string, _ uint64) error {
	return nil
}

func (fvmSimulator *FvmSimulator) GasAvailable(_ string, _ uint64) (uint64, error) {
	return 0, nil
}
