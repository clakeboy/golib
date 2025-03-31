package explain

import (
	"fmt"
	"strings"
	"testing"
)

var content = `
func (r *ResData) ToJsonString() string {
	data := r.ToJson()
	if data == nil {
		return ""
	}
	return string(data)
}

		<React.Fragment>
                <CTable absolute x={'20px'} y={'116px'} width={'808px'} height={'392px'} scroll={true} bordered={true} headerTheme={'light'} hover={true} select={true} sm={true} fontSm={true}
                        onSelectPage={(page) => {
							if (page < 1) {
								console.log(1111);
							}
                            this.getData(page);
							return <div>sdfsdf</div>
                        }}
                        page={this.state.page}
                        dataCount={this.state.data_count} data={this.state.table_data} onClick={this.selectClickHandler}>
                    <Table.Header text={'Seq No.'} field={'ams_query_master_id'} width={'52px'} />
                    <Table.Header text={'Master B/L No.'} field={'master_bl_no'} width={'136px'} />
                    <Table.Header text={'Status'} field={'status_desc'} width={'80px'} />
                    <Table.Header text={'Mode'} field={'mode_code'} width={'36px'} />
                    <Table.Header text={'ETA'} field={'arrival_date'} width={'80px'} onFormat={(val, row) => { return val; }} />
                    <Table.Header text={'Code'} field={'last_disposition_code'} width={'36px'} />
                    <Table.Header text={'Qty'} field={'manifest_quantity'} width={'56px'} />
                    <Table.Header text={'HBLs'} field={'hbl_count'} width={'40px'} />
                    <Table.Header text={'Departure'} field={'departure_port_name'} width={'76px'} />
                    <Table.Header text={'Discharge'} field={'arrival_port_name'} width={'76px'} />
                    <Table.Header text={'Vessel / Flight Info'} field={'conveyance_name'} width={'136px'} />
                    <Table.Header text={'Voy No.'} field={'manifest_number'} width={'64px'} />
                    <Table.Header text={'Firms'} field={'FIRMS_code'} width={'48px'} />
                    <Table.Header text={'ISF'} field={'isf_sw'} width={'28px'} />
                    <Table.Header text={'Last Request'} field={'last_modified_date'} width={'108px'} />
                    <Table.Header text={'ETD'} field={'departure_date'} width={'80px'} onFormat={(val, row) => { return val;  }} />
                    <Table.Header text={'File No.'} field={'file_doc_number'} width={'96px'} />
                    <Table.Header text={'I.T. No.'} field={'inbond_number'} width={'80px'} />
                    <Table.Header text={'I.T. Origin'} field={'inbond_origin_port_name'} width={'88px'} />
                    <Table.Header text={'I.T. Dest'} field={'inbond_dest_port_name'} width={'88px'} />
                    <Table.Header text={'Last Request'} field={'modified_by_code'} width={'80px'} />
                    <Table.Header text={'Group'} field={'group_code'} width={'64px'} />
                    <Table.Header text={'Br'} field={'branch_code'} width={'36px'} />
                </CTable>
                <ConditionForm ref={c => this.condition = c} onChange={this.dataChangeHandler}>
                    <Input field={'master_bl_no'} size={'xs'} absolute={true} x={'100px'} y={'48px'} width={'112px'} data={this.state.form.master_bl_no} />
                    <Label sm={true} text={'Master No.'} absolute={true} x={'20px'} y={'48px'} width={'72px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />
                    <Input size={'xs'} absolute={true} x={'100px'} y={'76px'} width={'112px'} comboData={[{ 'status_id': 0, 'status_desc': 'Open' }]} data={this.state.form.status_desc}
                           combo={{
                               searchColumn: 'status_desc',
                               filterColumns: [
                                   { field: 'status_desc', width: '100px' },
                                   'status_id'
                               ]
                           }} />
                    <Label sm={true} text={'Status'} absolute={true} x={'20px'} y={'76px'} width={'72px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />
                    <Input field={'arrival_date$d1'} size={'xs'} absolute={true} x={'616px'} y={'48px'} width={'100px'} calendar={true} data={this.state.form.arrival_date$d1} />
                    <Label sm={true} text={'E T A'} absolute={true} x={'540px'} y={'48px'} width={'68px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />
                    <Input field={'arrival_date$d2'} size={'xs'} absolute={true} x={'724px'} y={'48px'} width={'100px'} calendar={true} data={this.state.form.arrival_date$d2} />
                    <Input field={'MOT'} size={'xs'} absolute={true} x={'316px'} y={'48px'} width={'48px'} data={this.state.form.MOT} />
                    <Label sm={true} text={'Mode'} absolute={true} x={'232px'} y={'48px'} width={'76px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />
                    <Input field={'last_disposition_code'} size={'xs'} absolute={true} x={'316px'} y={'76px'} width={'48px'} data={this.state.form.last_disposition_code} />
                    <Label sm={true} text={'Event Code'} absolute={true} x={'232px'} y={'76px'} width={'76px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />
                    <Input field={'last_modified_date$d1'} size={'xs'} absolute={true} x={'616px'} y={'76px'} width={'100px'} calendar={true} data={this.state.form.last_modified_date$d1} />
                    <Label sm={true} text={'Last Event'} absolute={true} x={'540px'} y={'76px'} width={'68px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />
                    <Input field={'last_modified_date$d2'} size={'xs'} absolute={true} x={'724px'} y={'76px'} width={'100px'} calendar={true} data={this.state.form.last_modified_date$d2} />
                    <Input field={'group_code'} size={'xs'} absolute={true} x={'444px'} y={'48px'} width={'76px'} data={this.state.form.group_code} />
                    <Label sm={true} text={'Group'} absolute={true} x={'384px'} y={'48px'} width={'52px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />
                    <Input field={'branch_code'} size={'xs'} absolute={true} x={'444px'} y={'76px'} width={'76px'} data={this.state.form.branch_code} />
                    <Label sm={true} text={'Branch'} absolute={true} x={'384px'} y={'76px'} width={'52px'} height={'20px'} color={'rgb(0,0,0)'} backColor={'rgb()'} />

                    <Button size={'xs'} absolute={true} x={'20px'} y={'8px'} width={'52px'} height={'24px'} onClick={
                        this.queryHandler
                    }>Query</Button>
                    <Button size={'xs'} absolute={true} x={'76px'} y={'8px'} width={'52px'} height={'24px'} onClick={
                        this.clearHandler
                    }>Clear</Button>
                </ConditionForm>
                <WModal ref={c => this.modal = c} fade center={true} />
            </React.Fragment>`

func TestJSX(t *testing.T) {
	nodes := JSXString(content)
	compileStr := content
	for _, n := range nodes {
		if n.NodeName == "Button" && !n.EndNode {
			n.properties = append(n.properties, &JSXNodeProperty{
				Name:  "field",
				Value: "{'newbutton'}",
			})
			compileStr = strings.Replace(compileStr, string(n.Original), n.Compile(), 1)
		}
	}

	fmt.Println(compileStr)
}

func TestJSXProperty(t *testing.T) {
	str := []byte("calke")
	fmt.Printf("%X", str)
}
