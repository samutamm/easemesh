package installation

import (
	"fmt"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/pkg/errors"
)

type InstallStage interface {
	Do(*installbase.StageContext, Installation) error
	Clear(*installbase.StageContext) error
}

// Installation represents installing or clearing infrastructure
// components for the EaseMesh
type Installation interface {
	DoInstallStage(*installbase.StageContext) error
	ClearResource(*installbase.StageContext)
}

type installation struct {
	stages []InstallStage
	step   int
}

func New(stages ...InstallStage) Installation {
	return &installation{stages: stages, step: 0}
}

func (i *installation) DoInstallStage(context *installbase.StageContext) error {
	if i.step >= len(i.stages) {
		return nil
	}
	current := i.step
	i.step++
	return i.stages[current].Do(context, i)
}

func (i *installation) ClearResource(context *installbase.StageContext) {
	for _, f := range context.ClearFuncs {
		err := f(context)
		if err != nil {
			common.OutputErrorf("clear resource error:%s", err)
		}
	}
}

type InstallFunc func(*installbase.StageContext) error

type HookFunc InstallFunc
type ClearFunc HookFunc
type PreCheckFunc HookFunc

type DescribeFunc func(*installbase.StageContext, installbase.InstallPhase) string

func Wrap(preCheckFunc HookFunc, installFunc InstallFunc, clearFunc HookFunc, description DescribeFunc) InstallStage {
	return &baseInstallStage{preCheck: PreCheckFunc(preCheckFunc), installFunc: installFunc, clearFunc: ClearFunc(clearFunc), description: description}
}

type baseInstallStage struct {
	preCheck    PreCheckFunc
	installFunc InstallFunc
	clearFunc   ClearFunc
	description DescribeFunc
}

var _ InstallStage = &baseInstallStage{}

func (b *baseInstallStage) Do(context *installbase.StageContext, install Installation) error {
	fmt.Printf("%s\n", b.description(context, installbase.BeginPhase))
	if b.preCheck != nil {
		if err := b.preCheck(context); err != nil {
			return errors.Wrap(err, "pre check installation condition failed")
		}
	}
	err := b.installFunc(context)
	context.ClearFuncs = append(context.ClearFuncs, b.clearFunc)
	if err != nil {
		return errors.Wrap(err, "invoke install func error")
	}

	fmt.Printf("Install successfully end, following resource are deployed successfully: %s\n", b.description(context, installbase.EndPhase))
	return install.DoInstallStage(context)
}

func (b *baseInstallStage) Clear(context *installbase.StageContext) error {
	// Do nothing
	if b.clearFunc != nil {
		return b.clearFunc(context)
	}
	return nil
}