import React from 'react'
import { Layout, Menu } from 'antd';
import { Link, connect, Loading, Dispatch, ILayoutModelState } from 'umi';
import { MenuInfo } from 'rc-menu/lib/interface';
import { AppstoreOutlined, HomeOutlined, CommentOutlined } from '@ant-design/icons';

const { Header, Content } = Layout;
interface IPageProps {
    model: ILayoutModelState;
    loading: boolean;
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
                            <Menu.Item key="/items" icon={<AppstoreOutlined />}>Items</Menu.Item>
                            <Menu.Item key="/reviews" icon={<CommentOutlined/>}>Reviews</Menu.Item>
                            {/* 
                            <SubMenu key="/reviews" icon={<MailOutlined />} title="Reviews">
                                <Menu.ItemGroup key="g1" title="Item 1">
                                    <Menu.Item key="1">Option 1</Menu.Item>
                                    <Menu.Item key="2">Option 2</Menu.Item>
                                </Menu.ItemGroup>
                                <Menu.ItemGroup key="g2" title="Item 2">
                                    <Menu.Item key="3">Option 3</Menu.Item>
                                    <Menu.Item key="4">Option 4</Menu.Item>
                                </Menu.ItemGroup>
                            </SubMenu>
                            <SubMenu key="sub2" icon={<AppstoreOutlined />} title="Navigation Two">
                                <Menu.Item key="5">Option 5</Menu.Item>
                                <Menu.Item key="6">Option 6</Menu.Item>
                                <SubMenu key="sub3" title="Submenu">
                                    <Menu.Item key="7">Option 7</Menu.Item>
                                    <Menu.Item key="8">Option 8</Menu.Item>
                                </SubMenu>
                            </SubMenu>
                            <SubMenu key="sub4" icon={<SettingOutlined />} title="Navigation Three">
                                <Menu.Item key="9">Option 9</Menu.Item>
                                <Menu.Item key="10">Option 10</Menu.Item>
                                <Menu.Item key="11">Option 11</Menu.Item>
                                <Menu.Item key="12">Option 12</Menu.Item>
                            </SubMenu> */}
                        </Menu>
                    </Header>
                    <Content className="content">
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