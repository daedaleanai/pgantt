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
import moment from 'moment'

import {
  dateRangeSet, zoomSet, showTasksOutsideTimescaleSet, showTasksClosedSet,
  showTasksUnscheduledSet
} from '../actions/settings';

// Trickery required to have the date picker allow picking a range.
const { RangePicker } = DatePicker;

// Controls for the Gantt chart allowing basic filtering of the tasks, defining
// a date range and zooming.
class GanttToolbar extends Component {

  // Handles the user setting the date range.
  onRangeChange = (dateRange, dateString) => {
    let start = null;
    let end = null;
    if (dateRange != null) {
      start = dateRange[0];
      end = dateRange[1];
    }
    this.props.dateRangeSet(start, end);
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
              checked={this.props.showTasksOutsideTimescale}
              onChange={(e) => this.props.showTasksOutsideTimescaleSet(e.target.checked)}
            >
              Show Tasks Outside of the Timescale
            </Checkbox>,
            <Checkbox
              key="e0b1a9db-14a1-4b10-81ce-45f0d454895a"
              checked={this.props.showTasksUnscheduled}
              onChange={(e) => this.props.showTasksUnscheduledSet(e.target.checked)}
            >
              Show Unscheduled Tasks
            </Checkbox>,

            <Checkbox
              key="bbf1de5c-1f5b-415e-b759-d0eac641ca30"
              checked={this.props.showTasksClosed}
              onChange={(e) => this.props.showTasksClosedSet(e.target.checked)}
            >
              Show Closed Tasks
            </Checkbox>,
            <RangePicker
              key="0ca5b8f8-bb61-423b-9b1e-909cdf4bff83"
              onChange={this.onRangeChange}
              value={this.props.dateRange}
            />,
            <Radio.Group
              key="151ea1b3-52e5-4886-8b09-4f13551f3433"
              options={options}
              onChange={(e) => this.props.zoomSet(e.target.value)}
              value={this.props.zoom}
              optionType="button"
            />,
          ]}
        >
        </PageHeader>
      </div>
    );
  }
}

// Builds props for the `GanttToolbar` component out of the Redux store state.
function mapStateToProps(state, ownProps) {
  const proj = state.projects.filter(proj => proj.phid === ownProps.phid);
  return {
    projectName: proj.length !== 0 ? proj[0].name : "",
    showTasksOutsideTimescale: state.settings.showTasksOutsideTimescale,
    showTasksUnscheduled: state.settings.showTasksUnscheduled,
    showTasksClosed: state.settings.showTasksClosed,
    dateRange: [moment(state.settings.startDate), moment(state.settings.endDate)],
    zoom: state.settings.zoom,
  };
}

// Create functions that dispatch actions to the Redux store.
function mapDispatchToProps(dispatch) {
  return {
    showTasksOutsideTimescaleSet: (setting) => dispatch(showTasksOutsideTimescaleSet(setting)),
    showTasksUnscheduledSet: (setting) => dispatch(showTasksUnscheduledSet(setting)),
    showTasksClosedSet: (setting) => dispatch(showTasksClosedSet(setting)),
    dateRangeSet: (start, end) => dispatch(dateRangeSet(start, end)),
    zoomSet: (zoom) => dispatch(zoomSet(zoom))
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(GanttToolbar);
