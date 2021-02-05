import React from 'react'
import { Layout, Menu, Card, Breadcrumb, } from 'antd';
import { connect, Loading, Dispatch, ILayoutModelState, Link } from 'umi';
import { MenuInfo } from 'rc-menu/lib/interface';
import { AppstoreOutlined, HomeOutlined, CommentOutlined, AmazonOutlined, ContainerOutlined } from '@ant-design/icons';

const { SubMenu } = Menu;
const { Header, Content } = Layout;
const _breadcrumbNameMap: { [key: string]: string; } = {
    '/amazon': 'Amazon',
    '/amazon/reviews': 'Reviews',
    '/amazon/items': 'Items',
    '/wayfair': 'Wayfair',
    '/wayfair/reviews': 'Reviews',
    '/wayfair/items': 'Items',
};

interface IPageProps {
    model: ILayoutModelState;
    loading: boolean;
    location: any;
    dispatch: Dispatch;
}

class AppLayout extends React.Component<IPageProps> {
    handleClick = (e: MenuInfo) => {
        const { dispatch } = this.props;
        dispatch({
            type: 'layout/navigate', payload: {
                path: e.key,
            }
        });
    };

    buildBreadcrumbs = () => {
        const { location } = this.props;
        const pathSnippets = location.pathname.split('/').filter((i: any) => i);
        const extraBreadcrumbItems = pathSnippets.map((_: any, index: any) => {
            const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;
            const link = (url: string) => {
                if (url == "/amazon" || url == "/wayfair") {
                    return <label>{_breadcrumbNameMap[url]}</label>;
                } else {
                    return <Link to={url}>{_breadcrumbNameMap[url]}</Link>;
                }
            }
            return (
                <Breadcrumb.Item key={url}>
                    { link(url)}
                </Breadcrumb.Item>
            );
        });

        const breadcrumbItems = [
            <Breadcrumb.Item key="home">
                <Link to="/">Home</Link>
            </Breadcrumb.Item>,
        ].concat(extraBreadcrumbItems);

        return breadcrumbItems;
    }

    render() {
        const { model } = this.props;
        return (
            <Layout style={{ minHeight: '100vh' }}>
                <Layout className="layout">
                    <Header className="header">
                        <Menu
                            onClick={this.handleClick}
                            defaultSelectedKeys={model.selectedPathKeys}
                            mode="horizontal"
                            theme="dark"
                        >
                            <Menu.Item key="/" icon={<HomeOutlined />}>Overview</Menu.Item>
                            <SubMenu key="Amazon" icon={<AmazonOutlined />} title="Amazon">
                                <Menu.Item key="/amazon/items" icon={<ContainerOutlined />}>Items</Menu.Item>
                                <Menu.Item key="/amazon/reviews" icon={<CommentOutlined />}>Reviews</Menu.Item>
                            </SubMenu>
                            <SubMenu key="Wayfair" icon={<AppstoreOutlined />} title="Wayfair">
                                <Menu.Item key="/wayfair/items" icon={<ContainerOutlined />}>Items</Menu.Item>
                                <Menu.Item key="/wayfair/reviews" icon={<CommentOutlined />}>Reviews</Menu.Item>
                            </SubMenu>
                        </Menu>
                    </Header>
                    <Content className="content">
                        <Card>
                            <Breadcrumb>{this.buildBreadcrumbs()}</Breadcrumb>
                        </Card>
                        {this.props.children}
                    </Content>
                </Layout>
            </Layout>
        );
    }
}

export default connect(({ layout, loading }: { layout: ILayoutModelState; loading: Loading }) => ({
    model: layout,
    loading: loading.models.layout,
}))(AppLayout);