package miner

import (
	"errors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/dline"
	"golang.org/x/xerrors"
)

// Returns deadline-related calculations for a deadline in some proving period and the current epoch.
func NewDeadlineInfo(periodStart abi.ChainEpoch, deadlineIdx uint64, currEpoch abi.ChainEpoch) *dline.Info {
	return dline.NewInfo(periodStart, deadlineIdx, currEpoch, WPoStPeriodDeadlines, WPoStProvingPeriod, WPoStChallengeWindow, WPoStChallengeLookback, FaultDeclarationCutoff)
}

func QuantSpecForDeadline(di *dline.Info) builtin.QuantSpec {
	return builtin.NewQuantSpec(WPoStProvingPeriod, di.Last())
}

// FindSector returns the deadline and partition index for a sector number.
// It returns an error if the sector number is not tracked by deadlines.
func FindSector(store adt.Store, deadlines *Deadlines, sectorNum abi.SectorNumber) (uint64, uint64, error) {
	for dlIdx := range deadlines.Due {
		dl, err := deadlines.LoadDeadline(store, uint64(dlIdx))
		if err != nil {
			return 0, 0, err
		}

		partitions, err := adt.AsArray(store, dl.Partitions, DeadlinePartitionsAmtBitwidth)
		if err != nil {
			return 0, 0, err
		}
		var partition Partition

		partIdx := uint64(0)
		stopErr := errors.New("stop")
		err = partitions.ForEach(&partition, func(i int64) error {
			found, err := partition.Sectors.IsSet(uint64(sectorNum))
			if err != nil {
				return err
			}
			if found {
				partIdx = uint64(i)
				return stopErr
			}
			return nil
		})
		if err == stopErr {
			return uint64(dlIdx), partIdx, nil
		} else if err != nil {
			return 0, 0, err
		}

	}
	return 0, 0, xerrors.Errorf("sector %d not due at any deadline", sectorNum)
}
