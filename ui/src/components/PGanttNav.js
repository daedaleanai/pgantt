//------------------------------------------------------------------------------
// Copyright (C) 2021 Daedalean AG
//
// This file is part of PGantt.
//
// PGantt is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 2 of the License, or
// (at your option) any later version.
//
// PGantt is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PGantt.  If not, see <https://www.gnu.org/licenses/>.
//------------------------------------------------------------------------------

import React, { Component } from 'react';
import { Menu, message } from 'antd';
import { MailOutlined, AppstoreOutlined, SettingOutlined } from '@ant-design/icons';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import { projectsSet } from '../actions/projects';
import { projectsGet } from '../utils/api';

const { SubMenu } = Menu;

const styles = {
  logo: {
    float: 'left',
    fontSize: 'x-large'
  }
};

class PGanttNav extends Component {
  state = {
    current: 'pgantt',
  };

  handleClick = e => {
    console.log('click ', e);
    this.setState({ current: e.key });
  };

  componentDidMount() {
    projectsGet()
      .then(data => this.props.projectsSet(data.Data))
      .catch(msg => message.error(msg.toString()));
  }

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
          <Link to='/'>
            PGantt
          </Link>
        </Menu.Item>
        <SubMenu key="projects" icon={<SettingOutlined />} title="Projects">
          {this.props.projects.map(project => (
            <Menu.Item key={project.Phid}>
              <Link to={'/project/' + project.Phid}>
                {project.Name}
              </Link>
            </Menu.Item>
          ))}
        </SubMenu>
      </Menu>
    );
  }
}

function mapStateToProps(state) {
  return {
    ...state
  };
}

function mapDispatchToProps(dispatch) {
  return {
    projectsSet: (data) => dispatch(projectsSet(data))
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(PGanttNav);
