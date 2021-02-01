import React from 'react'
import { Table } from 'antd';
import moment from 'moment'
import { connect, IReviewListModelState, Loading, Dispatch } from 'umi';

interface IPageProps {
  model: IReviewListModelState;
  loading: boolean;
  dispatch: Dispatch;
}

class ReviewsPage<T extends IPageProps> extends React.Component<T> {

  componentDidMount() {
    const { dispatch } = this.props;

    dispatch({
      type: 'reviewList/getReviews'
    });
  }

  _columns: any[] = [
    {
      title: 'ID',
      dataIndex: 'ID',
      sorter: (a: any, b: any) => a.ID.localeCompare(b.ID),
    },
    {
      title: 'Location',
      dataIndex: 'Location',
      sorter: (a: any, b: any) => a.Location.localeCompare(b.Location),
    },
    {
      title: 'CustomerName',
      dataIndex: 'CustomerName',
      sorter: (a: any, b: any) => a.CustomerName.localeCompare(b.CustomerName),
    },
    // {
    //   title: 'Title',
    //   dataIndex: 'Title',
    //   key: "Title",
    //   sorter: (a: any, b: any) => a.Title.localeCompare(b.Title),
    // },
    {
      title: 'IsVerified',
      dataIndex: 'IsVerified',
      sorter: (a: any, b: any) => a.IsVerified > b.IsVerified,
      render: (_: any, x: any) => <label>{x.IsVerified ? "TRUE" : "FALSE"}</label>,
    },
    {
      title: 'Rating',
      dataIndex: 'Rating',
      sorter: (a: any, b: any) => a.Rating - b.Rating,
    },
    {
      title: 'CreatedOn',
      dataIndex: 'CreatedOn',
      sorter: (a: any, b: any) => a.CreatedOn.localeCompare(b.CreatedOn),
      render: (_: any, x: any) => <label>{moment(x.CreatedOn).format("MM/DD/YYYY")}</label>,
    },
  ];

  render() {
    const { model, loading } = this.props
    return (
      <Table dataSource={model.reviews} columns={this._columns} rowKey="ID" loading={loading} />
    );
  }
}

export default connect(({ reviewList, loading }: { reviewList: IReviewListModelState; loading: Loading }) => ({
  model: reviewList,
  loading: loading.models.reviewList,
}))(ReviewsPage as React.ComponentClass<any>);