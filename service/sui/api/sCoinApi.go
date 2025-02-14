package api

import (
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
)

var (
	SCALLOPPACKAGE = "0x80ca577876dec91ae6d22090e56c39bc60dce9086ab0729930c6900bc4162b4c"
	SCOINTREASURY = "0x5c1678c8261ac9eec024d4d630006a9f55c80dc0b1aa38a003fcb1d425818c6b"

	SCALLOPVERSION = "0x07871c4b3c847a0f674510d4978d5cf6f960452795e8ff6f189fd2088a3f6ac7"
	SCALLOPMARKETOBJECT = "0xa757975255146dc9686aa823b7838b507f315d704f428cbadad2f4ea061939d9"
	SCALLOPMINTPACKAGE = "0x3fc1f14ca1017cff1df9cd053ce1f55251e9df3019d728c7265f028bb87f0f97"
)

func MintSCoin(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, coinType, underlyingCoinType string, coinArgument *sui_types.Argument) (*sui_types.Argument,error) {
	scaPackage, err := sui_types.NewObjectIdFromHex(SCALLOPPACKAGE)
	if err != nil {
		return nil, err
	}

	module := move_types.Identifier("s_coin_converter")
	function := move_types.Identifier("mint_s_coin")
	sCoinStructTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: sCoinStructTag,
	}
	underlyingCoinStructTag, err := GetStructTag(underlyingCoinType)
	if err != nil {
		return nil, err
	}
	type2Tag := move_types.TypeTag{
		Struct: underlyingCoinStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag, type2Tag)

	scaTreasuryCallArg,err := GetObjectArg(client, SCOINTREASURY, false, SCALLOPPACKAGE, "s_coin_converter", "mint_s_coin")
	if err != nil {
		return nil, err
	}
	scaTreasuryArgument, err := ptb.Input(sui_types.CallArg{Object: scaTreasuryCallArg})
	if err != nil {
		return nil, err
	}

	marketCoin, err := Mint(ptb, client, underlyingCoinType, coinArgument)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument

	arguments = append(arguments, scaTreasuryArgument, *marketCoin)
	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *scaPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}

func Mint(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, underlyingCoinType string, coinArgument *sui_types.Argument) (*sui_types.Argument,error) {
	nemoPackage, err := sui_types.NewObjectIdFromHex(SCALLOPMINTPACKAGE)
	if err != nil {
		return nil, err
	}
	module := move_types.Identifier("mint")
	function := move_types.Identifier("mint")

	underlyingCoinStructTag, err := GetStructTag(underlyingCoinType)
	if err != nil {
		return nil, err
	}
	type1Tag := move_types.TypeTag{
		Struct: underlyingCoinStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, type1Tag)

	versionCallArg,err := GetObjectArg(client, SCALLOPVERSION, false, SCALLOPMINTPACKAGE, "mint", "mint")
	if err != nil {
		return nil, err
	}

	marketObjectCallArg,err := GetObjectArg(client, SCALLOPMARKETOBJECT, false, SCALLOPMINTPACKAGE, "mint", "mint")
	if err != nil {
		return nil, err
	}

	clockCallArg,err := GetObjectArg(client, constant.CLOCK, false, SCALLOPMINTPACKAGE, "mint", "mint")
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: versionCallArg}, sui_types.CallArg{Object: marketObjectCallArg})
	var arguments []sui_types.Argument
	for _, v := range callArgs {
		argument, err := ptb.Input(v)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}

	arguments = append(arguments, *coinArgument)
	clockArgument,err := ptb.Input(sui_types.CallArg{Object: clockCallArg})
	arguments = append(arguments, clockArgument)

	command := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *nemoPackage,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return &command, nil
}
