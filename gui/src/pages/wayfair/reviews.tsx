import React from 'react'
import { Table, Card, Form, Input, Button, DatePicker } from 'antd';
import moment from 'moment'
import { connect, IReviewListModelState, Loading, Dispatch } from 'umi';
import { CheckCircleTwoTone } from '@ant-design/icons'

interface IPageProps {
  model: IReviewListModelState;
  loading: boolean;
  dispatch: Dispatch;
}

class ReviewsPage<T extends IPageProps> extends React.Component<T> {

  componentDidMount() {
    this.getReviews();
  }

  getReviews = () => {
    const { dispatch } = this.props;

    dispatch({
      type: 'wayfairReviewList/getReviews'
    });
  };
  onSKUChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairReviewList/setState',
      payload: {
        sku: e.target.value,
      },
    });
  };

  onItemNoChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairReviewList/setState',
      payload: {
        itemNo: e.target.value,
      },
    });
  };

  onFromDateChanged = (e: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairReviewList/setState',
      payload: {
        fromDate: e.format("YYYY-MM-DD"),
      },
    });
  };

  onShowSizeChange = (oldSize: number, newSize: number) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairReviewList/setState',
      payload: {
        pageSize: newSize,
      },
    });
  };

  export = () => {
    const { dispatch } = this.props;
    dispatch({
      type: 'wayfairReviewList/export',
    });
  };

  disabledDate = (current: any) => {
    return current < moment().add(-1, "month").add(-1, "d") || current > moment();
  }

  _columns: any[] = [
    {
      title: 'SKU',
      dataIndex: 'SKU',
      defaultSortOrder: 'ascend',
      width: 100,
      sorter: (a: any, b: any) => a.SKU.localeCompare(b.SKU),
    },
    {
      title: 'ItemNo',
      dataIndex: 'CustomerNo',
      width: 100,
      sorter: (a: any, b: any) => a.CustomerNo.localeCompare(b.CustomerNo),
    },
    {
      title: 'StripInfo',
      dataIndex: 'StripInfo',
      sorter: (a: any, b: any) => a.Title.localeCompare(b.Title),
    },
    {
      title: 'Rating',
      dataIndex: 'Rating',
      align: "center",
      width: 80,
      sorter: (a: any, b: any) => a.Rating - b.Rating,
    },
    {
      title: 'IsVerified',
      dataIndex: 'IsVerified',
      align: "center",
      width: 80,
      sorter: (a: any, b: any) => a.IsVerified > b.IsVerified,
      render: (_: any, x: any) => <label>{x.IsVerified ? <CheckCircleTwoTone twoToneColor="#52c41a" /> : ""}</label>,
    },
    {
      title: 'Location',
      dataIndex: 'Location',
      width: 120,
      ellipsis: true,
      sorter: (a: any, b: any) => a.Location.localeCompare(b.Location),
    },
    {
      title: 'CustomerName',
      dataIndex: 'CustomerName',
      width: 150,
      ellipsis: true,
      sorter: (a: any, b: any) => a.CustomerName.localeCompare(b.CustomerName),
    },
    {
      title: 'CreatedOn',
      dataIndex: 'CreatedOn',
      width: 100,
      sorter: (a: any, b: any) => a.CreatedOn.localeCompare(b.CreatedOn),
      render: (_: any, x: any) => <label>{moment(x.CreatedOn).format("MM/DD/YYYY")}</label>,
    },
  ];

  render() {
    const { model, loading } = this.props
    return (
      <div>
        <Card style={{ margin: "5px 0" }}>
          <Form
            layout="inline"
            onFinish={this.getReviews}
            initialValues={{
              fromDate: moment(model.fromDate),
            }}
          >
            <Form.Item name="sku">
              <Input placeholder="SKU" onChange={this.onSKUChanged} />
            </Form.Item>

            <Form.Item name="itemNo">
              <Input placeholder="ItemNo" onChange={this.onItemNoChanged} />
            </Form.Item>

            <Form.Item name="fromDate">
              <DatePicker
                format={"MM/DD/YYYY"}
                disabledDate={this.disabledDate}
                onChange={this.onFromDateChanged}
              />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit">Search</Button>
            </Form.Item>
            <Button htmlType="button" onClick={this.export}>Export</Button>
          </Form>
        </Card>
        <Table
          dataSource={model.reviews}
          columns={this._columns}
          size="small"
          rowKey="WayfairID"
          loading={loading}
          pagination={{
            total: model.totalCount,
            pageSize: model.pageSize,
            onShowSizeChange: this.onShowSizeChange,
            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} items`,
          }}
        />
      </div>
    );
  }
}

export default connect(({ wayfairReviewList, loading }: { wayfairReviewList: IReviewListModelState; loading: Loading }) => ({
  model: wayfairReviewList,
  loading: loading.models.wayfairReviewList,
}))(ReviewsPage as React.ComponentClass<any>);