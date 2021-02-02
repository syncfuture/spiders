import React from 'react'
import { Table, Divider, Select, Button, Input, Form, Card } from 'antd';
import { connect, IItemListModelState, Dispatch, Loading } from 'umi';

const { Option } = Select;

interface IPageProps {
  model: IItemListModelState;
  loading: boolean;
  dispatch: Dispatch;
}

class ItemsPage<T extends IPageProps> extends React.Component<T> {
  componentDidMount() {
    this.getItems();
  }

  getItems = () => {
    const { dispatch } = this.props;
    dispatch({
      type: 'itemList/getItems',
    });
  };

  // loadMore = () => {
  //   const { dispatch } = this.props;
  //   dispatch({
  //     type: 'itemList/loadMore',
  //   });
  // };

  onASINChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'itemList/setState',
      payload: {
        asin: e.target.value,
      },
    });
  };

  onItemNoChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'itemList/setState',
      payload: {
        itemNo: e.target.value,
      },
    });
  };

  onStatusChanged = (newValue: number) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'itemList/setState',
      payload: {
        status: newValue,
      },
    });
  };

  search = (values: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'itemList/search',
      payload: {
        ...values
      },
    });
  };

  scrape = (values: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'itemList/scrape',
      payload: {
        ...values
      },
    });
  };

  _columns: any[] = [
    {
      title: 'ASIN',
      dataIndex: 'ASIN',
      defaultSortOrder: 'ascend',
      sorter: (a: any, b: any) => a.ASIN.localeCompare(b.ASIN),
    },
    {
      title: 'ItemNo',
      dataIndex: 'ItemNo',
      sorter: (a: any, b: any) => a.ItemNo.localeCompare(b.ItemNo),
    },
    {
      title: 'Status',
      dataIndex: 'Status',
      sorter: (a: any, b: any) => a.Status - b.Status,
    },
  ];


  render() {
    const { model, loading } = this.props
    return (
      <div>
        <Card>
          <Form
            name="basic"
            layout="inline"
            initialValues={{ status: model.status.toString() }}
            onFinish={this.getItems}
          >
            <Form.Item name="asin">
              <Input placeholder="ASIN" onChange={this.onASINChanged} />
            </Form.Item>

            <Form.Item name="itemNo">
              <Input placeholder="ItemNo" onChange={this.onItemNoChanged} />
            </Form.Item>

            <Form.Item name="status">
              <Select style={{ width: 120 }} onChange={this.onStatusChanged}>
                <Option value="-1">All</Option>
                <Option value="0">Pending</Option>
                <Option value="1">Finished</Option>
                <Option value="2">Error</Option>
              </Select>
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit">Search</Button>
            </Form.Item>
            <Button onClick={this.scrape}>Scrape</Button>
          </Form>
        </Card>
        <Table
          dataSource={model.items}
          columns={this._columns}
          rowKey="ASIN"
          loading={loading}
          pagination={{
            total: model.totalCount,
            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} items`,
          }}
        // footer={() => <div> <Button type="link" onClick={this.loadMore}>Load more...</Button></div>}
        />
      </div>
    );
  }
}

export default connect(({ itemList, loading }: { itemList: IItemListModelState; loading: Loading }) => ({
  model: itemList,
  loading: loading.models.itemList,
}))(ItemsPage as React.ComponentClass<any>);