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
import moment from 'moment';

import PGanttNav from './PGanttNav';
import Gantt from './Gantt';
import GanttToolbar from './GanttToolbar';

class PGanttApp extends Component {
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

  constructor() {
    super();

    let start = moment();
    let end = moment();
    if (start.day() !== 0) {
      start = start.day(0); // previous Sunday
    }
    if (end.day() !== 6) {
      end = end.day(6); // next Saturday
    }

    this.state = {
      startDate: start,
      endDate: end,
      currentZoom: "Days"
    };
  }

  render() {
    const { currentZoom, startDate, endDate } = this.state;
    console.log("zoom", this.state);

    return (
      <div className="box">
        <div className="row header">
          <PGanttNav />
          <GanttToolbar
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
        </div>
        <div className="row footer">
        </div>
      </div>
    );
  }
}

export default PGanttApp;
