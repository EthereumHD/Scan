package plotter

import (
	"fmt"
	"go-web3/dto"
	"go-web3/providers"
)

type Plotter struct {
	provider providers.ProviderInterface
}

func NewPlotter(provider providers.ProviderInterface) *Plotter {
	plotter := new(Plotter)
	plotter.provider = provider
	return plotter
}

func (p *Plotter) Start() (bool, error) {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_start", params)

	if err != nil {
		fmt.Printf("start plot , http req err : %s\n", err.Error())
		return false, err
	}

	response, err := pointer.ToBoolean()
	if err != nil {
		fmt.Printf("start plot , http rsp err : %s\n", err.Error())
		return false, err
	}
	return response, nil
}

func (p *Plotter) Stop() (bool, error) {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_stop", params)

	if err != nil {
		fmt.Printf("stop plot , http req err : %s\n", err.Error())
		return false, err
	}

	response, err := pointer.ToBoolean()
	if err != nil {
		fmt.Printf("stop plot , http rsp err : %s\n", err.Error())
		return false, err
	}
	return response, nil
}

func (p *Plotter) Ploting() (bool, error) {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_ploting", params)

	if err != nil {
		fmt.Printf("get ploting , http req err : %s\n", err.Error())
		return false, err
	}

	response, err := pointer.ToBoolean()
	if err != nil {
		fmt.Printf("get ploting , http rsp err : %s\n", err.Error())
		return false, err
	}
	return response, nil
}

func (p *Plotter) GetPlotdataDir() (string, error) {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_getPlotdataDir", params)

	if err != nil {
		fmt.Printf("GetPlotdataDir , http req err : %s\n", err.Error())
		return "", err
	}

	response, err := pointer.ToString()
	if err != nil {
		fmt.Printf("GetPlotdataDir , http rsp err %s\n", err.Error())
		return "", err
	}
	return response, nil
}

func (p *Plotter) GetSeed() (string, error) {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_getSeed", params)

	if err != nil {
		fmt.Printf("GetSeed , http req err : %s\n", err.Error())
		return "", err
	}

	response, err := pointer.ToString()
	if err != nil {
		fmt.Printf("GetSeed , http rsp err : %s\n", err.Error())
		return "", err
	}
	return response, nil
}

func (p *Plotter) SetPlotdataDir(dir string) error {
	params := []string{}
	params = append(params, dir)
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_setPlotdataDir", params)

	if err != nil {
		fmt.Printf("SetPlotdataDir , http req err : %s\n", err.Error())
		return err
	}

	return nil
}

func (p *Plotter) SetSeed(addr string) error {
	params := []string{}
	params = append(params, addr)
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_setSeed", params)

	if err != nil {
		fmt.Printf("SetSeed , http req err : %s\n", err.Error())
		return err
	}

	return nil
}

func (p *Plotter) CommitWork(start uint64, quantity uint64) (string, error) {
	params := []uint64{}
	params = append(params, start)
	params = append(params, quantity)
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_commitWork", params)

	if err != nil {
		fmt.Printf("CommitWork , http req err : %s\n", err.Error())
		return "", err
	}

	response, err := pointer.ToString()
	if err != nil {
		fmt.Printf("CommitWork , http rsp err %s\n", err.Error())
		return "", err
	}

	return response, nil
}

func (p *Plotter) GetTaskIds() ([]string, error) {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_getTaskIds", params)

	if err != nil {
		fmt.Printf("GetTaskIds , http req err : %s\n", err.Error())
		return nil, err
	}

	response, err := pointer.ToStringArray()
	if err != nil {
		fmt.Printf("GetTaskIds , http rsp err %s\n", err.Error())
		return nil, err
	}
	return response, nil
}

func (p *Plotter) GetTaskProgress(tid string) (float64, error) {
	params := []string{}
	params = append(params, tid)
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_getTaskProgress", params)

	if err != nil {
		fmt.Printf("GetTaskProgress , http req err : %s\n", err.Error())
		return 0, err
	}

	response, err := pointer.ToFloat()
	if err != nil {
		fmt.Printf("GetTaskProgress , http rsp err %s\n", err.Error())
		return 0, err
	}
	return response, nil
}

func (p *Plotter) RemovePlotDataById(id string) error {
	params := []string{}
	params = append(params, id)
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_removePlotDataById", params)

	if err != nil {
		fmt.Printf("RemovePlotDataById , http req err : %s\n", err.Error())
		return err
	}
	return nil
}

func (p *Plotter) GetPlotDatalist() ([]string, error) {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_getPlotDatalist", params)

	if err != nil {
		fmt.Printf("GetPlotDatalist , http req err : %s\n", err.Error())
		return nil, err
	}
	response, err := pointer.ToStringArray()
	if err != nil {
		fmt.Printf("GetPlotDatalist , http rsp err %s\n", err.Error())
		return nil, err
	}

	return response, nil
}

func (p *Plotter) ClearPlotData() error {
	params := []string{}
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(&pointer, "plotter_clearPlotData", params)

	if err != nil {
		fmt.Printf("ClearPlotData , http req err : %s\n", err.Error())
		return err
	}
	return nil
}
