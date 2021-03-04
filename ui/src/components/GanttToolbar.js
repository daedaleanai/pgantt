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
import { PageHeader, Radio, Checkbox, DatePicker } from 'antd';
import { connect } from 'react-redux';

const { RangePicker } = DatePicker;

class GanttToolbar extends Component {
  state = {
    currentZoom: 'Days',
    showOutsideTimescale: true,
    showClosed: false
  };

  handleZoomChange = (e) => {
    this.setState({
      currentZoom: e.target.value,
    });

    if (this.props.onZoomChange) {
      this.props.onZoomChange(e.target.value);
    }
  }

  onToggleOutsideTimescale = e => {
    this.setState({
      showOutsideTimescale: e.target.checked,
    });

    if (this.props.onToggleOutsideTimescale)  {
      this.props.onToggleOutsideTimescale(e.target.checked);
    }
  };

  onToggleClosed = e => {
    this.setState({
      showClosed: e.target.checked,
    });

    if (this.props.onToggleClosed)  {
      this.props.onToggleClosed(e.target.checked);
    }
  };

  onRangeChange = (date, dateString) => {
    this.setState({
      dateRange: date
    });

    if (this.props.onRangeChange) {
      if (date) {
        this.props.onRangeChange(date[0], date[1]);
      } else {
        this.props.onRangeChange(null, null);
      }
    }
  }

  render() {
    const options = [
      { label: 'Days', value: 'Days' },
      { label: 'Months', value: 'Months' },
    ];

    return (
      <div className="zoom-bar">
        <PageHeader
          ghost={false}
          title={this.props.projectName}
          extra={[
            <Checkbox
              key="45462ce1-2d60-4f3d-8fa5-265a024724c8"
              checked={this.state.showOutsideTimescale}
              onChange={this.onToggleOutsideTimescale}
            >
              Show Tasks Outside of the Timescale
            </Checkbox>,
            <Checkbox
              key="bbf1de5c-1f5b-415e-b759-d0eac641ca30"
              checked={this.state.showClosed}
              onChange={this.onToggleClosed}
            >
              Show Closed Tasks
            </Checkbox>,
            <RangePicker
              key="0ca5b8f8-bb61-423b-9b1e-909cdf4bff83"
              onChange={this.onRangeChange}
              value={this.state.dateRange}
            />,
            <Radio.Group
              key="151ea1b3-52e5-4886-8b09-4f13551f3433"
              options={options}
              onChange={this.handleZoomChange}
              value={this.state.currentZoom}
              optionType="button"
            />,
          ]}
        >
        </PageHeader>
      </div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const proj = state.projects.filter(proj => proj.phid === ownProps.phid);
  return {
    projectName: proj.length !== 0 ? proj[0].name : ""
  };
}

function mapDispatchToProps(dispatch) {
  return {};
}

export default connect(mapStateToProps, mapDispatchToProps)(GanttToolbar);
