package database

import (
	"encoding/json"
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
)

// SaveShieldPool allows to save for the given height the given shieldtypes pool
func (db *Db) SaveShieldPool(pool *types.ShieldPool) error {
	stmt := `
INSERT INTO shield_pool (pool_id, from_address, shield, native_deposit, foreign_deposit, sponsor, sponsor_address, description, shield_limit, pause, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (pool_id) DO UPDATE 
    SET from_address = excluded.from_address, 
	shield = excluded.shield, 
	native_deposit = excluded.native_deposit, 
	foreign_deposit = excluded.foreign_deposit, 
	description = excluded.description, 
	shield_limit = excluded.shield_limit, 
    height = excluded.height
WHERE shield_pool.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		pool.PoolID,
		pool.FromAddress,
		pq.Array(dbtypes.NewDbCoins(pool.Shield)),
		pq.Array(dbtypes.NewDbCoins(pool.NativeDeposit)),
		pq.Array(dbtypes.NewDbCoins(pool.ForeignDeposit)),
		pool.Sponsor,
		pool.SponsorAddr,
		pool.Description,
		pool.ShieldLimit.Int64(),
		pool.Pause,
		pool.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield pool: %s", err)
	}

	return nil
}

// UpdatePoolPauseStatus updates the pool pause status
func (db *Db) UpdatePoolPauseStatus(poolID uint64, pause bool) error {
	stmt := `UPDATE shield_pool SET pause = $1 WHERE pool_id = %2`

	_, err := db.Sql.Exec(stmt, pause, poolID)
	if err != nil {
		return fmt.Errorf("error while updating shield pool pause status: %s", err)
	}

	return nil
}

// UpdatePoolSponsor updates the pool sponsor address
func (db *Db) UpdatePoolSponsor(poolID uint64, sponsor string, sponsorAddress string) error {
	stmt := `UPDATE shield_pool SET sponsor = $1 AND sponsor_address = $2 WHERE pool_id = %3`

	_, err := db.Sql.Exec(stmt, sponsor, sponsorAddress, poolID)
	if err != nil {
		return fmt.Errorf("error while updating shield pool sponsor: %s", err)
	}

	return nil
}

// SaveShieldPurchase allows to save shield purchase for the given height
func (db *Db) SaveShieldPurchase(shield *types.ShieldPurchase) error {
	stmt := `
INSERT INTO shield_purchase (pool_id, from_address, shield, description, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (from_address) DO UPDATE 
    SET pool_id = excluded.pool_id, 
	shield = excluded.shield, 
	description = excluded.description, 
    height = excluded.height
WHERE shield_purchase.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		shield.PoolID,
		shield.FromAddress,
		pq.Array(dbtypes.NewDbCoins(shield.Shield)),
		shield.Description,
		shield.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield purchase: %s", err)
	}

	return nil
}

// SaveShieldPoolParams allows to save shield pool params
func (db *Db) SaveShieldPoolParams(params *types.ShieldPoolParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling shield pool params: %s", err)
	}

	stmt := `
INSERT INTO shield_pool_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE shield_pool_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing shield pool params: %s", err)
	}

	return nil
}

// SaveShieldClaimProposalParams allows to save shield claim proposal params
func (db *Db) SaveShieldClaimProposalParams(params *types.ShieldClaimProposalParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling shield claim proposal params: %s", err)
	}

	stmt := `
INSERT INTO shield_claim_proposal_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE shield_claim_proposal_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing shield claim proposal params: %s", err)
	}

	return nil
}

// SaveShieldProvider allows to save the shield provider for the given height
func (db *Db) SaveShieldProvider(provider *types.ShieldProvider) error {
	stmt := `
INSERT INTO shield_provider (address, collateral, delegation_bonded, native_rewards, 
    foreign_rewards,total_locked, withdrawing, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (address) DO UPDATE 
    SET collateral = excluded.collateral, 
	delegation_bonded = excluded.delegation_bonded, 
	native_rewards = excluded.native_rewards, 
	foreign_rewards = excluded.foreign_rewards, 
	total_locked = excluded.total_locked, 
	withdrawing = excluded.withdrawing, 
    height = excluded.height
WHERE shield_provider.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		provider.Address,
		provider.Collateral,
		provider.DelegationBonded,
		pq.Array(dbtypes.NewDbDecCoins(provider.NativeRewards)),
		pq.Array(dbtypes.NewDbDecCoins(provider.ForeignRewards)),
		provider.TotalLocked,
		provider.Withdrawing,
		provider.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield provider: %s", err)
	}

	return nil
}

// SaveShieldPurchaseList allows to save the shield purchase record for the given height
func (db *Db) SaveShieldPurchaseList(list *types.ShieldPurchaseList) error {
	stmt := `
INSERT INTO shield_purchase_list (purchase_id, pool_id, purchaser, deletion_time, protection_end_time, 
    foreign_service_fees, native_service_fees, shield, description, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (purchase_id) DO UPDATE 
    SET pool_id = excluded.pool_id, 
	purchaser = excluded.purchaser, 
	deletion_time = excluded.deletion_time, 
	protection_end_time = excluded.protection_end_time, 
	foreign_service_fees = excluded.foreign_service_fees, 
	native_service_fees = excluded.native_service_fees,
	shield = excluded.shield, 
	description = excluded.description, 
    height = excluded.height
WHERE shield_purchase_list.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		list.PurchaseID,
		list.PoolID,
		list.Purchaser,
		list.DeletionTime,
		list.ProtectionEndTime,
		pq.Array(dbtypes.NewDbDecCoins(list.ForeignServiceFees)),
		pq.Array(dbtypes.NewDbDecCoins(list.NativeServiceFees)),
		list.Shield.Int64(),
		list.Description,
		list.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield purchase list: %s", err)
	}

	return nil
}

// SaveShieldWithdraw allows to save the shield withdraw for the given height
func (db *Db) SaveShieldWithdraw(withdraw *types.ShieldWithdraw) error {
	stmt := `
INSERT INTO shield_withdraws (address, amount, completion_time, height)
VALUES ($1, $2, $3, $4)`

	_, err := db.Sql.Exec(stmt,
		withdraw.Address,
		withdraw.Amount,
		withdraw.CompletionTime,
		withdraw.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield withdraw: %s", err)
	}

	return nil
}
