package shield

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"

	tmtypes "github.com/tendermint/tendermint/types"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "shield").Msg("parsing genesis")

	// Read the genesis state
	var genState shieldtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[shieldtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling shield state: %s", err)
	}

	// Save shield pools
	err = m.saveShieldPools(doc, genState.Pools)
	if err != nil {
		return fmt.Errorf("error while storing shield genesis pools: %s", err)
	}

	// Save the shield providers
	err = m.saveShieldProviders(doc, genState.Providers)
	if err != nil {
		return fmt.Errorf("error while storing shield genesis providers: %s", err)
	}

	// Save shield purchase list
	err = m.savePurchaseList(doc, genState.PurchaseLists)
	if err != nil {
		return fmt.Errorf("error while storing shield genesis purchase list: %s", err)
	}

	// Save shield withdraws
	err = m.saveShieldWithdraws(doc, genState.Withdraws)
	if err != nil {
		return fmt.Errorf("error while storing shield withdraws: %s", err)
	}

	// Save pool params
	err = m.db.SaveShieldPoolParams(types.NewShieldPoolParams(genState.PoolParams, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis shield pool params: %s", err)
	}

	// Save claim proposal params
	err = m.db.SaveShieldClaimProposalParams(types.NewShieldClaimProposalParams(genState.ClaimProposalParams, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing shield claim proposal params: %s", err)
	}
	return nil
}

// saveShieldPools stores the shield pools present inside the given genesis state
func (m *Module) saveShieldPools(doc *tmtypes.GenesisDoc, pools []shieldtypes.Pool) error {
	for _, pool := range pools {
		poolRecord := types.NewShieldPool(pool.Id, "", sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, pool.Shield)), nil, nil,
			pool.Sponsor, pool.SponsorAddr, pool.Description, pool.ShieldLimit, !pool.Active, doc.InitialHeight)
		err := m.db.SaveShieldPool(poolRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

// saveShieldProviders stores the shield providers present inside the given genesis state
func (m *Module) saveShieldProviders(doc *tmtypes.GenesisDoc, providers []shieldtypes.Provider) error {
	for _, provider := range providers {
		providerRecord := types.NewShieldProvider(provider.Address, provider.Collateral.Int64(), provider.DelegationBonded.Int64(),
			provider.Rewards.Native, provider.Rewards.Foreign, provider.TotalLocked.Int64(), provider.Withdrawing.Int64(), doc.InitialHeight)
		err := m.db.SaveShieldProvider(providerRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

// savePurchaseList stores the shield purchase record inside the given genesis state
func (m *Module) savePurchaseList(doc *tmtypes.GenesisDoc, list []shieldtypes.PurchaseList) error {
	for _, purchase := range list {
		for _, entry := range purchase.Entries {
			purchaseRecord := types.NewShieldPurchaseList(entry.PurchaseId, purchase.PoolId, purchase.Purchaser, entry.DeletionTime, entry.ProtectionEndTime,
				entry.ServiceFees.Foreign, entry.ServiceFees.Native, entry.Shield, entry.Description, doc.InitialHeight)
			err := m.db.SaveShieldPurchaseList(purchaseRecord)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// saveShieldWithdraws stores the shield withdraws present inside the given genesis state
func (m *Module) saveShieldWithdraws(doc *tmtypes.GenesisDoc, withdraws []shieldtypes.Withdraw) error {
	for _, withdraw := range withdraws {
		withdrawRecord := types.NewShieldWithdraw(withdraw.Address, withdraw.Amount.Int64(), withdraw.CompletionTime, doc.InitialHeight)
		err := m.db.SaveShieldWithdraw(withdrawRecord)
		if err != nil {
			return err
		}
	}

	return nil
}
