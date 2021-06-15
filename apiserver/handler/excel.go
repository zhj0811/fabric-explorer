package handler

import (
	"bytes"
	"mime/multipart"

	"github.com/tealeg/xlsx"
	"github.com/jzsg/fabric-explorer/apiserver/common"
	"github.com/jzsg/fabric-explorer/common/define"
)

func ParsePolicies(fh *multipart.FileHeader) ([]define.Policy, int, error) {
	f, errCode, err := newFileFormFileHeader(fh)
	if err != nil {
		return nil, errCode, err
	}
	var res []define.Policy
	for _, sheet := range f.Sheets {
		for rowIndex, row := range sheet.Rows {
			//跳过第一行表头信息
			if len(row.Cells) != 11 {
				return nil, common.ExcelFormatErr, common.ErrExcelFormat
			}
			if rowIndex == 0 {
				continue
			}
			policy := define.Policy{
				Number:   row.Cells[0].Value,
				Type:     row.Cells[1].Value,
				Insured:  row.Cells[2].Value,
				USCC:     row.Cells[3].Value,
				StartAt:  row.Cells[4].Value,
				ExpireAt: row.Cells[5].Value,
				Insurer:  row.Cells[7].Value,
				Amount:   row.Cells[8].Value,
				Premium:  row.Cells[9].Value,
				Rate:     row.Cells[10].Value,
				Content:  row.Cells[11].Value,
			}
			pooled := row.Cells[6].Value
			if pooled == "是" {
				policy.Pooled = true
			} else if pooled != "否" {
				return nil, common.ExcelFormatErr, common.ErrExcelFormat
			}

			res = append(res, policy)
		}
	}
	return res, common.Success, nil
}

func ParseServices(fh *multipart.FileHeader) ([]define.Service, int, error) {
	f, errCode, err := newFileFormFileHeader(fh)
	if err != nil {
		return nil, errCode, err
	}
	var services []define.Service
	for _, sheet := range f.Sheets {
		for rowIndex, row := range sheet.Rows {
			//跳过第一行表头信息
			if len(row.Cells) != 5 {
				return nil, common.ExcelFormatErr, common.ErrExcelFormat
			}
			if rowIndex == 0 {
				continue
			}
			service := define.Service{
				Number:       row.Cells[0].Value,
				Insured:      row.Cells[1].Value,
				Date:         row.Cells[2].Value,
				Organization: row.Cells[3].Value,
				Type:         row.Cells[4].Value,
			}
			services = append(services, service)
		}
	}
	return services, common.Success, nil
}

func newFileFormFileHeader(fh *multipart.FileHeader) (*xlsx.File, int, error) {
	//if !strings.HasSuffix(fh.Filename, ".xlsx") {
	//	return common.ExcelFormatErr, common.ErrExcelFormat
	//}

	var buf bytes.Buffer
	file, err := fh.Open()
	if err != nil {
		return nil, common.OpenFileErr, err
	}
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, common.OpenFileErr, err
	}
	//bytes := buf.Bytes()
	f, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		return nil, common.OpenFileErr, err
	}
	return f, common.Success, nil
}
