
import React, { Component } from 'react';
import { Menu } from 'antd';
import { MailOutlined, AppstoreOutlined, SettingOutlined } from '@ant-design/icons';

const { SubMenu } = Menu;

const styles = {
  logo: {
    float: 'left',
    fontSize: 'x-large'
  }
};

class PGanttNav extends Component {
  state = {
    current: 'mail',
  };

  handleClick = e => {
    console.log('click ', e);
    this.setState({ current: e.key });
  };

  render() {
    const { current } = this.state;
    return (
      <Menu
        onClick={this.handleClick}
        selectedKeys={[current]}
        mode="horizontal"
        theme='dark'
      >
        <Menu.Item key="pgantt" style={styles.logo}>
          PGantt
        </Menu.Item>
        <SubMenu key="projects" icon={<SettingOutlined />} title="Projects">
          <Menu.Item key="setting:1">Option 1</Menu.Item>
          <Menu.Item key="setting:2">Option 2</Menu.Item>
          <Menu.Item key="setting:3">Option 3</Menu.Item>
          <Menu.Item key="setting:4">Option 4</Menu.Item>
        </SubMenu>
      </Menu>
    );
  }
}

export default PGanttNav;
