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
      title: 'ID',
      dataIndex: 'reviewId',
      defaultSortOrder: 'descend',
      width: 100,
      sorter: (a: any, b: any) => a.reviewId - b.reviewId,
    },
    {
      title: 'SKU',
      dataIndex: 'sku',
      width: 100,
      sorter: (a: any, b: any) => a.sku.localeCompare(b.sku),
    },
    {
      title: 'Rating',
      dataIndex: 'ratingStars',
      align: "center",
      width: 80,
      sorter: (a: any, b: any) => a.ratingStars - b.ratingStars,
    },
    {
      title: 'IsVerified',
      dataIndex: 'hasVerifiedBuyerStatus',
      align: "center",
      width: 80,
      sorter: (a: any, b: any) => a.hasVerifiedBuyerStatus > b.hasVerifiedBuyerStatus,
      render: (_: any, x: any) => <label>{x.hasVerifiedBuyerStatus ? <CheckCircleTwoTone twoToneColor="#52c41a" /> : ""}</label>,
    },
    {
      title: 'Language',
      dataIndex: 'languageCode',
      width: 120,
      ellipsis: true,
      sorter: (a: any, b: any) => a.languageCode.localeCompare(b.languageCode),
    },
    {
      title: 'Reviewer',
      dataIndex: 'reviewerName',
      width: 150,
      ellipsis: true,
      sorter: (a: any, b: any) => a.reviewerName.localeCompare(b.reviewerName),
    },
    {
      title: 'CreatedOn',
      dataIndex: 'createdOn',
      width: 100,
      sorter: (a: any, b: any) => a.createdOn.localeCompare(b.createdOn),
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
          rowKey="reviewId"
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