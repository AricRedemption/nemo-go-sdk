package api

import (
	"errors"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
	"nemo-go-sdk/service/sui/common/models"
)

func GetPriceVoucherFromXOracle(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "oracle"
	functionName := "get_price_voucher_from_x_oracle"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}

	structTag, err := GetStructTag(nemoConfig.UnderlyingCoinType)
	if err != nil {
		return nil, err
	}
	typeTag := move_types.TypeTag{
		Struct: structTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag, typeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	scallopVersionCallArg,err := GetObjectArg(client, SCALLOPVERSION, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	scallopMarketCallArg,err := GetObjectArg(client, SCALLOPMARKETOBJECT, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	clockCallArg,err := GetObjectArg(client, constant.CLOCK, false, nemoConfig.NemoContract, "oracle", "get_price_voucher_from_x_oracle")
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: scallopVersionCallArg}, sui_types.CallArg{Object: scallopMarketCallArg}, sui_types.CallArg{Object: syStateCallArg}, sui_types.CallArg{Object: clockCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func GetPriceVoucherFromVolo(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "oracle"
	functionName := "get_price_voucher_from_volo"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	nativePoolCallArg,err := GetObjectArg(client, nemoConfig.NativePool, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	metadataCallArg,err := GetObjectArg(client, nemoConfig.Metadata, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: nativePoolCallArg}, sui_types.CallArg{Object: metadataCallArg}, sui_types.CallArg{Object: syStateCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func GetPriceVoucherFromSpring(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "oracle"
	functionName := "get_price_voucher_from_spring"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	lstInfoCallArg,err := GetObjectArg(client, nemoConfig.LstInfo, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: lstInfoCallArg}, sui_types.CallArg{Object: syStateCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func GetPriceVoucherFromAftermath(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "oracle"
	functionName := "get_price_voucher_from_aftermath"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	lstInfoCallArg,err := GetObjectArg(client, nemoConfig.LstInfo, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: lstInfoCallArg}, sui_types.CallArg{Object: syStateCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func GetPriceVoucherFromHasui(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoConfig.NemoContract)
	if err != nil {
		return nil, err
	}

	moduleName := "oracle"
	functionName := "get_price_voucher_from_haSui"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(nemoConfig.SyCoinType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	priceOracleCallArg,err := GetObjectArg(client, nemoConfig.PriceOracle, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	haedalStakeingCallArg,err := GetObjectArg(client, nemoConfig.HaedalStakeing, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateCallArg,err := GetObjectArg(client, nemoConfig.SyState, false, nemoConfig.NemoContract, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: priceOracleCallArg}, sui_types.CallArg{Object: haedalStakeingCallArg}, sui_types.CallArg{Object: syStateCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackageId,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func GetPriceVoucher(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoConfig *models.NemoConfig) (*sui_types.Argument,error){
	if constant.IsScallopCoin(nemoConfig.CoinType){
		return GetPriceVoucherFromXOracle(ptb, client, nemoConfig)
	}else if constant.IsVSui(nemoConfig.CoinType){
		return GetPriceVoucherFromVolo(ptb, client, nemoConfig)
	}else if constant.IsSpringSui(nemoConfig.CoinType){
		return GetPriceVoucherFromSpring(ptb, client, nemoConfig)
	}else if constant.IsAfSui(nemoConfig.CoinType) {
		return GetPriceVoucherFromAftermath(ptb, client, nemoConfig)
	}else if constant.IsHaSui(nemoConfig.CoinType) {
		return GetPriceVoucherFromHasui(ptb, client, nemoConfig)
	}
	return nil, errors.New("coinType oracle not support！")
}