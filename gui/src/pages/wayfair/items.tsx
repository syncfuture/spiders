import React from 'react'
import { Table, Select, Button, Input, Form, Card, Progress, Space } from 'antd';
import { connect, IWayfairItemListModelState, Dispatch, Loading } from 'umi';

const { Option } = Select;

interface IPageProps {
  model: IWayfairItemListModelState;
  loading: boolean;
  dispatch: Dispatch;
  timeout: NodeJS.Timeout;
}

class ItemsPage<T extends IPageProps> extends React.Component<T> {
  componentDidMount() {
    this.getItems();
  }

  getItems = () => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairItemList/getItems',
    });
  };

  // loadMore = () => {
  //   const { dispatch } = this.props;
  //   dispatch({
  //     type: 'wayfairItemList/loadMore',
  //   });
  // };

  onSKUChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairItemList/setState',
      payload: {
        sku: e.target.value,
      },
    });
  };

  onItemNoChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairItemList/setState',
      payload: {
        itemNo: e.target.value,
      },
    });
  };

  onStatusChanged = (newValue: number) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairItemList/setState',
      payload: {
        status: newValue,
      },
    });
  };

  scrape = (values: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairItemList/scrape',
      payload: {
        ...values
      },
    });

    // this.updateScrapeStatus();
  };

  cancel = () => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairItemList/cancel',
    });
  };

  // updateScrapeStatus = () => {
  //   const { dispatch, model } = this.props;
  //   let { timeout } = this.props;

  //   dispatch({
  //     type: 'wayfairItemList/getScrapeStatus',
  //   });

  //   if (!model.running) {
  //     clearTimeout(timeout);
  //   } else {
  //     timeout = setTimeout(this.updateScrapeStatus, 1000);
  //   }
  // };

  onShowSizeChange = (oldSize: number, newSize: number) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairItemList/getScrapeStatus',
    });
  };

  _columns: any[] = [
    {
      title: 'SKU',
      dataIndex: 'SKU',
      defaultSortOrder: 'ascend',
      sorter: (a: any, b: any) => a.SKU.localeCompare(b.SKU),
    },
    {
      title: 'ItemNOs',
      dataIndex: 'ItemNOs',
      sorter: (a: any, b: any) => a.ItemNo.localeCompare(b.ItemNo),
    },
    {
      title: 'Status',
      dataIndex: 'Status',
      sorter: (a: any, b: any) => a.Status - b.Status,
    },
  ];


  render() {
    const { model, loading } = this.props;
    const { Current, TotalCount } = model.scrapeStatus;
    let percent = TotalCount > 0 ? Math.floor((Current / TotalCount) * 100) : 0;

    return (
      <div>
        {/* <Card>
          <Breadcrumb>
            <Breadcrumb.Item>Home</Breadcrumb.Item>
            <Breadcrumb.Item>Wayfair</Breadcrumb.Item>
            <Breadcrumb.Item>Items</Breadcrumb.Item>
          </Breadcrumb>
        </Card> */}

        <Card style={{ margin: "5px 0" }}>
          <Form
            layout="inline"
            initialValues={{ status: model.status.toString() }}
            onFinish={this.getItems}
          >
            <Form.Item name="sku">
              <Input placeholder="SKU" onChange={this.onSKUChanged} />
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
            <Space>
              <Button onClick={this.scrape} disabled={model.running}>Scrape</Button>
              <Button htmlType="button" disabled={!model.running} onClick={this.cancel}>Cancel</Button>
            </Space>
          </Form>
          {model.running ? <Progress percent={percent} size="small" /> : null}
        </Card>
        <Table
          dataSource={model.items}
          columns={this._columns}
          size="small"
          rowKey="SKU"
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

export default connect(({ wayfairItemList, loading }: { wayfairItemList: IWayfairItemListModelState; loading: Loading }) => ({
  model: wayfairItemList,
  loading: loading.models.wayfairItemList,
}))(ItemsPage as React.ComponentClass<any>);