import React from 'react'
import { Card } from 'antd';

interface IPageProps {
  model: any;
}

class IndexPage<T extends IPageProps> extends React.Component<T> {

  render() {
    return (
      <Card style={{ margin: 5 }}>Building...</Card>
    );
  }
}

export default IndexPage