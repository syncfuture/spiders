import React from 'react'

interface IPageProps {
  model: any;
}

class IndexPage<T extends IPageProps> extends React.Component<T> {

  render() {
    return (
      <div>Index</div>
    );
  }
}

export default IndexPage