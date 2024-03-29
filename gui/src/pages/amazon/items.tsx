import React from 'react'
import { Table, Select, Button, Input, Form, Card } from 'antd';
import { connect, IAmazonItemListModelState, Dispatch, Loading } from 'umi';

const { Option } = Select;

interface IPageProps {
  model: IAmazonItemListModelState;
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
      type: 'amazonItemList/getItems',
    });
  };

  // loadMore = () => {
  //   const { dispatch } = this.props;
  //   dispatch({
  //     type: 'amazonItemList/loadMore',
  //   });
  // };

  onASINChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'amazonItemList/setState',
      payload: {
        asin: e.target.value,
      },
    });
  };

  onItemNoChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'amazonItemList/setState',
      payload: {
        itemNo: e.target.value,
      },
    });
  };

  onStatusChanged = (newValue: number) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'amazonItemList/setState',
      payload: {
        status: newValue,
      },
    });
  };

  scrape = (values: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'amazonItemList/scrape',
      payload: {
        ...values
      },
    });
  };

  onShowSizeChange = (oldSize: number, newSize: number) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'amazonItemList/setState',
      payload: {
        pageSize: newSize,
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
        {/* <Card>
          <Breadcrumb>
            <Breadcrumb.Item>Home</Breadcrumb.Item>
            <Breadcrumb.Item>Amazon</Breadcrumb.Item>
            <Breadcrumb.Item>Items</Breadcrumb.Item>
          </Breadcrumb>
        </Card> */}

        <Card style={{ margin: "5px 0" }}>
          <Form
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
                <Option value="">All</Option>
                <Option value="0">Pending</Option>
                <Option value="1">Finished</Option>
                <Option value="-1">Error</Option>
                <Option value="404">NotFound</Option>
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
          size="small"
          rowKey="ASIN"
          loading={loading}
          pagination={{
            total: model.totalCount,
            pageSize: model.pageSize,
            onShowSizeChange: this.onShowSizeChange,
            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} items`,
          }}
        // footer={() => <div> <Button type="link" onClick={this.loadMore}>Load more...</Button></div>}
        />
      </div>
    );
  }
}

export default connect(({ amazonItemList, loading }: { amazonItemList: IAmazonItemListModelState; loading: Loading }) => ({
  model: amazonItemList,
  loading: loading.models.amazonItemList,
}))(ItemsPage as React.ComponentClass<any>);