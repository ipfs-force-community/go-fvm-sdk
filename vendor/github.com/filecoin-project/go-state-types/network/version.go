package network

import "math"

// Enumeration of network upgrades where actor behaviour can change (without necessarily
// vendoring and versioning the whole actor codebase).
type Version uint

const (
	Version0  = Version(iota) // genesis    (specs-actors v0.9.3)
	Version1                  // breeze     (specs-actors v0.9.7)
	Version2                  // smoke      (specs-actors v0.9.8)
	Version3                  // ignition   (specs-actors v0.9.11)
	Version4                  // actors v2  (specs-actors v2.0.3)
	Version5                  // tape       (specs-actors v2.1.0)
	Version6                  // kumquat    (specs-actors v2.2.0)
	Version7                  // calico     (specs-actors v2.3.2)
	Version8                  // persian    (post-2.3.2 behaviour transition)
	Version9                  // orange     (post-2.3.2 behaviour transition)
	Version10                 // trust      (specs-actors v3.0.1)
	Version11                 // norwegian  (specs-actors v3.1.0)
	Version12                 // turbo      (specs-actors v4.0.0)
	Version13                 // hyperdrive (specs-actors v5.0.1)
	Version14                 // chocolate (specs-actors v6.0.0)
	Version15                 // ???
	Version16                 // ???
	Version17                 // ???

	// VersionMax is the maximum version number
	VersionMax = Version(math.MaxUint32)
)
