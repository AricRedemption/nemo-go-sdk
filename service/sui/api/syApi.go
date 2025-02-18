package api

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
)

var (
	SYSTATE = "0xdc04b5fffe78fae13e967c5943ea6b543637df8955afca1e89a70d0cf5a1a0c2"
	PRICEORACLECONFIG = "0x8dc043ba780bc9f5b4eab09c4e6d82d7af295e5c5ab32be5c27d9933fb02421b"
	MARKETGLOBALCONFIG = "0x9bdde7b16ccaa212b80cb3ae8d644aa1c7f65fd12764ce9bc267fe28de72b54d"
	MARKETSTATECONFIG = "0x876d6f16bb3c7439cd2490b13dbbd249618e5d10bef53a00fac800f60fd86acf"
	YIELDFACTORYCONFIG = "0x0f3e1b1922a2445a4ed5ec936a348cf6bfe50f829b92da0ba9ed3490ae1f1439"
)

func Deposit(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoPackage, coinType, syType string, coinArgument *sui_types.Argument) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoPackage)
	if err != nil {
		return nil, err
	}

	moduleName := "sy"
	functionName := "deposit"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}

	structTag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	typeTag := move_types.TypeTag{
		Struct: structTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, typeTag, syTypeTag)

	var arguments []sui_types.Argument
	versionArgument,err := GetObjectArgument(ptb, client, VERSION, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateArgument,err := GetObjectArgument(ptb, client, SYSTATE, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}
	arguments = append(arguments, versionArgument, *coinArgument, syStateArgument)

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

func SeedLiquidity() {}

func SwapExactPtForSy(ptb *sui_types.ProgrammableTransactionBuilder, blockClient *sui.ISuiAPI, client *client.Client, currentNemoPackage ,pyState, syType, ownerAddress string, nemoPackageList []string, oracleArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(currentNemoPackage)
	if err != nil {
		return nil, err
	}

	moduleName := "market"
	functionName := "swap_exact_pt_for_sy"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	versionArgument,err := GetObjectArgument(ptb, client, VERSION, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateInfo, err := GetObjectFieldByObjectId(client, pyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoPackageList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockClient, client, expectPyPositionTypeList, syType, maturity, ownerAddress)
	if err != nil {
		return nil, err
	}
	if pyPosition == ""{
		return nil, errors.New("without pyPosition！")
	}

	var pyPositionArgument *sui_types.Argument
	argument, err := GetObjectArgument(ptb, client, pyPosition, false, currentNemoPackage, moduleName, functionName)
	if err != nil{
		return nil, err
	}
	pyPositionArgument = &argument

	pyStateArgument,err := GetObjectArgument(ptb, client, pyState, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, MARKETSTATECONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, MARKETGLOBALCONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	ptAmountIn := CreatePureU64CallArg(3)
	ptAmountArgument,err := ptb.Input(ptAmountIn)
	if err != nil {
		return nil, err
	}
	syOut := CreatePureU64CallArg(1)
	syOutArgument,err := ptb.Input(syOut)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	resultArg := &sui_types.Argument{
		NestedResult: &struct {
			Result1 uint16
			Result2 uint16
		}{Result1: *oracleArgument.Result, Result2: 0},
	}

	arguments = append(arguments, versionArgument, ptAmountArgument, syOutArgument, *pyPositionArgument, pyStateArgument, *resultArg, marketGlobalConfigArgument, marketStateArgument, clockArgument)
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

func SyRedeem(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoPackage, coinType, syType string, argument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoPackage)
	if err != nil {
		return nil, err
	}

	moduleName := "sy"
	functionName := "redeem"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	arg0Tag, err := GetStructTag(coinType)
	if err != nil {
		return nil, err
	}
	arg0TypeTag := move_types.TypeTag{
		Struct: arg0Tag,
	}
	arg1Tag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	arg1TypeTag := move_types.TypeTag{
		Struct: arg1Tag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, arg0TypeTag, arg1TypeTag)

	versionArgument,err := GetObjectArgument(ptb, client, VERSION, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	syStateArgument,err := GetObjectArgument(ptb, client, SYSTATE, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, versionArgument, *argument, syStateArgument)
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

func CreatePureU64CallArg(value uint64) sui_types.CallArg {
	// 创建一个8字节的缓冲区
	buf := make([]byte, 8)

	// 使用 binary.LittleEndian 将 uint64 写入字节数组
	binary.LittleEndian.PutUint64(buf, value)

	// 返回构造的 CallArg
	return sui_types.CallArg{
		Pure: &buf,
		Object: nil,  // Pure 类型不需要 Object
	}
}

func SwapExactYtForSy(ptb *sui_types.ProgrammableTransactionBuilder, blockClient *sui.ISuiAPI, client *client.Client, currentNemoPackage ,pyState, syType, ownerAddress string, nemoPackageList []string, oracleArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(currentNemoPackage)
	if err != nil {
		return nil, err
	}

	moduleName := "router"
	functionName := "swap_exact_yt_for_sy"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)

	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	versionArgument,err := GetObjectArgument(ptb, client, VERSION, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateInfo, err := GetObjectFieldByObjectId(client, pyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoPackageList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockClient, client, expectPyPositionTypeList, syType, maturity, ownerAddress)
	if err != nil {
		return nil, err
	}
	if pyPosition == ""{
		return nil, errors.New("without pyPosition！")
	}

	var pyPositionArgument *sui_types.Argument
	argument, err := GetObjectArgument(ptb, client, pyPosition, false, currentNemoPackage, moduleName, functionName)
	if err != nil{
		return nil, err
	}
	pyPositionArgument = &argument

	pyStateArgument,err := GetObjectArgument(ptb, client, pyState, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	yieldFactoryArgument,err := GetObjectArgument(ptb, client, YIELDFACTORYCONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, MARKETSTATECONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, MARKETGLOBALCONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	ytAmountIn := CreatePureU64CallArg(4000000)
	ytAmountArgument,err := ptb.Input(ytAmountIn)
	if err != nil {
		return nil, err
	}
	syOut := CreatePureU64CallArg(10)
	syOutArgument,err := ptb.Input(syOut)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	resultArg := &sui_types.Argument{
		NestedResult: &struct {
			Result1 uint16
			Result2 uint16
		}{Result1: *oracleArgument.Result, Result2: 0},
	}

	arguments = append(arguments, versionArgument, ytAmountArgument, syOutArgument, *pyPositionArgument, pyStateArgument, *resultArg, yieldFactoryArgument, marketGlobalConfigArgument, marketStateArgument, clockArgument)
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

func SwapExactSyForYt(ptb *sui_types.ProgrammableTransactionBuilder, blockClient *sui.ISuiAPI, client *client.Client, currentNemoPackage ,pyState, syType, ownerAddress string, nemoPackageList []string, approxYtOut, netSyTokenization, minYtOut uint64, oracleArgument, coinArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(currentNemoPackage)
	if err != nil {
		return nil, err
	}

	moduleName := "router"
	functionName := "swap_exact_sy_for_yt"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	versionArgument,err := GetObjectArgument(ptb, client, VERSION, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	pyStateInfo, err := GetObjectFieldByObjectId(client, pyState)
	if err != nil{
		return nil, err
	}
	maturity := pyStateInfo["expiry"].(string)

	expectPyPositionTypeList := make([]string, 0)
	for _, pkg := range nemoPackageList{
		expectPyPositionTypeList = append(expectPyPositionTypeList, fmt.Sprintf("%v::py_position::PyPosition", pkg))
	}

	pyPosition,err := GetOwnerObjectByType(blockClient, client, expectPyPositionTypeList, syType, maturity, ownerAddress)
	if err != nil {
		return nil, err
	}
	var pyPositionArgument *sui_types.Argument
	if pyPosition == ""{
		pyPositionArgument, err = InitPyPosition(ptb, client, currentNemoPackage, syType)
		if err != nil{
			return nil, err
		}
	}else {
		argument, err := GetObjectArgument(ptb, client, pyPosition, false, currentNemoPackage, moduleName, functionName)
		if err != nil{
			return nil, err
		}
		pyPositionArgument = &argument
	}

	pyStateArgument,err := GetObjectArgument(ptb, client, pyState, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	yieldFactoryArgument,err := GetObjectArgument(ptb, client, YIELDFACTORYCONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, MARKETGLOBALCONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, MARKETSTATECONFIG, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, currentNemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	approxYtOutCallArg := CreatePureU64CallArg(approxYtOut)
	approxYtOutArgument,err := ptb.Input(approxYtOutCallArg)
	if err != nil {
		return nil, err
	}

	netSyTokenizationCallArg := CreatePureU64CallArg(netSyTokenization)
	netSyTokenizationArgument,err := ptb.Input(netSyTokenizationCallArg)
	if err != nil {
		return nil, err
	}

	minYtOutCallArg := CreatePureU64CallArg(minYtOut)
	minYtOutArgument,err := ptb.Input(minYtOutCallArg)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, versionArgument, minYtOutArgument, approxYtOutArgument, netSyTokenizationArgument, *coinArgument, *oracleArgument, *pyPositionArgument, pyStateArgument, yieldFactoryArgument, marketGlobalConfigArgument, marketStateArgument, clockArgument)
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

func GetApproxYtOutForNetSyInInternal(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoPackage, syType, pyState, marketGlobalConfig, marketState string, netSyIn, minYtOut uint64, oracleArgument *sui_types.Argument) (*sui_types.Argument,error){
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoPackage)
	if err != nil {
		return nil, err
	}

	moduleName := "offchain"
	functionName := "get_approx_yt_out_for_net_sy_in_internal"
	module := move_types.Identifier(moduleName)
	function := move_types.Identifier(functionName)
	syStructTag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	syTypeTag := move_types.TypeTag{
		Struct: syStructTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, syTypeTag)

	pyStateArgument,err := GetObjectArgument(ptb, client, pyState, false, nemoPackage, moduleName, functionName)

	marketGlobalConfigArgument,err := GetObjectArgument(ptb, client, marketGlobalConfig, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	marketStateArgument,err := GetObjectArgument(ptb, client, marketState, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	clockArgument,err := GetObjectArgument(ptb, client, constant.CLOCK, false, nemoPackage, moduleName, functionName)
	if err != nil {
		return nil, err
	}

	netSyInArg := CreatePureU64CallArg(netSyIn)
	netSyInAtgument,err := ptb.Pure(netSyInArg)
	if err != nil {
		return nil, err
	}
	minYtOutArg := CreatePureU64CallArg(minYtOut)
	minYtOutAtgument,err := ptb.Pure(minYtOutArg)
	if err != nil {
		return nil, err
	}

	var arguments []sui_types.Argument
	arguments = append(arguments, netSyInAtgument, minYtOutAtgument, *oracleArgument, pyStateArgument, pyStateArgument, marketStateArgument, marketGlobalConfigArgument, clockArgument)
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