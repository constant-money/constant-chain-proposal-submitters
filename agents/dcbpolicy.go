package agents

import (
	"proposalsubmitters/entities"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
)

func buildCrowdsalesSellBond(
	burnAmount uint64,
	constantPrice uint64,
	blockHeight uint64,
	bonds []*entities.DCBBondInfo,
) ([]component.SaleData, error) {
	// Check if there's a bond that can reduce burnAmount of Constant
	var bondToSell *entities.DCBBondInfo
	for _, b := range bonds {
		// TODO(0xbunyip): pick price between b.Price and b.BuyBackPrice
		if b.Price > 0 && b.Price*b.Amount >= burnAmount {
			bondToSell = b
			break
		}
	}
	if bondToSell == nil {
		return nil, nil
	}

	sale := component.SaleData{
		EndBlock:         blockHeight + 1000,
		BuyingAsset:      common.ConstantID,
		BuyingAmount:     burnAmount,
		DefaultBuyPrice:  constantPrice,
		SellingAsset:     bondToSell.BondID,
		SellingAmount:    burnAmount / bondToSell.Price,
		DefaultSellPrice: bondToSell.Price,
	}
	return []component.SaleData{sale}, nil
}

func buildTradeSellBond(
	burnAmount uint64,
	blockHeight uint64,
	bonds []*entities.DCBBondInfo,
) ([]*component.TradeBondWithGOV, error) {
	// Check if there's a bond that can reduce burnAmount of Constant
	var bondToSell *entities.DCBBondInfo
	for _, b := range bonds {
		if b.Maturity < blockHeight && b.BuyBack*b.Amount >= burnAmount {
			bondToSell = b
			break
		}
	}
	if bondToSell == nil {
		return nil, nil
	}

	trade := &component.TradeBondWithGOV{
		BondID: &bondToSell.BondID,
		Amount: burnAmount / bondToSell.BuyBack,
		Buy:    false,
	}
	return []*component.TradeBondWithGOV{trade}, nil
}

func buildSpendReserve(
	burnAmount uint64,
	constantPrice uint64,
	blockHeight uint64,
	dr *DataRequester,
) (map[common.Hash]*component.SpendReserveData, error) {
	// TODO(@0xbunyip): choose between ETH and USD
	price, err := dr.AssetPrice(common.ETHAssetID)
	if err != nil {
		return nil, err
	}
	reserve := map[common.Hash]*component.SpendReserveData{
		common.ETHAssetID: &component.SpendReserveData{
			EndBlock:        blockHeight + 1000,
			ReserveMinPrice: price,
			Amount:          burnAmount,
		},
	}
	return reserve, nil
}