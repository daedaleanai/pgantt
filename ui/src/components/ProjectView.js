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
import { connect } from 'react-redux';

import Gantt from './Gantt';
import GanttToolbar from './GanttToolbar';
import WrongRoute from './WrongRoute';

class ProjectView extends Component {
  state = {
    startDate: null,
    endDate: null,
    currentZoom: "Days"
  };

  handleZoomChange = (zoom) => {
    this.setState({
      currentZoom: zoom
    });
  }

  handleRangeChange = (startDate, endDate) => {
    this.setState({
      startDate: startDate,
      endDate: endDate
    });
  }

  render() {
    if (!this.props.projectExists) {
      return (<WrongRoute/>);
    }

    const { currentZoom, startDate, endDate } = this.state;

    return (
      <div className="row content">
        <div className="box">
        <div className="row header">
          <GanttToolbar
            phid={this.props.match.params.phid}
            zoom={currentZoom}
            onZoomChange={this.handleZoomChange}
            onRangeChange={this.handleRangeChange}
          />
        </div>
        <div className="row content">
          <Gantt
            zoom={currentZoom}
            startDate={startDate}
            endDate={endDate}
          />
        </div></div>
        </div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const proj = state.projects.filter(proj => proj.phid === ownProps.match.params.phid);
  return {
    projectExists: proj.length !== 0
  };
}

function mapDispatchToProps(dispatch) {
  return {};
}

export default connect(mapStateToProps, mapDispatchToProps)(ProjectView);
