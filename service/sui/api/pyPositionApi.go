package api

import (
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"nemo-go-sdk/service/sui/common/constant"
	"strings"
)

var (
	VERSION = "0x4000b5c20e70358a42ae45421c96d2f110817d75b80df30dad5b5d4f1fdad6af"
	PYSTATE = "0x4f3cc83141fa233ee733074f90c8b955cb76030b4ae119bd5d40f9bbace97cde"
)

func InitPyPosition(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, nemoPackage, syType string) (*sui_types.Argument,error) {
	nemoPackageId, err := sui_types.NewObjectIdFromHex(nemoPackage)
	if err != nil {
		return nil, err
	}

	module := move_types.Identifier("py")
	function := move_types.Identifier("init_py_position")
	structTag, err := GetStructTag(syType)
	if err != nil {
		return nil, err
	}
	typeTag := move_types.TypeTag{
		Struct: structTag,
	}
	typeArguments := make([]move_types.TypeTag, 0)
	typeArguments = append(typeArguments, typeTag)

	versionCallArg,err := GetObjectArg(client, VERSION, false, nemoPackage, "py", "init_py_position")
	if err != nil {
		return nil, err
	}

	pyStateCallArg,err := GetObjectArg(client, PYSTATE, false, nemoPackage, "py", "init_py_position")
	if err != nil {
		return nil, err
	}

	clockCallArg,err := GetObjectArg(client, constant.CLOCK, false, nemoPackage, "py", "init_py_position")
	if err != nil {
		return nil, err
	}

	callArgs := make([]sui_types.CallArg, 0)
	callArgs = append(callArgs, sui_types.CallArg{Object: versionCallArg}, sui_types.CallArg{Object: pyStateCallArg}, sui_types.CallArg{Object: clockCallArg})
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

func GetObjectArgument(ptb *sui_types.ProgrammableTransactionBuilder, client *client.Client, shareObject string, isCoin bool, contractPackage, module, function string) (sui_types.Argument, error){
	arg, err := GetObjectArg(client, shareObject, isCoin, contractPackage, module, function)
	if err != nil{
		return sui_types.Argument{}, err
	}
	return ptb.Input(sui_types.CallArg{Object: arg})
}

func GetObjectArg(client *client.Client, shareObject string, isCoin bool, contractPackage, module, function string) (*sui_types.ObjectArg, error) {
	hexObject, err := sui_types.NewObjectIdFromHex(shareObject)
	if err != nil {
		return nil, err
	}
	sourceObjectData, err := GetObjectMetadata(client, shareObject)
	if err != nil {
		return nil, err
	}

	if !isCoin && sourceObjectData.Data.Owner.AddressOwner == nil{
		return &sui_types.ObjectArg{
			SharedObject: &struct {
				Id                   sui_types.ObjectID
				InitialSharedVersion sui_types.SequenceNumber
				Mutable              bool
			}{
				Id:                   *hexObject,
				InitialSharedVersion: *sourceObjectData.Data.Owner.Shared.InitialSharedVersion,
				Mutable:              GetObjectMutable(client, *sourceObjectData.Data.Type, contractPackage, module, function),
			},
		}, nil
	}
	return &sui_types.ObjectArg{
		ImmOrOwnedObject: &sui_types.ObjectRef{
			ObjectId: sourceObjectData.Data.ObjectId,
			Version: sourceObjectData.Data.Version.Uint64(),
			Digest: sourceObjectData.Data.Digest,
		},
	}, nil
}

func GetStructTag(syType string) (*move_types.StructTag, error) {
	elements := strings.Split(syType, "::")
	if len(elements) != 3 {
		return nil, errors.New("syType invalid！")
	}

	addr, err := sui_types.NewAddressFromHex(elements[0])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("init syType address error: %v", err))
	}

	structTag := &move_types.StructTag{
		Address:    *addr,
		Module:     move_types.Identifier(elements[1]),
		Name:       move_types.Identifier(elements[2]),
		TypeParams: []move_types.TypeTag{},
	}

	return structTag, nil
}
