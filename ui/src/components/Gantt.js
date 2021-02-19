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
import { gantt } from 'dhtmlx-gantt';
import moment from 'moment';

const data = {
    data: [
        { id: 1, text: 'Task #1', start_date: '15-04-2019', duration: 3, progress: 0.6 },
        { id: 2, text: 'Task #2', start_date: '18-04-2019', duration: 3, progress: 0.4 }
    ],
    links: [
        { id: 1, source: 1, target: 2, type: '0' }
    ]
};

class Gantt extends Component {
  componentDidMount() {
    gantt.init(this.ganttContainer);
    gantt.config.show_tasks_outside_timescale = true;
    gantt.parse(data);
    this.initGanttDataProcessor();
  }

  setZoom(value) {
    switch (value) {
    case 'Days':
      gantt.config.min_column_width = 70;
      gantt.config.scale_unit = 'week';
      gantt.config.date_scale = 'Week %W';
      gantt.config.subscales = [
        { unit: 'day', step: 1, date: '%d %M' }
      ];
      gantt.config.scale_height = 60;
      break;
    case 'Months':
      gantt.config.min_column_width = 70;
      gantt.config.scale_unit = 'month';
      gantt.config.date_scale = '%F';
      gantt.config.scale_height = 60;
      gantt.config.subscales = [
        { unit:'week', step:1, date:'Week %W' }
      ];
      break;
    default:
      break;
    }
  }

  setRange(start, end) {
    let s = start;
    let e = end;
    if (start == null || end === null) {
      s = moment();
      e = moment();
      if (s.day() !== 0) {
        s = s.day(0); // previous Sunday
      }
      if (e.day() !== 6) {
        e = e.day(6); // next Saturday
      }

    }
    gantt.config.start_date = s.toDate();
    gantt.config.end_date = e.toDate();
  }

  initGanttDataProcessor() {
    const onDataUpdated = this.props.onDataUpdated;
    this.dataProcessor = gantt.createDataProcessor((entityType, action, item, id) => {
      return new Promise((resolve, reject) => {
        if (onDataUpdated) {
         onDataUpdated(entityType, action, item, id);
        }
        return resolve();
      });
    });
  }

  componentWillUnmount() {
    if (this.dataProcessor) {
      this.dataProcessor.destructor();
      this.dataProcessor = null;
    }
  }

  shouldComponentUpdate(nextProps) {
    if (this.props.zoom !== nextProps.zoom) {
      return true;
    }

    if (this.props.startDate !== nextProps.startDate) {
      return true;
    }

    if (this.props.endDate !== nextProps.endDate) {
      return true;
    }

    return false;
  }

  componentDidUpdate() {
    gantt.render();
  }

  render() {
    const { zoom, startDate, endDate } = this.props;
    this.setRange(startDate, endDate);
    this.setZoom(zoom);
    return (
      <div
        ref={ (input) => { this.ganttContainer = input; } }
        style={ { width: '100%', height: '100%' } }
      ></div>
    );
  }
}

export default Gantt;
