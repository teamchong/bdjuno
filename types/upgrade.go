package types

// UpgradeParams contains the info of the on-chain upgrade of the x/upgrade module
type UpgradeParams struct {
	BinaryVersion string
	UpgradeInfo   string
	UpgradeHeight string
	UpgradeStatus string
}

// NewUpgradeParams allows to build a new UpgradeParams instance
func NewUpgradeParams(
	binaryVersion string,
	upgradeInfo string,
	upgradeHeight string,
	upgradeStatus string) UpgradeParams {
	return UpgradeParams{
		BinaryVersion: binaryVersion,
		UpgradeInfo:   upgradeInfo,
		UpgradeHeight: upgradeHeight,
		UpgradeStatus: upgradeStatus,
	}
}
