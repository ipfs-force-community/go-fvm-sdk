package simulated

func (fvmSimulator *FvmSimulator) ChargeGas(_ string, _ uint64) error {
	return nil
}

func (fvmSimulator *FvmSimulator) AvailableGas(_ string, _ uint64) (uint64, error) {
	return 0, nil
}
